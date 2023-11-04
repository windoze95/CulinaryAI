package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	goaway "github.com/TwiN/go-away"
	"github.com/asaskevich/govalidator"
	"github.com/windoze95/saltybytes-api/internal/config"
	"github.com/windoze95/saltybytes-api/internal/models"
	"github.com/windoze95/saltybytes-api/internal/openai"
	"github.com/windoze95/saltybytes-api/internal/repository"
	"github.com/windoze95/saltybytes-api/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Cfg  *config.Config
	Repo *repository.UserRepository
}

// Constructor function for initializing a new UserService
func NewUserService(cfg *config.Config, repo *repository.UserRepository) *UserService {
	return &UserService{
		Cfg:  cfg,
		Repo: repo,
	}
}

func (s *UserService) CreateUser(username, firstName, email, password string) (*models.User, error) {
	// // Validate username
	// if err := s.ValidateUsername(username); err != nil {
	// 	return err
	// }

	// // Validate password
	// if err := s.ValidatePassword(password); err != nil {
	// 	return err
	// }

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	hashedPasswordStr := string(hashedPassword)

	// Create User and UserSettings
	user := &models.User{
		Username:  username,
		FirstName: firstName,
		Email:     email,
		Auth: models.UserAuth{
			HashedPassword: hashedPasswordStr,
			AuthType:       models.Standard,
		},
		Subscription: models.Subscription{
			SubscriptionTier: models.Free,
			ExpiresAt:        time.Now().AddDate(0, 1, 0), // One month from now
		},
		Settings: models.UserSettings{
			KeepScreenAwake: true, // Default value
			// UsePersonalAPIKey:  false, // Default value
			// EncryptedOpenAIKey: "",    // Default value
		},
		GuidingContent: models.GuidingContent{
			// UnitSystem: models.USCustomary, // Default value
			Requirements: s.Cfg.DefaultRequirements,
		},
		CollectedRecipes: []models.Recipe{},
	}

	// settings := &models.UserSettings{}
	// gc := &models.GuidingContent{}
	// gc.UnitSystem = 1 // Default value

	// if err := s.Repo.CreateUser(user); err != nil {
	// 	if pgErr, ok := err.(*pq.Error); ok {
	// 		if pgErr.Code == "23505" { // Unique constraint violation
	// 			if strings.Contains(pgErr.Error(), "username") {
	// 				return fmt.Errorf("username already in use")
	// 			} else if strings.Contains(pgErr.Error(), "email") {
	// 				return fmt.Errorf("email already in use")
	// 			}
	// 		}
	// 	}
	// 	return fmt.Errorf("error creating user: %v", err)
	// }

	user, err = s.Repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) LoginUser(username, password string) (*models.User, error) {
	user, err := s.Repo.GetUserAuthByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Auth.HashedPassword), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Clear the hashed password before returning the user
	// user.HashedPassword = ""

	// user = util.StripSensitiveUserData(user)

	return util.StripSensitiveUserData(user), nil
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	return s.Repo.GetUserByID(userID)
}

// func (s *UserService) GetPreloadedUserByID(userID uint) (*models.User, error) {
// 	return s.Repo.GetPreloadedUserByID(userID)
// }

func (s *UserService) VerifyOpenAIKeyInUserSettings(user *models.User) (bool, error) {
	// Decrypt the OpenAI key
	decryptedKey, err := util.DecryptOpenAIKey(s.Cfg.Env.OpenAIKeyEncryptionKey.Value(), user.Settings.EncryptedOpenAIKey)
	if err != nil {
		return false, fmt.Errorf("failed to decrypt OpenAI key: %v", err)
	}

	// Verify the OpenAI key
	isValid, err := openai.VerifyOpenAIKey(decryptedKey)
	if err != nil {
		return false, fmt.Errorf("failed to verify OpenAI key: %v", err)
	}

	return isValid, nil
}

func (s *UserService) UpdateUserSettings(user *models.User, newOpenAIKey string) (bool, error) {
	// Encrypt the OpenAI key
	encryptedOpenAIKey, err := util.EncryptOpenAIKey(s.Cfg.Env.OpenAIKeyEncryptionKey.Value(), newOpenAIKey)
	if err != nil {
		return false, err
	}

	// Check if the OpenAI key has changed
	openAIKeyChanged := encryptedOpenAIKey != user.Settings.EncryptedOpenAIKey
	if openAIKeyChanged {
		if err := s.Repo.UpdateUserSettingsOpenAIKey(user.ID, encryptedOpenAIKey); err != nil {
			return false, err
		}
	}
	return openAIKeyChanged, nil
}

func (s *UserService) UpdateGuidingContent(user *models.User, updatedGC *models.GuidingContent) error {
	return s.Repo.UpdateGuidingContent(user.ID, updatedGC)
}

// // VerifyRecaptcha verifies the provided reCAPTCHA response
// func (s *UserService) VerifyRecaptcha(recaptchaResponse string) error {
// 	secretKey := s.Cfg.Env.RecaptchaSecretKey.Value()

// 	// Google reCAPTCHA API endpoint for server-side verification
// 	apiURL := "https://www.google.com/recaptcha/api/siteverify"

// 	response, err := http.PostForm(apiURL, url.Values{"secret": {secretKey}, "response": {recaptchaResponse}})
// 	if err != nil {
// 		return errors.New("Failed to verify reCAPTCHA: " + err.Error())
// 	}
// 	defer response.Body.Close()

// 	var result map[string]interface{}
// 	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
// 		return errors.New("Failed to read reCAPTCHA response: " + err.Error())
// 	}

// 	if success, ok := result["success"].(bool); !ok || !success {
// 		return errors.New("reCAPTCHA verification failed")
// 	}

// 	return nil
// }

func (s *UserService) ValidateUsername(username string) error {
	// exists, err := s.Repo.UsernameExists(username)
	// if err != nil {
	// 	return fmt.Errorf("error checking username: %v", err)
	// }
	// if exists {
	// 	return fmt.Errorf("username is already taken")
	// }

	minLength := 3
	if len(username) < minLength {
		return fmt.Errorf("username must be at least %d characters", minLength)
	}

	if !govalidator.IsAlphanumeric(username) {
		return fmt.Errorf("username can only contain alphanumeric characters")
	}

	var forbiddenUsernames = []string{
		"admin",
		"administrator",
		"root",
		"julian",
		"awfulbits",
		// "windoze95",
		"yana",
		"russianminx",
		"russianminxx",
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
		"saltybytes",
		"saltybytes_ai",
		"saltybytes-ai",
		"saltybytesadmin",
		"saltybytes_admin",
		"saltybytes-admin",
		"saltybytesroot",
		"saltybytes_root",
		"saltybytes-root",
	}

	lowercaseUsername := strings.ToLower(username)
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

func (s *UserService) ValidateEmail(email string) error {
	if !govalidator.IsEmail(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (s *UserService) ValidatePassword(password string) error {
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
