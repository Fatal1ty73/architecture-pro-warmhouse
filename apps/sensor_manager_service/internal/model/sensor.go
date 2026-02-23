package model

import "time"

type Sensor struct {
	ID           int64     `json:"id"`
	HouseID      int64     `json:"house_id"`
	Name         string    `json:"name"`
	TypeID       int64     `json:"type_id"`
	Location     string    `json:"location"`
	SerialNumber string    `json:"serial_number"`
	Value        float64   `json:"value"`
	Unit         string    `json:"unit"`
	Status       string    `json:"status"`
	LastUpdated  time.Time `json:"last_updated"`
	CreatedAt    time.Time `json:"created_at"`
}

type SensorCreate struct {
	HouseID      int64   `json:"house_id" binding:"required"`
	Name         string  `json:"name" binding:"required,max=100"`
	TypeID       int64   `json:"type_id" binding:"required"`
	Location     string  `json:"location" binding:"required,max=100"`
	SerialNumber string  `json:"serial_number" binding:"required,max=100"`
	Value        float64 `json:"value"`
	Unit         string  `json:"unit" binding:"max=20"`
	Status       string  `json:"status" binding:"required"`
}

type SensorUpdate struct {
	Name         *string  `json:"name"`
	TypeID       *int64   `json:"type_id"`
	Location     *string  `json:"location"`
	SerialNumber *string  `json:"serial_number"`
	Value        *float64 `json:"value"`
	Unit         *string  `json:"unit"`
	Status       *string  `json:"status"`
}

type SensorList struct {
	Items  []Sensor `json:"items"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Total  int      `json:"total"`
}
