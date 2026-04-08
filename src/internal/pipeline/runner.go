package pipeline

import (
	"context"
	"fmt"
	"log"

	"github.com/attila/focalytics/internal/aggregate"
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
		_ = r.progress.Publish(progress.Event{Kind: progress.EventKindStageStart, Stage: stage.Name(), Message: "stage started", CurrentPath: request.ArchiveRoot})
		result, err := stage.Run(ctx, runContext)
		results = append(results, result)
		if err != nil {
			_ = r.progress.Publish(progress.Event{Kind: progress.EventKindWarning, Stage: stage.Name(), Message: err.Error(), CurrentPath: request.ArchiveRoot, Warnings: 1})
			return app.RunResult{ExitCode: r.exitPolicy.RuntimeFailure, StageResults: results}, fmt.Errorf("stage %s failed: %w", stage.Name(), err)
		}
		if result.Fatal || result.Status == app.StageStatusFailure {
			return app.RunResult{ExitCode: r.exitPolicy.RuntimeFailure, StageResults: results}, fmt.Errorf("stage %s reported a fatal failure", result.StageName)
		}
		_ = r.progress.Publish(progress.Event{Kind: progress.EventKindStageEnd, Stage: stage.Name(), Message: "stage complete", CurrentPath: request.ArchiveRoot})
	}

	completionNote := ""
	if artifact, ok := runContext.Artifact(app.ArtifactAggregateResult); ok {
		if aggregateResult, ok := artifact.(aggregate.Result); ok {
			completionNote = aggregateResult.Filter.CompletionNote()
			if completionNote != "" {
				_ = r.progress.Publish(progress.Event{Kind: progress.EventKindStatus, Message: completionNote, CurrentPath: request.ArchiveRoot})
			}
		}
	}

	_ = r.progress.Publish(progress.Event{Kind: progress.EventKindStatus, Message: "run complete", CurrentPath: request.ArchiveRoot})

	return app.RunResult{ExitCode: r.exitPolicy.Success, StageResults: results, CompletionNote: completionNote}, nil
}
