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
		Run:   startJob,
		Args:  cobra.MinimumNArgs(1),
	}
)

func startJob(cobraCmd *cobra.Command, args []string) {
	cmd := strings.Join(args, " ")

	jobID, err := makeClient().StartJob(cmd)

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	fmt.Printf("Started job : %v\n", *jobID)
}
