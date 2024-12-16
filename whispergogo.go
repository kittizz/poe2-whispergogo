package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Message struct {
	ChatID  string `json:"chat_id"`
	Message string `json:"message"`
}

func NewWhisperGOGO() *WhisperGOGO {
	return &WhisperGOGO{
		baseURL: WHISPERGOGO_BASE_URL,
		client:  &http.Client{},
	}
}

type WhisperGOGO struct {
	baseURL string
	client  *http.Client
}

func (t *WhisperGOGO) SendMessage(chatID string, message string) error {
	payload := Message{
		ChatID:  chatID,
		Message: message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling json: %v", err)
	}

	resp, err := t.client.Post(
		t.baseURL+"/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

type VerifyChatResponse struct {
	ChatID string `json:"chat_id"`
	Valid  bool   `json:"valid"`
}

func (t *WhisperGOGO) VerifyChat(chatID string) (*VerifyChatResponse, error) {
	resp, err := t.client.Get(fmt.Sprintf("%s/verify/%s", t.baseURL, url.PathEscape(chatID)))
	if err != nil {
		return nil, fmt.Errorf("error verifying chat: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chat verification failed with status code: %d", resp.StatusCode)
	}

	var verifyResp VerifyChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &verifyResp, nil
}
