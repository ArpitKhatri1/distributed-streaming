package handlers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type PreSignedRequestType struct {
	FileName string `json:"filename"`
	FileType string `json:"filetype"`
}

type PreSignedReponseType struct {
	URL string `json:"url"`
}

func GetPresignedURL(c *gin.Context) {
	var req PreSignedRequestType

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	bucketName := os.Getenv("BACKBLAZE_BUCKET_NAME")
	endpoint := os.Getenv("BACKBLAZE_BUCKET_ENDPOINT")
	accessKey := os.Getenv("BACKBLAZE_BUCKET_KEY_ID")
	secretKey := os.Getenv("BACKBLAZE_BUCKET_KEY")
	region := os.Getenv("BACKBLAZE_BUCKET_REGION")

	if bucketName == "" || endpoint == "" || accessKey == "" || secretKey == "" || region == "" {
		c.JSON(500, gin.H{"error": "Missing Backblaze environment variables"})
		return
	}

	s3Client := s3.New(s3.Options{
		Credentials:      credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:           region,
		EndpointResolver: s3.EndpointResolverFromURL(endpoint),
	})

	presignClient := s3.NewPresignClient(s3Client)

	presignedURL, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(strconv.FormatInt(time.Now().UnixNano(), 10) + req.FileName), //avoid same file name overwriting
		ContentType: aws.String(req.FileType),
	}, s3.WithPresignExpires(1*time.Minute))

	if err != nil {
		c.JSON(500, map[string]interface{}{"error": "Failed to create presigned URL"})
		return
	}
	fmt.Println(presignedURL.URL)

	c.JSON(200, PreSignedReponseType{URL: presignedURL.URL})
}
