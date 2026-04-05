package progress

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func TestTextSinkPublishesCompactStatusLine(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := TextSink{Writer: buffer, IncludeStatus: true}

	err := sink.Publish(Event{
		Kind:                EventKindStatus,
		Stage:               "discovery",
		Message:             "candidate discovered",
		CurrentPath:         "archive/2024/photo.jpg",
		FilesSeen:           12,
		DirectoriesSeen:     3,
		CandidatesFound:     4,
		Warnings:            1,
		ThroughputPerSecond: 6.5,
	})
	if err != nil {
		t.Fatalf("expected publish success: %v", err)
	}

	output := buffer.String()
	for _, fragment := range []string{"candidate discovered", "stage=discovery", "path=archive/2024/photo.jpg", "files=12", "dirs=3", "candidates=4", "warnings=1", "throughput=6.50/s"} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("expected fragment %q in output %q", fragment, output)
		}
	}
}

func TestTextSinkAlwaysPublishesWarnings(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := TextSink{Writer: buffer}

	err := sink.Publish(Event{Kind: EventKindWarning, CurrentPath: "archive/2024/photo.jpg", Message: "embedded metadata unavailable"})
	if err != nil {
		t.Fatalf("expected publish success: %v", err)
	}

	output := buffer.String()
	for _, fragment := range []string{"warning", "path=archive/2024/photo.jpg", "embedded metadata unavailable"} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("expected fragment %q in output %q", fragment, output)
		}
	}
}

func TestTextSinkPublishesMetricAndStageMessagesWhenEnabled(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := TextSink{Writer: buffer, IncludeStatus: true}

	if err := sink.Publish(Event{Kind: EventKindStageStart, Stage: "metadata", Message: "stage started"}); err != nil {
		t.Fatalf("publish stage event: %v", err)
	}
	if err := sink.Publish(Event{Kind: EventKindMetric, Stage: "metadata", Message: "metadata progress", ProcessedCount: 5, TotalCount: 10, Warnings: 2}); err != nil {
		t.Fatalf("publish metric event: %v", err)
	}

	output := buffer.String()
	for _, fragment := range []string{"stage started stage=metadata", "metadata progress stage=metadata processed=5 total=10 warnings=2"} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("expected fragment %q in output %q", fragment, output)
		}
	}
}

func TestTUISinkPublishesToChannel(t *testing.T) {
	events := make(chan Event, 1)
	sink := TUISink{Events: events}
	expected := Event{Kind: EventKindStatus, Stage: "discovery", FilesSeen: 4}

	if err := sink.Publish(expected); err != nil {
		t.Fatalf("publish failed: %v", err)
	}
	got := <-events
	if got.Stage != expected.Stage || got.FilesSeen != expected.FilesSeen || got.Kind != expected.Kind {
		t.Fatalf("unexpected event: got %+v want %+v", got, expected)
	}
}

func TestWaitForEventReturnsEventMessage(t *testing.T) {
	events := make(chan Event, 1)
	events <- Event{Kind: EventKindStatus, Stage: "discovery", FilesSeen: 7}

	msg := waitForEvent(events)()
	event, ok := msg.(eventMsg)
	if !ok {
		t.Fatalf("expected eventMsg, got %T", msg)
	}
	if !event.ok || event.event.FilesSeen != 7 {
		t.Fatalf("unexpected event message: %+v", event)
	}
}

func TestTUIModelInitAndUpdatePaths(t *testing.T) {
	model := NewTUIModel(nil)
	if cmd := model.Init(); cmd == nil {
		t.Fatal("expected init command")
	}

	updated, cmd := model.Update(spinner.TickMsg{})
	if cmd == nil {
		t.Fatal("expected spinner update command")
	}
	if _, ok := updated.(TUIModel); !ok {
		t.Fatalf("expected TUIModel, got %T", updated)
	}

	updated, cmd = model.Update(eventMsg{ok: false})
	if _, ok := updated.(TUIModel); !ok {
		t.Fatalf("expected TUIModel on close, got %T", updated)
	}
	if cmd == nil {
		t.Fatal("expected quit command")
	}

	updated, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	if _, ok := updated.(TUIModel); !ok {
		t.Fatalf("expected TUIModel on default update, got %T", updated)
	}
	if cmd != nil {
		t.Fatalf("expected nil command for default update, got %v", cmd)
	}
}
