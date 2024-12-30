package huggingface

import (
	"Orbyters/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetGPTNeoResponse(prompt string) (string, error) {
	apiToken := config.HuggingFaceKey
	url := config.HugginFaceUrl

	body := map[string]interface{}{
		"inputs": prompt,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	var response []map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	if len(response) > 0 && response[0]["generated_text"] != nil {
		return response[0]["generated_text"].(string), nil
	}

	return "", fmt.Errorf("no generated text found in response")
}
