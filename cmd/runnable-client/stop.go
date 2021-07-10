package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdStop = cobra.Command{
		Use:   "stop [jobID]",
		Short: "Stops the job",
		Run:   stopJob,
		Args:  cobra.ExactArgs(1),
	}
)

func stopJob(cobraCmd *cobra.Command, args []string) {
	jobID := args[0]

	resp, err := client().R().
		Post(serverAddress + "/job/" + jobID + "/stop")

	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Response Status Code: %v\n", resp.StatusCode())
	fmt.Printf("Response Status: %v\n", resp.Status())
	fmt.Printf("Response body: %v\n", resp)
}
