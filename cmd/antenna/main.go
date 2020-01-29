package main

import (
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "antenna [options] [command]",
	}
	cmd.PersistentFlags().StringP("config", "c", "", "config file")
	cmd.MarkPersistentFlagFilename("config", "yaml", "yml")
	cmd.AddCommand(
		newAppCommand(),
		newServeCommand(),
	)
	cmd.Execute()
}
