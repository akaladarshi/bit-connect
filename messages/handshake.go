package messages

import (
	"fmt"
	"io"
	"time"

	"github.com/akaladarshi/bit-connect/common"
	"github.com/akaladarshi/bit-connect/configs"
)

const (
	version                 = "version"
	fullNodeServices uint64 = 1

	defaultLastHeight = 0
	defaultAgent      = ""
)

type HandshakeMsg struct {
	ProtocolVersion int32
	Services        uint64
	Timestamp       int64
	AddrRecv        NetAddr
	AddrFrom        NetAddr
	Nonce           uint64
	UserAgent       string
	LastHeight      int32
	Relay           bool // we can remove this field as it is not used.
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
		ProtocolVersion: common.LatestProtocolVersion,
		Services:        fullNodeServices,
		Timestamp:       time.Unix(time.Now().Unix(), 0).Unix(),
		AddrRecv:        addrRecv,
		AddrFrom:        addrSender,
		Nonce:           generateNonce(),
		UserAgent:       defaultAgent,
		LastHeight:      defaultLastHeight,
		Relay:           false,
	}, nil
}

func (h *HandshakeMsg) Encode(w io.Writer) error {
	err := EncodeData(w, h.ProtocolVersion, h.Services, h.Timestamp)
	if err != nil {
		return err
	}

	err = h.AddrRecv.Encode(w)
	if err != nil {
		return fmt.Errorf("failed to encode receiver address: %w", err)
	}

	err = h.AddrFrom.Encode(w)
	if err != nil {
		return fmt.Errorf("failed to encode sender address: %w", err)
	}

	err = EncodeData(w, h.Nonce, h.UserAgent, h.LastHeight, h.Relay)
	if err != nil {
		return err
	}

	return nil
}

func (h *HandshakeMsg) Decode(r io.Reader) error {
	err := DecodeData(r, &h.ProtocolVersion, &h.Services, &h.Timestamp)
	if err != nil {
		return err
	}

	senderAddr := &NetAddr{}
	err = senderAddr.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode sender address: %w", err)
	}

	h.AddrFrom = *senderAddr

	receiverAddr := &NetAddr{}
	err = receiverAddr.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode receiver address: %w", err)
	}

	h.AddrRecv = *receiverAddr

	return DecodeData(r, &h.Nonce, &h.UserAgent, &h.LastHeight, &h.Relay)
}

func (h *HandshakeMsg) GetCommand() string {
	return version
}
