package cmd

import (
	"fmt"

	"github.com/larsks/blogtool/version"
	"github.com/spf13/cobra"
)

func NewCmdVersion() *cobra.Command {
	cmd := cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", version.BuildVersion)
			fmt.Printf("Build ref: %s\n", version.BuildRef)
			fmt.Printf("Build date: %s\n", version.BuildDate)
		},
	}

	return &cmd
}
