package Dto

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
