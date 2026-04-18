package Dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequestDto struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponseDto struct {
	Object string `json:"object"`
	Model  string `json:"model"`
}

type MessageDto struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChoiceDto struct {
	Index        int        `json:"index"`
	FinishReason string     `json:"finish_reason"`
	Message      MessageDto `json:"message"`
}

type ChatCompletionResponseDto struct {
	ID      string      `json:"id"`
	Object  string      `json:"object"`
	Model   string      `json:"model"`
	Created int64       `json:"created"`
	Choices []ChoiceDto `json:"choices"`
}

func ParseChatRequest(body []byte) (ChatRequestDto, error) {
	var zero ChatRequestDto
	if len(strings.TrimSpace(string(body))) == 0 {
		return zero, errors.New("request body is empty")
	}
	var req ChatRequestDto
	if err := json.Unmarshal(body, &req); err != nil {
		return zero, fmt.Errorf("invalid JSON: %w", err)
	}
	return req, nil
}

func (r *ChatRequestDto) Validate() error {
	if strings.TrimSpace(r.Model) == "" {
		return errors.New("model is required")
	}
	if len(r.Messages) == 0 {
		return errors.New("messages must be a non-empty array")
	}
	for i, m := range r.Messages {
		if strings.TrimSpace(m.Role) == "" {
			return fmt.Errorf("messages[%d].role is required", i)
		}
		if strings.TrimSpace(m.Content) == "" {
			return fmt.Errorf("messages[%d].content is required", i)
		}
	}
	return nil
}

func ParseAndValidateChatRequest(body []byte) (ChatRequestDto, error) {
	req, err := ParseChatRequest(body)
	if err != nil {
		return ChatRequestDto{}, err
	}
	if err := req.Validate(); err != nil {
		return ChatRequestDto{}, err
	}
	return req, nil
}
