package progress

import (
	"strings"
	"testing"
)

func TestTUIModelTracksStageAndMetrics(t *testing.T) {
	model := NewTUIModel(nil)
	model.applyEvent(Event{Kind: EventKindStageStart, Stage: "discovery", Message: "stage started"})
	model.applyEvent(Event{Kind: EventKindStatus, Stage: "discovery", FilesSeen: 12, CandidatesFound: 4, DirectoriesSeen: 3, Warnings: 1, ThroughputPerSecond: 6.5})
	model.applyEvent(Event{Kind: EventKindStageEnd, Stage: "discovery", Message: "stage complete"})
	model.applyEvent(Event{Kind: EventKindStageStart, Stage: "metadata", Message: "stage started"})
	model.applyEvent(Event{Kind: EventKindMetric, Stage: "metadata", ProcessedCount: 8, TotalCount: 10, Warnings: 2})

	view := model.View()
	for _, fragment := range []string{"✓ Discovery", "Metadata", "files=12 candidates=4 dirs=3 warnings=1 rate=6.50/s", "processed=8/10 warnings=2"} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("expected fragment %q in view %q", fragment, view)
		}
	}
}

func TestTUIModelRetainsRecentWarnings(t *testing.T) {
	model := NewTUIModel(nil)
	model.applyEvent(Event{Kind: EventKindWarning, Stage: "metadata", Message: "first"})
	model.applyEvent(Event{Kind: EventKindWarning, Stage: "metadata", Message: "second"})
	model.applyEvent(Event{Kind: EventKindWarning, Stage: "metadata", Message: "third"})
	model.applyEvent(Event{Kind: EventKindWarning, Stage: "metadata", Message: "fourth"})

	view := model.View()
	if strings.Contains(view, "first") {
		t.Fatalf("expected oldest warning to be trimmed, got %q", view)
	}
	for _, fragment := range []string{"Metadata: second", "Metadata: third", "Metadata: fourth"} {
		if !strings.Contains(view, fragment) {
			t.Fatalf("expected fragment %q in view %q", fragment, view)
		}
	}
}

func TestTUIModelShowsCompletionNote(t *testing.T) {
	model := NewTUIModel(nil)
	model.applyEvent(Event{Kind: EventKindStatus, Message: "Phone filter active: excluded 2 phone-made photos from gear and technical insights; timeline and total photos still reflect the full archive."})

	view := model.View()
	if !strings.Contains(view, "Phone filter active: excluded 2 phone-made photos") {
		t.Fatalf("expected completion note in view %q", view)
	}
}
