package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "bit-p2p",
		Short: "bitcoin p2p handlers client",
	}

	rootCmd.AddCommand(
		NewRunCmd(),
	)

	return rootCmd
}
