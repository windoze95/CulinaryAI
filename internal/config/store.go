package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// SSMService struct to hold the SSM client.
type SSMService struct {
	client *ssm.Client
}

// NewSSMService initializes a new SSMService with the provided AWS configuration.
func NewSSMService(region, accessKeyID, secretAccessKey string) (*SSMService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""))),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading configuration: %v", err)
	}

	ssmClient := ssm.NewFromConfig(cfg)
	return &SSMService{client: ssmClient}, nil
}

// GetParameter retrieves a parameter from AWS SSM Parameter Store.
func (s *SSMService) GetParameter(paramName string, withDecryption bool) (*ssm.GetParameterOutput, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: &withDecryption,
	}
	return s.client.GetParameter(context.TODO(), input)
}

// GetSecureStrings retrieves a list of SecureString parameters from AWS SSM Parameter Store.
func (s *SSMService) GetSecureParameterList(path string) ([]string, error) {
	var apiKeys []string

	withDecryption := true
	isRecursive := true
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		WithDecryption: &withDecryption,
		Recursive:      &isRecursive,
	}

	paginator := ssm.NewGetParametersByPathPaginator(s.client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("error fetching parameters page: %v", err)
		}

		for _, param := range page.Parameters {
			log.Printf("parameter name: %v", *param.Value)
			if param.Type == types.ParameterTypeSecureString {
				keys := strings.Split(*param.Value, ",")
				apiKeys = append(apiKeys, keys...)
			}
		}
	}

	return apiKeys, nil
}

// GetAllParameters retrieves all parameters from AWS SSM Parameter Store.
func (s *SSMService) GetOpenaiPromptsFromParameters(path string) (*OpenaiPrompts, error) {
	isRecursive := true
	input := &ssm.GetParametersByPathInput{
		Path: aws.String(path),
		// WithDecryption: true,
		Recursive: &isRecursive,
	}

	var prompts OpenaiPrompts
	paginator := ssm.NewGetParametersByPathPaginator(s.client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("error fetching parameters page: %v", err)
		}

		for _, param := range page.Parameters {
			if err := mapParameterToStructField(&prompts, param); err != nil {
				return nil, err
			}
		}
	}
	return &prompts, nil
}

// mapParameterToStructField maps a parameter to the corresponding field in the OpenaiPrompts struct
func mapParameterToStructField(prompts *OpenaiPrompts, param types.Parameter) error {
	paramJSON := fmt.Sprintf(`{"%s": %q}`, *param.Name, *param.Value)
	return json.Unmarshal([]byte(paramJSON), prompts)
}
