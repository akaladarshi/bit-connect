package message

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/akaladarshi/bit-connect/common"
)

const (
	CommandSize   = 12
	MaxHeaderSize = 24
)

// header represents the header of a message
type header struct {
	Magic       uint32            // network name (4 bytes)
	Command     [CommandSize]byte // message name (12 bytes)
	PayloadSize uint32            // size of payload (4 bytes)
	Checksum    [4]byte           // checksum of payload (4 bytes)
}

// newHeader creates a new header message
func newHeader(network uint32, cmd string, payload []byte) (*header, error) {
	if len(cmd) > CommandSize {
		return nil, fmt.Errorf("command size exceeds limit: %d", len(cmd))
	}

	var command [CommandSize]byte
	copy(command[:], cmd)

	return &header{
		Magic:       network,
		Command:     command,
		PayloadSize: uint32(len(payload)),
		Checksum:    common.PayloadHash(payload),
	}, nil
}

// encode encodes the header message
func (h *header) encode(w io.Writer) error {
	return EncodeData(w, h.Magic, h.Command, h.PayloadSize, h.Checksum)
}

// decode decodes the header message
func (h *header) decode(r io.Reader) error {
	return DecodeData(r, &h.Magic, &h.Command, &h.PayloadSize, &h.Checksum)
}

// validate validates the header message
func (h *header) validate() error {
	// right now we only support regtest network
	if h.Magic != common.Regtest {
		return fmt.Errorf("unsupported network: %d", h.Magic)
	}

	cmd := string(bytes.TrimRight(h.Command[:], "\x00"))
	if !utf8.ValidString(cmd[:]) {
		return fmt.Errorf("invalid command: %s", cmd)
	}

	if h.PayloadSize > MaxPayloadSize {
		return fmt.Errorf("payload size exceeds limit: expected %d, got %d", MaxPayloadSize, h.PayloadSize)
	}

	return nil
}
