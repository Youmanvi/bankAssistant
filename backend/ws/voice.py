"""
LLM WebSocket endpoint for real-time voice interaction with Retell AI.
"""

import asyncio
import logging
from concurrent.futures import TimeoutError as ConnectionTimeoutError
from fastapi import APIRouter, WebSocket, WebSocketDisconnect
from models import ResponseRequiredRequest, ConfigResponse
from constants import RETELL_DEFAULT_CONFIG
from services.llm_service import LlmService
from db.db_service import db_service

logger = logging.getLogger(__name__)
router = APIRouter(tags=["LLM"])


@router.websocket("/llm-websocket/{call_id}")
async def websocket_handler(websocket: WebSocket, call_id: str):
    """
    WebSocket endpoint for real-time LLM interaction with Retell AI.

    Handles:
    - call_details: Initialize LLM client with user info
    - ping_pong: Keep-alive messages
    - update_only: State updates without response needed
    - response_required: User input requiring LLM response
    - reminder_required: Timeout reminder needed

    Args:
        websocket: The WebSocket connection
        call_id: Unique identifier for the call
    """
    try:
        await websocket.accept()
        llm_client = None
        response_id = 0

        # Send config to Retell
        config = ConfigResponse(
            response_type="config",
            config=RETELL_DEFAULT_CONFIG,
        )
        await websocket.send_json(config.model_dump())

        async def handle_message(request_json: dict):
            nonlocal response_id, llm_client

            interaction_type = request_json.get("interaction_type")
            logger.debug(f"Received interaction: {interaction_type}")

            # Handle call_details: Initialize LLM client
            if interaction_type == "call_details":
                try:
                    call_data = request_json.get("call", {})
                    from_number = call_data.get("from_number", "")

                    # Format phone number
                    if from_number:
                        formatted_number = (
                            "+1-" + from_number[2:5] + "-" +
                            from_number[5:8] + "-" + from_number[8:]
                        )
                        # Get user info from database
                        user = db_service.get_user(formatted_number)
                        if user:
                            llm_client = LlmService()
                            llm_client.set_user_info(user.get("name"), formatted_number)

                            # Send greeting
                            first_event = llm_client.draft_begin_message()
                            await websocket.send_json(first_event.model_dump())
                        else:
                            logger.warning(f"User not found: {formatted_number}")
                            llm_client = LlmService()
                            llm_client.set_user_info("there")
                            first_event = llm_client.draft_begin_message()
                            await websocket.send_json(first_event.model_dump())
                except Exception as e:
                    logger.error(f"Error handling call_details: {e}")

            # Handle ping_pong: Keep-alive
            elif interaction_type == "ping_pong":
                await websocket.send_json(
                    {
                        "response_type": "ping_pong",
                        "timestamp": request_json.get("timestamp"),
                    }
                )

            # Handle update_only: No response needed
            elif interaction_type == "update_only":
                pass

            # Handle response_required or reminder_required
            elif interaction_type in ["response_required", "reminder_required"]:
                response_id = request_json.get("response_id", 0)
                logger.info(
                    f"Response required: interaction_type={interaction_type}, "
                    f"response_id={response_id}"
                )

                if not llm_client:
                    logger.warning("LLM client not initialized")
                    return

                # Create request object
                try:
                    request_obj = ResponseRequiredRequest(
                        interaction_type=interaction_type,
                        response_id=response_id,
                        transcript=request_json.get("transcript", []),
                    )

                    # Stream response from LLM
                    async for event in llm_client.draft_response(request_obj):
                        await websocket.send_json(event.model_dump())
                        if request_obj.response_id < response_id:
                            logger.debug("New response needed, abandoning current")
                            break

                except Exception as e:
                    logger.error(f"Error generating response: {e}")

        # Main WebSocket loop
        async for data in websocket.iter_json():
            asyncio.create_task(handle_message(data))

    except WebSocketDisconnect:
        logger.info(f"LLM WebSocket disconnected: {call_id}")
    except ConnectionTimeoutError:
        logger.warning(f"Connection timeout: {call_id}")
    except Exception as e:
        logger.error(f"Error in LLM WebSocket: {e}")
        try:
            await websocket.close(1011, "Server error")
        except Exception as close_error:
            logger.error(f"Error closing WebSocket: {close_error}")
    finally:
        logger.info(f"LLM WebSocket closed: {call_id}")
