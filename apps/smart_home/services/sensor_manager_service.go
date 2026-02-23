package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"smarthome/models"
)

type SensorManagerService struct {
	BaseURL    string
	HTTPClient *http.Client
}

type SensorListResponse struct {
	Items []models.Sensor `json:"items"`
}

func NewSensorManagerService(baseURL string) *SensorManagerService {
	if baseURL == "" {
		return nil
	}
	return &SensorManagerService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *SensorManagerService) GetSensors() ([]models.Sensor, error) {
	resp, err := s.HTTPClient.Get(fmt.Sprintf("%s/api/v1/sensors", s.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var out SensorListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out.Items, nil
}

func (s *SensorManagerService) GetSensorByID(id string) (*models.Sensor, error) {
	resp, err := s.HTTPClient.Get(fmt.Sprintf("%s/api/v1/sensors/%s", s.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var out models.Sensor
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
