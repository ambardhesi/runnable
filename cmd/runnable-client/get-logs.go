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

	logs, err := makeClient().GetLogs(jobID)

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	fmt.Printf("logs : %v\n", *logs)
}
