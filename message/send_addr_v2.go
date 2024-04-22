package message

import (
	"io"
)

const (
	sendAddrV2 = "sendaddrv2"
)

// SendAddrV2Msg represents a send addr v2 message
// only sent by nodes which support protocol version 70016 or higher
type SendAddrV2Msg struct{}

// Encode encodes the send addr v2 message
func (s *SendAddrV2Msg) Encode(w io.Writer) error {
	return nil
}

// Decode decodes the send addr v2 message
func (s *SendAddrV2Msg) Decode(r io.Reader) error {
	return nil
}

// GetCommand returns the send addr v2 message command
func (s *SendAddrV2Msg) GetCommand() string {
	return sendAddrV2
}
