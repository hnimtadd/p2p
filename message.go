package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/hnimtadd/p2p/encoding"
)

type (
	MessageType byte
	Message     struct {
		Payload []byte
		Header  MessageType
	}
	DecodedMessage struct {
		From NetAddr
		Data any
	}
	HandleFunc func(*RPC) (*DecodedMessage, error)
)

func (m *Message) Bytes() []byte {
	buf := new(bytes.Buffer)
	if err := NewMessageGobEncoder(buf).Encode(m); err != nil {
		panic(fmt.Sprintf("cannot encode message to bytes, err: %v", err.Error()))
	}
	return buf.Bytes()
}

func (m *Message) FromBytes(b []byte) error {
	buf := bytes.NewReader(b)
	if err := NewMessageGobDecoder(buf).Decode(m); err != nil {
		return err
	}
	return nil
}

type MessageGobEncoder struct {
	w io.Writer
}

func NewMessageGobEncoder(w io.Writer) encoding.Encoder[*Message] {
	return &MessageGobEncoder{
		w: w,
	}
}

func (e *MessageGobEncoder) Encode(msg *Message) error {
	if err := gob.NewEncoder(e.w).Encode(msg); err != nil {
		return err
	}
	return nil
}

type MessageGobDecoder struct {
	r io.Reader
}

func NewMessageGobDecoder(r io.Reader) encoding.Decoder[*Message] {
	return &MessageGobDecoder{
		r: r,
	}
}

func (e *MessageGobDecoder) Decode(msg *Message) error {
	if err := gob.NewDecoder(e.r).Decode(msg); err != nil {
		return err
	}
	return nil
}
