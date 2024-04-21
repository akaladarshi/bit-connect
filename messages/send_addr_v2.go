package messages

import "io"

const (
	sendAddrV2 = "sendAddrV2"
)

type SendAddrV2Msg struct {
	header *Header
}

// NewSendAddrV2Msg creates a new send addr v2 message
func NewSendAddrV2Msg(network uint32) *SendAddrV2Msg {
	var command [CommandSize]byte
	copy(command[:], sendAddrV2)

	return &SendAddrV2Msg{
		header: &Header{
			Magic:   network,
			Command: command,
		},
	}
}

func (s *SendAddrV2Msg) Encode(w io.Writer) error {
	return nil
}

func (s *SendAddrV2Msg) Decode(r io.Reader) error {
	header := &Header{}
	err := DecodeData(r, &header.Magic, &header.Command)
	if err != nil {
		return err
	}

	s.header = header

	return nil
}

func (s *SendAddrV2Msg) GetCommand() string {
	return sendAddrV2
}
