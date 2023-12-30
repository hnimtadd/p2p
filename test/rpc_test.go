package test

import (
	"testing"

	"github.com/hnimtadd/p2p"
	"github.com/stretchr/testify/assert"
)

func TestRPCWithMessage(t *testing.T) {
	msg := p2p.Message{
		Header:  1,
		Payload: []byte("helloworld"),
	}
	rpc := p2p.RPC{
		From:    nil,
		Payload: msg.Bytes(),
	}

	buf := rpc.Bytes()

	encodedRPC := new(p2p.RPC)
	assert.Nil(t, encodedRPC.FromBytes(buf))
	assert.Equal(t, rpc, *encodedRPC)

	encodedMessage := new(p2p.Message)
	assert.Nil(t, encodedMessage.FromBytes(encodedRPC.Payload))
	assert.Equal(t, msg, *encodedMessage)
}
