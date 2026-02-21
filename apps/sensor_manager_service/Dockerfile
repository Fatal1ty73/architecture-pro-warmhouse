FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/sensor-manager ./main.go

FROM alpine:3.20
RUN adduser -D -g '' appuser
USER appuser
COPY --from=builder /bin/sensor-manager /usr/local/bin/sensor-manager
EXPOSE 8082
ENTRYPOINT ["sensor-manager"]
