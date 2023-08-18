package main

import (
	"context"
	"errors"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type ChatService struct {
	Client    *openai.Client
	PrePrompt string
}

func NewChatService(decryptedAPIKey string, prePrompt string) (*ChatService, error) {
	// decryptedAPIKey, err := Decrypt(NewCipherConfig(), encryptedAPIKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to decrypt API key: %v", err)
	// }

	return &ChatService{
		Client:    openai.NewClient(decryptedAPIKey),
		PrePrompt: prePrompt,
	}, nil
}

func (cs *ChatService) CreateChatCompletion(userPrompt string) (string, error) {
	resp, err := cs.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: cs.PrePrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", errors.New("OpenAI API returned an empty message")
	}

	return resp.Choices[0].Message.Content, nil
}
