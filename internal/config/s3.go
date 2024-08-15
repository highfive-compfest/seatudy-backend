package config

import (
	"fmt"

	"mime/multipart"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

)
type FileUploader interface {
    UploadFile(key string, fileHeader *multipart.FileHeader) (string, error)
}	
type S3FileUploader struct {
	S3Service *s3.S3
}



// InitializeS3 initializes the S3 session and returns the S3 service client
func InitializeS3() (*S3FileUploader, error){
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Env.AwsRegion),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating AWS session: %v", err)
	}

	return &S3FileUploader{
		S3Service: s3.New(sess),
	}, nil

	
}

// UploadFile uploads a file to the specified S3 bucket
// UploadFile uploads a file to the specified S3 bucket
func (uploader *S3FileUploader) UploadFile(key string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("unable to open file %q: %v", fileHeader.Filename, err)
	}
	defer file.Close()

	encodedKey := url.PathEscape(key)

	_, err = uploader.S3Service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(Env.AwsBucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("unable to upload %q to %q: %v", fileHeader.Filename, Env.AwsBucketName, err)
	}

	// Construct the permanent URL
	urlStr := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", Env.AwsBucketName, aws.StringValue(uploader.S3Service.Config.Region), encodedKey)
	return urlStr, nil
}
