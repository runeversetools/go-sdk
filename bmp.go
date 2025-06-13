package rvtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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

	resp, err := httpClient.Do(req)
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
