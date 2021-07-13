package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cmdStart = cobra.Command{
		Use:   "start command [args...]",
		Short: "Starts a job on the Runnable server",
		RunE:  startJob,
		Args:  cobra.MinimumNArgs(1),
	}
)

func startJob(cobraCmd *cobra.Command, args []string) error {
	cmd := strings.Join(args, " ")

	jobID, err := makeClient().StartJob(cmd)

	if err != nil {
		return nil
	}

	fmt.Printf("Started job : %v\n", *jobID)
	return nil
}
