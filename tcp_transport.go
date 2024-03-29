package p2p

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
)

var (
	ErrPeerExisted    = errors.New("duplicate peer")
	ErrPeerInvalid    = errors.New("peer with type invalid")
	ErrPeerNotExisted = errors.New("given peer not found")
)

// ################### TRANSPORT ########################

type TCPTransport struct {
	addr   NetAddr
	ln     net.Listener
	peers  map[NetAddr]*TCPPeer // map nodeId with peer present for that node
	peerCh chan Peer
	errCh  chan error
	infoCh chan string
	rpcCh  chan *RPC
	nodeId string
	mu     sync.RWMutex
}

func NewTCPTransport(
	addr NetAddr,
	// peerCh chan<- Peer,
	// errCh chan<- error,
	// notiCh chan<- string,
	// rpcCh chan<- *RPC,
) (Transport, error) {
	ln, err := net.Listen("tcp", addr.(string))
	if err != nil {
		return nil, fmt.Errorf("tcp transport: cannot init tcp transport, err: %v", err.Error())
	}
	transport := &TCPTransport{
		peers:  make(map[NetAddr]*TCPPeer, 1024),
		peerCh: make(chan Peer, 1024),
		errCh:  make(chan error, 1024),
		infoCh: make(chan string, 1024),
		rpcCh:  make(chan *RPC, 1024),
		// peerCh: peerCh,
		// errCh:  errCh,
		// notiCh: notiCh,
		// rpcCh:  rpcCh,
		ln:   ln,
		addr: addr,
	}
	return transport, nil
}

// ###################### Interface implement ####################

func (t *TCPTransport) Addr() NetAddr {
	return t.addr
}

func (t *TCPTransport) Start() {
	go t.Loop()
}

func (t *TCPTransport) ConsumePeer() <-chan Peer {
	return t.peerCh
}

func (t *TCPTransport) ConsumeRPC() <-chan *RPC {
	return t.rpcCh
}

func (t *TCPTransport) ConsumeError() <-chan error {
	return t.errCh
}

func (t *TCPTransport) ConsumeInfo() <-chan string {
	return t.infoCh
}

func (t *TCPTransport) Peers() []*PeerInfo {
	peerInfos := []*PeerInfo{}
	for _, peer := range t.peers {
		peerInfos = append(peerInfos, peer.Info())
	}
	return peerInfos
}

func (t *TCPTransport) PeerCount() int {
	return len(t.peers)
}

func (t *TCPTransport) AddPeer(peer Peer) error {
	tcpPeer, ok := peer.(*TCPPeer)
	if !ok {
		return ErrPeerInvalid
	}
	return t.connect(tcpPeer)
}

func (t *TCPTransport) RemovePeer(peer Peer) error {
	tcpPeer, ok := peer.(*TCPPeer)
	if !ok {
		return ErrPeerInvalid
	}
	t.disconect(tcpPeer)
	return nil
}

// Broadcast wraps payload in RPC and send to all peer
func (t *TCPTransport) Broadcast(payload []byte) error {
	t.mu.Lock()
	rpc := &RPC{
		From:    t.Addr(),
		Payload: payload,
	}
	defer t.mu.Unlock()
	for _, peer := range t.peers {
		if err := peer.Accept(rpc.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// Send wrap payload in RPC and send to given peer
func (t *TCPTransport) Send(to NetAddr, payload []byte) error {
	t.mu.Lock()
	peer, ok := t.peers[to]
	t.mu.Unlock()
	if !ok {
		return ErrPeerNotExisted
	}
	rpc := &RPC{
		From:    t.Addr(),
		Payload: payload,
	}
	if err := peer.Accept(rpc.Bytes()); err != nil {
		return err
	}
	return nil
}

func (t *TCPTransport) Dial(addr NetAddr) (Peer, error) {
	conn, err := net.Dial("tcp", addr.(string))
	if err != nil {
		return nil, fmt.Errorf("cannot dial to given addr %s, err: %s", addr.(string), err.Error())
	}

	peerInfo, err := TCPHandshake(t, conn)
	if err != nil {
		return nil, err
	}
	peer, err := NewTCPPeer(conn, peerInfo.NodeID, false, t.rpcCh, t.infoCh)
	if err != nil {
		return nil, err
	}
	err = t.connect(peer)
	if err != nil {
		return peer, fmt.Errorf("cannot make peer with addr %s, err: %s", addr.(string), err.Error())
	}
	return peer, nil
}

// ###################### End-Interface implement ####################

func (t *TCPTransport) Accept(conn net.Conn) error {
	peerInfo, err := TCPHandshakeReply(conn, t)
	if err != nil {
		return err
	}
	peer, err := NewTCPPeer(conn, peerInfo.NodeID, true, t.rpcCh, t.infoCh)
	if err != nil {
		return err
	}
	if err := t.connect(peer); err != nil {
		return err
	}
	return nil
}

func (t *TCPTransport) Loop() {
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			t.errCh <- err
			return
		}
		err = t.Accept(conn)
		if err != nil {
			t.errCh <- err
		}
	}
}

func (t *TCPTransport) connect(p *TCPPeer) error {
	t.mu.Lock()
	_, ok := t.peers[p.Addr()]
	t.mu.Unlock()
	if ok {
		return ErrPeerExisted
	}
	t.mu.RLock()
	t.peers[p.Addr()] = p
	t.mu.RUnlock()
	t.peerCh <- p
	return nil
}

func (t *TCPTransport) disconect(p *TCPPeer) {
	t.mu.RLock()
	delete(t.peers, p.Addr())
	t.mu.RUnlock()
	go p.Stop()
}

// ################### END-TRANSPORT ########################

//##################- PEER -#########################

type TCPPeer struct {
	addr    NetAddr
	conn    net.Conn
	rpcCh   chan<- *RPC
	infoCh  chan<- string
	nodeID  string
	inbound bool
}

func NewTCPPeer(
	conn net.Conn,
	nodeId string,
	inbound bool,
	rpcCh chan<- *RPC,
	infoCh chan<- string,
) (*TCPPeer, error) {
	peer := &TCPPeer{
		nodeID:  nodeId,
		rpcCh:   rpcCh,
		infoCh:  infoCh,
		inbound: inbound,
		conn:    conn,
	}
	return peer, nil
}

// ################## Peer-interface-implement ################

func (p *TCPPeer) Addr() NetAddr {
	return NetAddr(p.nodeID)
}

func (p *TCPPeer) Info() *PeerInfo {
	return &PeerInfo{
		Addr:   p.Addr(),
		NodeID: p.nodeID,
	}
}

// Accept send plain payload to conn, rpc shold be encoded to bytes before send
func (p *TCPPeer) Accept(payload []byte) error {
	return Send(p.conn, payload)
}

// ################## End-Peer-interface-implement ################

func (p TCPPeer) String() string {
	inbound := "outbound"
	if p.inbound {
		inbound = "inbound"
	}
	return fmt.Sprintf("[NODE] addr: %s, direction: %s", p.Addr(), inbound)
}

// start  should be call after peer processed and we want to work with peer
func (p *TCPPeer) Start() {
	go p.loop()
}

func (p *TCPPeer) Stop() {
	p.conn.Close()
}

func (p *TCPPeer) loop() {
	defer p.conn.Close()
loop:
	for {
		payload, err := Receive(p.conn)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNABORTED) {
				break loop
			}
			p.infoCh <- fmt.Sprintf("cannot receive bytes from conn, err: %v", err.Error())
			break loop
		}
		// buf[:n] is encoded rpc
		rpc := new(RPC)
		if err := rpc.FromBytes(payload); err != nil {
			continue
		}
		p.rpcCh <- rpc
	}
}

//################## END-PEER #########################
