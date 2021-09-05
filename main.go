package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getGitTopdir() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(out), "\n"), err
}

func readConfigFile() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.SetConfigName(".blogtool")

	if repodir, err := getGitTopdir(); err == nil {
		log.Debug().Str("repodir", repodir).Msgf("looking for config in repodir")
		viper.AddConfigPath(repodir)
	}
	viper.AddConfigPath(filepath.Join(homedir, ".config"))

	viper.ReadInConfig()

	return nil
}

func NewCmdNew() *cobra.Command {
	cmd := cobra.Command{
		Use:  "new",
		Args: cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			maxlen := viper.GetInt("max-slug-len")

			tags, err := cmd.Flags().GetStringSlice("tag")
			if err != nil {
				return err
			}

			categories, err := cmd.Flags().GetStringSlice("category")
			if err != nil {
				return err
			}

			date, err := cmd.Flags().GetString("date")
			if err != nil {
				return err
			}

			post := Post{
				Metadata: Metadata{
					Title:      args[0],
					Tags:       tags,
					Categories: categories,
					Date:       date,
				},
			}
			slug := post.Slug(maxlen)
			if err := os.MkdirAll(slug, 0o777); err != nil {
				return err
			}

			if err := post.WriteToFile(filepath.Join(slug, "index.md")); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringSliceP("tag", "t", nil, "Specify tags for post")
	cmd.Flags().StringSliceP("category", "c", nil, "Specify category for post")
	cmd.Flags().StringP("date", "d", "", "Specify post date")
	cmd.Flags().StringP("slug", "s", "", "Specify post slug")
	cmd.Flags().Int("max-slug-len", 30, "Set maximum length of slug")

	viper.BindPFlag("max-slug-len", cmd.Flags().Lookup("max-slug-len"))

	return &cmd
}

func NewCmdParse() *cobra.Command {
	cmd := cobra.Command{
		Use: "parse",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				post, err := ReadPostFromFile(arg)
				if err != nil {
					return err
				}
				fmt.Printf("%+v\n", post)

				post.WriteToFile(fmt.Sprintf("%s.out", arg))
			}

			return nil
		},
	}

	return &cmd
}

func NewCmdRoot() *cobra.Command {
	cmd := cobra.Command{
		Use: "blogtool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := readConfigFile(); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(NewCmdNew())
	cmd.AddCommand(NewCmdParse())
	return &cmd
}

func main() {
	cobra.CheckErr(NewCmdRoot().Execute())
}
