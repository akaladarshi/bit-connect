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

func NewHandshakeConfig(senderAddress, receiverAddress string) (*HandshakeConfig, error) {
	if senderAddress == "" || receiverAddress == "" {
		return nil, ErrInvalidHandshakeConfig
	}

	return &HandshakeConfig{
		SenderAddress:   senderAddress,
		ReceiverAddress: receiverAddress,
	}, nil
}
