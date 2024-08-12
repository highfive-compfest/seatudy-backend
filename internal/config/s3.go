package config

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3Svc *s3.S3

// InitializeS3 initializes the S3 session and returns the S3 service client
func InitializeS3() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Env.AwsRegion),
	})
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	S3Svc = s3.New(sess)
}

// UploadFile uploads a file to the specified S3 bucket
// UploadFile uploads a file to the specified S3 bucket
func UploadFile(key string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("unable to open file %q: %v", fileHeader.Filename, err)
	}
	defer file.Close()

	encodedKey := url.PathEscape(key)

	_, err = S3Svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(Env.AwsBucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("unable to upload %q to %q: %v", fileHeader.Filename, Env.AwsBucketName, err)
	}

	// Construct the permanent URL
	urlStr := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", Env.AwsBucketName, aws.StringValue(S3Svc.Config.Region), encodedKey)
	return urlStr, nil
}
