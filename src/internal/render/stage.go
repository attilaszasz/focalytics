package render

import (
	"context"
	"fmt"

	"github.com/attila/focalytics/internal/aggregate"
	"github.com/attila/focalytics/internal/app"
)

type Stage struct {
	service Service
}

func NewStage(service Service) Stage {
	return Stage{service: service}
}

func (s Stage) Name() string {
	return "render"
}

func (s Stage) Run(_ context.Context, runContext app.RunContext) (app.StageResult, error) {
	artifact, ok := runContext.Artifact(app.ArtifactAggregateResult)
	if !ok {
		return app.StageResult{StageName: s.Name(), Status: app.StageStatusFailure, Fatal: true, ErrorMessage: "aggregate result missing"}, fmt.Errorf("aggregate result missing")
	}
	aggregateResult, ok := artifact.(aggregate.Result)
	if !ok {
		return app.StageResult{StageName: s.Name(), Status: app.StageStatusFailure, Fatal: true, ErrorMessage: "aggregate artifact type mismatch"}, fmt.Errorf("aggregate artifact type mismatch")
	}

	result, err := s.service.Generate(aggregateResult, runContext.Request.ArchiveRoot, runContext.Request.Stdout)
	if err != nil {
		return app.StageResult{StageName: s.Name(), Status: app.StageStatusFailure, Fatal: true, ErrorMessage: err.Error()}, err
	}
	runContext.SetArtifact(app.ArtifactRenderResult, result)

	return app.StageResult{StageName: s.Name(), Status: app.StageStatusSuccess}, nil
}
