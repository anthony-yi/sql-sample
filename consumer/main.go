package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/sirupsen/logrus"
)

func main() {
	// Define command line flags
	queueURL := flag.String("queue-url", "", "SQS Url")
	numThreads := flag.Int("threads", 2, "Number of consumers")
	maxMessages := flag.Int("max-messages", 1, "Max number of messages to receive")
	waitTimeSeconds := flag.Int("wait-time-seconds", 10, "Wait time in seconds")
	flag.Parse()

	// Load the AWS configuration from the currently active AWS profile
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error loading AWS config: %v", err)
	}

	// Create an SQS client
	client := sqs.NewFromConfig(cfg)

	// Create a wait group to wait for all consumer threads to finish
	var wg sync.WaitGroup

	// Create a channel to signal the consumer threads to stop
	stopCh := make(chan struct{})

	// Create a logger instance
	logger := logrus.New()

	// Subscribe to the FIFO queue with multiple consumer threads
	for i := 0; i < *numThreads; i++ {
		wg.Add(1)
		go startConsumer(i, client, *queueURL, *maxMessages, *waitTimeSeconds, &wg, stopCh, logger)
	}

	// Wait for a key press to exit
	fmt.Println("Press Enter to stop consumers...")
	fmt.Scanln()

	// Signal consumer threads to stop
	close(stopCh)

	// Wait for consumer threads to finish
	wg.Wait()
}

func startConsumer(id int, client *sqs.Client, queueURL string, maxMessages, waitTimeSeconds int, wg *sync.WaitGroup, stopCh <-chan struct{}, logger *logrus.Logger) {
	defer wg.Done()

	for {
		select {
		case <-stopCh:
			return
		default:

			input := &sqs.ReceiveMessageInput{
				QueueUrl:            &queueURL,
				MaxNumberOfMessages: int32(maxMessages),
				WaitTimeSeconds:     int32(waitTimeSeconds),
				AttributeNames: []types.QueueAttributeName{
					types.QueueAttributeNameAll,
				},
			}

			// Receive messages from the FIFO queue
			resp, err := client.ReceiveMessage(context.TODO(), input)
			if err != nil {
				logger.Error("Error receiving messages:", err)
				continue
			}

			// Process the received messages
			for _, msg := range resp.Messages {
				messageBody := *msg.Body
				messageGroupID := ""
				messageDeduplicationID := ""
				if msg.Attributes != nil {
					if attr, exists := msg.Attributes["MessageGroupId"]; exists {
						messageGroupID = attr
					}
					if attr, exists := msg.Attributes["MessageDeduplicationId"]; exists {
						messageDeduplicationID = attr
					}
				}
				// Add a timestamp to the log message
				logMessage := fmt.Sprintf("Consumer# %d Received message: %s (MessageGroupId: %s, MessageDeduplicationId: %s)", id, messageBody, messageGroupID, messageDeduplicationID)
				fmt.Println(logMessage)

				// Simulate processing time
				time.Sleep(2 * time.Second)

				// Delete the message from the queue
				_, err := client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
					QueueUrl:      &queueURL,
					ReceiptHandle: msg.ReceiptHandle,
				})
				if err != nil {
					logger.Error("Error deleting message:", err)
				}
			}
		}
	}
}
