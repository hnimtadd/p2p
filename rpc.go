package p2p

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/hnimtadd/p2p/encoding"
)

type RPC struct {
	From    NetAddr
	Payload []byte
}

func (rpc *RPC) Bytes() []byte {
	buf := new(bytes.Buffer)
	if err := NewRPCGobEncoder(buf).Encode(rpc); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (rpc *RPC) FromBytes(buf []byte) error {
	r := bytes.NewReader(buf)
	if err := NewRPCGobDecoder(r).Decode(rpc); err != nil {
		return err
	}
	return nil
}

type RPCGobEncoder struct {
	w io.Writer
}

func NewRPCGobEncoder(w io.Writer) encoding.Encoder[*RPC] {
	return &RPCGobEncoder{
		w: w,
	}
}

func (enc *RPCGobEncoder) Encode(rpc *RPC) error {
	if err := gob.NewEncoder(enc.w).Encode(rpc); err != nil {
		return err
	}
	return nil
}

type RPCGobDecoder struct {
	r io.Reader
}

func NewRPCGobDecoder(r io.Reader) encoding.Decoder[*RPC] {
	return &RPCGobDecoder{
		r: r,
	}
}

func (dec *RPCGobDecoder) Decode(rpc *RPC) error {
	if err := gob.NewDecoder(dec.r).Decode(rpc); err != nil {
		return err
	}
	return nil
}
