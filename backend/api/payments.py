"""
Payments API endpoints.
"""

import logging
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from db.db_service import db_service

logger = logging.getLogger(__name__)
router = APIRouter(tags=["Payments"], prefix="/payments")


class PaymentRequest(BaseModel):
    """Payment request schema."""
    from_account: str
    to_account: str
    amount: float


@router.get("/")
async def list_payments():
    """
    Get all payments.

    Returns:
        dict: All payments
    """
    try:
        payments = db_service.get_payments()
        return {"payments": payments, "count": len(payments)}
    except Exception as e:
        logger.error(f"Error retrieving payments: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve payments")


@router.get("/{payment_id}")
async def get_payment(payment_id: str):
    """
    Get a specific payment by ID.

    Args:
        payment_id: The payment ID

    Returns:
        dict: Payment details
    """
    try:
        payment = db_service.get_payment(payment_id)
        if not payment:
            raise HTTPException(status_code=404, detail="Payment not found")
        return {"payment_id": payment_id, "details": payment}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error retrieving payment {payment_id}: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve payment")


@router.post("/transfer")
async def transfer_funds(request: PaymentRequest):
    """
    Transfer funds between two accounts.

    Args:
        request: Payment request with from_account, to_account, and amount

    Returns:
        dict: Payment details with payment_id
    """
    try:
        # Validate accounts exist
        from_account = db_service.get_account(request.from_account)
        to_account = db_service.get_account(request.to_account)

        if not from_account:
            raise HTTPException(status_code=404, detail="Source account not found")
        if not to_account:
            raise HTTPException(status_code=404, detail="Destination account not found")

        # Validate amount
        if request.amount <= 0:
            raise HTTPException(status_code=400, detail="Amount must be positive")

        # Check sufficient funds
        if from_account.get("balance", 0) < request.amount:
            raise HTTPException(status_code=400, detail="Insufficient funds")

        # Perform transfer
        payment_id = db_service.transfer_funds(
            request.from_account, request.to_account, request.amount
        )

        if not payment_id:
            raise HTTPException(status_code=500, detail="Transfer failed")

        return {
            "status": "success",
            "payment_id": payment_id,
            "from_account": request.from_account,
            "to_account": request.to_account,
            "amount": request.amount,
        }

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error processing transfer: {e}")
        raise HTTPException(status_code=500, detail="Failed to process transfer")
