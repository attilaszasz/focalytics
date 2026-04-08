package aggregate

import (
	"context"
	"fmt"

	"github.com/attila/focalytics/internal/app"
	"github.com/attila/focalytics/internal/metadata"
)

type Stage struct {
	service Service
}

func NewStage(service Service) Stage {
	return Stage{service: service}
}

func (s Stage) Name() string {
	return "aggregate"
}

func (s Stage) Run(_ context.Context, runContext app.RunContext) (app.StageResult, error) {
	artifact, ok := runContext.Artifact(app.ArtifactMetadataResult)
	if !ok {
		return app.StageResult{StageName: s.Name(), Status: app.StageStatusFailure, Fatal: true, ErrorMessage: "metadata result missing"}, fmt.Errorf("metadata result missing")
	}
	metadataResult, ok := artifact.(metadata.Result)
	if !ok {
		return app.StageResult{StageName: s.Name(), Status: app.StageStatusFailure, Fatal: true, ErrorMessage: "metadata artifact type mismatch"}, fmt.Errorf("metadata artifact type mismatch")
	}

	result := s.service.Aggregate(metadataResult, runContext.Request.IgnorePhonePhotos)
	runContext.SetArtifact(app.ArtifactAggregateResult, result)

	return app.StageResult{StageName: s.Name(), Status: app.StageStatusSuccess}, nil
}
