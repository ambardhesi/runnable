package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdGetLogs = cobra.Command{
		Use:   "logs [jobId]",
		Short: "Gets the job logs for the given ID.",
		Run:   getJobLogs,
		Args:  cobra.ExactArgs(1),
	}
)

func getJobLogs(cobraCmd *cobra.Command, args []string) {
	jobID := args[0]

	resp, err := client().R().
		Get(serverAddress + "/job/" + jobID + "/logs")

	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Response Status Code: %v\n", resp.StatusCode())
	fmt.Printf("Response Status: %v\n", resp.Status())
	fmt.Printf("Response body: %v\n", resp)
}
