package main

import (
	"os"

	"github.com/larsks/blogtool/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cobra.CheckErr(cmd.NewCmdRoot().Execute())
}
