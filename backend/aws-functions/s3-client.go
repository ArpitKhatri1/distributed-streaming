package awsfunctions

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func CreateS3Client() {
	bucketName := os.Getenv("BACKBLAZE_BUCKET_NAME")
	endpoint := os.Getenv("BACKBLAZE_BUCKET_ENDPOINT")
	accessKey := os.Getenv("BACKBLAZE_BUCKET_KEY_ID")
	secretKey := os.Getenv("BACKBLAZE_BUCKET_KEY")
	region := os.Getenv("BACKBLAZE_BUCKET_REGION")

	if bucketName == "" || endpoint == "" || accessKey == "" || secretKey == "" || region == "" {
		log.Fatal("S3 configuration cannot be empty")
	}

	S3Client = s3.New(s3.Options{
		Credentials:      credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:           region,
		EndpointResolver: s3.EndpointResolverFromURL(endpoint), // Endpoint resolver is deprecated fix it.
		UsePathStyle:     true,
	})

}
