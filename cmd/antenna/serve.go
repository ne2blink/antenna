package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the telegram bot",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("helle world")
		},
	}
}
