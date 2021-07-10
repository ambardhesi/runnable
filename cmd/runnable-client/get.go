package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdGet = cobra.Command{
		Use:   "get [jobId]",
		Short: "Gets a job for the given ID.",
		Run:   getJob,
		Args:  cobra.ExactArgs(1),
	}
)

func getJob(cobraCmd *cobra.Command, args []string) {
	jobID := args[0]

	resp, err := client().R().
		Get(serverAddress + "/job/" + jobID)

	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Response Status Code: %v\n", resp.StatusCode())
	fmt.Printf("Response Status: %v\n", resp.Status())
	fmt.Printf("Response body: %v\n", resp)
}
