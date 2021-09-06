package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/larsks/blogtool/git"
	"github.com/larsks/blogtool/post"
)

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

			post := post.Post{
				Metadata: post.Metadata{
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

			use_git, err := cmd.Flags().GetBool("git")
			if err != nil {
				return err
			}

			if err := os.MkdirAll(slug, 0o777); err != nil {
				return err
			}

			if err := post.WriteToFile(filepath.Join(slug, "index.md")); err != nil {
				return err
			}

			if use_git {
				start_branch, err := cmd.Flags().GetString("start-branch")
				if err != nil {
					return err
				}

				log.Debug().Str("slug", slug).Msgf("creating new branch")
				if err := git.CreateBranch(fmt.Sprintf("draft/%s", slug), start_branch); err != nil {
					return err
				}

				log.Debug().Str("slug", slug).Msgf("adding post")
				if err := git.AddFiles(filepath.Join(slug, "index.md")); err != nil {
					return err
				}

				log.Debug().Str("slug", slug).Msgf("committing changes")
				if err := git.Commit(fmt.Sprintf("Add %s", slug)); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringSliceP("tag", "t", nil, "Specify tags for post")
	cmd.Flags().StringSliceP("category", "c", nil, "Specify category for post")
	cmd.Flags().StringP("date", "d", "", "Specify post date")
	cmd.Flags().StringP("slug", "s", "", "Specify post slug")
	cmd.Flags().BoolP("git", "g", false, "Create new git branch for post")
	cmd.Flags().StringP("start-branch", "b", "master", "Name of start branch")

	cmd.Flags().Int("max-slug-len", 30, "Set maximum length of slug")
	cmd.Flags().MarkHidden("max-slug-len") //nolint

	if err := viper.BindPFlag("start-branch", cmd.Flags().Lookup("start-branch")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("max-slug-len", cmd.Flags().Lookup("max-slug-len")); err != nil {
		panic(err)
	}

	return &cmd
}
