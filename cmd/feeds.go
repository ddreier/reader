package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"reader/data/storm"
)

var feedsCmd = &cobra.Command{
	Use: "feeds",
}

var feedsListCmd = &cobra.Command{
	Use: "list",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := storm.NewStormFeedsDatabase("feeds.db")
		if err != nil {
			return err
		}

		feeds, err := db.GetFeedList()
		if err != nil {
			return err
		}

		for _, f := range feeds {
			fmt.Printf("%s: %s", f.Name, f.Addr.String())
		}

		return nil
	},
}

var feedsAddCmd = &cobra.Command{
	Use:  "add <feed name> <feed URL>",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := storm.NewStormFeedsDatabase("feeds.db")
		if err != nil {
			return err
		}

		u, err := url.Parse(args[1])
		if err != nil {
			return err
		}

		_, err = db.AddFeed(args[0], *u)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	feedsCmd.AddCommand(feedsListCmd)
	feedsCmd.AddCommand(feedsAddCmd)
}
