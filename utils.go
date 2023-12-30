package p2p

import (
	"fmt"
	"net"
)

func Read(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func Write(conn net.Conn, payload []byte) error {
	n, err := conn.Write(payload)
	if err != nil {
		return err
	}

	if n != len(payload) {
		return fmt.Errorf("given payload with length %d, written %d", len(payload), n)
	}
	return nil
}
