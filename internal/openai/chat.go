package openai

import (
	"errors"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"github.com/windoze95/saltybytes-api/internal/models"
	"github.com/windoze95/saltybytes-api/internal/util"
)

// GenerateNewRecipe generates a new recipe.
func generateRecipeWithChat(r *RecipeManager) error {
	// New recipe, there shouldn't be a history
	if r.RecipeHistoryMessages != nil || len(r.RecipeHistoryMessages) > 0 {
		return errors.New("RecipeHistoryMessages was not empty")
	}

	// Build the chat completion message stream
	sysPromptTemplate := r.Cfg.OpenaiPrompts.GenNewRecipeSys
	userPromptTemplate := r.Cfg.OpenaiPrompts.GenNewRecipeUser
	sysPrompt := r.Cfg.OpenaiPrompts.FillSysPrompt(sysPromptTemplate, r.UnitSystem, r.Requirements)
	userPrompt := r.Cfg.OpenaiPrompts.FillUserPrompt(userPromptTemplate, r.UserPrompt)
	chatCompletionMessages := []openai.ChatCompletionMessage{
		createSysMsg(sysPrompt),
		createUserMsg(userPrompt),
	}

	// Create the request
	recipeDefReplyRequest, err := createRecipeDefRequest(chatCompletionMessages)
	if err != nil {
		return err
	}

	// Perform the chat completion
	resp, err := createChatCompletionWithRetry(recipeDefReplyRequest, r.Cfg)
	if err != nil {
		return fmt.Errorf("failed to create chat completion: %v", err)
	}

	// Get the recipe def
	recipeDefJSON := resp.Choices[0].Message.FunctionCall.Arguments
	if len(resp.Choices) == 0 || recipeDefJSON == "" {
		return errors.New("OpenAI API returned an empty message")
	}

	// Deserialize the recipe def
	var functionCallArgument models.RecipeDef
	if err = util.DeserializeFromJSONString(recipeDefJSON, &functionCallArgument); err != nil {
		return fmt.Errorf("failed to deserialize FunctionCallArgument: %v", err)
	}

	// Set the recipe def
	r.RecipeDef = &functionCallArgument

	// Set the next history message
	r.NextRecipeHistoryMessage = models.RecipeHistoryMessage{
		UserPrompt:     r.UserPrompt,
		RecipeResponse: functionCallArgument,
		RecipeType:     models.RecipeTypeChat,
	}

	return nil
}
