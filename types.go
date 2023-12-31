package p2p

import "fmt"

type NodeInformation struct {
	Addr   NetAddr
	NodeID string
}

type PeerInfo struct {
	Addr   NetAddr
	NodeID string
}

func (p PeerInfo) String() string {
	return fmt.Sprintf("[nodeID: %s Addr: %s]", p.NodeID, p.Addr)
}
