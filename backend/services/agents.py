"""
Agent Service for multi-agent orchestration using OpenAI Swarm.
"""

from typing import Dict, List, Generator
import logging
from swarm import Swarm
from openai import OpenAI
from agents.triage_agent import TriageAgent
from agents.accounts_agent import (
    AccountsAgent,
    handle_account_balance,
    retrieve_bank_statement,
)
from agents.payments_agent import (
    PaymentsAgent,
    transfer_funds,
    schedule_payment,
    cancel_payment,
)
from agents.applications_agent import (
    ApplicationsAgent,
    apply_for_loan,
    apply_for_credit_card,
)


logger = logging.getLogger(__name__)


class AgentService:
    """Orchestrates multiple agents using OpenAI Swarm framework."""

    OPENAI_MODEL = "gpt-4o-mini"

    def __init__(self):
        """Initialize the Agent Service with all available agents."""
        self.client = Swarm()
        self.openai_client = OpenAI()
        self.messages = []

        # Initialize agents with their respective functions
        self.triage_agent = TriageAgent(
            [
                self.transfer_to_accounts,
                self.transfer_to_payments,
                self.transfer_to_applications,
            ]
        )

        self.accounts_agent = AccountsAgent(
            transfer_to_payments=self.transfer_to_payments,
            handle_account_balance=handle_account_balance,
            retrieve_bank_statement=retrieve_bank_statement,
        )

        self.payments_agent = PaymentsAgent(
            transfer_back_to_triage=self.transfer_back_to_triage,
            transfer_funds=transfer_funds,
            schedule_payment=schedule_payment,
            cancel_payment=cancel_payment,
        )

        self.applications_agent = ApplicationsAgent(
            transfer_back_to_triage=self.transfer_back_to_triage,
            apply_for_loan=apply_for_loan,
            apply_for_credit_card=apply_for_credit_card,
        )

        self.current_agent = self.triage_agent

    # ========================================================================
    # Agent Transfer Methods
    # ========================================================================

    def transfer_to_accounts(self, context_variables: Dict, user_message: str):
        """Transfer control to the Accounts Agent."""
        logger.info("Transferring to Accounts Agent")
        self.current_agent = self.accounts_agent
        return self.accounts_agent

    def transfer_to_payments(self, context_variables: Dict, user_message: str):
        """Transfer control to the Payments Agent."""
        logger.info("Transferring to Payments Agent")
        self.current_agent = self.payments_agent
        return self.payments_agent

    def transfer_to_applications(self, context_variables: Dict, user_message: str):
        """Transfer control to the Applications Agent."""
        logger.info("Transferring to Applications Agent")
        self.current_agent = self.applications_agent
        return self.applications_agent

    def transfer_back_to_triage(self, context_variables: Dict, response: str):
        """Transfer control back to the Triage Agent."""
        logger.info("Transferring back to Triage Agent")
        self.current_agent = self.triage_agent
        return self.triage_agent

    # ========================================================================
    # Swarm Execution
    # ========================================================================

    def run(
        self, messages: List[Dict[str, str]], stream: bool = False
    ) -> Generator:
        """
        Execute the swarm with the provided messages.

        Args:
            messages: List of message dicts for the LLM
            stream: Whether to stream the response

        Yields:
            Response chunks from the swarm
        """
        self.messages.extend(messages)
        logger.debug(f"Current message history: {self.messages}")

        response = self.client.run(
            agent=self.current_agent, messages=self.messages, stream=stream
        )

        return response

    def reset_messages(self) -> None:
        """Reset the message history."""
        self.messages = []

    def reset_agent(self) -> None:
        """Reset to the triage agent."""
        self.current_agent = self.triage_agent

    def get_current_agent_name(self) -> str:
        """Get the name of the current agent."""
        agent_names = {
            self.triage_agent: "TriageAgent",
            self.accounts_agent: "AccountsAgent",
            self.payments_agent: "PaymentsAgent",
            self.applications_agent: "ApplicationsAgent",
        }
        return agent_names.get(self.current_agent, "UnknownAgent")
