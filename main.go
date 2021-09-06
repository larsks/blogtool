package main

import (
	"github.com/larsks/blogtool/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.NewCmdRoot().Execute())
}
