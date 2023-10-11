package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type Config struct {
	Title string `json:"title"`
	Env   Env    `json:"env"`
}

type EnvVar string

// Value returns the value of the environment variable in heroku
func (e EnvVar) Value() string {
	return os.Getenv(string(e))
}

type Env struct {
	Port                   EnvVar `json:"port"`
	DatabaseUrl            EnvVar `json:"databaseUrl"`
	OpenAIKeyEncryptionKey EnvVar `json:"openAIKeyEncryptionKey"`
	JwtSecretKey           EnvVar `json:"jwtSecretKey"`
	PublicOpenAIKey        EnvVar `json:"publicOpenAIKey"`
	AWSRegion              EnvVar `json:"awsRegion"`
	AWSAccessKeyID         EnvVar `json:"awsAccessKeyId"`
	AWSSecretAccessKey     EnvVar `json:"awsSecretAccessKey"`
	S3Bucket               EnvVar `json:"s3Bucket"`
	RecaptchaSecretKey     EnvVar `json:"recaptchaSecretKey"`
	FacebookClientID       EnvVar `json:"facebookClientId"`
	FacebookClientSecret   EnvVar `json:"facebookClientSecret"`
	FacebookRedirectURL    EnvVar `json:"facebookRedirectUrl"`
}

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
