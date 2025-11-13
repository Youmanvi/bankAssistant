"""
Loan and Credit Card Applications API endpoints.
"""

import logging
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

logger = logging.getLogger(__name__)
router = APIRouter(tags=["Applications"], prefix="/applications")


class LoanApplicationRequest(BaseModel):
    """Loan application request schema."""
    user_id: str
    loan_amount: float
    loan_purpose: str
    term_years: int


class CreditCardApplicationRequest(BaseModel):
    """Credit card application request schema."""
    user_id: str
    card_type: str
    credit_limit: float


@router.post("/loan")
async def apply_for_loan(request: LoanApplicationRequest):
    """
    Submit a loan application.

    Args:
        request: Loan application details

    Returns:
        dict: Application status and details
    """
    try:
        # Validate input
        if request.loan_amount <= 0:
            raise HTTPException(status_code=400, detail="Loan amount must be positive")
        if request.term_years <= 0:
            raise HTTPException(status_code=400, detail="Term must be positive")

        # Generate application ID
        application_id = f"LOAN_{request.user_id}_{__import__('time').time()}"

        return {
            "status": "submitted",
            "application_id": application_id,
            "user_id": request.user_id,
            "loan_amount": request.loan_amount,
            "loan_purpose": request.loan_purpose,
            "term_years": request.term_years,
            "message": "Your loan application has been submitted successfully.",
        }

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error processing loan application: {e}")
        raise HTTPException(
            status_code=500, detail="Failed to process loan application"
        )


@router.post("/credit-card")
async def apply_for_credit_card(request: CreditCardApplicationRequest):
    """
    Submit a credit card application.

    Args:
        request: Credit card application details

    Returns:
        dict: Application status and details
    """
    try:
        # Validate input
        if request.credit_limit <= 0:
            raise HTTPException(status_code=400, detail="Credit limit must be positive")

        # Generate application ID
        application_id = f"CARD_{request.user_id}_{__import__('time').time()}"

        return {
            "status": "submitted",
            "application_id": application_id,
            "user_id": request.user_id,
            "card_type": request.card_type,
            "credit_limit": request.credit_limit,
            "message": "Your credit card application has been submitted successfully.",
        }

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error processing credit card application: {e}")
        raise HTTPException(
            status_code=500, detail="Failed to process credit card application"
        )


@router.get("/{application_id}")
async def get_application_status(application_id: str):
    """
    Get the status of an application.

    Args:
        application_id: The application ID

    Returns:
        dict: Application status
    """
    try:
        # For now, return a mock status
        # In a real implementation, this would query a database
        return {
            "application_id": application_id,
            "status": "under_review",
            "message": "Your application is being reviewed.",
        }
    except Exception as e:
        logger.error(f"Error retrieving application status: {e}")
        raise HTTPException(
            status_code=500, detail="Failed to retrieve application status"
        )
