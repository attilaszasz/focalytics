//go:build integration

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/attila/focalytics/internal/app"
)

func TestExecuteIntegrationWithLocalDirectory(t *testing.T) {
	archiveRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(archiveRoot, "nested"), 0o755); err != nil {
		t.Fatalf("create nested dir: %v", err)
	}
	for path := range map[string]struct{}{
		filepath.Join(archiveRoot, "nested", "frame.jpg"):  {},
		filepath.Join(archiveRoot, "nested", "frame.xmp"):  {},
		filepath.Join(archiveRoot, "nested", "ignore.txt"): {},
		filepath.Join(archiveRoot, "cover.jpeg"):           {},
	} {
		if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Execute([]string{archiveRoot}, bytes.NewBuffer(nil), stdout, stderr)
	if exitCode != app.DefaultExitPolicy().Success {
		t.Fatalf("expected success exit code, got %d", exitCode)
	}
	output := stdout.String()
	for _, fragment := range []string{"image\tcover.jpeg", "image\tnested/frame.jpg", "sidecar\tnested/frame.xmp"} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("expected stdout to contain %q, got %q", fragment, output)
		}
	}
	if strings.Contains(output, "ignore.txt") {
		t.Fatalf("expected unsupported files to be skipped, got %q", output)
	}
	stderrOutput := stderr.String()
	if !strings.Contains(stderrOutput, "throughput=") {
		t.Fatalf("expected progress metrics on stderr, got %q", stderrOutput)
	}
	if !strings.Contains(stderrOutput, "embedded metadata unavailable") {
		t.Fatalf("expected metadata warning on stderr, got %q", stderrOutput)
	}
}
