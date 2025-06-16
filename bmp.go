package rvtools

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type BmpReportData string

func (client *Client) BmpFeedback(valid bool, reportData BmpReportData) error {
	reqBytes, err := json.Marshal(struct {
		Valid      bool          `json:"valid"`
		ReportData BmpReportData `json:"reportData"`
	}{
		Valid:      valid,
		ReportData: reportData,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", client.Host+"/bmp/feedback", bytes.NewReader(reqBytes))
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", client.ApiKey)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error string `json:"error"`
		}

		return fmt.Errorf("BmpFeedback: %s", errorResponse.Error)
	}

	return nil
}

func WithBmpDCIScript(scriptBytes []byte) func(map[string]any) {
	return func(request map[string]any) {
		request["dciScriptBase64"] = base64.StdEncoding.EncodeToString(scriptBytes)
	}
}

func WithBmpMockedDCIScript() func(map[string]any) {
	return func(request map[string]any) {
		request["dciScriptBase64"] = "mock"
	}
}

func WithBmpParams(paramsBytes []byte) func(map[string]any) {
	return func(request map[string]any) {
		request["paramsBase64"] = base64.StdEncoding.EncodeToString(paramsBytes)
	}
}

func WithBmpLanguage(language string) func(map[string]any) {
	return func(request map[string]any) {
		request["language"] = language
	}
}

func WithBmpAppVersion(appVersion string) func(map[string]any) {
	return func(request map[string]any) {
		request["appVersion"] = appVersion
	}
}

func WithBmpAppVersionCode(appVersionCode string) func(map[string]any) {
	return func(request map[string]any) {
		request["appVersionCode"] = appVersionCode
	}
}

func WithBmpSession(session map[string]any) func(map[string]any) {
	return func(request map[string]any) {
		request["session"] = session
	}
}

func WithBmpOption(key string, value any) func(map[string]any) {
	return func(request map[string]any) {
		request[key] = value
	}
}
