package s3helper

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Bucket defines the expected JSON structure for S3 buckets.
// An S3 bucket can be used for source (input) and destination (output) files.
type S3Bucket struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

/*
S3Download downloads a file from an S3 bucket/key.
*/
func S3Download(file *os.File, bucket, key string) error {
	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	log.Println("Downloaded", numBytes, "bytes")
	return nil
}

/*
S3Upload uploads a file to an S3 bucket.
*/
func S3Upload(file *os.File, bucket, key string) error {
	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	log.Println("Successfully uploaded to", result.Location)
	return nil
}

// ParseFilenameFromKey parses the S3 filename from the key.
func ParseFilenameFromKey(key string) string {
	keySlice := strings.Split(key, "/")
	return keySlice[len(keySlice)-1]
}
