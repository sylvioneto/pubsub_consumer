package main

import "testing"

var invalidMessage string = `
{
	"bad": "message"
}
`

var validMessage string = `
{
	"fromCustomer": "a",
	"fromAccount": "a",
	"toCustomer": "a",
	"toAccount": "a",
	"amount": 0
}
`

func TestValidateFail(t *testing.T) {
	msg := PubSubMessage{ID: "bad", Data: []byte(invalidMessage)}
	err := msg.validate()
	if err == nil {
        t.Errorf("TestValidateFail err=%s; want Error", err)
	}
}

func TestValidateSuccess(t *testing.T) {
	msg := PubSubMessage{ID: "good", Data: []byte(validMessage)}
	err := msg.validate()
	if err != nil {
        t.Errorf("TestValidateSuccess err=%s; want nil", err)
	}
}
