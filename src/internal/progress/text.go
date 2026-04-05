package progress

import (
	"fmt"
	"io"
)

type TextSink struct {
	Writer io.Writer
}

func (t TextSink) Publish(event Event) error {
	if t.Writer == nil {
		return nil
	}

	_, err := fmt.Fprintf(
		t.Writer,
		"%s path=%s files=%d dirs=%d candidates=%d warnings=%d throughput=%.2f/s\n",
		event.Message,
		event.CurrentPath,
		event.FilesSeen,
		event.DirectoriesSeen,
		event.CandidatesFound,
		event.Warnings,
		event.ThroughputPerSecond,
	)
	return err
}
