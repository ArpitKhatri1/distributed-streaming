package handlers

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	awsfunctions "github.com/ArpitKhatri1/distributed-streaming/aws-functions"
	"github.com/aws/aws-sdk-go-v2/aws"
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

	presignClient := s3.NewPresignClient(awsfunctions.S3Client)
	safeFileName := url.PathEscape(req.FileName)
	key := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + safeFileName

	presignedURL, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(req.FileType),
	}, s3.WithPresignExpires(1*time.Minute))

	if err != nil {
		c.JSON(500, map[string]interface{}{"error": "Failed to create presigned URL"})
		return
	}
	fmt.Println(presignedURL.URL)

	c.JSON(200, PreSignedReponseType{URL: presignedURL.URL})

}
