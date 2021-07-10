package main

import (
	"fmt"
	"strings"

	"github.com/ambardhesi/runnable/internal/server"
	"github.com/spf13/cobra"
)

var (
	cmdStart = cobra.Command{
		Use:   "start command [args...]",
		Short: "Starts a job on the Runnable server",
		Run:   startJob,
		Args:  cobra.MinimumNArgs(1),
	}
)

func startJob(cobraCmd *cobra.Command, args []string) {
	cmd := strings.Join(args, " ")
	request := server.StartJobRequest{
		Command: cmd,
	}

	resp, err := client().R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(serverAddress + "/job")

	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Response Status Code: %v\n", resp.StatusCode())
	fmt.Printf("Response Status: %v\n", resp.Status())
	fmt.Printf("Response body: %v\n", resp)
}
