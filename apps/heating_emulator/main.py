from datetime import datetime, timezone
from typing import Dict

from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI(title="heating-emulator", version="1.0.0")

# serial_number -> heating enabled/disabled
STATE: Dict[str, bool] = {}


class HeatingRequest(BaseModel):
    set_value: int = 0


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/api/v1/heating/{serial_number}/power")
def set_heating(serial_number: str, payload: HeatingRequest) -> dict[str, object]:
    is_enabled = payload.set_value > 0
    STATE[serial_number] = is_enabled
    return {
        "serial_number": serial_number,
        "heating_enabled": is_enabled,
        "status": "active" if is_enabled else "inactive",
        "updated_at": datetime.now(timezone.utc),
    }


@app.get("/api/v1/heating/{serial_number}/status")
def get_heating_status(serial_number: str) -> dict[str, object]:
    is_enabled = STATE.get(serial_number, False)
    return {
        "serial_number": serial_number,
        "heating_enabled": is_enabled,
        "status": "active" if is_enabled else "inactive",
    }
