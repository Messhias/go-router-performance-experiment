package Dto

type ChatRequestDto struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
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
	Created int64       `json:"created"`
	Model   string      `json:"model"`
	Choices []ChoiceDto `json:"choices"`
}
