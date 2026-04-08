package progress

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type TUISink struct {
	Events chan<- Event
}

func (t TUISink) Publish(event Event) error {
	if t.Events == nil {
		return nil
	}
	select {
	case t.Events <- event:
	default:
	}
	return nil
}

type stageStatus string

const (
	stagePending   stageStatus = "pending"
	stageActive    stageStatus = "active"
	stageCompleted stageStatus = "completed"
)

var stageOrder = []string{"discovery", "metadata", "aggregate", "render"}

var stageLabels = map[string]string{
	"discovery": "Discovery",
	"metadata":  "Metadata",
	"aggregate": "Aggregate",
	"render":    "Render",
}

type eventMsg struct {
	event Event
	ok    bool
}

type TUIModel struct {
	spinner           spinner.Model
	events            <-chan Event
	stages            map[string]stageStatus
	currentStage      string
	latest            Event
	discoverySnapshot Event
	metadataProcessed int
	metadataTotal     int
	completionNote    string
	warnings          []string
}

func NewTUIModel(events <-chan Event) TUIModel {
	spin := spinner.New()
	spin.Spinner = spinner.Line

	stages := make(map[string]stageStatus, len(stageOrder))
	for _, stage := range stageOrder {
		stages[stage] = stagePending
	}

	return TUIModel{spinner: spin, events: events, stages: stages}
}

func (m TUIModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, waitForEvent(m.events))
}

func waitForEvent(events <-chan Event) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-events
		return eventMsg{event: event, ok: ok}
	}
}

func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case eventMsg:
		if !msg.ok {
			return m, tea.Quit
		}
		m.applyEvent(msg.event)
		return m, waitForEvent(m.events)
	default:
		return m, nil
	}
}

func (m *TUIModel) applyEvent(event Event) {
	m.latest = event
	if event.Stage != "" {
		m.currentStage = event.Stage
	}

	switch event.Kind {
	case EventKindStageStart:
		for stage := range m.stages {
			if m.stages[stage] == stageActive {
				m.stages[stage] = stageCompleted
			}
		}
		m.stages[event.Stage] = stageActive
	case EventKindStageEnd:
		m.stages[event.Stage] = stageCompleted
	case EventKindMetric:
		m.metadataProcessed = event.ProcessedCount
		m.metadataTotal = event.TotalCount
	case EventKindWarning:
		m.warnings = append(m.warnings, fmt.Sprintf("%s: %s", stageLabels[event.Stage], event.Message))
		if len(m.warnings) > 3 {
			m.warnings = m.warnings[len(m.warnings)-3:]
		}
	case EventKindStatus:
		if event.Stage == "discovery" {
			m.discoverySnapshot = event
		}
		if event.Stage == "metadata" && event.TotalCount > 0 {
			m.metadataProcessed = event.ProcessedCount
			m.metadataTotal = event.TotalCount
		}
		if event.Stage == "" && event.Message != "" && event.Message != "run started" && event.Message != "run complete" {
			m.completionNote = event.Message
		}
	}
}

func (m TUIModel) View() string {
	var lines []string
	lines = append(lines, "focalytics progress")

	for _, stage := range stageOrder {
		status := m.stages[stage]
		prefix := "○"
		switch status {
		case stageActive:
			prefix = m.spinner.View()
		case stageCompleted:
			prefix = "✓"
		}

		line := fmt.Sprintf("%s %s", prefix, stageLabels[stage])
		if metrics := m.metricsForStage(stage); metrics != "" {
			line = line + "  " + metrics
		}
		lines = append(lines, line)
	}

	if len(m.warnings) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Warnings:")
		for _, warning := range m.warnings {
			lines = append(lines, "- "+warning)
		}
	}
	if m.completionNote != "" {
		lines = append(lines, "")
		lines = append(lines, m.completionNote)
	}

	return strings.Join(lines, "\n") + "\n"
}

func (m TUIModel) metricsForStage(stage string) string {
	if stage == "discovery" && (m.discoverySnapshot.FilesSeen > 0 || m.discoverySnapshot.CandidatesFound > 0 || m.discoverySnapshot.DirectoriesSeen > 0) {
		return fmt.Sprintf("files=%d candidates=%d dirs=%d warnings=%d rate=%.2f/s", m.discoverySnapshot.FilesSeen, m.discoverySnapshot.CandidatesFound, m.discoverySnapshot.DirectoriesSeen, m.discoverySnapshot.Warnings, m.discoverySnapshot.ThroughputPerSecond)
	}
	if stage == "metadata" && m.metadataTotal > 0 {
		return fmt.Sprintf("processed=%d/%d warnings=%d", m.metadataProcessed, m.metadataTotal, m.latest.Warnings)
	}
	return ""
}
