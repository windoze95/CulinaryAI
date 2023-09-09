package s3

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/windoze95/culinaryai/internal/config"
)

// UploadRecipeImageToS3 uploads a given byte array to an S3 bucket and returns the location URL.
func UploadRecipeImageToS3(cfg config.Config, imgBytes []byte) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Env.AWSRegion.Value()),
		Credentials: credentials.NewStaticCredentials(cfg.Env.AWSAccessKeyID.Value(), cfg.Env.AWSSecretAccessKey.Value(), ""),
	}))

	uploader := s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cfg.Env.S3Bucket.Value()),
		Key:    aws.String(cfg.Env.S3Key.Value()),
		Body:   bytes.NewReader(imgBytes),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	return result.Location, nil
}
