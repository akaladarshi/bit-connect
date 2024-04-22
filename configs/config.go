package configs

import (
	"errors"
)

var (
	ErrInvalidHandshakeConfig = errors.New("invalid handshake config")
)

// HandshakeConfig represents a handshake configuration
type HandshakeConfig struct {
	SenderAddress   string
	ReceiverAddress string
}

// NewHandshakeConfig creates a new handshake configuration
func NewHandshakeConfig(localAddress, remotePeerAddress string) (*HandshakeConfig, error) {
	if localAddress == "" || remotePeerAddress == "" {
		return nil, ErrInvalidHandshakeConfig
	}

	return &HandshakeConfig{
		SenderAddress:   localAddress,
		ReceiverAddress: remotePeerAddress,
	}, nil
}
