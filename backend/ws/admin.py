"""
Admin dashboard WebSocket and API endpoints.
"""

import logging
from fastapi import APIRouter, WebSocket, WebSocketDisconnect, Query
from fastapi.responses import JSONResponse
from ws.connection_manager import connection_manager
from db.db_service import db_service
from constants import WebSocketEventType

logger = logging.getLogger(__name__)
router = APIRouter(tags=["Admin"], prefix="/admin")


@router.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket, client_id: str = Query(None)):
    """
    WebSocket endpoint for the admin dashboard.

    Handles client communication for real-time updates:
    - get_db: Retrieve database state
    - get_calls: Retrieve all call records
    - get_all_dbs: Retrieve combined database and calls

    Args:
        websocket: The WebSocket connection
        client_id: Unique identifier for the client
    """
    if not client_id:
        logger.warning("WebSocket connection attempt without client_id")
        await websocket.close(code=4001, reason="client_id required")
        return

    # Register the connection
    await connection_manager.connect(websocket, client_id)
    logger.info(f"Admin client connected: {client_id}")

    try:
        while True:
            data = await websocket.receive_json()
            event = data.get("event")
            logger.debug(f"Received admin event: {event}")

            try:
                if event == WebSocketEventType.GET_DB.value:
                    # Send database state
                    db_state = db_service.get_db()
                    message = {
                        "event": WebSocketEventType.DB_RESPONSE.value,
                        "data": db_state,
                    }
                    await connection_manager.send_personal_message(message, websocket)

                elif event == WebSocketEventType.GET_CALLS.value:
                    # Send all calls
                    calls = db_service.get_all_calls()
                    message = {
                        "event": WebSocketEventType.CALLS_RESPONSE.value,
                        "data": calls,
                    }
                    await connection_manager.send_personal_message(message, websocket)

                elif event == WebSocketEventType.GET_ALL_DBS.value:
                    # Send combined database and calls
                    db_state = db_service.get_db()
                    calls = db_service.get_all_calls()
                    message = {
                        "event": WebSocketEventType.COMBINED_RESPONSE.value,
                        "calls": calls,
                        "db": db_state,
                    }
                    await connection_manager.send_personal_message(message, websocket)

                else:
                    logger.warning(f"Unknown event: {event}")

            except Exception as e:
                logger.error(f"Error handling event {event}: {e}")
                error_message = {
                    "event": "error",
                    "error": str(e),
                }
                await connection_manager.send_personal_message(
                    error_message, websocket
                )

    except WebSocketDisconnect:
        logger.info(f"Admin client disconnected: {client_id}")
        await connection_manager.disconnect(client_id)
    except Exception as e:
        logger.error(f"WebSocket error: {e}")
        await connection_manager.disconnect(client_id)


@router.get("/database")
async def get_database():
    """
    REST endpoint to get the current database state.

    Returns:
        JSONResponse: Current database state
    """
    try:
        db_state = db_service.get_db()
        return JSONResponse(content=db_state)
    except Exception as e:
        logger.error(f"Error retrieving database: {e}")
        return JSONResponse(
            status_code=500,
            content={"error": "Failed to retrieve database"},
        )


@router.get("/calls")
async def get_calls():
    """
    REST endpoint to get all call records.

    Returns:
        JSONResponse: All call records
    """
    try:
        calls = db_service.get_all_calls()
        return JSONResponse(content=calls)
    except Exception as e:
        logger.error(f"Error retrieving calls: {e}")
        return JSONResponse(
            status_code=500,
            content={"error": "Failed to retrieve calls"},
        )


@router.get("/status")
async def get_status():
    """
    Get the status of admin services.

    Returns:
        JSONResponse: Service status
    """
    return JSONResponse(
        content={
            "status": "healthy",
            "active_connections": connection_manager.get_active_connections_count(),
        }
    )
