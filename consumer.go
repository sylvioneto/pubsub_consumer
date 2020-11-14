package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
func ProcessLog(ctx context.Context, m PubSubMessage) int {
	log.Printf("messageId: %s", string(m.ID))
	err := m.validate()
	if err != nil {
		log.Printf("messageId 1731286397516536: %s", err)
		return http.StatusBadRequest
	}
	return http.StatusOK
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
