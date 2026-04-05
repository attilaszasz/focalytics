package progress

type EventKind string

const (
	EventKindStatus     EventKind = "status"
	EventKindWarning    EventKind = "warning"
	EventKindMetric     EventKind = "metric"
	EventKindStageStart EventKind = "stage_start"
	EventKindStageEnd   EventKind = "stage_end"
)

type Event struct {
	Stage               string
	Kind                EventKind
	Message             string
	CurrentPath         string
	FilesSeen           int
	Warnings            int
	CandidatesFound     int
	DirectoriesSeen     int
	ProcessedCount      int
	TotalCount          int
	ThroughputPerSecond float64
}
