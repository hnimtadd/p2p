package p2p

type Transport interface {
	Addr() NetAddr
	Start()
	// Dial connect to peer at addr, make handshake and add connection to this node's peer list if success.
	Dial(addr NetAddr) (Peer, error)
	// ConsumePeer return receive only peerCh of this transport, when this transport accept new per connection, peer will be sent to this channel
	ConsumePeer() <-chan Peer
	// ConsumeRPC return rpc receive only channel of this transport, when there are new rpc from the connection, rpc will be sent to this channel
	ConsumeRPC() <-chan *RPC
	// ConsumeError return error receive only channel of this tranpsport, when there are error from underlying transport, error will be sent to this channel
	ConsumeError() <-chan error
	// ConsumeInfo return string receive only channel of this tranport, when there are any message from tranpoort, message will be sent to this channel
	ConsumeInfo() <-chan string
	// Broadcast wraps payload in RPC and send to all connected peers.
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
