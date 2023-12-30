package p2p

type Transport interface {
	Addr() NetAddr
	Start()
	Dial(NetAddr) error
	Broadcast(payload []byte) error
	AddPeer(Peer) error
	RemovePeer(Peer) error
	Peers() []*PeerInfo
	PeerCount() int
}

type Peer interface {
	Addr() NetAddr
	Info() *PeerInfo
	Accept(payload []byte) error
}
