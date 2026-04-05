package app

import (
	"bytes"
	"log"
	"testing"

	"github.com/attila/focalytics/internal/progress"
)

func TestNewRunContextAssignsDependencies(t *testing.T) {
	request := ScanRequest{ArchiveRoot: "/tmp/archive"}
	logger := log.New(&bytes.Buffer{}, "", 0)
	sink := progress.NoopSink{}

	runContext := NewRunContext(request, DefaultExitPolicy(), sink, logger)

	if runContext.Request.ArchiveRoot != request.ArchiveRoot {
		t.Fatalf("unexpected archive root: got %q want %q", runContext.Request.ArchiveRoot, request.ArchiveRoot)
	}
	if runContext.ProgressSink == nil {
		t.Fatal("expected progress sink to be set")
	}
	if runContext.Logger != logger {
		t.Fatal("expected logger to be retained")
	}
}
