package progress

import "testing"

func TestNoopSinkPublishReturnsNil(t *testing.T) {
	err := NoopSink{}.Publish(Event{Kind: EventKindStatus})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
