import logging
from typing import Optional

# Configure logging
def setup_logging(logger_name: Optional[str] = None) -> logging.Logger:
    logging.basicConfig(
        level=logging.INFO,
        format="%(levelname)s\t%(asctime)s:\t%(message)s",
        datefmt="%Y-%m-%d %H:%M:%S"
    )
    return logging.getLogger(logger_name or "ark-api")
