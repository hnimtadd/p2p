package p2p

type Transport interface {
	Addr() NetAddr
	Start()
	// Dial connect to peer at addr, make handshake and add connection to this node's peer list if success.
	Dial(addr NetAddr) error
	// Broadcast wraps payload in RPC and send to all peers.
	Broadcast(payload []byte) error
	// Send wrap payload in RPC and send to given peer.
	Send(to NetAddr, payload []byte) error
	// AddPeer adds make handshake with peer and connection to this node's peers list.
	AddPeer(peer Peer) error
	// RemovePeer removes given peer from this node's peers list.
	RemovePeer(Peer) error
	// Peers return all connected peers's infomation.
	Peers() []*PeerInfo
	// PeerCount return number of connected peers of this node.
	PeerCount() int
}

type Peer interface {
	Addr() NetAddr
	Info() *PeerInfo
	Start()
	Stop()
	// Accept write given payload into this peer's connection.
	Accept(payload []byte) error
}
