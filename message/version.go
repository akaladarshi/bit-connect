package message

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

type Version struct {
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

// NewVersionMsg creates a new version message
func NewVersionMsg(cfg *configs.HandshakeConfig) (Message, error) {
	addrRecv, err := NewNetAddr(fullNodeServices, cfg.ReceiverAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create receiver net address: %w", err)
	}

	addrSender, err := NewNetAddr(fullNodeServices, cfg.SenderAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create sender net address: %w", err)
	}

	return &Version{
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

func (h *Version) Encode(w io.Writer) error {
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

func (h *Version) Decode(r io.Reader) error {
	err := DecodeData(r, &h.ProtocolVersion, &h.Services, &h.Timestamp)
	if err != nil {
		return err
	}

	h.AddrFrom = NetAddr{}
	err = h.AddrFrom.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode sender address: %w", err)
	}

	h.AddrRecv = NetAddr{}
	err = h.AddrRecv.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode receiver address: %w", err)
	}

	return DecodeData(r, &h.Nonce, &h.UserAgent, &h.LastHeight, &h.Relay)
}

func (h *Version) GetCommand() string {
	return version
}
