package cmd

import (
	"time"

	"github.com/araddon/dateparse"
	"github.com/spf13/cobra"
)

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
