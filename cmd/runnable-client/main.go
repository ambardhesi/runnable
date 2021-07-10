package main

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	// TODO Make this configurable (as a CLI flag)
	serverAddress = "http://localhost:8080"
)

var (
	rootCmd = cobra.Command{
		Use:   "runnable-client",
		Short: "Client for the the Runnable job manager server.",
	}
)

func init() {
	rootCmd.AddCommand(&cmdStart, &cmdStop, &cmdGet, &cmdGetLogs)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
