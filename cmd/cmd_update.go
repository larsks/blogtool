package cmd

import (
	"path/filepath"

	"github.com/larsks/blogtool/post"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

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
				post, err := post.ReadPostFromFile(postIndex)
				plog := log.With().Str("slug", arg).Logger()

				if err != nil {
					return err
				}

				if date != "" {
					plog.Info().Str("date", date).Msg("Setting date")
					post.Date = date
				}

				if len(tags) > 0 {
					plog.Info().Msg("Setting tags")
					if appendValues {
						post.Tags = append(post.Tags, tags...)
					} else {
						post.Tags = tags
					}

					post.Tags = uniqueValues(post.Tags)
				}

				if len(categories) > 0 {
					plog.Info().Msg("Setting categories")
					if appendValues {
						post.Categories = append(post.Categories, categories...)
					} else {
						post.Tags = categories
					}

					post.Categories = uniqueValues(post.Categories)
				}

				plog.Info().Msg("writing updated post")
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
