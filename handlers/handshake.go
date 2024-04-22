package handlers

import (
	"fmt"
	"io"
	"net"

	"github.com/akaladarshi/bit-connect/common"
	"github.com/akaladarshi/bit-connect/configs"
	"github.com/akaladarshi/bit-connect/message"
	"github.com/rs/zerolog/log"
)

// HandshakeHandler represents a handshake handler
type HandshakeHandler struct {
	cfg *configs.HandshakeConfig
}

// NewHandshakeHandler creates a new handshake handler
func NewHandshakeHandler(cfg *configs.HandshakeConfig) *HandshakeHandler {
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
func (h *HandshakeHandler) Handle(conn net.Conn) error {
	log.Info().Msg("sending version message to remote peer")

	// initiate the handshake by sending the version message to the peer
	err := h.sendVersionMessage(conn)
	if err != nil {
		return fmt.Errorf("failed to send handshake message: %w", err)
	}

	// read the message sent by the peer in response to the version message
	err = h.readRemotePeerMessages(conn)
	if err != nil {
		return fmt.Errorf("failed to read handshake message: %w", err)
	}

	log.Info().Msg("sending version ack message to remote peer")

	// send the version ack message to the peer
	err = h.sendVersionAckMessage(conn)
	if err != nil {
		return fmt.Errorf("failed to send version ack message: %w", err)
	}

	return nil
}

// all the messages are in the format of bitcoin message
func (h *HandshakeHandler) readRemotePeerMessages(r io.Reader) error {
	versionMsg, err := readVersionMsg(r)
	if err != nil {
		return fmt.Errorf("failed to read version message: %w", err)
	}

	log.Info().Msg("received version message from remote peer")

	// only decode the addr v2 message if the protocol version is the latest (70016)
	// node before 70016 does not send addr v2 message
	if versionMsg.ProtocolVersion >= common.LatestProtocolVersion {
		log.Info().Msgf("remote peer protocol version is %d, reading send addr v2 message", common.LatestProtocolVersion)
		err = readSendAddrV2Msg(r)
		if err != nil {
			return fmt.Errorf("failed to read send addr v2 message: %w", err)
		}
	}

	// decode the version ack message
	err = readVersionAckMsg(r)
	if err != nil {
		return fmt.Errorf("failed to read version ack message: %w", err)

	}

	log.Info().Msg("received version ack message from remote peer")
	return nil
}

func readVersionMsg(r io.Reader) (*message.Version, error) {
	bitcoinMsg, err := readBitcoinMessage(r)
	if err != nil {
		return nil, err
	}

	versionMsg := message.Version{}
	// verify the received message
	err = bitcoinMsg.LoadAndVerifyOriginalMsg(&versionMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to load and verify received message: %w", err)
	}

	return &versionMsg, nil
}

func readSendAddrV2Msg(r io.Reader) error {
	bitcoinMsg, err := readBitcoinMessage(r)
	if err != nil {
		return err
	}

	err = bitcoinMsg.LoadAndVerifyOriginalMsg(&message.SendAddrV2Msg{})
	if err != nil {
		return fmt.Errorf("failed to load and verify received message: %w", err)
	}

	return nil
}

func readVersionAckMsg(r io.Reader) error {
	bitcoinMsg, err := readBitcoinMessage(r)
	if err != nil {
		return err
	}

	err = bitcoinMsg.LoadAndVerifyOriginalMsg(&message.VersionAckMsg{})
	if err != nil {
		return fmt.Errorf("failed to load and verify received message: %w", err)
	}

	return nil
}

func readBitcoinMessage(r io.Reader) (*message.BitcoinMsg, error) {
	bitcoinMsg := message.BitcoinMsg{}
	err := bitcoinMsg.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode bitcoin bitcoinMsg: %w", err)
	}

	return &bitcoinMsg, nil
}

func (h *HandshakeHandler) sendVersionMessage(w io.Writer) error {
	versionMsg, err := message.NewVersionMsg(h.cfg.SenderAddress, h.cfg.ReceiverAddress)
	if err != nil {
		return fmt.Errorf("failed to create new version message: %w", err)
	}

	return sendBitcoinMessage(w, versionMsg)
}

func (h *HandshakeHandler) sendVersionAckMessage(w io.Writer) error {
	return sendBitcoinMessage(w, &message.VersionAckMsg{})
}

func sendBitcoinMessage(w io.Writer, msg message.Message) error {
	bitcoinMsg, err := message.NewBitCoinMsg(msg)
	if err != nil {
		return fmt.Errorf("failed to create new bitcoin message: %w", err)
	}

	// encode and write the message to the writer
	err = bitcoinMsg.Encode(w)
	if err != nil {
		return fmt.Errorf("failed to encode bitcoin message: %w", err)
	}

	return nil
}
