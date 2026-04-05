package metadata

import (
	"context"
	"fmt"

	"github.com/attila/focalytics/internal/app"
	"github.com/attila/focalytics/internal/discovery"
)

type Stage struct {
	service Service
}

func NewStage(service Service) Stage {
	return Stage{service: service}
}

func (s Stage) Name() string {
	return "metadata"
}

func (s Stage) Run(_ context.Context, runContext app.RunContext) (app.StageResult, error) {
	artifact, ok := runContext.Artifact(app.ArtifactDiscoveryResult)
	if !ok {
		return app.StageResult{StageName: s.Name(), Status: app.StageStatusFailure, Fatal: true, ErrorMessage: "discovery result missing"}, fmt.Errorf("discovery result missing")
	}
	discoveryResult, ok := artifact.(discovery.Result)
	if !ok {
		return app.StageResult{StageName: s.Name(), Status: app.StageStatusFailure, Fatal: true, ErrorMessage: "discovery artifact type mismatch"}, fmt.Errorf("discovery artifact type mismatch")
	}

	result, err := s.service.Recover(discoveryResult, runContext.ProgressSink)
	if err != nil {
		return app.StageResult{StageName: s.Name(), Status: app.StageStatusFailure, Fatal: true, ErrorMessage: err.Error()}, err
	}
	runContext.SetArtifact(app.ArtifactMetadataResult, result)

	return app.StageResult{StageName: s.Name(), Status: app.StageStatusSuccess}, nil
}
