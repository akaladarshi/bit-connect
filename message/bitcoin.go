package message

import (
	"bytes"
	"fmt"
	"io"

	"github.com/akaladarshi/bit-connect/common"
)

// BitcoinMsg represents a bitcoin message
type BitcoinMsg struct {
	header  *header
	payload []byte
}

// NewBitCoinMsg creates a new bitcoin message
func NewBitCoinMsg(msg Message) (Message, error) {
	var buff bytes.Buffer
	err := msg.Encode(&buff)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %w", err)
	}

	payload := buff.Bytes()

	head, err := newHeader(common.Regtest, msg.GetCommand(), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create header message: %w", err)
	}

	err = head.validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate header: %w", err)
	}

	return &BitcoinMsg{
		header:  head,
		payload: payload,
	}, nil
}

// Encode encodes the bitcoin message
func (m *BitcoinMsg) Encode(w io.Writer) error {
	var headerBuffer = bytes.NewBuffer(make([]byte, 0, MaxHeaderSize))
	err := m.header.encode(headerBuffer)
	if err != nil {
		return fmt.Errorf("failed to encode header: %w", err)
	}

	// write the header to the writer
	_, err = w.Write(headerBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// write the payload to the writer
	_, err = w.Write(m.payload)
	if err != nil {
		return fmt.Errorf("failed to write payload: %w", err)
	}

	return nil
}

// Decode decodes the bitcoin message
func (m *BitcoinMsg) Decode(r io.Reader) error {
	m.header = &header{}
	err := m.header.decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode header: %w", err)
	}

	err = m.header.validate()
	if err != nil {
		return fmt.Errorf("failed to validate header: %w", err)
	}

	payload := make([]byte, m.header.PayloadSize)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return fmt.Errorf("failed to read payload: %w", err)
	}

	checksum := common.PayloadHash(payload)
	// checksum is the first 4 bytes of the hash
	if !bytes.Equal(checksum[:], m.header.Checksum[:]) {
		return fmt.Errorf("invalid checksum: expected %x, got %x", m.header.Checksum, checksum)
	}

	m.payload = payload
	return nil
}

// GetCommand returns the command of the message in string format
func (m *BitcoinMsg) GetCommand() string {
	return string(bytes.TrimRight(m.header.Command[:], "\x00"))
}

// GetPayload returns the payload of the message
func (m *BitcoinMsg) GetPayload() []byte {
	return m.payload
}
