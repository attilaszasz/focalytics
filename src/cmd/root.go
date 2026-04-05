package cmd

import (
	"io"
	"log"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/attila/focalytics/internal/aggregate"
	"github.com/attila/focalytics/internal/app"
	"github.com/attila/focalytics/internal/discovery"
	"github.com/attila/focalytics/internal/metadata"
	"github.com/attila/focalytics/internal/pipeline"
	"github.com/attila/focalytics/internal/progress"
	"github.com/attila/focalytics/internal/render"
)

type IOStreams struct {
	In          io.Reader
	Out         io.Writer
	ErrOut      io.Writer
	Interactive bool
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
	interactive := supportsInteractiveOutput(stdout, stderr)
	streams := IOStreams{
		In:          stdin,
		Out:         stdout,
		ErrOut:      stderr,
		Interactive: interactive,
	}
	exitPolicy := app.DefaultExitPolicy()
	sink := progress.Sink(progress.TextSink{Writer: stderr})
	var programWait func()
	if interactive {
		events := make(chan progress.Event, 512)
		sink = progress.TUISink{Events: events}
		programWait = runInteractiveProgress(events, stderr)
	}
	runner := newRunner(sink, stderr, exitPolicy)
	command := NewRootCommand(runner, exitPolicy, streams)
	args = normalizeArgs(args)
	command.SetArgs(args)

	if err := command.Execute(); err != nil {
		if programWait != nil {
			programWait()
		}
		_, _ = io.WriteString(stderr, err.Error()+"\n")
		return app.ExitCodeForError(err, exitPolicy)
	}
	if programWait != nil {
		programWait()
	}

	return exitPolicy.Success
}

func newRunner(sink progress.Sink, stderr io.Writer, exitPolicy app.ExitPolicy) *pipeline.Runner {
	return pipeline.NewRunner(
		[]pipeline.Stage{
			discovery.NewStage(discovery.NewService()),
			metadata.NewStage(metadata.NewService()),
			aggregate.NewStage(aggregate.NewService()),
			render.NewStage(render.NewService()),
		},
		sink,
		log.New(stderr, "", 0),
		exitPolicy,
	)
}

func normalizeArgs(args []string) []string {
	if len(args) > 0 && args[0] != "run" && args[0] != "help" && args[0] != "completion" && args[0][0] != '-' {
		return append([]string{"run"}, args...)
	}
	return args
}

func supportsInteractiveOutput(stdout, stderr io.Writer) bool {
	stdoutFile, ok := stdout.(interface{ Fd() uintptr })
	if !ok {
		return false
	}
	stderrFile, ok := stderr.(interface{ Fd() uintptr })
	if !ok {
		return false
	}
	return term.IsTerminal(int(stdoutFile.Fd())) && term.IsTerminal(int(stderrFile.Fd()))
}

func runInteractiveProgress(events chan progress.Event, stderr io.Writer) func() {
	program := tea.NewProgram(progress.NewTUIModel(events), tea.WithInput(nil), tea.WithOutput(stderr))
	var once sync.Once
	done := make(chan struct{})
	go func() {
		_, _ = program.Run()
		close(done)
	}()
	return func() {
		once.Do(func() {
			close(events)
			<-done
		})
	}
}
