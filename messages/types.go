package messages

import (
	"encoding/binary"
	"io"
	"net"
)

// NetAddr represents a network address
// read more: https://en.bitcoin.it/wiki/Protocol_documentation#Network_address
type NetAddr struct {
	Services uint64   // same as services in the version message
	IP       [16]byte // only supports IPV4 in IPv6 format
	Port     uint16   // port number
}

// NewNetAddr creates a new netAddr
func NewNetAddr(services uint64, addr string) (NetAddr, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return NetAddr{}, err
	}

	return NetAddr{
		Services: services,
		IP:       [16]byte(tcpAddr.IP.To16()),
		Port:     uint16(tcpAddr.Port),
	}, nil
}

func (n *NetAddr) Encode(w io.Writer) error {
	err := EncodeData(w, n.Services, n.IP)
	if err != nil {
		return err
	}

	var portBuf = make([]byte, 2)
	// big endian is used for port
	binary.BigEndian.PutUint16(portBuf, n.Port)
	_, err = w.Write(portBuf)
	if err != nil {
		return err
	}

	return nil
}

func (n *NetAddr) Decode(r io.Reader) error {
	err := DecodeData(r, &n.Services, &n.IP)
	if err != nil {
		return err
	}

	var portBuf = make([]byte, 2)
	_, err = io.ReadFull(r, portBuf[:])
	if err != nil {
		return err
	}

	// big endian is used for port
	n.Port = binary.BigEndian.Uint16(portBuf)

	return nil
}
