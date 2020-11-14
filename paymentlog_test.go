package paymentlog

import (
	"testing"
)

var invalidMessage string = `
{
	"bad": "message"
}
`

var validMessage string = `
{
	"fromCustomer": "Jon",
	"fromAccount": "0010001",
	"toCustomer": "Cam",
	"toAccount": "0010002",
	"amount": 10.11
}
`

func TestValidate_withInvalidMsg(t *testing.T) {
	msg := PubSubMessage{Data: []byte(invalidMessage)}
	err := msg.validate()
	if err == nil {
		t.Errorf("Got err=%s; want errors", err)
	}
}

func TestValidate_withValidMsg(t *testing.T) {
	msg := PubSubMessage{Data: []byte(validMessage)}
	err := msg.validate()
	if err != nil {
		t.Errorf("Got err=%s; want no errors", err)
	}
}
