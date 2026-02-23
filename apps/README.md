# Smart Home Microservices MVP

## Prerequisites

- Docker and Docker Compose

## Getting Started

### Option 1: Using Docker Compose (Recommended)

The easiest way to start the application is to use Docker Compose:

```bash
./init.sh
```

This script will:

1. Build and start the PostgreSQL and application containers
2. Wait for the services to be ready
3. Display information about how to access the API

Alternatively, you can run Docker Compose directly:

```bash
docker-compose up -d
```

Legacy API remains on `http://localhost:8080/api/v1` (unchanged).
New WebApp entry point is `http://localhost:8084/api/v1`.

### Option 2: Manual setup

If you prefer to run the application without Docker:

1. Start the PostgreSQL database:

```bash
docker-compose up -d postgres
```

2. Build and run services independently:

```bash
# legacy monolith
cd smart_home && go build -o smarthome && ./smarthome

# sensor manager
cd ../sensor_manager_service && go run main.go

# device handle service
cd ../device_handle_service && uvicorn main:app --host 0.0.0.0 --port 8083

# heating emulator
cd ../heating_emulator && uvicorn main:app --host 0.0.0.0 --port 8086

# webapp
cd ../webapp && mvn spring-boot:run
```

## API Testing

A Postman collection is provided for testing the API. Import the `smarthome-api.postman_collection.json` file into Postman to get started.

## Services and Ports

- `webapp` (Java Spring): `http://localhost:8084`
- `sensor-manager-service` (Go): `http://localhost:8082`
- `device-handle-service` (Python FastAPI): `http://localhost:8083`
- `heating-emulator` (Python FastAPI): `http://localhost:8086`
- `temperature-api` (Python FastAPI): `http://localhost:8081`
- `smart-home` legacy monolith (Go): `http://localhost:8080`
- `postgres`: `localhost:5432`

## API Scope

OpenAPI contract is described in `../swagger-api.yaml`.
`WebApp` exposes:
- Auth: `/api/v1/auth/login`, `/api/v1/auth/logout`
- Users: `/api/v1/users...`
- Houses: `/api/v1/houses...`
- Sensors: proxied to `SensorManagerService`
- Devices: proxied to `DeviceHandleService`
