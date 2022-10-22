package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"reader/data/storm"
	"reader/http"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Setup signal handlers.
		ctx, cancel := context.WithCancel(context.Background())
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() { <-c; cancel() }()

		fdb, err := storm.NewStormFeedsDatabase("feeds.db")
		if err != nil {
			return err
		}
		udb, err := storm.NewStormItemsDatabase("unread.db")
		if err != nil {
			return err
		}

		s := http.NewServer()
		s.Addr = ":8888"
		s.Feeds = fdb
		s.UnreadItems = udb

		if err := s.Open(); err != nil {
			return err
		}

		// Wait for CTRL-C.
		<-ctx.Done()

		return nil
	},
}
