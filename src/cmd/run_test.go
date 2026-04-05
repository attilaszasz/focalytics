package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/attila/focalytics/internal/app"
)

func TestNewRunCommandAcceptsValidArchiveRoot(t *testing.T) {
	archiveRoot := t.TempDir()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	runner := &fakeRunner{result: app.RunResult{ExitCode: app.DefaultExitPolicy().Success}}

	command := NewRunCommand(runner, app.DefaultExitPolicy(), IOStreams{Out: stdout, ErrOut: stderr})
	command.SetArgs([]string{archiveRoot})

	if err := command.Execute(); err != nil {
		t.Fatalf("expected valid archive root to execute cleanly: %v", err)
	}
	if !runner.called {
		t.Fatal("expected runner to be invoked")
	}
	if runner.request.ArchiveRoot != archiveRoot {
		t.Fatalf("unexpected archive root: got %q want %q", runner.request.ArchiveRoot, archiveRoot)
	}
}

func TestNewRunCommandRejectsInvalidArchiveRoot(t *testing.T) {
	command := NewRunCommand(&fakeRunner{}, app.DefaultExitPolicy(), IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"/path/that/does/not/exist"})

	err := command.Execute()
	if err == nil {
		t.Fatal("expected invalid archive root to return an error")
	}
	if app.ExitCodeForError(err, app.DefaultExitPolicy()) != app.DefaultExitPolicy().InvalidInput {
		t.Fatalf("expected invalid-input exit code, got %d", app.ExitCodeForError(err, app.DefaultExitPolicy()))
	}
}

func TestExecuteEmitsProgressToStderr(t *testing.T) {
	archiveRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(archiveRoot, "photo.jpg"), []byte("test"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	workDir := t.TempDir()
	originalWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("chdir workdir: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalWorkingDirectory)
	}()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Execute([]string{archiveRoot}, bytes.NewBuffer(nil), stdout, stderr)
	if exitCode != app.DefaultExitPolicy().Success {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "candidate discovered") {
		t.Fatalf("expected progress output, got %q", stderr.String())
	}
}

func TestNewRunCommandPreservesRunnerStdout(t *testing.T) {
	archiveRoot := t.TempDir()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	runner := &reportingRunner{path: "/tmp/focalytics_report_20260405_1105.html"}

	command := NewRunCommand(runner, app.DefaultExitPolicy(), IOStreams{Out: stdout, ErrOut: stderr})
	command.SetArgs([]string{archiveRoot})

	if err := command.Execute(); err != nil {
		t.Fatalf("expected valid archive root to execute cleanly: %v", err)
	}
	if !strings.Contains(stdout.String(), runner.path) {
		t.Fatalf("expected stdout to contain runner report path, got %q", stdout.String())
	}
}

type reportingRunner struct {
	path string
}

func (r *reportingRunner) Run(_ context.Context, request app.ScanRequest) (app.RunResult, error) {
	_, _ = fmt.Fprintf(request.Stdout, "report\t%s\n", r.path)
	return app.RunResult{ExitCode: app.DefaultExitPolicy().Success}, nil
}
