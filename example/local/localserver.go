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
	errCh := make(chan error, 1)
	peerCh := make(chan p2p.Peer)
	infoCh := make(chan string)
	rpcCh := make(chan *p2p.RPC)

	tcpNode, err := p2p.NewTCPTransport(addr, peerCh, errCh, infoCh, rpcCh)
	if err != nil {
		log.Panicf("cannot init tcpTransport, err: %v", err)
	}
	tcpNode.Start()

	go func() {
		for _, addr := range seed {
			if err := tcpNode.Dial(addr); err != nil {
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
				fmt.Print(*info)
			}
			fmt.Println()
		case err := <-errCh:
			log.Panic(err)
		case peer := <-peerCh:
			fmt.Printf("[NODE %s] new peer: %s\n", addr, peer)
		case info := <-infoCh:
			fmt.Printf("[NODE %s] info:%s\n", addr, info)
		case rpc := <-rpcCh:
			fmt.Printf("[NODE %s] rcp: %v\n", addr, rpc)
		}
	}
}
