package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	goaway "github.com/TwiN/go-away"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
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

func collectRecipe(db *gorm.DB, userID uint, recipeID uint) error {
	var user User
	var recipe Recipe

	// Find the user and recipe
	if err := db.First(&user, userID).Error; err != nil {
		return err
	}
	if err := db.First(&recipe, recipeID).Error; err != nil {
		return err
	}

	// Add the recipe to the user's collected recipes
	user.CollectedRecipes = append(user.CollectedRecipes, recipe)

	// Save the user
	if err := db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

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

func signupUserHandler(c *gin.Context) {
	var newUser struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Returns error if a required field is not included
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, gin.H{"message": "User signed up successfully"})
}

var forbiddenUsernames = []string{
	"admin",
	"administrator",
	"root",
	"julian",
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
	userID := session.Values["user_id"]
	if userID == nil {
		return nil
	}

	user := &User{}
	if err := db.Where("id = ?", userID).First(user).Error; err != nil {
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
