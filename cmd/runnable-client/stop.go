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

	err := makeClient().StopJob(jobID)

	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
}
