package configs

import (
	"errors"
)

var (
	ErrInvalidHandshakeConfig = errors.New("invalid handshake config")
)

type HandshakeConfig struct {
	SenderAddress   string
	ReceiverAddress string
}

func NewHandshakeConfig(localAddress, remotePeerAddress string) (*HandshakeConfig, error) {
	if localAddress == "" || remotePeerAddress == "" {
		return nil, ErrInvalidHandshakeConfig
	}

	return &HandshakeConfig{
		SenderAddress:   localAddress,
		ReceiverAddress: remotePeerAddress,
	}, nil
}
