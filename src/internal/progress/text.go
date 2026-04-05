package progress

import (
	"fmt"
	"io"
)

type TextSink struct {
	Writer        io.Writer
	IncludeStatus bool
}

func (t TextSink) Publish(event Event) error {
	if t.Writer == nil {
		return nil
	}

	switch event.Kind {
	case EventKindWarning:
		_, err := fmt.Fprintf(t.Writer, "warning path=%s message=%s\n", event.CurrentPath, event.Message)
		return err
	case EventKindStageStart, EventKindStageEnd:
		if !t.IncludeStatus {
			return nil
		}
		_, err := fmt.Fprintf(t.Writer, "%s stage=%s\n", event.Message, event.Stage)
		return err
	case EventKindMetric:
		if !t.IncludeStatus {
			return nil
		}
		_, err := fmt.Fprintf(t.Writer, "%s stage=%s processed=%d total=%d warnings=%d\n", event.Message, event.Stage, event.ProcessedCount, event.TotalCount, event.Warnings)
		return err
	case EventKindStatus:
		if !t.IncludeStatus {
			return nil
		}
		_, err := fmt.Fprintf(
			t.Writer,
			"%s stage=%s path=%s files=%d dirs=%d candidates=%d warnings=%d throughput=%.2f/s\n",
			event.Message,
			event.Stage,
			event.CurrentPath,
			event.FilesSeen,
			event.DirectoriesSeen,
			event.CandidatesFound,
			event.Warnings,
			event.ThroughputPerSecond,
		)
		return err
	default:
		return nil
	}
}
