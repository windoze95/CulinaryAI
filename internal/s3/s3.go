package s3

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/windoze95/saltybytes-api/internal/config"
)

// UploadRecipeImageToS3 uploads a given byte array to an S3 bucket and returns the location URL.
func UploadRecipeImageToS3(cfg *config.Config, imgBytes []byte, s3Key string) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Env.AWSRegion.Value()),
		Credentials: credentials.NewStaticCredentials(cfg.Env.AWSAccessKeyID.Value(), cfg.Env.AWSSecretAccessKey.Value(), ""),
	}))

	uploader := s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cfg.Env.S3Bucket.Value()),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(imgBytes),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	return result.Location, nil
}

// DeleteRecipeImageFromS3 deletes a given image from an S3 bucket.
func DeleteRecipeImageFromS3(cfg *config.Config, s3Key string) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Env.AWSRegion.Value()),
		Credentials: credentials.NewStaticCredentials(cfg.Env.AWSAccessKeyID.Value(), cfg.Env.AWSSecretAccessKey.Value(), ""),
	}))

	deleter := s3.New(sess)

	_, err := deleter.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(cfg.Env.S3Bucket.Value()),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %v", err)
	}

	return nil
}

// GenerateS3Key generates the S3 key for a recipe image, given the recipe ID.
func GenerateS3Key(recipeID uint) string {
	return fmt.Sprintf("recipes/%d/images/recipe_image_%d.jpg", recipeID, recipeID)
}
