package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	println("Teststing consumer")
	err := test("spedroza-sandbox", "test")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func test(projectID string, topicID string) error {

	// client
	log.Println("Create PubSub client")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	// create subscription
	log.Printf("Get topic %v", topicID)
	topic := client.Topic(topicID)
	subID := "sub-test"
	log.Printf("Creating sub %v", subID)
	sub, err := client.CreateSubscription(context.Background(), subID,
		pubsub.SubscriptionConfig{Topic: topic})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			log.Println("Sub already exists")
			sub = client.Subscription(subID)
		} else {
			return fmt.Errorf("Subscription: %v", err)
		}
	}

	// pull messages
	log.Println("Pulling messages")
	received := 0
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		log.Printf("messageId: %s - data: %s", m.ID, m.Data)
		msg := PubSubMessage{ID: m.ID, Data: m.Data}
		err := msg.validate()
		if err != nil {
			log.Printf("validate: %s", err)
		}
		m.Ack()
		received++
		if received == 10 {
			cancel()
		}
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}
