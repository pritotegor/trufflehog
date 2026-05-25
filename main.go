package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/sirupsen/logrus"

	"github.com/trufflesecurity/trufflehog/v3/pkg/cmd"
)

func init() {
	// Maximize CPU utilization for parallel scanning.
	// Use all available CPUs since this runs on a dedicated scan machine
	// and we don't need to reserve one for system responsiveness.
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Set up context with cancellation for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logrus.WithField("signal", sig).Info("Received signal, shutting down...")
		cancel()
	}()

	// Build and run the CLI application.
	app := kingpin.New("trufflehog", "Find credentials all over the place.").Author("TruffleSecurity")
	app.HelpFlag.Short('h')
	app.Version(versionString())

	if err := cmd.Run(ctx, app, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// versionString returns the version information for the binary.
// Build-time variables are injected via ldflags:
//
//	-X main.version=<version>
//	-X main.buildDate=<date>
//	-X main.buildCommit=<commit>
var (
	version     = "dev"
	buildDate   = "unknown"
	buildCommit = "unknown"
)

func versionString() string {
	return fmt.Sprintf("%s (commit=%s, built=%s)", version, buildCommit, buildDate)
}
