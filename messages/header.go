package messages

import (
	"fmt"
	"io"
)

const CommandSize = 12

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

func (h *Header) Encode(w io.Writer) error {
	err := EncodeData(w, h.Magic, h.Command, h.PayloadSize, h.Checksum)
	if err != nil {
		return fmt.Errorf("failed to encode header: %w", err)
	}

	return nil
}
