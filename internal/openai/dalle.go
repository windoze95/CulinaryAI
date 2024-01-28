package openai

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/windoze95/saltybytes-api/internal/config"
)

// generateRecipeImage generates an image using DALL-E based on the prompt in RecipeManager.RecipeDef.ImagePrompt,
// then assigns the image bytes to RecipeManager.ImageBytes.
func generateRecipeImage(r *RecipeManager) error {
	// Tests for the presence of a prompt
	if r.RecipeDef.ImagePrompt == "" {
		return errors.New("ImagePrompt is nil")
	}

	imageBytes, err := createImage(r.RecipeDef.ImagePrompt, r.Cfg)
	if err != nil {
		log.Printf("error: failed to create recipe image completion: %v", err)
		return err
	}

	r.ImageBytes = imageBytes

	return nil
}

// createImage generates an image using DALL-E based on the provided prompt.
func createImage(prompt string, cfg *config.Config) ([]byte, error) {
	maxRetries := 3
	var respBase64 openai.ImageResponse
	var err error

	for i := 0; i < maxRetries; i++ {
		c, err := newOpenaiClient(cfg)
		if err != nil {
			log.Printf("error: failed to create image service: %v", err)
			return nil, err
		}

		respBase64, err = c.Client.CreateImage(
			context.Background(),
			openai.ImageRequest{
				Prompt:         prompt,
				Size:           openai.CreateImageSize512x512,
				ResponseFormat: openai.CreateImageResponseFormatB64JSON,
				N:              1,
			},
		)

		if err == nil {
			break
		}

		shouldRetry, waitTime, noRetryErr := handleAPIError(err)
		if !shouldRetry {
			return nil, noRetryErr
		}

		// Wait before next retry
		time.Sleep(waitTime * time.Duration(i))
	}

	if err != nil {
		return nil, fmt.Errorf("exhausted maximum retries. Exiting. CreateImage error: %v", err)
	}

	if len(respBase64.Data) == 0 || respBase64.Data[0].B64JSON == "" {
		return nil, errors.New("openAI API returned an empty image")
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %v", err)
	}

	return imgBytes, nil
}
