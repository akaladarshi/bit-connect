package messages

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

// Header represents the header of a message
type Header struct {
	Magic       uint32            // network name (4 bytes)
	Command     [CommandSize]byte // message name (12 bytes)
	PayloadSize uint32            // size of payload (4 bytes)
	Checksum    [4]byte           // checksum of payload (4 bytes)
}

// NewHeader creates a new header message
func NewHeader(network uint32, cmd string, payloadSize uint32, checksum [4]byte) (*Header, error) {
	if len(cmd) > CommandSize {
		return nil, fmt.Errorf("command size exceeds limit: %d", len(cmd))
	}

	var command [CommandSize]byte
	copy(command[:], cmd)

	return &Header{
		Magic:       network,
		Command:     command,
		PayloadSize: payloadSize,
		Checksum:    checksum,
	}, nil
}

// Encode encodes the header message
func (h *Header) Encode(w io.Writer) error {
	return EncodeData(w, h.Magic, h.Command, h.PayloadSize, h.Checksum)
}

// Decode decodes the header message
func (h *Header) Decode(r io.Reader) error {
	return DecodeData(r, &h.Magic, &h.Command, &h.PayloadSize, &h.Checksum)
}

// Validate validates the header message
func (h *Header) Validate() error {
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
