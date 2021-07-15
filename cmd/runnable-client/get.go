package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdGet = cobra.Command{
		Use:   "get [jobId]",
		Short: "Gets a job for the given ID.",
		RunE:  getJob,
		Args:  cobra.ExactArgs(1),
	}
)

func getJob(cobraCmd *cobra.Command, args []string) error {
	jobID := args[0]

	job, err := makeClient().Get(jobID)

	if err != nil {
		return err
	}

	fmt.Printf("Job : %v\n", *job)
	return nil
}
