package p2p

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"

	"github.com/hnimtadd/dpoker/network/utils"
)

type handshake struct {
	FromID string `json:"id"`
}

func handshakeToBytes(h *handshake) []byte {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(h); err != nil {
		log.Panicf("handshake: cannot encode handshake to bytes, err: %s", err.Error())
	}
	return buf.Bytes()
}

func handshakeFromBytes(payload []byte) *handshake {
	handshake := new(handshake)
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(handshake); err != nil {
		log.Panicf("handshake: cannot decode handshake fromto bytes, err: %s", err.Error())
	}
	return handshake
}

// DefaultTCPHandshake should be call after dial to other peer from fromNode
// This should ask the other node the id of the node and then create peer with this id
func TCPHandshake(fromNode Transport, peer net.Conn) (NodeInformation, error) {
	syn := handshake{
		FromID: fromNode.Addr().(string),
	}
	synBytes := handshakeToBytes(&syn)
	err := utils.Write(peer, synBytes)
	if err != nil {
		return NodeInformation{}, err
	}

	saBytes, err := utils.Read(peer)
	if err != nil {
		return NodeInformation{}, err
	}
	saHandshake := handshakeFromBytes(saBytes)
	return NodeInformation{
		NodeID: saHandshake.FromID,
	}, nil
}

func TCPHandshakeReply(peer net.Conn, toNode Transport) (NodeInformation, error) {
	sBytes, err := utils.Read(peer)
	if err != nil {
		return NodeInformation{}, err
	}
	sHandshake := handshakeFromBytes(sBytes)

	saHandkshake := handshake{
		FromID: toNode.Addr().(string),
	}

	saBytes := handshakeToBytes(&saHandkshake)
	if err := utils.Write(peer, saBytes); err != nil {
		return NodeInformation{}, err
	}

	return NodeInformation{
		NodeID: sHandshake.FromID,
	}, nil
}
