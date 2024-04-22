package cmd

import (
	"fmt"
	"net"

	"github.com/akaladarshi/bit-connect/configs"
	"github.com/akaladarshi/bit-connect/handlers"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewConnectCommand() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "connect",
		Short: "initiates connection to a bitcoin node",
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeAddress := cmd.Flag(nodeAddressFlag).Value.String()
			if nodeAddress == "" {
				return fmt.Errorf("node address is required")
			}

			return runP2PHandshakeClient(nodeAddress, args)
		},
	}

	runCmd.Flags().String(nodeAddressFlag, "0.0.0.0:18444", "Bitcoin node address")
	return runCmd
}

func runP2PHandshakeClient(nodeAddress string, _ []string) error {
	log.Info().Msgf("🔗 connecting to remote peer at %s", nodeAddress)

	// connect to node address
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to node address: %w", err)
	}

	log.Info().Msg("🌐 connection established")

	// close connection when done
	defer closeConnection(conn)

	cfg, err := configs.NewHandshakeConfig(conn.LocalAddr().String(), conn.RemoteAddr().String())
	if err != nil {
		return fmt.Errorf("failed to create handshake config: %w", err)
	}

	log.Info().Msg("🤝 initiating handshake")

	// handshake with node
	err = handlers.NewHandshakeHandler(cfg).Handle(conn)
	if err != nil {
		log.Error().Msg("❌ handshake failed")

		return fmt.Errorf("failed to handshake with node: %w", err)
	}

	log.Info().Msg("✅ handshake completed")
	return nil
}

func closeConnection(conn net.Conn) {
	log.Info().Msg("🔌 closing connection")
	mustNotErr(conn.Close())
}

func mustNotErr(err error) {
	if err != nil {
		panic(err)
	}
}
