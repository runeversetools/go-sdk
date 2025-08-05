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

type IOSBmpSession struct {
	sessionData map[string]any
	client      *Client
}

func (s *IOSBmpSession) Udid() string {
	if udid, ok := s.sessionData["udid"].(string); ok {
		return udid
	}

	return ""
}

func (s *IOSBmpSession) StartMillis() int {
	if startMillis, ok := s.sessionData["startMillis"].(float64); ok {
		return int(startMillis)
	}

	return 0
}

type iosBmpInitOptions struct {
	IOSVersion    string `json:"iosVersion"`
	MinIOSVersion string `json:"minIosVersion"`
	MaxIOSVersion string `json:"maxIosVersion"`
	Model         string `json:"model"`
}

type iosBmpInitOption func(*iosBmpInitOptions)

func (client *Client) IOSBmpInit(opts ...iosBmpInitOption) (*IOSBmpDevice, *IOSBmpSession, error) {
	var options iosBmpInitOptions
	for _, option := range opts {
		option(&options)
	}

	req, err := http.NewRequest("GET", client.Host+"/bmp/ios/init", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error while creating request: %w", err)
	}

	query := req.URL.Query()

	if options.IOSVersion != "" {
		query.Set("iosVersion", options.IOSVersion)
	}

	if options.MinIOSVersion != "" {
		query.Set("minIosVersion", options.MinIOSVersion)
	}

	if options.MaxIOSVersion != "" {
		query.Set("maxIosVersion", options.MaxIOSVersion)
	}

	if options.Model != "" {
		query.Set("model", options.Model)
	}

	req.URL.RawQuery = query.Encode()

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error while sending request: %w", err)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error while reading response: %w", err)
	}

	var iosBmpInitData struct {
		Error   string         `json:"error"`
		Device  IOSBmpDevice   `json:"device"`
		Session map[string]any `json:"session"`
	}

	if err := json.Unmarshal(respBytes, &iosBmpInitData); err != nil {
		return nil, nil, fmt.Errorf("error while unmarshalling response: %w", err)
	}

	if iosBmpInitData.Error != "" {
		return nil, nil, &RemoteError{iosBmpInitData.Error}
	}

	return &iosBmpInitData.Device, &IOSBmpSession{
		sessionData: iosBmpInitData.Session,
		client:      client,
	}, nil
}

func (session *IOSBmpSession) Sensor(bmpVersion string, appPackage string, options ...func(map[string]any)) (string, BmpReportData, error) {
	requestData := map[string]any{
		"bmpVersion": bmpVersion,
		"appPackage": appPackage,
		"session":    session.sessionData,
	}

	for _, option := range options {
		option(requestData)
	}

	client := session.client

	reqBytes, err := json.Marshal(requestData)
	if err != nil {
		return "", "", fmt.Errorf("error while marshalling request: %w", err)
	}

	req, err := http.NewRequest("POST", client.Host+"/bmp/ios/sensor", bytes.NewReader(reqBytes))
	if err != nil {
		return "", "", fmt.Errorf("error while creating request: %w", err)
	}

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error while sending request: %w", err)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error while reading response: %w", err)
	}

	var iosBmpSensorResponse struct {
		Error      string         `json:"error"`
		Sensor     string         `json:"sensor"`
		Session    map[string]any `json:"session"`
		ReportData BmpReportData  `json:"reportData"`
	}

	if err := json.Unmarshal(respBytes, &iosBmpSensorResponse); err != nil {
		return "", "", fmt.Errorf("error while unmarshalling response: %w", err)
	}

	if iosBmpSensorResponse.Error != "" {
		return "", "", &RemoteError{iosBmpSensorResponse.Error}
	}

	session.sessionData = iosBmpSensorResponse.Session

	return iosBmpSensorResponse.Sensor, iosBmpSensorResponse.ReportData, nil
}

// Sets the iOS version for the BMP session.
func WithIOSVersion(version string) func(*iosBmpInitOptions) {
	return func(c *iosBmpInitOptions) {
		c.IOSVersion = version
	}
}

// Sets the minimum iOS version (>= version) for the BMP session.
func WithMinIOSVersion(version string) func(*iosBmpInitOptions) {
	return func(c *iosBmpInitOptions) {
		c.MinIOSVersion = version
	}
}

// Sets the maximum iOS version (<= version) for the BMP session.
func WithMaxIOSVersion(version string) func(*iosBmpInitOptions) {
	return func(c *iosBmpInitOptions) {
		c.MaxIOSVersion = version
	}
}

// Sets the device model for the BMP session.
func WithDeviceModel(model string) func(*iosBmpInitOptions) {
	return func(c *iosBmpInitOptions) {
		c.Model = model
	}
}
