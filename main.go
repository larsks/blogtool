package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/araddon/dateparse"
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

func dateFromFlags(cmd *cobra.Command, defaultIsToday bool) (string, error) {
	dateIn, err := cmd.Flags().GetString("date")
	if err != nil {
		return "", err
	}

	if dateIn == "" {
		if !defaultIsToday {
			return "", nil
		}

		dateIn = "today"
	}

	var ts time.Time

	if dateIn == "today" {
		ts = time.Now()
	} else if dateIn != "" {
		ts, err = dateparse.ParseStrict(dateIn)
		if err != nil {
			return "", err
		}
	}

	return ts.Format("2006-01-02"), nil
}

func NewCmdNew() *cobra.Command {
	cmd := cobra.Command{
		Use:  "new <title>",
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

			date, err := dateFromFlags(cmd, true)
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

			slug, err := cmd.Flags().GetString("slug")
			if err != nil {
				return err
			}
			if slug == "" {
				slug = post.Slug(maxlen)
			}
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

func uniqueValues(values []string) []string {
	seen := make(map[string]bool)
	var res []string

	for _, value := range values {
		if _, found := seen[value]; !found {
			seen[value] = true
			res = append(res, value)
		}
	}

	return res
}

func NewCmdUpdate() *cobra.Command {
	cmd := cobra.Command{
		Use: "update <slug>",
		RunE: func(cmd *cobra.Command, args []string) error {
			date, err := dateFromFlags(cmd, false)
			if err != nil {
				return err
			}

			appendValues, err := cmd.Flags().GetBool("append")
			if err != nil {
				return err
			}

			tags, err := cmd.Flags().GetStringSlice("tag")
			if err != nil {
				return err
			}

			categories, err := cmd.Flags().GetStringSlice("category")
			if err != nil {
				return err
			}

			for _, arg := range args {
				postIndex := filepath.Join(arg, "index.md")
				post, err := ReadPostFromFile(postIndex)

				if err != nil {
					return err
				}

				if date != "" {
					log.Info().Str("date", date).Msgf("Setting date")
					post.Date = date
				}

				if len(tags) > 0 {
					log.Info().Msgf("Setting tags")
					if appendValues {
						post.Tags = append(post.Tags, tags...)
					} else {
						post.Tags = tags
					}

					post.Tags = uniqueValues(post.Tags)
				}

				if len(categories) > 0 {
					log.Info().Msgf("Setting categories")
					if appendValues {
						post.Categories = append(post.Categories, categories...)
					} else {
						post.Tags = categories
					}

					post.Categories = uniqueValues(post.Categories)
				}

				if err := post.WriteToFile(postIndex); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolP("append", "a", false, "Append tags/categories")
	cmd.Flags().StringSliceP("tag", "t", nil, "Specify tags for post")
	cmd.Flags().StringSliceP("category", "c", nil, "Specify category for post")
	cmd.Flags().StringP("date", "d", "", "Specify post date")

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
	cmd.AddCommand(NewCmdUpdate())
	return &cmd
}

func main() {
	cobra.CheckErr(NewCmdRoot().Execute())
}
