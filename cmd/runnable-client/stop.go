package main

import (
	"github.com/spf13/cobra"
)

var (
	cmdStop = cobra.Command{
		Use:   "stop [jobID]",
		Short: "Stops the job",
		RunE:  stopJob,
		Args:  cobra.ExactArgs(1),
	}
)

func stopJob(cobraCmd *cobra.Command, args []string) error {
	jobID := args[0]

	err := makeClient().StopJob(jobID)

	if err != nil {
		return err
	}

	return nil
}
