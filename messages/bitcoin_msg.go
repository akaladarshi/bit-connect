package messages

import (
	"bytes"
	"fmt"
	"io"

	"github.com/akaladarshi/bit-connect/common"
)

// BitCoinMsg represents a bitcoin message
// TODO: replace payload with Message interface
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
	var headerBuffer = bytes.NewBuffer(make([]byte, 0, MaxHeaderSize))

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

func (m *BitCoinMsg) Decode(r io.Reader) error {
	header := &Header{}
	err := header.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode header: %w", err)
	}

	err = header.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate header: %w", err)
	}

	m.header = header

	payload := make([]byte, m.header.PayloadSize)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return fmt.Errorf("failed to read payload: %w", err)
	}

	checksum := common.PayloadHash(payload)
	// checksum is the first 4 bytes of the hash
	if !bytes.Equal(checksum[:][0:4], m.header.Checksum[:]) {
		return fmt.Errorf("invalid checksum: expected %x, got %x", m.header.Checksum, checksum)
	}

	m.payload = payload
	return nil
}

// GetCommand returns the command of the message in string format
func (m *BitCoinMsg) GetCommand() string {
	return string(bytes.TrimRight(m.header.Command[:], "\x00"))
}

// GetPayload returns the payload of the message
func (m *BitCoinMsg) GetPayload() []byte {
	return m.payload
}
