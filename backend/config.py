"""
Configuration management for the application.
"""

import os
from dotenv import load_dotenv
from typing import List

# Load environment variables from .env file
load_dotenv(override=True)


class Settings:
    """Application settings loaded from environment variables."""

    # API Configuration
    DEBUG: bool = os.getenv("DEBUG", "False").lower() == "true"
    APP_HOST: str = os.getenv("APP_HOST", "0.0.0.0")
    APP_PORT: int = int(os.getenv("APP_PORT", 8000))

    # API Keys and Credentials
    RETELL_API_KEY: str = os.getenv("RETELL_API_KEY", "")
    OPENAI_API_KEY: str = os.getenv("OPENAI_API_KEY", "")
    PINATA_API_KEY: str = os.getenv("PINATA_API_KEY", "")
    PINATA_API_SECRET: str = os.getenv("PINATA_API_SECRET", "")

    # CORS Configuration
    CORS_ORIGINS: List[str] = [
        "http://localhost:3000",
        "http://localhost:3001",
        "http://localhost:8000",
    ]

    # Database Configuration
    DATABASE_URL: str = os.getenv("DATABASE_URL", "sqlite:///./db.sqlite")
    USE_IN_MEMORY_DB: bool = os.getenv("USE_IN_MEMORY_DB", "true").lower() == "true"

    # Logging Configuration
    LOG_LEVEL: str = os.getenv("LOG_LEVEL", "INFO")
    LOG_FILE: str = os.getenv("LOG_FILE", "app.log")

    # LLM Configuration
    LLM_MODEL: str = os.getenv("LLM_MODEL", "gpt-4")
    LLM_TEMPERATURE: float = float(os.getenv("LLM_TEMPERATURE", 0.7))
    LLM_MAX_TOKENS: int = int(os.getenv("LLM_MAX_TOKENS", 1024))

    # Retell Configuration
    RETELL_AGENT_ID: str = os.getenv("RETELL_AGENT_ID", "")
    RETELL_PHONE_NUMBER: str = os.getenv("RETELL_PHONE_NUMBER", "")

    # Validate required settings
    @classmethod
    def validate(cls) -> bool:
        """Validate that all required settings are configured."""
        required = ["RETELL_API_KEY", "OPENAI_API_KEY"]
        missing = [key for key in required if not getattr(cls, key, None)]

        if missing:
            print(f"Warning: Missing required environment variables: {', '.join(missing)}")
            return False
        return True


# Global settings instance
settings = Settings()
