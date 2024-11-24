package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/simulot/immich-go/app"
	"github.com/simulot/immich-go/app/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// immich-go entry point
func main() {
	ctx := context.Background()
	err := immichGoMain(ctx)
	if err != nil {
		if e := context.Cause(ctx); e != nil {
			err = e
		}
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// makes immich-go breakable with ^C and run it
func immichGoMain(ctx context.Context) error {
	// Create a context with cancel function to gracefully handle Ctrl+C events
	ctx, cancel := context.WithCancelCause(ctx)

	// Handle Ctrl+C signal (SIGINT)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	// Watch for ^C to be pressed
	go func() {
		<-signalChannel
		fmt.Println("\nCtrl+C received. Shutting down...")
		cancel(errors.New("Ctrl+C received")) // Cancel the context when Ctrl+C is received
	}()

	return runImmichGo(ctx)
}

// Run immich-go
func runImmichGo(ctx context.Context) error {
	viper.SetEnvPrefix("IMMICHGO")

	// Create the application context

	// Add the root command
	c := &cobra.Command{
		Use:     "immich-go",
		Short:   "Immich-go is a command line application to interact with the Immich application using its API",
		Long:    `An alternative to the immich-CLI command that doesn't depend on nodejs installation. It tries its best for importing google photos takeout archives.`,
		Version: app.Version,
	}
	cobra.EnableTraverseRunHooks = true // doc: cobra/site/content/user_guide.md
	a := app.New(ctx, c)

	// add immich-go commands
	c.AddCommand(app.NewVersionCommand(ctx, a))
	cmd.AddCommands(c, ctx, a)

	// let's start
	err := c.ExecuteContext(ctx)
	if err != nil && a.Log().GetSLog() != nil {
		a.Log().Error(err.Error())
	}

	return err
}
