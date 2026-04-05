package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/attila/focalytics/internal/app"
)

func NewRunCommand(runner app.Runner, exitPolicy app.ExitPolicy, streams IOStreams) *cobra.Command {
	command := &cobra.Command{
		Use:           "run [archive-root]",
		Short:         "Run one archive scan",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE:          newRunHandler(runner, exitPolicy, streams),
	}

	command.SetOut(streams.Out)
	command.SetErr(streams.ErrOut)

	return command
}

func newRunHandler(runner app.Runner, exitPolicy app.ExitPolicy, streams IOStreams) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, args []string) error {
		request, err := buildScanRequest(args, streams)
		if err != nil {
			return app.InvalidInputError(err)
		}

		result, err := runner.Run(context.Background(), request)
		if err != nil {
			return app.RuntimeFailureError(err)
		}

		if result.ExitCode != exitPolicy.Success {
			return app.CommandErrorFromCode(result.ExitCode, fmt.Errorf("command exited with code %d", result.ExitCode))
		}

		return nil
	}
}

func buildScanRequest(args []string, streams IOStreams) (app.ScanRequest, error) {
	if len(args) != 1 {
		return app.ScanRequest{}, fmt.Errorf("exactly one archive root path is required")
	}

	archiveRoot := filepath.Clean(args[0])
	info, err := os.Stat(archiveRoot)
	if err != nil {
		return app.ScanRequest{}, fmt.Errorf("archive root %q is not accessible: %w", archiveRoot, err)
	}
	if !info.IsDir() {
		return app.ScanRequest{}, fmt.Errorf("archive root %q must be a directory", archiveRoot)
	}

	return app.ScanRequest{
		ArchiveRoot: archiveRoot,
		Interactive: false,
		Stdout:      streams.Out,
		Stderr:      streams.ErrOut,
	}, nil
}
