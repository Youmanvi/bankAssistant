"""
Users API endpoints.
"""

import logging
from fastapi import APIRouter, HTTPException
from db.db_service import db_service

logger = logging.getLogger(__name__)
router = APIRouter(tags=["Users"], prefix="/users")


@router.get("/")
async def list_users():
    """
    Get all users.

    Returns:
        dict: All users
    """
    try:
        users = db_service.get_users()
        return {"users": users, "count": len(users)}
    except Exception as e:
        logger.error(f"Error retrieving users: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve users")


@router.get("/{user_id}")
async def get_user(user_id: str):
    """
    Get a specific user by ID (phone number).

    Args:
        user_id: The user ID (phone number)

    Returns:
        dict: User details
    """
    try:
        user = db_service.get_user(user_id)
        if not user:
            raise HTTPException(status_code=404, detail="User not found")
        return {"user_id": user_id, "details": user}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error retrieving user {user_id}: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve user")


@router.get("/{user_id}/accounts")
async def get_user_accounts(user_id: str):
    """
    Get all accounts for a specific user.

    Args:
        user_id: The user ID (phone number)

    Returns:
        dict: User's accounts
    """
    try:
        user = db_service.get_user(user_id)
        if not user:
            raise HTTPException(status_code=404, detail="User not found")

        account_ids = user.get("accounts", [])
        accounts = {}

        for account_id in account_ids:
            account = db_service.get_account(account_id)
            if account:
                accounts[account_id] = account

        return {
            "user_id": user_id,
            "account_count": len(accounts),
            "accounts": accounts,
        }
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error retrieving accounts for user {user_id}: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve accounts")


@router.get("/{user_id}/profile")
async def get_user_profile(user_id: str):
    """
    Get the profile information for a specific user.

    Args:
        user_id: The user ID (phone number)

    Returns:
        dict: User profile details (excluding sensitive info)
    """
    try:
        user = db_service.get_user(user_id)
        if not user:
            raise HTTPException(status_code=404, detail="User not found")

        # Return non-sensitive user information
        return {
            "user_id": user_id,
            "name": user.get("name"),
            "email": user.get("email"),
            "phone": user.get("phone"),
            "address": user.get("address"),
        }
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error retrieving profile for user {user_id}: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve profile")
