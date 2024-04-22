package message

import (
	"fmt"
	"io"
	"time"

	"github.com/akaladarshi/bit-connect/common"
)

const (
	version                 = "version"
	fullNodeServices uint64 = 1

	defaultLastHeight = 0
	defaultAgent      = ""
)

// Version represents a version message
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
func NewVersionMsg(senderAddr string, receiverAddr string) (Message, error) {
	addrRecv, err := NewNetAddr(fullNodeServices, receiverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create receiver net address: %w", err)
	}

	addrSender, err := NewNetAddr(fullNodeServices, senderAddr)
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

// Encode encodes the version message
func (v *Version) Encode(w io.Writer) error {
	err := EncodeData(w, v.ProtocolVersion, v.Services, v.Timestamp)
	if err != nil {
		return err
	}

	err = v.AddrRecv.Encode(w)
	if err != nil {
		return fmt.Errorf("failed to encode receiver address: %w", err)
	}

	err = v.AddrFrom.Encode(w)
	if err != nil {
		return fmt.Errorf("failed to encode sender address: %w", err)
	}

	err = EncodeData(w, v.Nonce, v.UserAgent, v.LastHeight, v.Relay)
	if err != nil {
		return err
	}

	return nil
}

// Decode decodes the version message
func (v *Version) Decode(r io.Reader) error {
	err := DecodeData(r, &v.ProtocolVersion, &v.Services, &v.Timestamp)
	if err != nil {
		return err
	}

	v.AddrFrom = NetAddr{}
	err = v.AddrFrom.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode sender address: %w", err)
	}

	v.AddrRecv = NetAddr{}
	err = v.AddrRecv.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode receiver address: %w", err)
	}

	return DecodeData(r, &v.Nonce, &v.UserAgent, &v.LastHeight, &v.Relay)
}

// GetCommand returns the version message command
func (v *Version) GetCommand() string {
	return version
}
