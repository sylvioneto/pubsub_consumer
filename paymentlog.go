package paymentlog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"github.com/xeipuuv/gojsonschema"
)

// PubSubMessage represents a Pub/Sub message.
type PubSubMessage struct {
	ID   string
	Data []byte
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

// ProcessLog consumes a Pub/Sub message from Payments Audit Log topic
func ProcessLog(ctx context.Context, m PubSubMessage) error {
	err := m.validate()
	if err != nil {
		println("messageId 1731286397516536: %s", err)
		return err
	}
	err = m.save()
	if err != nil {
		println("messageId 1731286397516536: %s", err)
		return err
	}
	return nil
}

// validate the payload
func (msg *PubSubMessage) validate() error {
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

func (msg *PubSubMessage) save() error {
	bucketName := os.Getenv("BUCKET")

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	obj := client.Bucket(bucketName).Object(msg.ID)
	wc := obj.NewWriter(ctx)
	if _, err := wc.Write(msg.Data); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
