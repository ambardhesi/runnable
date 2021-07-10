package main

import (
	"fmt"
	"os"

	rclient "github.com/ambardhesi/runnable/internal/client"
	"github.com/go-resty/resty/v2"
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

func client() *resty.Client {
	client := resty.New()
	tlsConfig, err := rclient.GetTLSConfig(clientCertFile, clientKeyFile, caCertFile)
	if err != nil {
		fmt.Printf("Failed to create http client %v\n", err)
	}
	client.SetTLSClientConfig(tlsConfig)

	return client
}

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

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
