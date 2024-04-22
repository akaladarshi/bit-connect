package main

import (
	"context"
	"fmt"
	"os"

	"github.com/akaladarshi/bit-connect/cmd"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	rootCmd := cmd.NewRootCmd()
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
