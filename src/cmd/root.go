package cmd

import (
	"io"
	"log"

	"github.com/spf13/cobra"

	"github.com/attila/focalytics/internal/app"
	"github.com/attila/focalytics/internal/discovery"
	"github.com/attila/focalytics/internal/pipeline"
	"github.com/attila/focalytics/internal/progress"
)

type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

func NewRootCommand(runner app.Runner, exitPolicy app.ExitPolicy, streams IOStreams) *cobra.Command {
	runHandler := newRunHandler(runner, exitPolicy, streams)

	command := &cobra.Command{
		Use:           "focalytics [archive-root]",
		Short:         "Analyze a local photo archive",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE:          runHandler,
	}

	command.AddCommand(NewRunCommand(runner, exitPolicy, streams))
	command.SetOut(streams.Out)
	command.SetErr(streams.ErrOut)

	return command
}

func Execute(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	streams := IOStreams{
		In:     stdin,
		Out:    stdout,
		ErrOut: stderr,
	}
	exitPolicy := app.DefaultExitPolicy()
	sink := progress.TextSink{Writer: stderr}
	runner := pipeline.NewRunner([]pipeline.Stage{discovery.NewStage(discovery.NewService())}, sink, log.New(stderr, "", 0), exitPolicy)
	command := NewRootCommand(runner, exitPolicy, streams)
	if len(args) > 0 && args[0] != "run" && args[0] != "help" && args[0] != "completion" && args[0][0] != '-' {
		args = append([]string{"run"}, args...)
	}
	command.SetArgs(args)

	if err := command.Execute(); err != nil {
		_, _ = io.WriteString(stderr, err.Error()+"\n")
		return app.ExitCodeForError(err, exitPolicy)
	}

	return exitPolicy.Success
}
