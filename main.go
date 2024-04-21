package main

import (
	"context"
	"fmt"
	"os"

	"github.com/akaladarshi/bit-connect/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
