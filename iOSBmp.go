package rvtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IOSBmpDevice struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	IOSVersion string `json:"iosVersion"`
}

type iosBmpInitContext struct {
	IOSVersion    string `json:"iosVersion"`
	MinIOSVersion string `json:"minIosVersion"`
	MaxIOSVersion string `json:"maxIosVersion"`
}

type IOSBmpInitOption func(*iosBmpInitContext)

func (client *Client) IOSBmpInit(options ...IOSBmpInitOption) (*IOSBmpDevice, map[string]any, error) {
	var initContext iosBmpInitContext
	for _, option := range options {
		option(&initContext)
	}

	req, err := http.NewRequest("GET", client.Host+"/bmp/ios/init", nil)
	if err != nil {
		return nil, nil, err
	}

	query := req.URL.Query()

	query.Set("iosVersion", initContext.IOSVersion)
	if initContext.IOSVersion != "" {
		query.Set("iosVersion", initContext.IOSVersion)
	}

	if initContext.MinIOSVersion != "" {
		query.Set("minIosVersion", initContext.MinIOSVersion)
	}

	if initContext.MaxIOSVersion != "" {
		query.Set("maxIosVersion", initContext.MaxIOSVersion)
	}

	req.URL.RawQuery = query.Encode()

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var iosBmpInitData struct {
		Error   string         `json:"error"`
		Device  IOSBmpDevice   `json:"device"`
		Session map[string]any `json:"session"`
	}

	if err := json.Unmarshal(respBytes, &iosBmpInitData); err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("IOSBmpGetSensor: %s", iosBmpInitData.Error)
	}

	return &iosBmpInitData.Device, iosBmpInitData.Session, nil
}

func (client *Client) IOSBmpGetSensor(bmpVersion string, appPackage string, additionalData map[string]any, session map[string]any) (string, map[string]any, BmpReportData, error) {
	requestData := map[string]any{
		"bmpVersion": bmpVersion,
		"appPackage": appPackage,
		"session":    session,
	}

	for key, value := range additionalData {
		requestData[key] = value
	}

	reqBytes, err := json.Marshal(requestData)
	if err != nil {
		return "", nil, nil, err
	}

	req, err := http.NewRequest("POST", client.Host+"/bmp/ios/sensor", bytes.NewReader(reqBytes))
	if err != nil {
		return "", nil, nil, err
	}

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", nil, nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, nil, err
	}

	var iosBmpSensorResponse struct {
		Error      string         `json:"error"`
		Sensor     string         `json:"sensor"`
		Session    map[string]any `json:"session"`
		ReportData BmpReportData  `json:"reportData"`
	}

	if err := json.Unmarshal(respBytes, &iosBmpSensorResponse); err != nil {
		return "", nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil, nil, fmt.Errorf("IOSBmpGetSensor: %s", iosBmpSensorResponse.Error)
	}

	return iosBmpSensorResponse.Sensor, iosBmpSensorResponse.Session, iosBmpSensorResponse.ReportData, nil
}

func WithIOSVersion(version string) func(*iosBmpInitContext) {
	return func(c *iosBmpInitContext) {
		c.IOSVersion = version
	}
}

func WithMinIOSVersion(version string) func(*iosBmpInitContext) {
	return func(c *iosBmpInitContext) {
		c.MinIOSVersion = version
	}
}

func WithMaxIOSVersion(version string) func(*iosBmpInitContext) {
	return func(c *iosBmpInitContext) {
		c.MaxIOSVersion = version
	}
}
