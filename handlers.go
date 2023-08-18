package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	goaway "github.com/TwiN/go-away"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/crypto/bcrypt"
)

type HTTPError struct {
	Code    int
	Message string
	Err     error
}

func (e *HTTPError) Error() string {
	return e.Message
}

func generateRecipeHandler(c *gin.Context) {
	// Retrieve the user from the context
	val, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No user information"})
		return
	}

	user, ok := val.(*User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User information is of the wrong type"})
		return
	}

	if err := db.Where("user_id = ?", user.ID).First(&user.Settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user settings: " + err.Error()})
		return
	}

	// Choose an api key
	var key string
	if user.Settings.EncryptedOpenAIKey != "" {
		decryptedKey, err := decryptOpenAIKey(user.Settings.EncryptedOpenAIKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decrypt API key: " + err.Error()})
		}
		key = decryptedKey
	} else {
		key = os.Getenv(gc.Env.PublicOpenAIKey)
	}

	// Initialize OpenAI client
	openaiClient := openai.NewClient(key)
	ctx := context.Background()

	// Setup request
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 20,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Lorem ipsum",
			},
		},
		Stream: true,
	}

	// Create stream
	stream, err := openaiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ChatCompletionStream error: %v", err)})
		return
	}
	defer stream.Close()

	// Stream the response back to the client
	c.Writer.Header().Set("Content-Type", "text/plain")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Stream error: %v", err)})
			return
		}

		// Write the content to the client
		_, err = c.Writer.Write([]byte(response.Choices[0].Delta.Content))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Write error: %v", err)})
			return
		}
	}
}

// func generateRecipeHandler(c *gin.Context) {
// 	// Retrieve the user from the context
// 	val, ok := c.Get("user")
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "No user information"})
// 		return
// 	}

// 	user, ok := val.(*User)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "User information is of the wrong type"})
// 		return
// 	}

// 	if err := db.Where("user_id = ?", user.ID).First(&user.Settings).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user settings: " + err.Error()})
// 		return
// 	}

// 	// Choose an api key
// 	var key string
// 	if user.Settings.EncryptedOpenAIKey != "" {
// 		decryptedKey, err := decryptOpenAIKey(user.Settings.EncryptedOpenAIKey)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decrypt API key: " + err.Error()})
// 		}
// 		key = decryptedKey
// 	} else {
// 		key = os.Getenv(gc.Env.PublicOpenAIKey)
// 	}

// 	// Parse the request body for the user's prompt
// 	var request struct {
// 		UserPrompt string `json:"userPrompt"`
// 	}

// 	if err := c.BindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
// 		return
// 	}

// 	// Create a new chat service instance with the user's decrypted key
// 	chatService, err := NewChatService(key, placeholderPrePrompt)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat service: " + err.Error()})
// 		return
// 	}

// 	// Create the chat completion with the user's prompt
// 	recipe, err := chatService.CreateChatCompletion(request.UserPrompt)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe: " + err.Error()})
// 		return
// 	}

// 	// Respond with the chat completion
// 	c.JSON(http.StatusOK, gin.H{"recipe": recipe})
// }

// Handler for collecting a recipe
func collectRecipeHandler(c *gin.Context) {
	userID := c.Param("id")
	recipeID := c.Param("recipe_id")

	var recipe Recipe
	if err := db.Where("id = ?", recipeID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	if recipe.DeletedAt != nil {
		recipe.GeneratedBy = nil
		recipe.DeletedAt = nil
	}

	var user User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the recipe still exists just before trying to modify it
	if err := db.Where("id = ?", recipeID).First(&Recipe{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe no longer exists"})
		return
	}

	user.CollectedRecipes = append(user.CollectedRecipes, recipe)

	db.Save(&user)
	db.Save(&recipe)

	c.JSON(http.StatusOK, gin.H{"message": "Recipe collected"})
}

func getSettingsHandler(c *gin.Context) {
	// Retrieve the user from the session
	val, ok := c.Get("user")
	// val, ok := session.Values["user"]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No user information"})
		return
	}

	user, ok := val.(*User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User information is of the wrong type"})
		return
	}

	if err := db.Where("user_id = ?", user.ID).First(&user.Settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user settings: " + err.Error()})
	}

	// Check the validity of the OpenAI key by making a test API call
	isValid, err := verifyOpenAIKey(user.Settings.EncryptedOpenAIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isValid {
		c.HTML(http.StatusOK, "settings.tmpl", gin.H{"isValid": false, "user": user})
		return
	}

	// Render the settings modal template with valid key and user data
	c.HTML(http.StatusOK, "settings.tmpl", gin.H{"isValid": true, "user": user})
}

func verifyOpenAIKey(encryptedOpenAIKey string) (bool, error) {
	// Set as invalid if no key exists yet
	if encryptedOpenAIKey == "" {
		return false, nil
	}

	// Decrypt the OpenAI key
	key, err := decryptOpenAIKey(encryptedOpenAIKey)
	if err != nil {
		return false, errors.New("failed to decrypt OpenAI key:" + err.Error())
	}

	// Set up OpenAI client with the given key
	client := openai.NewClient(key)
	ctx := context.Background()

	// Maximum number of retries
	const maxRetries = 3

	// Delay between retries
	const retryDelay = 10 * time.Second

	// Attempt the verification with retries
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Make a test API call using a minimal completion request
		req := openai.CompletionRequest{
			Model:     openai.GPT3Ada,
			MaxTokens: 5,
			Prompt:    "Test",
		}
		_, err := client.CreateCompletion(ctx, req)

		// Check for specific API errors
		e := &openai.APIError{}
		if errors.As(err, &e) {
			switch e.HTTPStatusCode {
			case 401:
				// Invalid auth or key (do not retry)
				return false, nil
			case 429:
				// Rate limiting or engine overload (wait and retry)
				time.Sleep(retryDelay)
				continue
			case 500:
				// OpenAI server error (retry)
				continue
			default:
				// Unhandled error (do not retry)
				// return false, err
				return true, err
			}
		}

		// If the call was successful, the key is valid
		if err == nil {
			return true, nil
		}
	}

	// If all attempts failed, return false
	return false, errors.New("failed to verify OpenAI key after multiple attempts")
}

func updateUserSettingsHandler(c *gin.Context) {
	// Retrieve the user from the session
	val, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No user information"})
		return
	}

	user, ok := val.(*User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User information is of the wrong type"})
		return
	}

	// Parse the new OpenAI key from the request body
	var newSettings struct {
		OpenAIKey string `json:"apikey"`
	}
	if err := c.ShouldBindJSON(&newSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var openAIKeyChanged = newSettings.OpenAIKey != ""

	// Check if the OpenAI key has been entered
	if openAIKeyChanged {
		// Encrypt the OpenAI key before storing
		encryptedOpenAIKey, err := encryptOpenAIKey(newSettings.OpenAIKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt OpenAI key: " + err.Error()})
			return
		}

		// Update the user's OpenAI key in the UserSettings
		// user.Settings.EncryptedOpenAIKey = encryptedOpenAIKey
		// if err := db.Model(&user.Settings).Update("EncryptedOpenAIKey", user.Settings.EncryptedOpenAIKey).Error; err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		// 	return
		// }
		if err := db.Model(&UserSettings{}).Where("user_id = ?", user.ID).Update("EncryptedOpenAIKey", encryptedOpenAIKey).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
			return
		}
	}

	// This won't seem as redundant when more settings are && added
	if openAIKeyChanged {
		c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "No changes made"})
	}
}

// Handler for logging in a user
func loginUserHandler(c *gin.Context) {
	var userCredentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Returns error if a required field is not included
	if err := c.ShouldBindJSON(&userCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("username = ?", userCredentials.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userCredentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	session := c.MustGet("session").(*sessions.Session)
	session.Values["user_id"] = user.ID
	// session.Values["user"] = user
	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}

func signupUserHandler(c *gin.Context) {
	var newUser struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Returns error if a required field is not included
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	if err := validateUsername(newUser.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user := User{
		Username:       newUser.Username,
		HashedPassword: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	userSettings := UserSettings{
		UserID: user.ID, // Use the ID of the newly created user
	}

	if err := db.Create(&userSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User signed up successfully"})
}

var forbiddenUsernames = []string{
	"admin",
	"administrator",
	"root",
	// "julian",
	"yana",
	"sys",
	"sysadmin",
	"system",
	"test",
	"testuser",
	"test-user",
	"test_user",
	"login",
	"logout",
	"register",
	"password",
	"user",
	"user123",
	"newuser",
	"yourapp",
	"yourcompany",
	"yourbrand",
	"support",
	"help",
	"faq",
	"culinaryai",
	"culinaryAI",
	"CulinaryAi",
	"CULINARYAI",
	"culinarya1",
	"cul1naryai",
	"culinaryal",
	"culnaryai",
	"culinary_ai",
	"culinary-ai",
	"culinaryaiadmin",
	"culinaryai_admin",
	"culinaryai-admin",
	"culinaryairoot",
	"culinaryai_root",
	"culinaryai-root",
}

func getSessionUser(c *gin.Context) *User {
	session := c.MustGet("session").(*sessions.Session)
	userID, ok := session.Values["user_id"].(uint) // Adjust the type as needed
	if !ok || userID == 0 {
		return nil
	}

	user := &User{}
	if err := db.Preload("Settings").Where("id = ?", userID).First(user).Error; err != nil {
		// If no user is found in the database, return nil
		return nil
	}

	return user
}

func validateUsername(username string) error {
	lowercaseUsername := strings.ToLower(username)

	var user User
	if err := db.Where("LOWER(username) = ?", lowercaseUsername).First(&user).Error; err == nil {
		return fmt.Errorf("username is already taken")
	}

	minLength := 3
	if len(username) < minLength {
		return fmt.Errorf("username must be at least %d characters", minLength)
	}

	if !govalidator.IsAlphanumeric(username) {
		return fmt.Errorf("username can only contain alphanumeric characters")
	}

	for _, forbiddenUsername := range forbiddenUsernames {
		if strings.EqualFold(lowercaseUsername, forbiddenUsername) {
			return fmt.Errorf("username '%s' is not allowed", username)
		}
	}

	// Profanity check using goaway library
	profanityDetector := goaway.NewProfanityDetector().WithSanitizeLeetSpeak(true).WithSanitizeSpecialCharacters(true).WithSanitizeAccents(false)
	if profanityDetector.IsProfane(username) {
		return fmt.Errorf("username contains inappropriate language")
	}

	// If we've passed all checks, the username is valid.
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	hasUppercase, _ := regexp.MatchString(`[A-Z]`, password)
	if !hasUppercase {
		return errors.New("password must contain at least one uppercase letter")
	}
	hasLowercase, _ := regexp.MatchString(`[a-z]`, password)
	if !hasLowercase {
		return errors.New("password must contain at least one lowercase letter")
	}
	hasNumber, _ := regexp.MatchString(`\d`, password)
	if !hasNumber {
		return errors.New("password must contain at least one digit")
	}
	hasSpecialChar, _ := regexp.MatchString(`[!@#$%^&*]`, password)
	if !hasSpecialChar {
		return errors.New("password must contain at least one special character")
	}
	return nil
}
