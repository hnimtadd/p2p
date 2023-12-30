package p2p

import "fmt"

type NodeInformation struct {
	NodeID string
}

type PeerInfo struct {
	NodeID string
}

func (p PeerInfo) String() string {
	return fmt.Sprintf("[nodeID: %s]", p.NodeID)
}
