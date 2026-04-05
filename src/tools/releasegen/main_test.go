package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/attila/focalytics/internal/release"
)

func TestMainDispatchesAssetNameCommand(t *testing.T) {
	originalArgs := os.Args
	originalStdout := stdoutWriter
	originalStderr := stderrWriter
	originalExit := exitWithCode
	t.Cleanup(func() {
		os.Args = originalArgs
		stdoutWriter = originalStdout
		stderrWriter = originalStderr
		exitWithCode = originalExit
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	stdoutWriter = stdout
	stderrWriter = stderr
	exitWithCode = func(code int) {
		t.Fatalf("unexpected exit: %d", code)
	}
	os.Args = []string{"releasegen", "asset-name", "--version", "v1.2.3", "--goos", "linux", "--goarch", "amd64", "--field", "archive"}

	main()

	if got := strings.TrimSpace(stdout.String()); got != "focalytics_v1.2.3_linux_amd64.tar.gz" {
		t.Fatalf("unexpected stdout: %s", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %s", stderr.String())
	}
}

func TestMainDispatchesErrorExit(t *testing.T) {
	originalArgs := os.Args
	originalStdout := stdoutWriter
	originalStderr := stderrWriter
	originalExit := exitWithCode
	t.Cleanup(func() {
		os.Args = originalArgs
		stdoutWriter = originalStdout
		stderrWriter = originalStderr
		exitWithCode = originalExit
	})

	stdoutWriter = &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	stderrWriter = stderr
	exitCode := -1
	exitWithCode = func(code int) {
		exitCode = code
	}
	os.Args = []string{"releasegen", "unknown"}

	main()

	if exitCode != 1 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "unknown subcommand") {
		t.Fatalf("expected unknown subcommand message, got %s", stderr.String())
	}
}

func TestRunAssetNameArchiveAndChecksums(t *testing.T) {
	archiveOutput := captureStdout(t, func() {
		if err := runAssetName([]string{"--version", "v1.2.3", "--goos", "linux", "--goarch", "amd64", "--field", "archive"}); err != nil {
			t.Fatalf("expected archive asset name: %v", err)
		}
	})
	if got := strings.TrimSpace(archiveOutput); got != "focalytics_v1.2.3_linux_amd64.tar.gz" {
		t.Fatalf("unexpected archive output: %s", got)
	}

	checksumOutput := captureStdout(t, func() {
		if err := runAssetName([]string{"--version", "v1.2.3", "--field", "checksums"}); err != nil {
			t.Fatalf("expected checksum name: %v", err)
		}
	})
	if got := strings.TrimSpace(checksumOutput); got != "focalytics_v1.2.3_checksums.txt" {
		t.Fatalf("unexpected checksum output: %s", got)
	}
}

func TestRunAssetNameRejectsUnknownField(t *testing.T) {
	err := runAssetName([]string{"--version", "v1.2.3", "--goos", "linux", "--goarch", "amd64", "--field", "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestRunVerifyUsesChecksumContract(t *testing.T) {
	checksumPath := writeChecksumManifest(t, "v1.2.3")
	if err := runVerify([]string{"--version", "v1.2.3", "--checksums", checksumPath}); err != nil {
		t.Fatalf("expected verify success: %v", err)
	}
}

func TestRunMetadataWritesPackageInputs(t *testing.T) {
	checksumPath := writeChecksumManifest(t, "v1.2.3")
	outDir := t.TempDir()

	if err := runMetadata([]string{"--repo", "attila/focalytics", "--version", "v1.2.3", "--checksums", checksumPath, "--out", outDir}); err != nil {
		t.Fatalf("expected metadata success: %v", err)
	}

	homebrewPath := filepath.Join(outDir, "homebrew-formula.json")
	homebrewContent, err := os.ReadFile(homebrewPath)
	if err != nil {
		t.Fatalf("expected homebrew output: %v", err)
	}
	var homebrew struct {
		ReleaseTag string `json:"release_tag"`
		Artifacts  []struct {
			URL string `json:"url"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal(homebrewContent, &homebrew); err != nil {
		t.Fatalf("expected valid homebrew json: %v", err)
	}
	if homebrew.ReleaseTag != "v1.2.3" {
		t.Fatalf("unexpected release tag: %s", homebrew.ReleaseTag)
	}
	if len(homebrew.Artifacts) == 0 || !strings.Contains(homebrew.Artifacts[0].URL, "/releases/download/v1.2.3/") {
		t.Fatalf("expected release asset url, got %+v", homebrew.Artifacts)
	}

	wingetPath := filepath.Join(outDir, "winget-manifests.json")
	if _, err := os.Stat(wingetPath); err != nil {
		t.Fatalf("expected winget output: %v", err)
	}
}

func TestLoadChecksumsRequiresPath(t *testing.T) {
	_, err := loadChecksums("")
	if err == nil {
		t.Fatal("expected missing path error")
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := stdoutWriter
	buffer := &bytes.Buffer{}
	stdoutWriter = buffer
	t.Cleanup(func() {
		stdoutWriter = original
	})

	fn()

	return buffer.String()
}

func writeChecksumManifest(t *testing.T, version string) string {
	t.Helper()
	lines := make([]string, 0, len(release.DefaultTargets()))
	for _, target := range release.DefaultTargets() {
		assetName, err := release.ReleaseAssetName(version, target)
		if err != nil {
			t.Fatalf("expected asset name: %v", err)
		}
		lines = append(lines, strings.Repeat("a", 64)+"  "+assetName)
	}

	path := filepath.Join(t.TempDir(), "checksums.txt")
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("expected checksum manifest: %v", err)
	}

	return path
}
