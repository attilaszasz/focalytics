package pipeline

import (
	"bytes"
	"context"
	"errors"
	"log"
	"testing"

	"github.com/attila/focalytics/internal/app"
	"github.com/attila/focalytics/internal/progress"
)

type fakeStage struct {
	name   string
	result app.StageResult
	err    error
}

func (f fakeStage) Name() string {
	return f.name
}

func (f fakeStage) Run(context.Context, app.RunContext) (app.StageResult, error) {
	return f.result, f.err
}

type recordingSink struct {
	events []progress.Event
}

func (r *recordingSink) Publish(event progress.Event) error {
	r.events = append(r.events, event)
	return nil
}

func TestNewRunnerDefaultsToNoopSink(t *testing.T) {
	runner := NewRunner(nil, nil, log.New(&bytes.Buffer{}, "", 0), app.DefaultExitPolicy())
	if runner.progress == nil {
		t.Fatal("expected nil sink to be replaced")
	}
}

func TestRunnerRunPublishesLifecycleEvents(t *testing.T) {
	sink := &recordingSink{}
	runner := NewRunner([]Stage{
		fakeStage{name: "noop", result: app.StageResult{StageName: "noop", Status: app.StageStatusSuccess}},
	}, sink, log.New(&bytes.Buffer{}, "", 0), app.DefaultExitPolicy())

	result, err := runner.Run(context.Background(), app.ScanRequest{ArchiveRoot: t.TempDir()})
	if err != nil {
		t.Fatalf("unexpected runner error: %v", err)
	}
	if result.ExitCode != app.DefaultExitPolicy().Success {
		t.Fatalf("unexpected exit code: got %d", result.ExitCode)
	}
	if len(sink.events) != 2 {
		t.Fatalf("expected start and completion events, got %d", len(sink.events))
	}
}

func TestRunnerRunReturnsRuntimeFailureOnStageError(t *testing.T) {
	sink := &recordingSink{}
	runner := NewRunner([]Stage{
		fakeStage{name: "broken", result: app.StageResult{StageName: "broken", Status: app.StageStatusFailure}, err: errors.New("broken stage")},
	}, sink, log.New(&bytes.Buffer{}, "", 0), app.DefaultExitPolicy())

	result, err := runner.Run(context.Background(), app.ScanRequest{ArchiveRoot: t.TempDir()})
	if err == nil {
		t.Fatal("expected stage error")
	}
	if result.ExitCode != app.DefaultExitPolicy().RuntimeFailure {
		t.Fatalf("unexpected exit code: got %d want %d", result.ExitCode, app.DefaultExitPolicy().RuntimeFailure)
	}
	if len(sink.events) == 0 {
		t.Fatal("expected warning event to be published")
	}
}

func TestRunnerRunReturnsRuntimeFailureOnFatalStageResult(t *testing.T) {
	runner := NewRunner([]Stage{
		fakeStage{name: "fatal", result: app.StageResult{StageName: "fatal", Status: app.StageStatusFailure, Fatal: true}},
	}, progress.NoopSink{}, log.New(&bytes.Buffer{}, "", 0), app.DefaultExitPolicy())

	result, err := runner.Run(context.Background(), app.ScanRequest{ArchiveRoot: t.TempDir()})
	if err == nil {
		t.Fatal("expected fatal stage result to fail")
	}
	if result.ExitCode != app.DefaultExitPolicy().RuntimeFailure {
		t.Fatalf("unexpected exit code: got %d want %d", result.ExitCode, app.DefaultExitPolicy().RuntimeFailure)
	}
}
