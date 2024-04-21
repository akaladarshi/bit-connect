package messages

import (
	"bytes"
	"fmt"
	"io"

	"github.com/akaladarshi/bit-connect/common"
)

// Message represents a message used in the p2p protocol
type Message interface {
	Encode() ([]byte, error)
	Decode([]byte) error
	Command() string
}

type BitCoinMsg struct {
	header  *Header
	payload []byte
}

func NewBitCoinMsg(header *Header, payload []byte) *BitCoinMsg {
	return &BitCoinMsg{
		header:  header,
		payload: payload,
	}
}

func (m *BitCoinMsg) Encode(w io.Writer) error {
	var headerBuffer = bytes.NewBuffer(make([]byte, 0, common.MaxHeaderSize))

	err := m.header.Encode(headerBuffer)
	if err != nil {
		return fmt.Errorf("failed to encode header: %w", err)
	}

	_, err = w.Write(headerBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	_, err = w.Write(m.payload)
	if err != nil {
		return fmt.Errorf("failed to write payload: %w", err)
	}

	return nil
}
