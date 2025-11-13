"""
Helper functions for validation and data access.
"""

import logging
from typing import List
from datetime import datetime
from database.db_service import db_service

logger = logging.getLogger(__name__)


def validate_account_id(account_id: str) -> bool:
    """
    Validate that an account ID exists in the database.

    Args:
        account_id: The account ID to validate

    Returns:
        bool: True if account exists, False otherwise
    """
    account = db_service.get_account(account_id)
    if account:
        logger.info(f"Account ID {account_id} is valid.")
        return True
    logger.warning(f"Account ID {account_id} is invalid.")
    return False


def validate_payment_id(payment_id: str) -> bool:
    """
    Validate that a payment ID exists in the database.

    Args:
        payment_id: The payment ID to validate

    Returns:
        bool: True if payment exists, False otherwise
    """
    payment = db_service.get_payment(payment_id)
    if payment:
        logger.info(f"Payment ID {payment_id} is valid.")
        return True
    logger.warning(f"Payment ID {payment_id} is invalid.")
    return False


def generate_payment_id() -> str:
    """
    Generate a new unique payment ID.

    Returns:
        str: New payment ID in format PAY###
    """
    payment_number = len(db_service.get_payments()) + 1
    payment_id = f"PAY{payment_number:03d}"
    logger.info(f"Generated new payment ID: {payment_id}")
    return payment_id


def validate_amount(amount: float) -> bool:
    """
    Validate that an amount is positive.

    Args:
        amount: The amount to validate

    Returns:
        bool: True if amount is positive, False otherwise
    """
    if amount <= 0:
        logger.warning(f"Invalid amount: {amount}. Amount must be positive.")
        return False
    logger.info(f"Amount {amount} is valid.")
    return True


def get_user_accounts(user_id: str) -> List[str]:
    """
    Get all account IDs for a user.

    Args:
        user_id: The user ID (phone number)

    Returns:
        List[str]: List of account IDs for the user
    """
    accounts = db_service.get_user_accounts(user_id)
    logger.info(f"User {user_id} has accounts: {accounts}")
    return accounts


# Legacy functions for backward compatibility
def get_db():
    """Legacy function: Get the entire database."""
    return db_service.get_db()


def set_db(new_db):
    """Legacy function: Set the entire database."""
    return db_service.set_db(new_db)
