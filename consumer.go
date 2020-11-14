package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/xeipuuv/gojsonschema"
)

// PubSubMessage represents a Pub/Sub message.
type PubSubMessage struct {
	ID   string
	Data []byte
}

// PaymentMessage represents a payment log message.
type PaymentMessage struct {
	Name string `json:"name"`
}

var schemaData = map[string]interface{}{
	"type":     "object",
	"required": []string{"fromCustomer", "toCustomer"},
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {
	name := string(m.Data) // Automatically decoded from base64.
	if name == "" {
		name = "World"
	}
	log.Printf("Hello, %s!", name)
	return nil
}

func (msg *PubSubMessage) validate() error {
	//log.Printf("This message is valid!")
	if !json.Valid([]byte(msg.Data)) {
		return fmt.Errorf("Invalid message")
	}

	var pMsg PaymentMessage
	err := json.Unmarshal(msg.Data, &pMsg)
	if err != nil {
		return fmt.Errorf("Unmarshal error %s", err)
	}

	schemaLoader := gojsonschema.NewGoLoader(schemaData)
	documentLoader := gojsonschema.NewBytesLoader(msg.Data)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("Schema validation error %s: ", err.Error())
	}

	if !result.Valid() {
		return fmt.Errorf("The document is not valid %s", result.Errors())
	}

	return nil
}
