package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/attila/focalytics/internal/app"
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
