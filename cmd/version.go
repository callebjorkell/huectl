package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of huectl",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version:    %s\n", version)
			fmt.Printf("Go Version: %s\n", runtime.Version())
			fmt.Printf("Built on:   %s\n", date)
			fmt.Printf("Commit:     %s\n", commit)
		},
	}
}