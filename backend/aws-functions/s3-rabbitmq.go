package awsfunctions

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func ConnectS3ToRabbitMQ() {
	ctx := context.TODO()
	bucket := os.Getenv("BACKBLAZE_BUCKET_NAME")

	input := &s3.PutBucketNotificationConfigurationInput{
		Bucket: &bucket,
		NotificationConfiguration: &types.NotificationConfiguration{
			QueueConfigurations: []types.QueueConfiguration{
				{
					Id: aws.String("rabbitmq-events"),
					Events: []types.Event{
						types.EventS3ObjectCreatedPut,
						types.EventS3ObjectCreatedPost,
					},
					QueueArn: aws.String("arn:minio:sqs::1:amqp"), // MinIO accepts AMQP targets like SQS ARN format
				},
			},
		},
	}

	_, err := S3Client.PutBucketNotificationConfiguration(ctx, input)
	if err != nil {
		log.Fatalf("Failed to initalise bucket notification %v", err)
	}

	fmt.Println("Added notification system")

}
