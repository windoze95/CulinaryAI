package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
)

type Config struct {
	Env                   Env `json:"env"`
	OpenaiPrompts         Prompts
	OpenaiKeys            []string
	CurrentOpenaiKeyIndex int
	Mutex                 sync.RWMutex
}

type Env struct {
	Port                   EnvVar `json:"port"`
	DatabaseUrl            EnvVar `json:"databaseUrl"`
	OpenAIKeyEncryptionKey EnvVar `json:"openAIKeyEncryptionKey"`
	JwtSecretKey           EnvVar `json:"jwtSecretKey"`
	PublicOpenAIKey        EnvVar `json:"publicOpenAIKey"`
	AWSRegion              EnvVar `json:"awsRegion"`
	AWSAccessKeyID         EnvVar `json:"awsAccessKeyID"`
	AWSSecretAccessKey     EnvVar `json:"awsSecretAccessKey"`
	S3Bucket               EnvVar `json:"s3Bucket"`
	OpenaiPromptsFilePath  EnvVar `json:"openaiPromptsFilePath"`
	OpenaiKeysFilePath     EnvVar `json:"openaiKeysFilePath"`
}

type EnvVar string

// Value returns the value of the environment variable
func (e EnvVar) Value() string {
	return os.Getenv(string(e))
}

type Prompts struct {
	GenNewRecipeSys              string `json:"genNewRecipeSys"`
	GenNewRecipeUser             string `json:"genNewRecipeUser"`
	GenNewVisionImportArgsSys    string `json:"genNewVisionImportArgsSys"`
	GenNewVisionImportArgsUser   string `json:"genNewVisionImportArgsUser"`
	GenNewVisionImportRecipeSys  string `json:"genNewVisionImportRecipeSys"`
	GenNewVisionImportRecipeUser string `json:"genNewVisionImportRecipeUser"`
	ReGenRecipeSys               string `json:"reGenRecipeSys"`
	ReGenRecipeUser              string `json:"reGenRecipeUser"`
}

type OpenaiPrompt string

const (
	GenNewRecipeSys              OpenaiPrompt = "GenNewRecipeSys"
	GenNewRecipeUser             OpenaiPrompt = "GenNewRecipeUser"
	GenNewVisionImportArgsSys    OpenaiPrompt = "GenNewVisionImportArgsSys"
	GenNewVisionImportArgsUser   OpenaiPrompt = "GenNewVisionImportArgsUser"
	GenNewVisionImportRecipeSys  OpenaiPrompt = "GenNewVisionImportRecipeSys"
	GenNewVisionImportRecipeUser OpenaiPrompt = "GenNewVisionImportRecipeUser"
	ReGenRecipeSys               OpenaiPrompt = "ReGenRecipeSys"
	ReGenRecipeUser              OpenaiPrompt = "ReGenRecipeUser"
)

// LoadConfig reads a JSON configuration file and returns a Config struct.
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	config.RefreshAPIKeys()

	return &config, nil
}

// CheckConfigFields validates that all fields in Config are populated
// and their Value method (if available) will not return an error.
func CheckConfigFields(config *Config) error {
	return checkFieldsRecursive(reflect.ValueOf(config))
}

// checkFieldsRecursive recursively checks each field.
func checkFieldsRecursive(v reflect.Value) error {
	// Dereference pointer values
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Iterate through each field
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// Check for zero values
		if isZeroValue(field) {
			return fmt.Errorf("$%s must be set", fieldType.Name)
		}

		// If it's an EnvVar, check its Value()
		if field.Type().String() == "config.EnvVar" {
			envVar := EnvVar(field.String())
			envVal := envVar.Value()
			if envVal == "" {
				return fmt.Errorf("value of $%s must be set", fieldType.Name)
			}
		}

		// Recursively check nested structs
		if field.Kind() == reflect.Struct {
			if err := checkFieldsRecursive(field); err != nil {
				return err
			}
		}
	}
	return nil
}

// isZeroValue checks if the value is a zero value for its type.
func isZeroValue(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

// RefreshPrompts reads the prompts file and updates the Config struct.
func (c *Config) RefreshPrompts() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	prompts, err := loadPromptsFromFile(c.Env.OpenaiPromptsFilePath.Value())
	if err != nil {
		log.Printf("Unable to refresh prompts: %v", err)
		return
	}

	c.OpenaiPrompts = *prompts
}

// RefreshAPIKeys reads the API keys file and updates the Config struct.
func (c *Config) RefreshAPIKeys() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	keys := loadAPIKeysFromFile(c.Env.OpenaiKeysFilePath.Value())
	c.OpenaiKeys = keys
	c.CurrentOpenaiKeyIndex = 0
}

// GetCurrentAPIKey returns the current API key and rotates to the next one.
func (c *Config) GetCurrentAPIKey() string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	key := c.OpenaiKeys[c.CurrentOpenaiKeyIndex]
	c.rotateAPIKey()
	return key
}

// rotateAPIKey rotates to the next API key.
func (c *Config) rotateAPIKey() {
	c.CurrentOpenaiKeyIndex = (c.CurrentOpenaiKeyIndex + 1) % len(c.OpenaiKeys)
}

// loadAPIKeysFromFile reads a JSON file with API keys and returns a slice of keys.
func loadAPIKeysFromFile(filePath string) []string {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Unable to read API keys file: %v", err)
		return []string{}
	}

	var data struct {
		APIKeys []string `json:"api_keys"`
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Printf("Unable to parse API keys file: %v", err)
		return []string{}
	}

	return data.APIKeys
}

// loadPromptsFromFile reads a JSON file with prompt templates and returns a Prompts struct.
func loadPromptsFromFile(filePath string) (*Prompts, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var prompts Prompts
	if err := json.Unmarshal(file, &prompts); err != nil {
		return nil, err
	}

	return &prompts, nil
}

// FetchSysPrompt fetches a system prompt and replaces placeholders.
func (p *Prompts) FetchSysPrompt(promptType string, unitSystem string, requirements string) string {
	prompt := ""

	switch promptType {
	case "GenNewRecipeSys":
		prompt = p.GenNewRecipeSys
	case "GenNewVisionImportArgsSys":
		prompt = p.GenNewVisionImportArgsSys
	case "GenNewVisionImportRecipeSys":
		prompt = p.GenNewVisionImportRecipeSys
	case "ReGenRecipeSys":
		prompt = p.ReGenRecipeSys
	}

	prompt = strings.Replace(prompt, "{unitSystem}", unitSystem, -1)
	prompt = strings.Replace(prompt, "{requirements}", requirements, -1)

	return prompt
}

// FetchUserPrompt fetches a user prompt and replaces placeholders.
func (p *Prompts) FetchUserPrompt(promptType string, userPrompt string) string {
	prompt := ""

	switch promptType {
	case "GenNewRecipeUser":
		prompt = p.GenNewRecipeUser
	case "GenNewVisionImportArgsUser":
		prompt = p.GenNewVisionImportArgsUser
	case "GenNewVisionImportRecipeUser":
		prompt = p.GenNewVisionImportRecipeUser
	case "ReGenRecipeUser":
		prompt = p.ReGenRecipeUser
	}

	prompt = strings.Replace(prompt, "{userPrompt}", userPrompt, -1)

	return prompt
}
