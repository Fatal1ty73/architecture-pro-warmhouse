package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type DeviceHandleService struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewDeviceHandleService(baseURL string) *DeviceHandleService {
	if baseURL == "" {
		return nil
	}
	return &DeviceHandleService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *DeviceHandleService) Forward(method, path string, query url.Values, body []byte) (int, []byte, error) {
	targetURL, err := url.Parse(fmt.Sprintf("%s%s", s.BaseURL, path))
	if err != nil {
		return 0, nil, err
	}
	if query != nil {
		targetURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequest(method, targetURL.String(), bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, payload, nil
}
