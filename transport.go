package p2p

type Transport interface {
	Addr() NetAddr
	Start()
	Dial(NetAddr) error
	Broadcast(payload []byte) error
}

type Peer interface {
	Addr() NetAddr
	Accept(payload []byte) error
}
