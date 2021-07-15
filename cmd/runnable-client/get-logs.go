package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdGetLogs = cobra.Command{
		Use:   "logs [jobId]",
		Short: "Gets the job logs for the given ID.",
		RunE:  getJobLogs,
		Args:  cobra.ExactArgs(1),
	}
)

func getJobLogs(cobraCmd *cobra.Command, args []string) error {
	jobID := args[0]

	logs, err := makeClient().GetLogs(jobID)

	if err != nil {
		return err
	}

	fmt.Printf("logs : %v\n", *logs)
	return nil
}
