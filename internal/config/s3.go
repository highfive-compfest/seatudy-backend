package config

import (
	"fmt"
	"log"
	"mime/multipart"
	

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
func UploadFile(bucket, key string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("unable to open file %q: %v", fileHeader.Filename, err)
	}
	defer file.Close()

	_, err = S3Svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("unable to upload %q to %q: %v", fileHeader.Filename, bucket, err)
	}

	// Construct the permanent URL
	urlStr := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, aws.StringValue(S3Svc.Config.Region), key)
	return urlStr, nil
}
