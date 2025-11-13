"""
Accounts API endpoints.
"""

import logging
from fastapi import APIRouter, HTTPException
from db.db_service import db_service

logger = logging.getLogger(__name__)
router = APIRouter(tags=["Accounts"], prefix="/accounts")


@router.get("/")
async def list_accounts():
    """
    Get all accounts.

    Returns:
        dict: All accounts
    """
    try:
        accounts = db_service.get_accounts()
        return {"accounts": accounts, "count": len(accounts)}
    except Exception as e:
        logger.error(f"Error retrieving accounts: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve accounts")


@router.get("/{account_id}")
async def get_account(account_id: str):
    """
    Get a specific account by ID.

    Args:
        account_id: The account ID

    Returns:
        dict: Account details
    """
    try:
        account = db_service.get_account(account_id)
        if not account:
            raise HTTPException(status_code=404, detail="Account not found")
        return {"account_id": account_id, "details": account}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error retrieving account {account_id}: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve account")


@router.get("/{account_id}/balance")
async def get_account_balance(account_id: str):
    """
    Get the balance of a specific account.

    Args:
        account_id: The account ID

    Returns:
        dict: Account balance
    """
    try:
        account = db_service.get_account(account_id)
        if not account:
            raise HTTPException(status_code=404, detail="Account not found")
        return {
            "account_id": account_id,
            "balance": account.get("balance", 0.0),
        }
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error retrieving balance for {account_id}: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve balance")


@router.get("/{account_id}/statements")
async def get_account_statements(account_id: str):
    """
    Get statements for a specific account.

    Args:
        account_id: The account ID

    Returns:
        dict: Account statements
    """
    try:
        account = db_service.get_account(account_id)
        if not account:
            raise HTTPException(status_code=404, detail="Account not found")
        return {
            "account_id": account_id,
            "statements": account.get("statements", {}),
        }
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error retrieving statements for {account_id}: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve statements")
