package cmd

import (
	"fmt"
	"net"

	"github.com/akaladarshi/bit-connect/configs"
	"github.com/akaladarshi/bit-connect/handlers"
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "run p2p handlers client",
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
	// TODO: parse address to check format
	// connect to node address
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to node address: %w", err)
	}

	// close connection when done
	defer closeConnection(conn)

	cfg, err := configs.NewHandshakeConfig(conn.LocalAddr().String(), conn.RemoteAddr().String())
	if err != nil {
		return fmt.Errorf("failed to create handshake config: %w", err)
	}

	// handshake with node
	err = handlers.NewHandshakeHandler(cfg).Handle(conn)
	if err != nil {
		return fmt.Errorf("failed to handshake with node: %w", err)
	}

	return nil
}

func closeConnection(conn net.Conn) {
	mustNotErr(conn.Close())
}

func mustNotErr(err error) {
	if err != nil {
		panic(err)
	}
}
