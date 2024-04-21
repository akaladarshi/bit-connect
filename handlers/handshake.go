package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/akaladarshi/bit-connect/common"
	"github.com/akaladarshi/bit-connect/configs"
	"github.com/akaladarshi/bit-connect/message"
)

// HandshakeHandler represents a handshake handler
type HandshakeHandler struct {
	cfg *configs.HandshakeConfig
}

// NewHandshakeHandler creates a new handshake handler
func NewHandshakeHandler(cfg *configs.HandshakeConfig) Handler {
	return &HandshakeHandler{
		cfg: cfg,
	}
}

/*
	Handshake protocol:
	local peer (L) -> remote peer (R)
	L -> R: Send version message with the local peer's version
	R -> L: Send version message back
	R -> L: Send sendaddrv2 message (only if the protocol version is 70016)
	R -> L: Send verack message
	L -> R: Send verack message after receiving version message from R
*/
// Handle handles the handshake
func (h *HandshakeHandler) Handle(rw net.Conn) error {
	versionMsg, err := message.NewVersionMsg(h.cfg)
	if err != nil {
		return fmt.Errorf("failed to create handshake message: %w", err)
	}

	// initiate the handshake by sending the version message to the peer
	err = sendMessage(rw, versionMsg)
	if err != nil {
		return fmt.Errorf("failed to send handshake message: %w", err)
	}

	// read the message sent by the peer in response to the version message
	err = h.readMessages(rw)
	if err != nil {
		return fmt.Errorf("failed to read handshake message: %w", err)
	}

	// send the version ack message to the peer
	versionAckMsg := message.NewVersionAckMsg(common.Regtest)
	err = sendMessage(rw, versionAckMsg)
	if err != nil {
		return fmt.Errorf("failed to send version ack message: %w", err)
	}

	// TODO: add logging
	fmt.Println("Handshake successful")

	return nil
}

func (h *HandshakeHandler) readMessages(r io.Reader) error {
	msg := message.BitcoinMsg{}
	err := msg.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode bitcoin msg: %w", err)
	}

	// decode the message payload
	handshakeMsg := message.Version{}
	err = handshakeMsg.Decode(bytes.NewReader(msg.GetPayload()))
	if err != nil {
		return fmt.Errorf("failed to decode handshake msg: %w", err)
	}

	// only decode the addr v2 message if the protocol version is the latest (70016)
	// node before 70016 does not send addr v2 message
	if handshakeMsg.ProtocolVersion == common.LatestProtocolVersion {
		sendV2 := message.SendAddrV2Msg{}
		err = sendV2.Decode(r)
		if err != nil {
			return fmt.Errorf("failed to decode send addr v2 msg: %w", err)
		}
	}

	// decode the version ack message
	versionAck := message.VersionAckMsg{}
	err = versionAck.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode version ack msg: %w", err)
	}

	return nil
}

func sendMessage(w io.Writer, msg message.Message) error {
	bitcoinMsg, err := message.NewBitCoinMsg(msg)
	if err != nil {
		return fmt.Errorf("failed to encode bitcoin message: %w", err)
	}

	// encode and write the message to the writer
	err = bitcoinMsg.Encode(w)
	if err != nil {
		return fmt.Errorf("failed to encode bitcoin message: %w", err)
	}

	return nil
}
