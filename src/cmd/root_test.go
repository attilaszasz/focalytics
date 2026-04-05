package cmd

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/attila/focalytics/internal/app"
	"github.com/attila/focalytics/internal/progress"
)

type fakeRunner struct {
	result  app.RunResult
	err     error
	called  bool
	request app.ScanRequest
}

func (f *fakeRunner) Run(_ context.Context, request app.ScanRequest) (app.RunResult, error) {
	f.called = true
	f.request = request
	return f.result, f.err
}

func TestNewRootCommandSetsExpectedUse(t *testing.T) {
	streams := IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}}
	command := NewRootCommand(&fakeRunner{}, app.DefaultExitPolicy(), streams)

	if got, want := command.Use, "focalytics [archive-root]"; got != want {
		t.Fatalf("unexpected use string: got %q want %q", got, want)
	}
	if len(command.Commands()) != 1 {
		t.Fatalf("expected one subcommand, got %d", len(command.Commands()))
	}
}

func TestExecuteRequiresArchiveRoot(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Execute(nil, bytes.NewBuffer(nil), stdout, stderr)

	if exitCode != app.DefaultExitPolicy().InvalidInput {
		t.Fatalf("unexpected exit code: got %d want %d", exitCode, app.DefaultExitPolicy().InvalidInput)
	}
	if stderr.Len() == 0 {
		t.Fatal("expected validation error on stderr")
	}
}

func TestSupportsInteractiveOutputRequiresTerminalWriters(t *testing.T) {
	stdout, err := os.CreateTemp(t.TempDir(), "stdout")
	if err != nil {
		t.Fatalf("create temp stdout: %v", err)
	}
	stderr, err := os.CreateTemp(t.TempDir(), "stderr")
	if err != nil {
		t.Fatalf("create temp stderr: %v", err)
	}
	defer func() {
		_ = stdout.Close()
		_ = stderr.Close()
	}()

	if supportsInteractiveOutput(stdout, stderr) {
		t.Fatal("expected regular files to be non-interactive")
	}
	if supportsInteractiveOutput(&bytes.Buffer{}, &bytes.Buffer{}) {
		t.Fatal("expected buffers to be non-interactive")
	}
}

func TestRunInteractiveProgressReturnsAfterClose(t *testing.T) {
	events := make(chan progress.Event, 1)
	buffer := &bytes.Buffer{}
	wait := runInteractiveProgress(events, buffer)
	events <- progress.Event{Kind: progress.EventKindStageStart, Stage: "discovery", Message: "stage started"}

	done := make(chan struct{})
	go func() {
		wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("interactive progress did not shut down")
	}
}
