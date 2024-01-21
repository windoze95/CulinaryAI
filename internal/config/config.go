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
	Env Env `json:"env"`
	// Prompts are actually the templates to construct the usable prompts.
	// Use the FillSysPrompt and FillUserPrompt methods to retrieve a prompt.
	OpenaiPrompts         OpenaiPrompts `json:"openai_prompts"`
	OpenaiKeys            []string      `json:"openai_keys"`
	CurrentOpenaiKeyIndex int
	Mutex                 sync.RWMutex
}

type Env struct {
	Port                   EnvVar `json:"port"`
	DatabaseUrl            EnvVar `json:"database_url"`
	OpenAIKeyEncryptionKey EnvVar `json:"openai_key_encryption_key"`
	JwtSecretKey           EnvVar `json:"jwt_secret_key"`
	PublicOpenAIKey        EnvVar `json:"public_openai_key"`
	AWSRegion              EnvVar `json:"aws_region"`
	AWSAccessKeyID         EnvVar `json:"aws_access_key_id"`
	AWSSecretAccessKey     EnvVar `json:"aws_secret_access_key"`
	S3Bucket               EnvVar `json:"s3_bucket"`
	OpenaiPromptsFilePath  EnvVar `json:"openai_prompts_file_path"`
	OpenaiKeysFilePath     EnvVar `json:"openai_keys_file_path"`
}

type EnvVar string

// Value returns the value of the environment variable
func (e EnvVar) Value() string {
	return os.Getenv(string(e))
}

// Prompts are actually the templates to construct the usable prompts.
// Use the FillSysPrompt and FillUserPrompt methods to retrieve a prompt.
type OpenaiPrompts struct {
	GenNewRecipeSys              OpenaiPromptTemplate `json:"gen_new_recipe_sys"`
	GenNewRecipeUser             OpenaiPromptTemplate `json:"gen_new_recipe_user"`
	GenNewVisionImportArgsSys    OpenaiPromptTemplate `json:"gen_new_vision_import_args_sys"`
	GenNewVisionImportArgsUser   OpenaiPromptTemplate `json:"gen_new_vision_import_args_user"`
	GenNewVisionImportRecipeSys  OpenaiPromptTemplate `json:"gen_new_vision_import_recipe_sys"`
	GenNewVisionImportRecipeUser OpenaiPromptTemplate `json:"gen_new_vision_import_recipe_user"`
	RegenRecipeSys               OpenaiPromptTemplate `json:"regen_recipe_sys"`
	RegenRecipeUser              OpenaiPromptTemplate `json:"regen_recipe_user"`
}

type OpenaiPromptTemplate string

// LoadConfig reads a JSON configuration file and returns a Config struct.
// Other (all untracked) config files are expected to maintain the same format as config.json.
// Example:
// {
//   "openai_prompts": {
//     "gen_new_recipe_sys": ...
// },
//
// {
//   "openai_keys": [
//     "key1",
//     ...
//   ]
// }
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
	config.RefreshPrompts()

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
		OpenaiKeys []string `json:"openai_keys"`
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Printf("Unable to parse API keys file: %v", err)
		return []string{}
	}

	return data.OpenaiKeys
}

// loadPromptsFromFile reads a JSON file with prompt templates and returns a Prompts struct.
func loadPromptsFromFile(filePath string) (*OpenaiPrompts, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data struct {
		OpenaiPrompts OpenaiPrompts `json:"openai_prompts"`
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &data.OpenaiPrompts, nil
}

// FillSysPrompt fetches a system prompt and replaces placeholders.
func (p *OpenaiPrompts) FillSysPrompt(promptTemplate OpenaiPromptTemplate, unitSystem string, requirements string) string {
	prompt := string(promptTemplate)

	prompt = strings.Replace(prompt, "{unitSystem}", unitSystem, -1)
	prompt = strings.Replace(prompt, "{requirements}", requirements, -1)

	return prompt
}

// FillUserPrompt fetches a user prompt and replaces placeholders.
func (p *OpenaiPrompts) FillUserPrompt(promptTemplate OpenaiPromptTemplate, userPrompt string) string {
	prompt := string(promptTemplate)

	prompt = strings.Replace(prompt, "{userPrompt}", userPrompt, -1)

	return prompt
}
