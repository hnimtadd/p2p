package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hnimtadd/p2p"
)

func main() {
	go MakeNodeAndRun(":3000", nil)
	time.Sleep(time.Second)
	go MakeNodeAndRun(":3001", []p2p.NetAddr{":3000"})
	select {}
}

func MakeNodeAndRun(addr p2p.NetAddr, seed []p2p.NetAddr) {
	fmt.Printf("[NODE %s] Start\n", addr)

	tcpNode, err := p2p.NewTCPTransport(addr)
	if err != nil {
		log.Panicf("cannot init tcpTransport, err: %v", err)
	}
	errCh := tcpNode.ConsumeError()
	peerCh := tcpNode.ConsumePeer()
	infoCh := tcpNode.ConsumeInfo()
	rpcCh := tcpNode.ConsumeRPC()
	tcpNode.Start()

	go func() {
		for _, addr := range seed {
			if _, err := tcpNode.Dial(addr); err != nil {
				log.Panicf("cannot dial to node, err: %v", err)
			}
		}
	}()
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			fmt.Printf("[NODE %s] current peers: ", addr)
			fmt.Printf("len: %d ", tcpNode.PeerCount())
			for _, info := range tcpNode.Peers() {
				if err := tcpNode.Send(info.Addr, []byte("hello")); err != nil {
					panic(err)
				}
			}
			fmt.Println("sended to peer")
		case err := <-errCh:
			log.Panic(err)
		case peer := <-peerCh:
			peer.Start()
			fmt.Printf("[NODE %s] new peer: %s\n", addr, peer)
		case info := <-infoCh:
			fmt.Printf("[NODE %s] info:%s\n", addr, info)
		case rpc := <-rpcCh:
			fmt.Printf("[NODE %s] rcp: %s\n", addr, string(rpc.Payload))
		}
	}
}
