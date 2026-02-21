package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TemperatureClient struct {
	BaseURL string
	Client  *http.Client
}

type TemperatureResponse struct {
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

func NewTemperatureClient(baseURL string) *TemperatureClient {
	return &TemperatureClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *TemperatureClient) FetchByLocation(location string) (*TemperatureResponse, error) {
	u := fmt.Sprintf("%s/temperature?location=%s", c.BaseURL, url.QueryEscape(location))
	resp, err := c.Client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("temperature API returned %d", resp.StatusCode)
	}

	var out TemperatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *TemperatureClient) FetchBySensorID(sensorID int64) (*TemperatureResponse, error) {
	u := fmt.Sprintf("%s/temperature/%d", c.BaseURL, sensorID)
	resp, err := c.Client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("temperature API returned %d", resp.StatusCode)
	}

	var out TemperatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}
