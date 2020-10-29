package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"reader/data"
	"reader/data/storm"
)

var itemsCmd = &cobra.Command{
	Use: "items",
}

var itemsListUnreadCmd = &cobra.Command{
	Use: "list-unread",
	RunE: func(cmd *cobra.Command, args []string) error {
		feedDb, err := storm.NewStormFeedsDatabase("feeds.db")
		if err != nil {
			return err
		}
		unreadDb, err := storm.NewStormItemsDatabase("unread.db")
		if err != nil {
			return err
		}

		feedList, err := feedDb.GetFeedList()
		if err != nil {
			return err
		}
		feedMap := make(map[uuid.UUID]data.Feed)
		for _, f := range feedList {
			feedMap[f.ID] = f
		}

		items, err := unreadDb.ListItems()
		if err != nil {
			return err
		}

		for _, i := range items {
			fmt.Printf("%s %s %s\n", feedMap[i.Feed].Name, i.PubDate, i.Title)
		}

		markRead, err := cmd.Flags().GetBool("mark-read")
		if err != nil {
			return err
		}
		if markRead {
			err = unreadDb.DeleteAllItems()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	itemsListUnreadCmd.Flags().Bool("mark-read", false, "Mark the items read")

	itemsCmd.AddCommand(itemsListUnreadCmd)
}
