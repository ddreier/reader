package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
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

		s := http.NewServer()
		s.Addr = ":8888"

		if err := s.Open(); err != nil {
			return err
		}

		// Wait for CTRL-C.
		<-ctx.Done()

		return nil
	},
}
