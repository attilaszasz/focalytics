package discovery

import (
	"context"

	"github.com/attila/focalytics/internal/app"
)

type Stage struct {
	service Service
}

func NewStage(service Service) Stage {
	return Stage{service: service}
}

func (s Stage) Name() string {
	return "discovery"
}

func (s Stage) Run(_ context.Context, runContext app.RunContext) (app.StageResult, error) {
	_, err := s.service.Discover(runContext.Request.ArchiveRoot, runContext.ProgressSink, runContext.Request.Stdout)
	if err != nil {
		return app.StageResult{StageName: s.Name(), Status: app.StageStatusFailure, Fatal: true, ErrorMessage: err.Error()}, err
	}

	return app.StageResult{StageName: s.Name(), Status: app.StageStatusSuccess}, nil
}
