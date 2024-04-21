package messages

import (
	"bytes"
	"fmt"
	"time"

	"github.com/akaladarshi/bit-connect/configs"
)

const (
	version                 = "version"
	fullNodeServices uint64 = 1

	defaultProtocolVersion = 70016
	defaultStartHeight     = 1
	defaultAgent           = ""
)

type HandshakeMsg struct {
	ProtocolVersion int32
	Services        uint64
	Timestamp       int64
	AddrRecv        NetAddr
	AddrFrom        NetAddr
	Nonce           uint64
	UserAgent       string
	Height          int32
	Relay           bool
}

// CreateHandshakeMsg creates a new version message
func CreateHandshakeMsg(cfg *configs.HandshakeConfig) (Message, error) {
	addrRecv, err := NewNetAddr(fullNodeServices, cfg.ReceiverAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create receiver net address: %w", err)
	}

	addrSender, err := NewNetAddr(fullNodeServices, cfg.SenderAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create sender net address: %w", err)
	}

	return &HandshakeMsg{
		ProtocolVersion: defaultProtocolVersion,
		Services:        fullNodeServices,
		Timestamp:       time.Unix(time.Now().Unix(), 0).Unix(),
		AddrRecv:        addrRecv,
		AddrFrom:        addrSender,
		Nonce:           generateNonce(),
		UserAgent:       defaultAgent,
		Height:          defaultStartHeight,
		Relay:           false,
	}, nil
}

func (h *HandshakeMsg) Encode() ([]byte, error) {
	var buf bytes.Buffer

	err := EncodeData(&buf, h.ProtocolVersion, h.Services, h.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to encode handshake data: %w", err)
	}

	err = h.AddrRecv.Encode(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to encode receiver address: %w", err)
	}

	err = h.AddrFrom.Encode(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to encode sender address: %w", err)
	}

	err = EncodeData(&buf, h.Nonce, h.UserAgent, h.Height)
	return buf.Bytes(), nil
}

func (h *HandshakeMsg) Decode(data []byte) error {
	return nil
}

func (h *HandshakeMsg) Command() string {
	return version
}
