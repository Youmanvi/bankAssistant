"""
Retell AI webhook and integration endpoints.
"""

import json
import os
import logging
from fastapi import APIRouter, Request, HTTPException
from fastapi.responses import JSONResponse
from retell import Retell

logger = logging.getLogger(__name__)
router = APIRouter(tags=["Retell"])

# Initialize Retell client
retell = Retell(api_key=os.environ.get("RETELL_API_KEY", ""))


@router.post("/webhook")
async def handle_webhook(request: Request):
    """
    Handle webhooks from Retell AI.

    Processes call events:
    - call_started
    - call_ended
    - call_analyzed

    Args:
        request: The incoming webhook request

    Returns:
        JSONResponse: Confirmation response
    """
    try:
        post_data = await request.json()

        # Verify the signature
        valid_signature = retell.verify(
            json.dumps(post_data, separators=(",", ":"), ensure_ascii=False),
            api_key=str(os.environ.get("RETELL_API_KEY", "")),
            signature=str(request.headers.get("X-Retell-Signature")),
        )

        if not valid_signature:
            logger.warning(
                f"Unauthorized webhook: {post_data.get('event')} - "
                f"{post_data.get('data', {}).get('call_id')}"
            )
            raise HTTPException(status_code=401, detail="Unauthorized")

        event = post_data.get("event")
        call_id = post_data.get("data", {}).get("call_id")

        # Log the event
        logger.info(f"Received Retell event: {event} - Call ID: {call_id}")

        # Handle specific events
        if event == "call_started":
            logger.info(f"Call started: {call_id}")
        elif event == "call_ended":
            logger.info(f"Call ended: {call_id}")
        elif event == "call_analyzed":
            logger.info(f"Call analyzed: {call_id}")
        else:
            logger.warning(f"Unknown event: {event}")

        return JSONResponse(status_code=200, content={"received": True})

    except json.JSONDecodeError:
        logger.error("Invalid JSON in webhook request")
        raise HTTPException(status_code=400, detail="Invalid JSON")
    except Exception as e:
        logger.error(f"Error processing webhook: {e}")
        raise HTTPException(status_code=500, detail="Internal Server Error")
