from datetime import datetime, timezone
from typing import Optional, Any

import os
import psycopg
import httpx
from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel, Field

app = FastAPI(title="device-handle-service", version="1.0.0")

DB_URL = os.getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smarthome")
HEATING_EMULATOR_URL = os.getenv("HEATING_EMULATOR_URL", "http://heating-emulator:8086")


class DeviceCreate(BaseModel):
    house_id: int
    name: str = Field(max_length=100)
    type_id: int
    location: str = Field(max_length=100)
    serial_number: str = Field(max_length=100)
    status: str = "inactive"


class DeviceUpdate(BaseModel):
    name: Optional[str] = Field(default=None, max_length=100)
    type_id: Optional[int] = None
    location: Optional[str] = Field(default=None, max_length=100)
    serial_number: Optional[str] = Field(default=None, max_length=100)
    status: Optional[str] = None


class DeviceExecute(BaseModel):
    set_value: Optional[int] = None


def as_dict(row: tuple[Any, ...]) -> dict[str, Any]:
    return {
        "id": row[0],
        "house_id": row[1],
        "name": row[2],
        "type_id": row[3],
        "location": row[4],
        "serial_number": row[5],
        "status": row[6],
        "last_updated": row[7],
        "created_at": row[8],
    }


def parse_pagination(limit: int = 50, offset: int = 0) -> tuple[int, int]:
    safe_limit = min(max(limit, 1), 200)
    safe_offset = max(offset, 0)
    return safe_limit, safe_offset


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.get("/api/v1/devices")
def list_devices(
    limit: int = Query(50),
    offset: int = Query(0),
    house_id: Optional[int] = Query(None),
    type_id: Optional[int] = Query(None),
    status: Optional[str] = Query(None),
) -> dict[str, Any]:
    limit, offset = parse_pagination(limit, offset)

    where = ["1=1"]
    args: list[Any] = []

    if house_id is not None:
        where.append("house_id = %s")
        args.append(house_id)
    if type_id is not None:
        where.append("type_id = %s")
        args.append(type_id)
    if status is not None:
        where.append("status = %s")
        args.append(status)

    where_clause = " AND ".join(where)

    with psycopg.connect(DB_URL) as conn:
        with conn.cursor() as cur:
            cur.execute(f"SELECT COUNT(*) FROM devices WHERE {where_clause}", args)
            total = cur.fetchone()[0]

            cur.execute(
                f"""
                SELECT id, house_id, name, type_id, location, serial_number, status, last_updated, created_at
                FROM devices
                WHERE {where_clause}
                ORDER BY id
                LIMIT %s OFFSET %s
                """,
                [*args, limit, offset],
            )
            items = [as_dict(r) for r in cur.fetchall()]

    return {"items": items, "limit": limit, "offset": offset, "total": total}


@app.get("/api/v1/devices/{device_id}")
def get_device(device_id: int) -> dict[str, Any]:
    with psycopg.connect(DB_URL) as conn:
        with conn.cursor() as cur:
            cur.execute(
                """
                SELECT id, house_id, name, type_id, location, serial_number, status, last_updated, created_at
                FROM devices
                WHERE id = %s
                """,
                (device_id,),
            )
            row = cur.fetchone()

    if row is None:
        raise HTTPException(status_code=404, detail={"code": "NOT_FOUND", "message": "device not found"})

    return as_dict(row)


@app.post("/api/v1/devices", status_code=201)
def create_device(payload: DeviceCreate) -> dict[str, Any]:
    now = datetime.now(timezone.utc)
    with psycopg.connect(DB_URL) as conn:
        with conn.cursor() as cur:
            cur.execute(
                """
                INSERT INTO devices (house_id, name, type_id, location, serial_number, status, last_updated, created_at)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING id, house_id, name, type_id, location, serial_number, status, last_updated, created_at
                """,
                (
                    payload.house_id,
                    payload.name,
                    payload.type_id,
                    payload.location,
                    payload.serial_number,
                    payload.status,
                    now,
                    now,
                ),
            )
            row = cur.fetchone()
            conn.commit()

    return as_dict(row)


@app.patch("/api/v1/devices/{device_id}")
def update_device(device_id: int, payload: DeviceUpdate) -> dict[str, Any]:
    values = payload.model_dump(exclude_none=True)
    if not values:
        raise HTTPException(status_code=400, detail={"code": "BAD_REQUEST", "message": "empty payload"})

    set_parts: list[str] = []
    args: list[Any] = []
    for key, value in values.items():
        set_parts.append(f"{key} = %s")
        args.append(value)
    set_parts.append("last_updated = %s")
    args.append(datetime.now(timezone.utc))
    args.append(device_id)

    with psycopg.connect(DB_URL) as conn:
        with conn.cursor() as cur:
            cur.execute(
                f"""
                UPDATE devices
                SET {', '.join(set_parts)}
                WHERE id = %s
                RETURNING id, house_id, name, type_id, location, serial_number, status, last_updated, created_at
                """,
                args,
            )
            row = cur.fetchone()
            conn.commit()

    if row is None:
        raise HTTPException(status_code=404, detail={"code": "NOT_FOUND", "message": "device not found"})

    return as_dict(row)


@app.delete("/api/v1/devices/{device_id}", status_code=204)
def delete_device(device_id: int) -> None:
    with psycopg.connect(DB_URL) as conn:
        with conn.cursor() as cur:
            cur.execute("DELETE FROM devices WHERE id = %s", (device_id,))
            if cur.rowcount == 0:
                raise HTTPException(status_code=404, detail={"code": "NOT_FOUND", "message": "device not found"})
            conn.commit()


@app.post("/api/v1/devices/{device_id}/execute")
def execute_device(device_id: int, payload: DeviceExecute) -> dict[str, Any]:
    with psycopg.connect(DB_URL) as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT serial_number FROM devices WHERE id = %s", (device_id,))
            row = cur.fetchone()
    if row is None:
        raise HTTPException(status_code=404, detail={"code": "NOT_FOUND", "message": "device not found"})
    serial_number = row[0]

    set_value = payload.set_value or 0
    try:
        with httpx.Client(timeout=5.0) as client:
            heating_resp = client.post(
                f"{HEATING_EMULATOR_URL}/api/v1/heating/{serial_number}/power",
                json={"set_value": set_value},
            )
            heating_resp.raise_for_status()
            heating_payload = heating_resp.json()
            desired_status = heating_payload.get("status", "inactive")
    except Exception as exc:
        raise HTTPException(
            status_code=502,
            detail={"code": "BAD_GATEWAY", "message": f"heating emulator unavailable: {exc}"},
        ) from exc

    with psycopg.connect(DB_URL) as conn:
        with conn.cursor() as cur:
            cur.execute(
                """
                UPDATE devices
                SET status = %s, last_updated = %s
                WHERE id = %s
                RETURNING id, house_id, name, type_id, location, serial_number, status, last_updated, created_at
                """,
                (desired_status, datetime.now(timezone.utc), device_id),
            )
            row = cur.fetchone()
            conn.commit()

    return as_dict(row)
