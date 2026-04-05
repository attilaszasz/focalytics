package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunReturnsSuccessForValidDirectory(t *testing.T) {
	archiveRoot := t.TempDir()
	stdout, err := os.CreateTemp(t.TempDir(), "stdout")
	if err != nil {
		t.Fatalf("create stdout temp file: %v", err)
	}
	t.Cleanup(func() { _ = stdout.Close() })

	stderr, err := os.CreateTemp(t.TempDir(), "stderr")
	if err != nil {
		t.Fatalf("create stderr temp file: %v", err)
	}
	t.Cleanup(func() { _ = stderr.Close() })

	exitCode := run([]string{filepath.Clean(archiveRoot)}, os.Stdin, stdout, stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: got %d want 0", exitCode)
	}
}

func TestMainUsesRunExitCode(t *testing.T) {
	originalArgs := os.Args
	originalExit := exitFunc
	t.Cleanup(func() {
		os.Args = originalArgs
		exitFunc = originalExit
	})

	archiveRoot := t.TempDir()
	os.Args = []string{"focalytics", archiveRoot}

	called := false
	exitCode := -1
	exitFunc = func(code int) {
		called = true
		exitCode = code
	}

	main()

	if !called {
		t.Fatal("expected main to call exit")
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: got %d want 0", exitCode)
	}
}
