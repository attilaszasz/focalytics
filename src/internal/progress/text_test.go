package progress

import (
	"bytes"
	"strings"
	"testing"
)

func TestTextSinkPublishesCompactStatusLine(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := TextSink{Writer: buffer}

	err := sink.Publish(Event{
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
	for _, fragment := range []string{"candidate discovered", "path=archive/2024/photo.jpg", "files=12", "dirs=3", "candidates=4", "warnings=1", "throughput=6.50/s"} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("expected fragment %q in output %q", fragment, output)
		}
	}
}
