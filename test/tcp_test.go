package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hnimtadd/p2p"
	"github.com/stretchr/testify/assert"
)

func TestTCP(t *testing.T) {
	n1, p1, e1, m1, r1 := MakeTCPTRansport(t, ":9999")
	n2, p2, e2, m2, r2 := MakeTCPTRansport(t, ":9998")

	go MonitorNode(n1.Addr().(string), p1, e1, m1, r1)
	go MonitorNode(n2.Addr().(string), p2, e2, m2, r2)

	peer2, err := n1.Dial(n2.Addr())
	assert.Nil(t, err)
	assert.NotNil(t, peer2)

	assert.Nil(t, n1.Send(n2.Addr(), []byte("hello")))
	time.Sleep(time.Second * 5)
}

func MonitorNode(
	prefix string,
	peerCh <-chan p2p.Peer,
	errCh <-chan error,
	msgCh <-chan string,
	rpcCh <-chan *p2p.RPC,
) {
	for {
		select {
		case peer := <-peerCh:
			fmt.Println(prefix, "new peer", peer)
		case err := <-errCh:
			fmt.Println(prefix, "new error", err)
		case msg := <-msgCh:
			fmt.Println(prefix, "new msg", msg)
		case rpc := <-rpcCh:
			fmt.Println(prefix, "new rpc", rpc)
		}
	}
}

func MakeTCPTRansport(
	t *testing.T,
	addr p2p.NetAddr,
) (p2p.Transport,
	<-chan p2p.Peer,
	<-chan error,
	<-chan string,
	<-chan *p2p.RPC,
) {
	node, err := p2p.NewTCPTransport(addr)
	assert.Nil(t, err)
	peerCh := node.ConsumePeer()
	msgCh := node.ConsumeInfo()
	errCh := node.ConsumeError()
	rpcCh := node.ConsumeRPC()
	go node.Start()
	return node, peerCh, errCh, msgCh, rpcCh
}
