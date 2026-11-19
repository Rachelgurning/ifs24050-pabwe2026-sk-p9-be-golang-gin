package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type LLM interface {
	Chat(ctx context.Context, system, user string) (string, error)
}

type OpenAIClient struct {
	Endpoint string
	Key      string
	Model    string
	HTTP     *http.Client
}

func NewOpenAIClient(endpoint, key, model string) *OpenAIClient {
	return &OpenAIClient{Endpoint: endpoint, Key: key, Model: model, HTTP: &http.Client{Timeout: 90 * time.Second}}
}

func (o *OpenAIClient) Configured() bool {
	k := strings.TrimSpace(o.Key)
	m := strings.TrimSpace(o.Model)
	return o.Endpoint != "" && k != "" && k != "key" && m != "" && m != "nama-model"
}

func (o *OpenAIClient) Chat(ctx context.Context, system, user string) (string, error) {
	if !o.Configured() {
		return "", fmt.Errorf("OpenAI belum dikonfigurasi")
	}
	payload, err := json.Marshal(map[string]any{
		"model":       o.Model,
		"temperature": 0.7,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(o.Endpoint, "/")+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.Key)
	resp, err := o.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("OpenAI mengembalikan status %d", resp.StatusCode)
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 || strings.TrimSpace(out.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("respons AI kosong")
	}
	return out.Choices[0].Message.Content, nil
}
