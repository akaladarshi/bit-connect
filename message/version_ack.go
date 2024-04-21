package message

import (
	"io"
)

const versionACK = "verack"

type VersionAckMsg struct {
	header *header
}

// NewVersionAckMsg creates a new version ack message
func NewVersionAckMsg(network uint32) *VersionAckMsg {
	var command [CommandSize]byte
	copy(command[:], versionACK)

	return &VersionAckMsg{
		header: &header{
			Magic:   network,
			Command: command,
		},
	}
}

// Encode encodes the version ack message
func (v *VersionAckMsg) Encode(w io.Writer) error {
	return nil
}

// Decode decodes the version ack message
func (v *VersionAckMsg) Decode(r io.Reader) error {
	v.header = &header{}
	return DecodeData(r, &v.header.Magic, &v.header.Command)
}

func (v *VersionAckMsg) GetCommand() string {
	return versionACK
}
