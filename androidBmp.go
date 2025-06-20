package rvtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AndroidBmpDevice struct {
	Manufacturer string `json:"manufacturer"`
	Hardware     string `json:"hardware"`
	Model        string `json:"model"`
	Release      string `json:"release"`
	Brand        string `json:"brand"`
	ID           string `json:"id"`
	HeightPixels int    `json:"heightPixels"`
	WidthPixels  int    `json:"widthPixels"`
}

type AndroidBmpSession struct {
	sessionData map[string]any
	client      *Client
}

func (s *AndroidBmpSession) AndroidId() string {
	if androidId, ok := s.sessionData["androidId"].(string); ok {
		return androidId
	}

	return ""
}

func (s *AndroidBmpSession) StartMillis() int {
	if startMillis, ok := s.sessionData["startMillis"].(float64); ok {
		return int(startMillis)
	}

	return 0
}

func (client *Client) AndroidBmpInit() (*AndroidBmpDevice, *AndroidBmpSession, error) {
	req, err := http.NewRequest("GET", client.Host+"/bmp/android/init", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error while creating request: %w", err)
	}

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error while sending request: %w", err)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error while reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error string `json:"error"`
		}

		if err := json.Unmarshal(respBytes, &errorResponse); err != nil {
			return nil, nil, fmt.Errorf("error while unmarshalling error response: %w", err)
		}

		return nil, nil, &RemoteError{errorResponse.Error}
	}

	var androidBmpInitData struct {
		Device  AndroidBmpDevice `json:"device"`
		Session map[string]any   `json:"session"`
	}

	if err := json.Unmarshal(respBytes, &androidBmpInitData); err != nil {
		return nil, nil, fmt.Errorf("error while unmarshalling response: %w", err)
	}

	return &androidBmpInitData.Device, &AndroidBmpSession{
		sessionData: androidBmpInitData.Session,
		client:      client,
	}, nil
}

func (session *AndroidBmpSession) Sensor(version string, appPackage string, opts ...func(map[string]any)) (string, BmpReportData, error) {
	requestData := map[string]any{
		"bmpVersion": version,
		"appPackage": appPackage,
		"session":    session.sessionData,
	}

	for _, option := range opts {
		option(requestData)
	}

	client := session.client

	reqBytes, err := json.Marshal(requestData)
	if err != nil {
		return "", "", fmt.Errorf("error while marshalling request: %w", err)
	}

	req, err := http.NewRequest("POST", client.Host+"/bmp/android/sensor", bytes.NewReader(reqBytes))
	if err != nil {
		return "", "", fmt.Errorf("error while creating request: %w", err)
	}

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error while sending request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error string `json:"error"`
		}

		return "", "", &RemoteError{errorResponse.Error}
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error while reading response: %w", err)
	}

	var androidBmpSensorResponse struct {
		Sensor     string         `json:"sensor"`
		Session    map[string]any `json:"session"`
		ReportData BmpReportData  `json:"reportData"`
	}

	if err := json.Unmarshal(respBytes, &androidBmpSensorResponse); err != nil {
		return "", "", fmt.Errorf("error while unmarshalling response: %w", err)
	}

	session.sessionData = androidBmpSensorResponse.Session

	return androidBmpSensorResponse.Sensor, androidBmpSensorResponse.ReportData, nil
}

func (client *Client) AndroidBmpGetSensorSessionless(version string, appPackage string, additionalData map[string]any) (string, *AndroidBmpDevice, map[string]any, BmpReportData, error) {
	requestData := map[string]any{
		"bmpVersion": version,
		"appPackage": appPackage,
	}

	for key, value := range additionalData {
		requestData[key] = value
	}

	reqBytes, err := json.Marshal(requestData)
	if err != nil {
		return "", nil, nil, "", err
	}

	req, err := http.NewRequest("POST", client.Host+"/bmp/android/sensor", bytes.NewReader(reqBytes))
	if err != nil {
		return "", nil, nil, "", err
	}

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return "", nil, nil, "", err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error string `json:"error"`
		}

		return "", nil, nil, "", fmt.Errorf("AndroidBmpGetSensor: %s", errorResponse.Error)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, nil, "", err
	}

	var androidBmpSensorResponse struct {
		Sensor     string           `json:"sensor"`
		Session    map[string]any   `json:"session"`
		Device     AndroidBmpDevice `json:"device"`
		ReportData BmpReportData    `json:"reportData"`
	}

	if err := json.Unmarshal(respBytes, &androidBmpSensorResponse); err != nil {
		return "", nil, nil, "", err
	}

	if androidBmpSensorResponse.Sensor == "" {
		return "", nil, nil, "", fmt.Errorf("AndroidBmpGetSensorSessionless: sensor is empty")
	}

	return androidBmpSensorResponse.Sensor, &androidBmpSensorResponse.Device, androidBmpSensorResponse.Session, androidBmpSensorResponse.ReportData, nil
}
