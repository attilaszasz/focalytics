package pipeline

import (
	"context"
	"fmt"
	"log"

	"github.com/attila/focalytics/internal/app"
	"github.com/attila/focalytics/internal/progress"
)

type Runner struct {
	stages     []Stage
	progress   progress.Sink
	logger     *log.Logger
	exitPolicy app.ExitPolicy
}

func NewRunner(stages []Stage, sink progress.Sink, logger *log.Logger, exitPolicy app.ExitPolicy) *Runner {
	if sink == nil {
		sink = progress.NoopSink{}
	}

	return &Runner{
		stages:     stages,
		progress:   sink,
		logger:     logger,
		exitPolicy: exitPolicy,
	}
}

func (r *Runner) Run(ctx context.Context, request app.ScanRequest) (app.RunResult, error) {
	runContext := app.NewRunContext(request, r.exitPolicy, r.progress, r.logger)
	results := make([]app.StageResult, 0, len(r.stages))

	_ = r.progress.Publish(progress.Event{Kind: progress.EventKindStatus, Message: "run started", CurrentPath: request.ArchiveRoot})

	for _, stage := range r.stages {
		result, err := stage.Run(ctx, runContext)
		results = append(results, result)
		if err != nil {
			_ = r.progress.Publish(progress.Event{Kind: progress.EventKindWarning, Message: err.Error(), CurrentPath: request.ArchiveRoot, Warnings: 1})
			return app.RunResult{ExitCode: r.exitPolicy.RuntimeFailure, StageResults: results}, fmt.Errorf("stage %s failed: %w", stage.Name(), err)
		}
		if result.Fatal || result.Status == app.StageStatusFailure {
			return app.RunResult{ExitCode: r.exitPolicy.RuntimeFailure, StageResults: results}, fmt.Errorf("stage %s reported a fatal failure", result.StageName)
		}
	}

	_ = r.progress.Publish(progress.Event{Kind: progress.EventKindStatus, Message: "run complete", CurrentPath: request.ArchiveRoot})

	return app.RunResult{ExitCode: r.exitPolicy.Success, StageResults: results}, nil
}
