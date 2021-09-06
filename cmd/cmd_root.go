package cmd

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/larsks/blogtool/git"
)

var verbosity int

func readConfigFile() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.SetConfigName(".blogtool")

	if repodir, err := git.GetTopdir(); err == nil {
		log.Debug().Str("repodir", repodir).Msgf("looking for config in repodir")
		viper.AddConfigPath(repodir)
	}
	viper.AddConfigPath(filepath.Join(homedir, ".config"))

	//nolint:errcheck
	viper.ReadInConfig()

	return nil
}

func setLogLevel() {
	var logLevel zerolog.Level
	switch verbosity {
	case 0:
		logLevel = zerolog.WarnLevel
	case 1:
		logLevel = zerolog.InfoLevel
	default:
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
}

func NewCmdRoot() *cobra.Command {
	cmd := cobra.Command{
		Use: "blogtool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			setLogLevel()

			if err := readConfigFile(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "Increase verbosity")

	cmd.AddCommand(NewCmdNew())
	cmd.AddCommand(NewCmdUpdate())
	cmd.AddCommand(NewCmdVersion())
	return &cmd
}
