package pgembed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type apiRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type apiResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`
}

// FetchEmbeddings fetches embeddings from the embedder service
func FetchEmbeddings(input []string, embedderUrl, model string) ([][]float32, error) {

	data := &apiRequest{
		Model: model,
		Input: input,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", embedderUrl, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	// req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad embedding response status code: %d", resp.StatusCode)
	}

	result := &apiResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Embeddings, nil
}
