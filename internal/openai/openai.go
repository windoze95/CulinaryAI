package openai

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/windoze95/saltybytes-api/internal/config"
	"github.com/windoze95/saltybytes-api/internal/models"
)

// OpenaiClient is a wrapper for the OpenAI API client.
type OpenaiClient struct {
	Client *openai.Client
}

// RecipeManager is a wrapper for the recipe generation process.
type RecipeManager struct {
	UserPrompt             string
	Requirements           string
	UnitSystem             string
	CreateType             models.RecipeType
	RecipeHistoryEntries   []models.RecipeHistoryEntry
	NextRecipeHistoryEntry models.RecipeHistoryEntry
	VisionImageURL         string
	ImageBytes             []byte
	Cfg                    *config.Config
	RecipeDef              *models.RecipeDef
}

// GenerateRecipeWithChat generates a new recipe using chat.
func (rm *RecipeManager) GenerateRecipeWithChat() error {
	return generateRecipeWithChat(rm)
}

// GenerateRecipeWithImportVision generates a new recipe using vision import.
func (rm *RecipeManager) GenerateRecipeWithImportVision() error {
	return generateRecipeWithImportVision(rm)
}

// GenerateRecipeImage generates an image using DALL-E based on the prompt in RecipeManager.RecipeDef.ImagePrompt,
// then assigns the image bytes to RecipeManager.ImageBytes.
func (rm *RecipeManager) GenerateRecipeImage() error {
	return generateRecipeImage(rm)
}

// newOpenaiClient creates a new OpenAI client.
func newOpenaiClient(cfg *config.Config) (*OpenaiClient, error) {
	return &OpenaiClient{
		Client: openai.NewClient(cfg.GetCurrentAPIKey()),
	}, nil
}

// createChatCompletionWithRetry creates a chat completion and retries if necessary.
func createChatCompletionWithRetry(chatCompletionRequest *openai.ChatCompletionRequest, cfg *config.Config) (*openai.ChatCompletionResponse, error) {
	maxRetries := 5
	var resp openai.ChatCompletionResponse
	var chatCompletionRespErr error
	for i := 0; i < maxRetries; i++ {
		c, err := newOpenaiClient(cfg)
		if err != nil {
			log.Printf("error: failed to create chat service: %v", err)
			return nil, err
		}

		resp, chatCompletionRespErr = c.Client.CreateChatCompletion(
			context.Background(),
			*chatCompletionRequest,
		)

		if chatCompletionRespErr == nil && len(resp.Choices) > 0 {
			break
		}

		shouldRetry, waitTime, noRetryErr := handleAPIError(chatCompletionRespErr)
		if !shouldRetry {
			return nil, fmt.Errorf("error: failed to create chat completion: %v", noRetryErr)
		}

		// Wait before next retry
		// Wait time increases slightly per iteration
		time.Sleep(waitTime * time.Duration(i))
	}
	if chatCompletionRespErr != nil {
		return nil, fmt.Errorf("error: failed to create chat completion: exhausted maximum retries. Exiting. ChatCompletion error: %v", chatCompletionRespErr)
	}

	return &resp, nil
}

// handleAPIError handles API errors and returns whether or not to retry, the wait time, and the error.
func handleAPIError(respErr error) (shouldRetry bool, waitTime time.Duration, err error) {
	e := &openai.APIError{}
	if errors.As(respErr, &e) {
		switch e.HTTPStatusCode {
		case 401:
			log.Printf("error: invalid auth or key. Will retry: %v", respErr)
			return true, 0, errors.New("invalid auth or key. Will retry")
			// We will rotate the keys on retry now
			// return false, 0, errors.New("invalid auth or key. Do not retry")
		case 429:
			return true, 2 * time.Second, errors.New("rate limiting or engine overload. Will retry")
		case 500:
			return true, 2 * time.Second, errors.New("openAI server error. Will retry")
		default:
			return false, 0, fmt.Errorf("unhandled error: %v", respErr)
		}
	}
	return false, 0, fmt.Errorf("unhandled error: %v", respErr)
}
