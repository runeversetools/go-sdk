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

type BmpReportData any

func (client *Client) AndroidBmpInit() (*AndroidBmpDevice, map[string]any, error) {
	req, err := http.NewRequest("GET", client.Host+"/bmp/android/init", nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error string `json:"error"`
		}

		if err := json.Unmarshal(respBytes, &errorResponse); err != nil {
			return nil, nil, fmt.Errorf("AndroidBmpInit: %s", errorResponse.Error)
		}

		return nil, nil, fmt.Errorf("AndroidBmpInit: %s", errorResponse.Error)
	}

	var androidBmpInitData struct {
		Device  AndroidBmpDevice `json:"device"`
		Session map[string]any   `json:"session"`
	}

	if err := json.Unmarshal(respBytes, &androidBmpInitData); err != nil {
		return nil, nil, err
	}

	return &androidBmpInitData.Device, androidBmpInitData.Session, nil
}

func (client *Client) AndroidBmpGetSensor(version string, appPackage string, language string, additionalData map[string]any, session map[string]any) (string, map[string]any, BmpReportData, error) {
	requestData := map[string]any{
		"bmpVersion": version,
		"appPackage": appPackage,
		"language":   language,
		"session":    session,
	}

	for key, value := range additionalData {
		requestData[key] = value
	}

	reqBytes, err := json.Marshal(requestData)
	if err != nil {
		return "", nil, nil, err
	}

	req, err := http.NewRequest("POST", client.Host+"/bmp/android/sensor", bytes.NewReader(reqBytes))
	if err != nil {
		return "", nil, nil, err
	}

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error string `json:"error"`
		}

		return "", nil, nil, fmt.Errorf("AndroidBmpGetSensor: %s", errorResponse.Error)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, nil, err
	}

	var androidBmpSensorResponse struct {
		Sensor     string         `json:"sensor"`
		Session    map[string]any `json:"session"`
		ReportData BmpReportData  `json:"reportData"`
	}

	if err := json.Unmarshal(respBytes, &androidBmpSensorResponse); err != nil {
		return "", nil, nil, err
	}

	if androidBmpSensorResponse.Sensor == "" {
		return "", nil, nil, fmt.Errorf("AndroidBmpGetSensorSessionless: sensor is empty")
	}

	return androidBmpSensorResponse.Sensor, androidBmpSensorResponse.Session, androidBmpSensorResponse.ReportData, nil
}

func (client *Client) AndroidBmpGetSensorSessionless(version string, appPackage string, language string, additionalData map[string]any) (string, *AndroidBmpDevice, map[string]any, BmpReportData, error) {
	requestData := map[string]any{
		"bmpVersion": version,
		"appPackage": appPackage,
		"language":   language,
	}

	for key, value := range additionalData {
		requestData[key] = value
	}

	reqBytes, err := json.Marshal(requestData)
	if err != nil {
		return "", nil, nil, nil, err
	}

	req, err := http.NewRequest("POST", client.Host+"/bmp/android/sensor", bytes.NewReader(reqBytes))
	if err != nil {
		return "", nil, nil, nil, err
	}

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", nil, nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error string `json:"error"`
		}

		return "", nil, nil, nil, fmt.Errorf("AndroidBmpGetSensor: %s", errorResponse.Error)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, nil, nil, err
	}

	var androidBmpSensorResponse struct {
		Sensor     string           `json:"sensor"`
		Session    map[string]any   `json:"session"`
		Device     AndroidBmpDevice `json:"device"`
		ReportData BmpReportData    `json:"reportData"`
	}

	if err := json.Unmarshal(respBytes, &androidBmpSensorResponse); err != nil {
		return "", nil, nil, nil, err
	}

	if androidBmpSensorResponse.Sensor == "" {
		return "", nil, nil, nil, fmt.Errorf("AndroidBmpGetSensorSessionless: sensor is empty")
	}

	return androidBmpSensorResponse.Sensor, &androidBmpSensorResponse.Device, androidBmpSensorResponse.Session, androidBmpSensorResponse.ReportData, nil
}
