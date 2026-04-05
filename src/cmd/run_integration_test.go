//go:build integration

package cmd

import (
	"bytes"
	"testing"

	"github.com/attila/focalytics/internal/app"
)

func TestExecuteIntegrationWithLocalDirectory(t *testing.T) {
	archiveRoot := t.TempDir()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Execute([]string{archiveRoot}, bytes.NewBuffer(nil), stdout, stderr)
	if exitCode != app.DefaultExitPolicy().Success {
		t.Fatalf("expected success exit code, got %d", exitCode)
	}
}
