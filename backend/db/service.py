"""
Database service for managing banking data.
Provides CRUD operations for users, accounts, and payments.
"""

from typing import Dict, Any, Optional, List
from datetime import datetime, timedelta
import random
from .mock_data import MOCK_DATABASE, MOCK_CALLS_SAMPLE_DATA


class DatabaseService:
    """Service for managing the in-memory banking database."""

    def __init__(self):
        """Initialize the database with mock data."""
        self.db: Dict[str, Any] = self._deep_copy(MOCK_DATABASE)
        self.calls_db: Dict[str, Any] = {}
        self._initialize_calls()

    def _deep_copy(self, data: Dict) -> Dict:
        """Create a deep copy of the database structure."""
        import copy
        return copy.deepcopy(data)

    def _initialize_calls(self) -> None:
        """Initialize the calls database with sample data."""
        for call in MOCK_CALLS_SAMPLE_DATA:
            call_with_timestamp = {
                **call,
                "time": (
                    datetime.now() - timedelta(minutes=random.randint(0, 120))
                ).isoformat(),
            }
            self.calls_db[call["id"]] = call_with_timestamp

    # ========================================================================
    # Database Getters
    # ========================================================================

    def get_db(self) -> Dict[str, Any]:
        """
        Retrieve the current state of the banking database.

        Returns:
            dict: The current banking database dictionary.
        """
        return self.db

    def get_users(self) -> Dict[str, Any]:
        """Get all users from the database."""
        return self.db.get("users", {})

    def get_user(self, user_id: str) -> Optional[Dict[str, Any]]:
        """Get a specific user by ID."""
        return self.db.get("users", {}).get(user_id)

    def get_accounts(self) -> Dict[str, Any]:
        """Get all accounts from the database."""
        return self.db.get("accounts", {})

    def get_account(self, account_id: str) -> Optional[Dict[str, Any]]:
        """Get a specific account by ID."""
        return self.db.get("accounts", {}).get(account_id)

    def get_payments(self) -> Dict[str, Any]:
        """Get all payments from the database."""
        return self.db.get("payments", {})

    def get_payment(self, payment_id: str) -> Optional[Dict[str, Any]]:
        """Get a specific payment by ID."""
        return self.db.get("payments", {}).get(payment_id)

    def get_all_calls(self) -> Dict[str, Any]:
        """Get all calls from the calls database."""
        return self.calls_db

    def get_call(self, call_id: str) -> Optional[Dict[str, Any]]:
        """Get a specific call by its ID."""
        return self.calls_db.get(call_id)

    # ========================================================================
    # Database Setters and Updates
    # ========================================================================

    def set_db(self, new_db: Dict[str, Any]) -> bool:
        """
        Update the banking database with a new state.

        Args:
            new_db (dict): The new banking database dictionary.

        Returns:
            bool: True if the database was successfully updated.
        """
        self.db = self._deep_copy(new_db)
        return True

    def update_user(self, user_id: str, user_data: Dict[str, Any]) -> bool:
        """Update a user's information."""
        if "users" not in self.db:
            self.db["users"] = {}
        self.db["users"][user_id] = user_data
        return True

    def update_account(self, account_id: str, account_data: Dict[str, Any]) -> bool:
        """Update an account's information."""
        if "accounts" not in self.db:
            self.db["accounts"] = {}
        self.db["accounts"][account_id] = account_data
        return True

    def update_payment(self, payment_id: str, payment_data: Dict[str, Any]) -> bool:
        """Update a payment's information."""
        if "payments" not in self.db:
            self.db["payments"] = {}
        self.db["payments"][payment_id] = payment_data
        return True

    def update_call(self, call_id: str, call_data: Dict[str, Any]) -> bool:
        """Update or insert a call record."""
        if call_id in self.calls_db:
            self.calls_db[call_id].update(call_data)
        else:
            self.calls_db[call_id] = call_data
        return True

    # ========================================================================
    # Database Deleters
    # ========================================================================

    def delete_call(self, call_id: str) -> bool:
        """Delete a call from the database."""
        if call_id in self.calls_db:
            del self.calls_db[call_id]
            return True
        return False

    def delete_payment(self, payment_id: str) -> bool:
        """Delete a payment from the database."""
        payments = self.db.get("payments", {})
        if payment_id in payments:
            del payments[payment_id]
            return True
        return False

    # ========================================================================
    # Business Logic Helpers
    # ========================================================================

    def get_user_accounts(self, user_id: str) -> List[str]:
        """Get all account IDs for a specific user."""
        user = self.get_user(user_id)
        if user:
            return user.get("accounts", [])
        return []

    def transfer_funds(
        self, from_account: str, to_account: str, amount: float
    ) -> Optional[str]:
        """
        Transfer funds from one account to another.

        Args:
            from_account: Source account ID
            to_account: Destination account ID
            amount: Amount to transfer

        Returns:
            payment_id if successful, None otherwise
        """
        from_acc = self.get_account(from_account)
        to_acc = self.get_account(to_account)

        # Validate accounts exist
        if not from_acc or not to_acc:
            return None

        # Validate sufficient funds
        if from_acc.get("balance", 0) < amount:
            return None

        # Perform transfer
        from_acc["balance"] -= amount
        to_acc["balance"] += amount

        # Create payment record
        payment_id = self._generate_payment_id()
        payment_data = {
            "from_account": from_account,
            "to_account": to_account,
            "amount": amount,
            "date": datetime.now().strftime("%Y-%m-%d"),
            "status": "Completed",
        }
        self.update_payment(payment_id, payment_data)

        return payment_id

    def _generate_payment_id(self) -> str:
        """Generate a new unique payment ID."""
        payments = self.get_payments()
        max_num = 0
        for payment_id in payments.keys():
            try:
                num = int(payment_id.replace("PAY", ""))
                max_num = max(max_num, num)
            except ValueError:
                pass
        return f"PAY{max_num + 1:03d}"

    # ========================================================================
    # Reset Functions (for testing)
    # ========================================================================

    def reset(self) -> None:
        """Reset database to initial state."""
        self.db = self._deep_copy(MOCK_DATABASE)
        self.calls_db.clear()
        self._initialize_calls()


# Global database service instance
db_service = DatabaseService()
