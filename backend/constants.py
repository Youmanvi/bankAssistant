"""
Application constants, enums, and configuration values.
"""

from enum import Enum

# ============================================================================
# Application Info
# ============================================================================

APP_NAME = "AI Banking Assistant"
APP_VERSION = "0.1.0"
API_PREFIX = "/api/v1"

# ============================================================================
# Enums
# ============================================================================


class WebSocketEvent(str, Enum):
    """Admin dashboard WebSocket events."""
    GET_DB = "get_db"
    GET_CALLS = "get_calls"
    GET_ALL_DBS = "get_all_dbs"
    DB_RESPONSE = "db_response"
    CALLS_RESPONSE = "calls_response"
    COMBINED_RESPONSE = "combined_response"


class Agent(str, Enum):
    """Available AI agents."""
    TRIAGE = "triage"
    ACCOUNTS = "accounts"
    PAYMENTS = "payments"
    APPLICATIONS = "applications"


class PaymentStatus(str, Enum):
    """Payment statuses."""
    PENDING = "Pending"
    COMPLETED = "Completed"
    SCHEDULED = "Scheduled"
    CANCELED = "Canceled"


class RetellEventType(str, Enum):
    """Retell AI events."""
    CALL_STARTED = "call_started"
    CALL_ENDED = "call_ended"
    CALL_ANALYZED = "call_analyzed"


# ============================================================================
# CORS Configuration
# ============================================================================

CORS_ORIGINS = [
    "http://localhost:3000",
    "http://localhost:3001",
    "http://localhost:8000",
]

CORS_ALLOW_CREDENTIALS = True
CORS_ALLOW_METHODS = ["*"]
CORS_ALLOW_HEADERS = ["*"]

# ============================================================================
# Server Configuration
# ============================================================================

DEFAULT_HOST = "0.0.0.0"
DEFAULT_PORT = 8000
DEFAULT_WORKERS = 1

# ============================================================================
# Logging
# ============================================================================

LOG_LEVEL = "INFO"
LOG_FORMAT = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"

# ============================================================================
# Timeouts & Limits
# ============================================================================

WEBSOCKET_TIMEOUT = 300  # 5 minutes
MAX_RESPONSE_ID = 999999
PDF_GENERATION_TIMEOUT = 30  # seconds

# ============================================================================
# Retell AI Configuration
# ============================================================================

RETELL_CONFIG = {
    "auto_reconnect": True,
    "call_details": True,
}
