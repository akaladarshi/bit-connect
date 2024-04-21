package messages

import (
	"io"
)

const (
	MaxPayloadSize = 32 * 1024 * 1024 // 32MB
)

// Message represents a message used in the p2p protocol
type Message interface {
	Encode(w io.Writer) error
	Decode(r io.Reader) error
	GetCommand() string
}
