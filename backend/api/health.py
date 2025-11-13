"""
Health check and status endpoints.
"""

from fastapi import APIRouter
from models import HealthResponse
from constants import APP_NAME, APP_VERSION

router = APIRouter(tags=["Health"])


@router.get("/health", response_model=HealthResponse)
async def health_check() -> HealthResponse:
    """
    Health check endpoint.

    Returns:
        HealthResponse: Current health status
    """
    return HealthResponse(
        status="healthy",
        version=APP_VERSION,
    )


@router.get("/status", response_model=dict)
async def status() -> dict:
    """
    Get application status.

    Returns:
        dict: Current status information
    """
    return {
        "app": APP_NAME,
        "version": APP_VERSION,
        "status": "running",
    }
