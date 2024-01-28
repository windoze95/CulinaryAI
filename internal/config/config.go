package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
)

// Config struct to hold the configuration.
type Config struct {
	Env Env `json:"env"`
	// Prompts are actually the templates to construct the usable prompts.
	// Use the FillSysPrompt and FillUserPrompt methods to retrieve a prompt.
	OpenaiPrompts         OpenaiPrompts `json:"openai_prompts"`
	OpenaiKeys            []string      `json:"openai_keys"`
	CurrentOpenaiKeyIndex int
	Mutex                 sync.RWMutex
}

// Env struct to hold the environment variables.
type Env struct {
	Port               EnvVar `json:"port"`
	DatabaseUrl        EnvVar `json:"database_url"`
	JwtSecretKey       EnvVar `json:"jwt_secret_key"`
	AWSRegion          EnvVar `json:"aws_region"`
	AWSAccessKeyID     EnvVar `json:"aws_access_key_id"`
	AWSSecretAccessKey EnvVar `json:"aws_secret_access_key"`
	S3Bucket           EnvVar `json:"s3_bucket"`
	IdHeader           EnvVar `json:"id_header"`
	OpenaiPromptsPath  EnvVar `json:"openai_prompts_path"`
	OpenaiKeysPath     EnvVar `json:"openai_keys_path"`
}

// EnvVar is a string that represents an environment variable.
type EnvVar string

// Value returns the value of the environment variable.
func (e EnvVar) Value() string {
	return os.Getenv(string(e))
}

// Prompts are actually the templates to construct the usable prompts.
// Use the FillSysPrompt and FillUserPrompt methods to retrieve a prompt.
type OpenaiPrompts struct {
	GenNewRecipeSys              OpenaiPromptTemplate `json:"/saltybytes/openai_prompts/gen_new_recipe_sys"`
	GenNewRecipeUser             OpenaiPromptTemplate `json:"/saltybytes/openai_prompts/gen_new_recipe_user"`
	GenNewVisionImportArgsSys    OpenaiPromptTemplate `json:"/saltybytes/openai_prompts/gen_new_vision_import_args_sys"`
	GenNewVisionImportArgsUser   OpenaiPromptTemplate `json:"/saltybytes/openai_prompts/gen_new_vision_import_args_user"`
	GenNewVisionImportRecipeSys  OpenaiPromptTemplate `json:"/saltybytes/openai_prompts/gen_new_vision_import_recipe_sys"`
	GenNewVisionImportRecipeUser OpenaiPromptTemplate `json:"/saltybytes/openai_prompts/gen_new_vision_import_recipe_user"`
	RegenRecipeSys               OpenaiPromptTemplate `json:"/saltybytes/openai_prompts/regen_recipe_sys"`
	RegenRecipeUser              OpenaiPromptTemplate `json:"/saltybytes/openai_prompts/regen_recipe_user"`
}

// OpenaiPromptTemplate is a string that represents an OpenAI prompt template.
type OpenaiPromptTemplate string

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

	return &config, nil
}

// CheckConfigFields validates that all fields in Config are populated
// and their Value method (if available) will not return an error.
func (c *Config) CheckConfigEnvFields() error {
	return checkFieldsRecursive(reflect.ValueOf(c.Env))
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

// LoadOpenaiKeys loads all OpenAI API keys from AWS SSM Parameter Store.
func (c *Config) LoadOpenaiKeys() error {
	// Initialize SSMService with AWS configuration
	ssmService, err := NewSSMService(c.Env.AWSRegion.Value(), c.Env.AWSAccessKeyID.Value(), c.Env.AWSSecretAccessKey.Value())
	if err != nil {
		return fmt.Errorf("failed to create SSM service: %v", err)
	}

	// Fetch Secure Strings
	apiKeys, err := ssmService.GetSecureParameterList(c.Env.OpenaiKeysPath.Value())
	if err != nil {
		return fmt.Errorf("failed to get all parameters: %v", err)
	}

	log.Printf("api keys path: %v", c.Env.OpenaiKeysPath.Value())
	log.Printf("api keys: %v", apiKeys)

	c.OpenaiKeys = apiKeys

	return nil
}

// LoadOpenaiPrompts loads all OpenAI prompts from AWS SSM Parameter Store.
func (c *Config) LoadOpenaiPrompts() error {
	// Initialize SSMService with AWS configuration
	ssmService, err := NewSSMService(c.Env.AWSRegion.Value(), c.Env.AWSAccessKeyID.Value(), c.Env.AWSSecretAccessKey.Value())
	if err != nil {
		return fmt.Errorf("failed to create SSM service: %v", err)
	}

	prompts, err := ssmService.GetOpenaiPromptsFromParameters(c.Env.OpenaiPromptsPath.Value())
	if err != nil {
		return fmt.Errorf("failed to get all parameters: %v", err)
	}

	c.OpenaiPrompts = *prompts

	return nil
}

// FillSysPrompt fetches a system prompt and replaces placeholders.
func (p *OpenaiPrompts) FillSysPrompt(promptTemplate OpenaiPromptTemplate, unitSystem string, requirements string) string {
	prompt := string(promptTemplate)

	sanitizedRequirements := strings.Replace(requirements, "`", "", -1)

	prompt = strings.Replace(prompt, "{unitSystem}", unitSystem, -1)
	prompt = strings.Replace(prompt, "{requirements}", sanitizedRequirements, 1)

	return prompt
}

// FillUserPrompt fetches a user prompt and replaces placeholders.
func (p *OpenaiPrompts) FillUserPrompt(promptTemplate OpenaiPromptTemplate, userPrompt string) string {
	prompt := string(promptTemplate)

	sanitizedUserPrompt := strings.Replace(userPrompt, "`", "", -1)

	prompt = strings.Replace(prompt, "{userPrompt}", sanitizedUserPrompt, 1)

	return prompt
}
