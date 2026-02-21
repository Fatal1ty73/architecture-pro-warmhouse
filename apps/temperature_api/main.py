from fastapi import FastAPI, Query, HTTPException
from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime, timezone
from enum import Enum
import random
import os
import uvicorn

app = FastAPI(title="temperature-api")

SENSOR_MAPPING = {
    "1": "Living Room",
    "2": "Bedroom",
    "3": "Kitchen",
}

LOCATION_MAPPING = {
    "Living Room": "1",
    "Bedroom": "2",
    "Kitchen": "3",
}


class SensorStatus(str, Enum):
    active = "active"
    inactive = "inactive"


class TemperatureResponse(BaseModel):
    value: float = Field(...)
    unit: str = Field(...)
    timestamp: datetime = Field(...)
    location: str
    status: SensorStatus
    sensor_id: str
    sensor_type: str
    description: str


def generate_random_temperature() -> float:
    return round(random.uniform(18.0, 30.0), 1)


@app.get("/temperature", response_model=TemperatureResponse)
def get_temperature(
        location: Optional[str] = Query(None),
        sensorId: Optional[str] = Query(None)
):
    location = (location or "").strip()
    sensor_id = (sensorId or "").strip()

    if not location and not sensor_id:
        raise HTTPException(
            status_code=400,
            detail="Either 'location' or 'sensorId' must be provided"
        )

    if not location:
        location = SENSOR_MAPPING.get(sensor_id, "Unknown")

    if not sensor_id:
        sensor_id = LOCATION_MAPPING.get(location, "0")

    status = (
        SensorStatus.active
        if sensor_id in SENSOR_MAPPING
        else SensorStatus.inactive
    )

    return TemperatureResponse(
        value=generate_random_temperature(),
        unit="C",
        timestamp=datetime.now(timezone.utc),
        location=location,
        status=status,
        sensor_id=sensor_id,
        sensor_type="temperature",
        description=f"Temperature reading from {location} sensor"
    )


def build_response(location: str, sensor_id: str) -> TemperatureResponse:
    status = (
        SensorStatus.active
        if sensor_id in SENSOR_MAPPING
        else SensorStatus.inactive
    )

    return TemperatureResponse(
        value=generate_random_temperature(),
        unit="C",
        timestamp=datetime.now(timezone.utc),
        location=location,
        status=status,
        sensor_id=sensor_id,
        sensor_type="temperature",
        description=f"Temperature reading from {location} sensor"
    )


@app.get("/temperature/{sensor_id}", response_model=TemperatureResponse)
def get_by_sensor_id(sensor_id: str):
    location = SENSOR_MAPPING.get(sensor_id)
    if not location:
        raise HTTPException(status_code=404, detail="Sensor not found")

    return build_response(location, sensor_id)

@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}
