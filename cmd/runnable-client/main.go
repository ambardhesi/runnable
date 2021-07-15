package main

import (
	"fmt"
	"os"

	httpClient "github.com/ambardhesi/runnable/internal/client"
	"github.com/spf13/cobra"
)

const (
	// TODO Make this configurable (as a CLI flag)
	serverAddress = "https://localhost:8080"
)

var (
	rootCmd = cobra.Command{
		Use:   "runnable-client",
		Short: "Client for the the Runnable job manager server.",
	}
	caCertFile     string
	clientCertFile string
	clientKeyFile  string
)

func init() {
	flags := rootCmd.PersistentFlags()

	flags.StringVar(&caCertFile, "ca", "", "Ca Cert file path")
	flags.StringVar(&clientCertFile, "cert", "", "Client cert file path")
	flags.StringVar(&clientKeyFile, "key", "", "Client key file path")

	for _, arg := range []string{"ca", "cert", "key"} {
		rootCmd.MarkPersistentFlagRequired(arg)
	}

	rootCmd.AddCommand(&cmdStart, &cmdStop, &cmdGet, &cmdGetLogs)
}

func makeClient() *httpClient.Client {
	config := httpClient.Config{
		ServerAddress:  serverAddress,
		CaCertFilePath: caCertFile,
		CertFilePath:   clientCertFile,
		KeyFilePath:    clientKeyFile,
	}
	client, err := httpClient.NewClient(config)
	if err != nil {
		fmt.Printf("Error creating client : %v\n", err)
		os.Exit(1)
	}

	return client
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
