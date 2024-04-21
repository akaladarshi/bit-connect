package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net"

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

func (h *HandshakeHandler) Handle(rw net.Conn) error {
	msg, err := msgs.CreateHandshakeMsg(h.cfg)
	if err != nil {
		return fmt.Errorf("failed to create handshake message: %w", err)
	}

	err = sendMessage(rw, msg)
	if err != nil {
		return fmt.Errorf("failed to send handshake message: %w", err)
	}

	err = h.readConn(rw)
	if err != nil {
		return fmt.Errorf("failed to read handshake message: %w", err)
	}

	versionAckMsg := msgs.NewVersionAckMsg(common.Regtest)
	err = sendMessage(rw, versionAckMsg)
	if err != nil {
		return fmt.Errorf("failed to send version ack message: %w", err)
	}

	fmt.Println("Handshake successful")

	return nil
}

func (h *HandshakeHandler) readConn(r io.Reader) error {
	msg := msgs.BitCoinMsg{}
	err := msg.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode bitcoin msg: %w", err)
	}

	handshakeMsg := msgs.HandshakeMsg{}
	err = handshakeMsg.Decode(bytes.NewReader(msg.GetPayload()))
	if err != nil {
		return fmt.Errorf("failed to decode handshake msg: %w", err)
	}

	if handshakeMsg.ProtocolVersion == common.LatestProtocolVersion {
		sendV2 := msgs.SendAddrV2Msg{}
		err = sendV2.Decode(r)
		if err != nil {
			return fmt.Errorf("failed to decode send addr v2 msg: %w", err)
		}
	}

	versionAck := msgs.VersionAckMsg{}
	err = versionAck.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode version ack msg: %w", err)
	}

	return nil
}

func sendMessage(w io.Writer, msg msgs.Message) error {
	var buff bytes.Buffer
	err := msg.Encode(&buff)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	payload := buff.Bytes()

	headerMsg, err := msgs.NewHeader(common.Regtest, msg.GetCommand(), uint32(len(payload)), [4]byte(common.PayloadHash(payload)))
	if err != nil {
		return fmt.Errorf("failed to create header message: %w", err)
	}

	err = msgs.NewBitCoinMsg(headerMsg, payload).Encode(w)
	if err != nil {
		return fmt.Errorf("failed to encode bitcoin message: %w", err)
	}

	return nil
}
