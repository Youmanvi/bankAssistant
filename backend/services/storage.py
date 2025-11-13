"""
Pinata Service for IPFS file storage integration.
"""

import os
import json
import requests
import logging
from typing import Optional, Dict
from config import settings


logger = logging.getLogger(__name__)


class PinataService:
    """Service for uploading files to IPFS via Pinata."""

    PINATA_API_URL = "https://api.pinata.cloud/pinning/pinFileToIPFS"
    PINATA_GATEWAY_BASE = "https://gateway.pinata.cloud/ipfs"

    def __init__(self):
        """Initialize the Pinata service with API credentials."""
        self.api_key = settings.PINATA_API_KEY
        self.api_secret = settings.PINATA_API_SECRET

    def upload_pdf(
        self, pdf_path: str, document_type: str
    ) -> Optional[Dict[str, str]]:
        """
        Upload a PDF file to Pinata IPFS.

        Args:
            pdf_path: Local path to the PDF file
            document_type: Type of document (LOAN, CARD, etc.)

        Returns:
            Dict with 'IpfsHash' and 'PinataURL' if successful, None otherwise
        """
        if not self.api_key or not self.api_secret:
            logger.error("Pinata API credentials not configured")
            return None

        if not os.path.isfile(pdf_path):
            logger.error(f"PDF file not found: {pdf_path}")
            return None

        try:
            with open(pdf_path, "rb") as pdf_file:
                return self._upload_file_to_pinata(pdf_file, document_type, pdf_path)
        except Exception as e:
            logger.error(f"Error uploading PDF to Pinata: {e}")
            return None

    def _upload_file_to_pinata(
        self, file_obj, document_type: str, pdf_path: str
    ) -> Optional[Dict[str, str]]:
        """
        Internal method to handle the actual file upload.

        Args:
            file_obj: File object to upload
            document_type: Type of document
            pdf_path: Original path for metadata

        Returns:
            Response dict with IPFS hash or None
        """
        # Prepare headers
        headers = {
            "pinata_api_key": self.api_key,
            "pinata_secret_api_key": self.api_secret,
        }

        # Prepare files
        files = {
            "file": (os.path.basename(pdf_path), file_obj, "application/pdf")
        }

        # Prepare metadata
        filename_parts = pdf_path.split("_")
        metadata = {
            "name": f"{document_type}_{filename_parts[-2]}_{filename_parts[-1]}",
            "keyvalues": {
                "type": document_type,
                "uploaded_at": str(__import__("datetime").datetime.now()),
            },
        }

        data = {"pinataMetadata": json.dumps(metadata)}

        try:
            response = requests.post(
                self.PINATA_API_URL, files=files, data=data, headers=headers
            )

            if response.status_code == 200:
                result = response.json()
                ipfs_hash = result.get("IpfsHash")
                pinata_url = f"{self.PINATA_GATEWAY_BASE}/{ipfs_hash}"

                logger.info(f"Successfully uploaded to IPFS: {ipfs_hash}")
                return {
                    "IpfsHash": ipfs_hash,
                    "PinataURL": pinata_url,
                    "status": "success",
                }
            else:
                logger.error(
                    f"Pinata upload failed: {response.status_code} - {response.text}"
                )
                return None

        except requests.RequestException as e:
            logger.error(f"Request error during Pinata upload: {e}")
            return None

    def get_gateway_url(self, ipfs_hash: str) -> str:
        """
        Generate a Pinata gateway URL for an IPFS hash.

        Args:
            ipfs_hash: The IPFS hash/CID

        Returns:
            Full Pinata gateway URL
        """
        return f"{self.PINATA_GATEWAY_BASE}/{ipfs_hash}"


# Global Pinata service instance
pinata_service = PinataService()
