package pipeline

import (
	"context"

	"github.com/attila/focalytics/internal/app"
)

type Stage interface {
	Name() string
	Run(ctx context.Context, runContext app.RunContext) (app.StageResult, error)
}
