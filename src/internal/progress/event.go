package progress

type EventKind string

const (
	EventKindStatus  EventKind = "status"
	EventKindWarning EventKind = "warning"
	EventKindMetric  EventKind = "metric"
)

type Event struct {
	Kind        EventKind
	Message     string
	CurrentPath string
	FilesSeen   int
	Warnings    int
}
