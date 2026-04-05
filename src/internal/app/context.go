package app

import (
	"context"
	"log"

	"github.com/attila/focalytics/internal/progress"
)

type RunContext struct {
	Request      ScanRequest
	ExitPolicy   ExitPolicy
	ProgressSink progress.Sink
	Logger       *log.Logger
	Artifacts    map[string]any
}

const (
	ArtifactDiscoveryResult = "discovery.result"
	ArtifactMetadataResult  = "metadata.result"
	ArtifactAggregateResult = "aggregate.result"
	ArtifactRenderResult    = "render.result"
)

type StageStatus string

const (
	StageStatusSkipped StageStatus = "skipped"
	StageStatusSuccess StageStatus = "success"
	StageStatusFailure StageStatus = "failure"
)

type StageResult struct {
	StageName    string
	Status       StageStatus
	Fatal        bool
	ErrorMessage string
}

type RunResult struct {
	ExitCode     int
	StageResults []StageResult
}

type Runner interface {
	Run(ctx context.Context, request ScanRequest) (RunResult, error)
}

func NewRunContext(request ScanRequest, exitPolicy ExitPolicy, sink progress.Sink, logger *log.Logger) RunContext {
	return RunContext{
		Request:      request,
		ExitPolicy:   exitPolicy,
		ProgressSink: sink,
		Logger:       logger,
		Artifacts:    map[string]any{},
	}
}

func (r RunContext) SetArtifact(key string, value any) {
	if r.Artifacts == nil {
		return
	}
	r.Artifacts[key] = value
}

func (r RunContext) Artifact(key string) (any, bool) {
	if r.Artifacts == nil {
		return nil, false
	}
	value, ok := r.Artifacts[key]
	return value, ok
}
