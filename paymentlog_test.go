package paymentlog

import (
	"net/http"
	"testing"
)

var invalidMessage string = `
{
	"bad": "message"
}
`

var validMessage string = `
{
	"fromCustomer": "Sylvio",
	"fromAccount": "0010001",
	"toCustomer": "Jessica",
	"toAccount": "0010002",
	"amount": 10.11
}
`

func TestValidate_withInvalidMsg(t *testing.T) {
	msg := PubSubMessage{ID: "bad", Data: []byte(invalidMessage)}
	err := msg.validate()
	if err == nil {
		t.Errorf("err=%s; want Error", err)
	}
}

func TestValidate_withValidMsg(t *testing.T) {
	msg := PubSubMessage{ID: "good", Data: []byte(validMessage)}
	err := msg.validate()
	if err != nil {
		t.Errorf("err=%s; want nil", err)
	}
}

func TestProcessLog_withValidMsg(t *testing.T) {
	msg := PubSubMessage{ID: "good", Data: []byte(validMessage)}
	httpStatus := ProcessLog(nil, msg)
	if httpStatus != http.StatusAccepted {
		t.Errorf("Status=%d; want %d", httpStatus, http.StatusAccepted)
	}
}

func TestProcessLog_withInvalidMsg(t *testing.T) {
	msg := PubSubMessage{ID: "bad", Data: []byte(invalidMessage)}
	httpStatus := ProcessLog(nil, msg)
	if httpStatus != http.StatusBadRequest {
		t.Errorf("Status=%d; want %d", httpStatus, http.StatusBadRequest)
	}
}
