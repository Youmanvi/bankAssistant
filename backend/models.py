"""
Data models and Pydantic schemas for the API and internal use.
"""

from typing import Any, Dict, List, Literal, Optional, Union
from pydantic import BaseModel


# ============================================================================
# Internal Data Models (Application State)
# ============================================================================

class User(BaseModel):
    """User profile information."""
    name: str
    accounts: List[str]
    ssn: str
    address: str
    date_of_birth: str
    email: str
    phone: str


class Account(BaseModel):
    """Bank account information."""
    balance: float
    statements: Dict[str, str]


class Payment(BaseModel):
    """Payment transaction information."""
    from_account: str
    to_account: str
    amount: float
    date: str
    status: str


# ============================================================================
# Retell AI WebSocket Models (Voice Interaction)
# ============================================================================

class Utterance(BaseModel):
    """Single message in a conversation."""
    role: Literal["agent", "user", "system"]
    content: str


# Retell → Server
class RetellMessage(BaseModel):
    """Base class for Retell messages."""
    interaction_type: str


class PingPongRequest(RetellMessage):
    """Keep-alive ping from Retell."""
    interaction_type: Literal["ping_pong"]
    timestamp: int


class CallDetailsRequest(RetellMessage):
    """Call setup details from Retell."""
    interaction_type: Literal["call_details"]
    call: Dict[str, Any]


class UpdateOnlyRequest(RetellMessage):
    """Transcript update without response."""
    interaction_type: Literal["update_only"]
    transcript: List[Utterance]


class ResponseRequiredRequest(RetellMessage):
    """User input requiring LLM response."""
    interaction_type: Literal["reminder_required", "response_required"]
    response_id: int
    transcript: List[Utterance]


# Server → Retell
class RetellResponse(BaseModel):
    """Base class for Retell responses."""
    response_type: str


class ConfigResponse(RetellResponse):
    """Server configuration for Retell."""
    response_type: Literal["config"] = "config"
    response_id: Optional[int] = None
    config: Dict[str, bool] = {
        "auto_reconnect": True,
        "call_details": True,
    }


class PingPongResponse(RetellResponse):
    """Keep-alive pong response."""
    response_type: Literal["ping_pong"] = "ping_pong"
    timestamp: int


class LLMResponse(RetellResponse):
    """LLM response to send to Retell."""
    response_type: Literal["response"] = "response"
    response_id: int
    content: str
    content_complete: bool
    end_call: Optional[bool] = False
    transfer_number: Optional[str] = None


# ============================================================================
# API Request/Response Models
# ============================================================================

class TransferRequest(BaseModel):
    """Request to transfer funds."""
    from_account: str
    to_account: str
    amount: float


class TransferResponse(BaseModel):
    """Response from fund transfer."""
    status: str
    payment_id: str
    from_account: str
    to_account: str
    amount: float


class LoanApplicationRequest(BaseModel):
    """Loan application request."""
    user_id: str
    loan_amount: float
    loan_purpose: str
    term_years: int


class LoanApplicationResponse(BaseModel):
    """Loan application response."""
    status: str
    application_id: str
    user_id: str
    loan_amount: float
    loan_purpose: str
    term_years: int
    message: str


class CreditCardApplicationRequest(BaseModel):
    """Credit card application request."""
    user_id: str
    card_type: str
    credit_limit: float


class CreditCardApplicationResponse(BaseModel):
    """Credit card application response."""
    status: str
    application_id: str
    user_id: str
    card_type: str
    credit_limit: float
    message: str


class HealthResponse(BaseModel):
    """Health check response."""
    status: str
    version: str


class StatusResponse(BaseModel):
    """Application status response."""
    app: str
    version: str
    status: str


class ErrorResponse(BaseModel):
    """Standard error response."""
    error: str
    detail: Optional[str] = None
    status_code: int


# ============================================================================
# Type Unions
# ============================================================================

RetellRequest = Union[
    PingPongRequest,
    CallDetailsRequest,
    UpdateOnlyRequest,
    ResponseRequiredRequest,
]

RetellMessageType = Union[
    ConfigResponse,
    PingPongResponse,
    LLMResponse,
]
