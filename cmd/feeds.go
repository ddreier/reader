package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
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
			fmt.Printf("%s %s: %s\n", f.ID.String(), f.Name, f.Addr.String())
		}

		return nil
	},
}

var feedsInfoCmd = &cobra.Command{
	Use:  "info <feed id>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := storm.NewStormFeedsDatabase("feeds.db")
		if err != nil {
			return err
		}

		id, err := uuid.Parse(args[0])
		if err != nil {
			return err
		}

		feed, err := db.GetFeedById(id)
		if err != nil {
			return err
		}

		fmt.Printf("ID:    %s\n", feed.ID)
		fmt.Printf("Name:  %s\n", feed.Name)
		fmt.Printf("URL:   %s\n", feed.Addr.String())
		fmt.Printf("Added: %s\n", feed.AddTime)
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

var feedsRemoveCmd = &cobra.Command{
	Use:  "remove <feed id>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := storm.NewStormFeedsDatabase("feeds.db")
		if err != nil {
			return err
		}

		id, err := uuid.Parse(args[0])
		if err != nil {
			return err
		}

		err = db.DeleteFeed(id)
		return err
	},
}

var feedsRefreshCmd = &cobra.Command{
	Use: "refresh",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := storm.NewStormFeedsDatabase("feeds.db")
		if err != nil {
			return err
		}

		feeds, err := db.GetFeedList()
		if err != nil {
			return err
		}

		fp := gofeed.NewParser()
		for i, f := range feeds {
			feed, err := fp.ParseURL(f.Addr.String())
			if err != nil {
				fmt.Printf("Error parsing feed %s (%s): %s\n", f.Name, f.Addr.String(), err)
				continue
			}

			if i > 0 {
				fmt.Println("----------------------------------")
			}
			fmt.Printf("Feed Title: %s\n", feed.Title)
			fmt.Printf("Feed Description: %s\n", feed.Description)
			fmt.Printf("Feed Updated: %s\n", feed.UpdatedParsed)
		}

		return nil
	},
}

func init() {
	feedsCmd.AddCommand(feedsListCmd)
	feedsCmd.AddCommand(feedsInfoCmd)
	feedsCmd.AddCommand(feedsAddCmd)
	feedsCmd.AddCommand(feedsRemoveCmd)
	feedsCmd.AddCommand(feedsRefreshCmd)
}
