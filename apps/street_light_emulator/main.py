from datetime import datetime, timezone
from typing import Dict

from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI(title="street-light-emulator", version="1.0.0")

# serial_number -> on/off
STATE: Dict[str, bool] = {}


class PowerRequest(BaseModel):
    set_value: int = 0


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/api/v1/street-lights/{serial_number}/power")
def set_power(serial_number: str, payload: PowerRequest) -> dict[str, object]:
    is_on = payload.set_value > 0
    STATE[serial_number] = is_on
    return {
        "serial_number": serial_number,
        "is_on": is_on,
        "status": "active" if is_on else "inactive",
        "updated_at": datetime.now(timezone.utc),
    }


@app.get("/api/v1/street-lights/{serial_number}/status")
def get_status(serial_number: str) -> dict[str, object]:
    is_on = STATE.get(serial_number, False)
    return {
        "serial_number": serial_number,
        "is_on": is_on,
        "status": "active" if is_on else "inactive",
    }
