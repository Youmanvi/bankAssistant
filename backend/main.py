"""
AI Banking Assistant - Main Application Entry Point

A voice-driven, multi-agent AI banking assistant powered by:
- FastAPI: Modern async web framework
- OpenAI Swarm: Multi-agent orchestration
- Retell AI: Voice interaction
- Pinata: IPFS decentralized storage
"""

import logging
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from dotenv import load_dotenv

from config import settings
from constants import (
    APP_NAME,
    APP_VERSION,
    CORS_ORIGINS,
    CORS_ALLOW_CREDENTIALS,
    CORS_ALLOW_METHODS,
    CORS_ALLOW_HEADERS,
)

# Import API routes (for Go to call)
from api import health, retell, accounts, payments, applications, users

# Import WebSocket handlers (internal communication)
from ws import voice, admin

from constants import API_PREFIX

# Load environment variables
load_dotenv(override=True)

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Create FastAPI app
app = FastAPI(
    title=APP_NAME,
    description="Voice-driven multi-agent AI banking assistant",
    version=APP_VERSION,
    docs_url="/docs",
    redoc_url="/redoc",
)

# Configure CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=CORS_ORIGINS,
    allow_credentials=CORS_ALLOW_CREDENTIALS,
    allow_methods=CORS_ALLOW_METHODS,
    allow_headers=CORS_ALLOW_HEADERS,
)


# ============================================================================
# Register API Routes (for Go to call)
# ============================================================================

# Health check & status
app.include_router(health.router)

# Retell AI webhook
app.include_router(retell.router)

# REST API v1 endpoints (Go orchestrator will call these)
app.include_router(accounts.router, prefix=API_PREFIX)
app.include_router(payments.router, prefix=API_PREFIX)
app.include_router(applications.router, prefix=API_PREFIX)
app.include_router(users.router, prefix=API_PREFIX)


# ============================================================================
# Register WebSocket Handlers (Internal Communication)
# ============================================================================

# Voice interaction with Retell AI
app.include_router(voice.router)

# Admin dashboard management
app.include_router(admin.router)


# ============================================================================
# Startup and Shutdown Events
# ============================================================================

@app.on_event("startup")
async def startup_event():
    """Handle application startup."""
    logger.info(f"{APP_NAME} v{APP_VERSION} starting up...")
    logger.info(f"CORS origins: {CORS_ORIGINS}")
    logger.info(f"API prefix: {API_PREFIX}")

    # Validate configuration
    if not settings.validate():
        logger.warning("Some environment variables may be missing")


@app.on_event("shutdown")
async def shutdown_event():
    """Handle application shutdown."""
    logger.info(f"{APP_NAME} shutting down...")


# ============================================================================
# Root Endpoint
# ============================================================================

@app.get("/")
async def root():
    """Root endpoint with API information."""
    return {
        "app": APP_NAME,
        "version": APP_VERSION,
        "description": APP_DESCRIPTION,
        "docs": "/docs",
        "health": "/health",
        "api_prefix": API_PREFIX,
    }


# ============================================================================
# Application Entry Point
# ============================================================================

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "main:app",
        host=settings.APP_HOST,
        port=settings.APP_PORT,
        reload=settings.DEBUG,
        log_level=settings.LOG_LEVEL.lower(),
    )
