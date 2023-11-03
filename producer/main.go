package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	// Define command line flags
	queueURL := flag.String("queue-url", "", "SQS Url")
	numGenerate := flag.Int("num", 10, "Number of messages to generate")
	randomizeIDs := flag.Bool("randomize-ids", false, "Randomize MessageGroupId and MessageDeduplicationId")
	flag.Parse()

	// Load the AWS configuration from the currently active AWS profile
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error loading AWS config: %v", err)
	}

	// Create an SQS client
	client := sqs.NewFromConfig(cfg)

	// Loop to send multiple dummy messages
	for i := 1; i <= *numGenerate; i++ {
		// Include the current timestamp in the message body
		now := time.Now()
		currentTime := now.Format(time.RFC3339)
		currentTimeUnix := now.Unix()
		messageBody := fmt.Sprintf("MessageID %d - Generated at: %s", currentTimeUnix, currentTime)

		// Determine the values for MessageGroupId and MessageDeduplicationId
		var messageGroupID, messageDeduplicationID string
		if *randomizeIDs {
			messageGroupID = fmt.Sprintf("Group%d", rand.Intn(100))           // Randomize MessageGroupId
			messageDeduplicationID = fmt.Sprintf("%d", time.Now().UnixNano()) // Randomize MessageDeduplicationId
		} else {
			messageGroupID = "Group1"                                 // Fixed value for MessageGroupId
			messageDeduplicationID = fmt.Sprintf("FixedDedupID%d", i) // Fixed value for MessageDeduplicationId
		}

		// Send a message to the SQS queue
		_, err := client.SendMessage(context.TODO(), &sqs.SendMessageInput{
			QueueUrl:               queueURL,
			MessageBody:            &messageBody,
			MessageGroupId:         &messageGroupID,
			MessageDeduplicationId: &messageDeduplicationID,
		})
		if err != nil {
			log.Printf("Error sending message %d: %v", i, err)
		} else {
			log.Printf("Sent message %s (MessageGroupId: %s, MessageDeduplicationId: %s)", messageBody, messageGroupID, messageDeduplicationID)
		}

		// Sleep to avoid exceeding the rate limit for sending messages
		time.Sleep(1000 * time.Millisecond)
	}
}
