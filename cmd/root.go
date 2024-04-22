package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "bit-connect",
		Short: "bitcoin connection handshake",
	}

	rootCmd.AddCommand(
		NewConnectCommand(),
	)

	return rootCmd
}
