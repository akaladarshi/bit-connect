package message

import (
	"io"
)

const versionACK = "verack"

type VersionAckMsg struct{}

// Encode encodes the version ack message
func (v *VersionAckMsg) Encode(_ io.Writer) error {
	return nil
}

// Decode decodes the version ack message
func (v *VersionAckMsg) Decode(_ io.Reader) error {
	return nil
}

// GetCommand returns the version ack message command
func (v *VersionAckMsg) GetCommand() string {
	return versionACK
}
