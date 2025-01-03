package huggingface

import (
	"Orbyters/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Choice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

type APIResponse struct {
	Choices []Choice `json:"choices"`
}

func GetMistralResponse(prompt string) (string, error) {
	apiToken := config.HuggingFaceKey
	url := config.HugginFaceUrl

	message := map[string]interface{}{
		"role":    "user",
		"content": prompt,
	}

	body := map[string]interface{}{
		"model":       config.ModelName,
		"messages":    []map[string]interface{}{message},
		"temperature": 0.5,
		"max_tokens":  2048,
		"top_p":       0.7,
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

	var apiResp APIResponse
	err = json.Unmarshal(respBody, &apiResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	if len(apiResp.Choices) > 0 && apiResp.Choices[0].Message.Content != "" {
		return apiResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no generated text found in response")
}
