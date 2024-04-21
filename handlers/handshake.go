package handlers

import (
	"fmt"
	"io"

	"github.com/akaladarshi/bit-connect/common"
	"github.com/akaladarshi/bit-connect/configs"
	msgs "github.com/akaladarshi/bit-connect/messages"
)

type HandshakeHandler struct {
	cfg *configs.HandshakeConfig
}

func NewHandshakeHandler(cfg *configs.HandshakeConfig) Handler {
	return &HandshakeHandler{
		cfg: cfg,
	}
}

func (h *HandshakeHandler) Handle(rw io.ReadWriteCloser) error {
	msg, err := msgs.CreateHandshakeMsg(h.cfg)
	if err != nil {
		return fmt.Errorf("failed to create handshake message: %w", err)
	}

	err = sendMessage(rw, msg)
	if err != nil {
		return fmt.Errorf("failed to send handshake message: %w", err)
	}

	_, err = h.readConn(rw)
	if err != nil {
		return fmt.Errorf("failed to read handshake message: %w", err)
	}

	return nil
}

func (h *HandshakeHandler) readConn(r io.Reader) ([]byte, error) {
	var buf [common.MaxHeaderSize]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return nil, fmt.Errorf("failed to read from connection: %w", err)
	}

	fmt.Println(buf)
	// net.TCPAddr{}
	// hr := bytes.NewReader(buf[:])
	return nil, nil
}

func sendMessage(w io.Writer, msg msgs.Message) error {
	payload, err := msg.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	headerMsg, err := msgs.NewHeader(common.Regtest, msg.Command(), uint32(len(payload)), [4]byte(common.PayloadHash(payload)))
	if err != nil {
		return fmt.Errorf("failed to create header message: %w", err)
	}

	err = msgs.NewBitCoinMsg(headerMsg, payload).Encode(w)
	if err != nil {
		return fmt.Errorf("failed to encode bitcoin message: %w", err)
	}

	return nil
}
