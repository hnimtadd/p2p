package test

import (
	"testing"

	"github.com/hnimtadd/p2p"
	"github.com/stretchr/testify/assert"
)

const (
	MessageTypeCustom p2p.MessageType = iota
)

func TestEncodingMessage(t *testing.T) {
	customMessage := p2p.Message{
		Payload: []byte("helloworld"),
		Header:  MessageTypeCustom,
	}

	buf := customMessage.Bytes()

	encodedMessage := new(p2p.Message)
	assert.Nil(t, encodedMessage.FromBytes(buf))
	assert.Equal(t, customMessage, *encodedMessage)
}
