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

var schemaData string = `
{	
	"$schema": "http://json-schema.org/draft/2019-09/schema#",
	"type": "object",
	"properties": {
	  "fromCustomer":    { "type": "string" },
	  "toCustomer":      { "type": "string" },
	  "fromAccount":   	 { "type": "string" },
	  "toAccount": 		 { "type": "string" },
	  "amount": 		 { "type": "number", "minimum": 0 },
	  "transactionDate": { "type": "string", "format": "date-time" }
	},
	"required": ["fromCustomer", "toCustomer", "fromAccount", "toAccount", "amount"]
}
`

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

	schemaLoader := gojsonschema.NewStringLoader(schemaData)
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

/*
func (msg *PubSubMessage) unmarshall() (PaymentMessage, error) {
	var pMsg PaymentMessage
	err := json.Unmarshal(msg.Data, &pMsg)
	if err != nil {
		return pMsg, fmt.Errorf("Unmarshal error %s", err)
	}
	return pMsg, nil
}
*/
