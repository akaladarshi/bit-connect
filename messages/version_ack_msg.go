package messages

import (
	"io"
)

const versionACK = "verack"

type VersionAckMsg struct {
	header *Header
}

// NewVersionAckMsg creates a new version ack message
func NewVersionAckMsg(network uint32) *VersionAckMsg {
	var command [CommandSize]byte
	copy(command[:], versionACK)

	return &VersionAckMsg{
		header: &Header{
			Magic:   network,
			Command: command,
		},
	}
}

func (v *VersionAckMsg) Encode(w io.Writer) error {
	return nil
}

func (v *VersionAckMsg) Decode(r io.Reader) error {
	header := &Header{}
	err := DecodeData(r, &header.Magic, &header.Command)
	if err != nil {
		return err
	}

	v.header = header
	return nil
}

func (v *VersionAckMsg) GetCommand() string {
	return versionACK
}
