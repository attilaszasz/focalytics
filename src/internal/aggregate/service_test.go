package aggregate

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/attila/focalytics/internal/app"
	"github.com/attila/focalytics/internal/metadata"
)

func TestAggregateTimelineSummariesDeterministic(t *testing.T) {
	service := NewService()
	metadataResult := metadata.Result{Facts: []metadata.Fact{
		{RelativePath: "a.jpg", CapturedAt: timePointer(time.Date(2020, time.January, 2, 10, 0, 0, 0, time.UTC))},
		{RelativePath: "b.jpg", CapturedAt: timePointer(time.Date(2021, time.May, 3, 12, 0, 0, 0, time.UTC))},
		{RelativePath: "c.jpg", CapturedAt: timePointer(time.Date(2020, time.January, 2, 14, 0, 0, 0, time.UTC))},
	}}

	first := service.Aggregate(metadataResult)
	second := service.Aggregate(metadataResult)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected deterministic aggregate output")
	}
	if first.DateSpan.FirstCapturedAt == nil || first.DateSpan.FirstCapturedAt.Year() != 2020 {
		t.Fatalf("unexpected first date span: %+v", first.DateSpan)
	}
	if first.DateSpan.LastCapturedAt == nil || first.DateSpan.LastCapturedAt.Year() != 2021 {
		t.Fatalf("unexpected last date span: %+v", first.DateSpan)
	}
	if len(first.Timeline.Years) != 2 || first.Timeline.Years[0].Key != "2020" || first.Timeline.Years[0].Count != 2 {
		t.Fatalf("unexpected year buckets: %+v", first.Timeline.Years)
	}
	if len(first.Timeline.Days) != 2 || first.Timeline.Days[0].Key != "2020-01-02" || first.Timeline.Days[0].Count != 2 {
		t.Fatalf("unexpected day buckets: %+v", first.Timeline.Days)
	}
}

func TestAggregateGearAndTechnicalSummaries(t *testing.T) {
	service := NewService()
	metadataResult := metadata.Result{Facts: []metadata.Fact{
		{
			CameraModel:             "Camera A",
			LensModel:               "Lens B",
			NormalizedFocalLengthMM: floatPointer(50),
			ApertureF:               floatPointer(2.8),
			ShutterSeconds:          floatPointer(1.0 / 125.0),
			ISO:                     intPointer(200),
		},
		{
			CameraModel:             "Camera B",
			LensModel:               "Lens A",
			NormalizedFocalLengthMM: floatPointer(85),
			ApertureF:               floatPointer(4),
			ShutterSeconds:          floatPointer(1.0 / 60.0),
			ISO:                     intPointer(100),
		},
		{
			CameraModel:             "Camera A",
			LensModel:               "Lens B",
			NormalizedFocalLengthMM: floatPointer(50),
			ApertureF:               floatPointer(2.8),
			ShutterSeconds:          floatPointer(1.0 / 125.0),
			ISO:                     intPointer(200),
		},
	}}

	result := service.Aggregate(metadataResult)
	if len(result.Gear.Cameras) != 2 || result.Gear.Cameras[0].Label != "Camera A" || result.Gear.Cameras[0].Count != 2 {
		t.Fatalf("unexpected camera ranking: %+v", result.Gear.Cameras)
	}
	if len(result.Gear.Lenses) != 2 || result.Gear.Lenses[0].Label != "Lens B" || result.Gear.Lenses[0].Count != 2 {
		t.Fatalf("unexpected lens ranking: %+v", result.Gear.Lenses)
	}
	if len(result.Technical.FocalLengths) != 2 || result.Technical.FocalLengths[0].Label != "50mm" {
		t.Fatalf("unexpected focal length buckets: %+v", result.Technical.FocalLengths)
	}
	if len(result.Technical.FocalLengthLenses) != 2 || result.Technical.FocalLengthLenses[0].BucketKey != "000500" {
		t.Fatalf("unexpected focal length lens summaries: %+v", result.Technical.FocalLengthLenses)
	}
	if len(result.Technical.FocalLengthLenses[0].Lenses) != 1 || result.Technical.FocalLengthLenses[0].Lenses[0].Label != "Lens B" || result.Technical.FocalLengthLenses[0].Lenses[0].Count != 2 {
		t.Fatalf("unexpected focal length lens breakdown: %+v", result.Technical.FocalLengthLenses[0].Lenses)
	}
	if len(result.Technical.Apertures) != 2 || result.Technical.Apertures[0].Label != "f/2.8" {
		t.Fatalf("unexpected aperture buckets: %+v", result.Technical.Apertures)
	}
	if len(result.Technical.ApertureLenses) != 2 || result.Technical.ApertureLenses[0].BucketKey != "000028" {
		t.Fatalf("unexpected aperture lens summaries: %+v", result.Technical.ApertureLenses)
	}
	if len(result.Technical.ApertureLenses[0].Lenses) != 1 || result.Technical.ApertureLenses[0].Lenses[0].Label != "Lens B" || result.Technical.ApertureLenses[0].Lenses[0].Count != 2 {
		t.Fatalf("unexpected aperture lens breakdown: %+v", result.Technical.ApertureLenses[0].Lenses)
	}
	if len(result.Technical.ShutterSpeeds) != 2 || result.Technical.ShutterSpeeds[0].Label != "1/125s" {
		t.Fatalf("unexpected shutter buckets: %+v", result.Technical.ShutterSpeeds)
	}
	if len(result.Technical.ISOs) != 2 || result.Technical.ISOs[0].Label != "ISO 100" {
		t.Fatalf("unexpected ISO buckets: %+v", result.Technical.ISOs)
	}
}

func TestAggregateWarningsAndExclusions(t *testing.T) {
	service := NewService()
	metadataResult := metadata.Result{
		Facts: []metadata.Fact{{
			Exclusions: []metadata.Exclusion{{Metric: metadata.MetricCameraModel, Reason: "camera missing"}, {Metric: metadata.MetricCameraModel, Reason: "camera missing"}, {Metric: metadata.MetricISO, Reason: "iso missing"}},
		}},
		Warnings: []metadata.Warning{{Path: "a.jpg", Message: "warn a"}, {Path: "b.jpg", Message: "warn b"}},
	}

	result := service.Aggregate(metadataResult)
	if result.WarningsTotal != 2 {
		t.Fatalf("expected warnings total 2, got %d", result.WarningsTotal)
	}
	if result.Totals.ExcludedMetrics != 3 {
		t.Fatalf("expected 3 excluded metrics, got %d", result.Totals.ExcludedMetrics)
	}
	if len(result.Exclusions) != 2 {
		t.Fatalf("expected 2 grouped exclusions, got %+v", result.Exclusions)
	}
	if result.Exclusions[0].Metric != metadata.MetricCameraModel || result.Exclusions[0].Count != 2 {
		t.Fatalf("unexpected exclusions: %+v", result.Exclusions)
	}
}

func TestStageStoresAggregateArtifact(t *testing.T) {
	runContext := app.NewRunContext(app.ScanRequest{}, app.DefaultExitPolicy(), nil, nil)
	runContext.SetArtifact(app.ArtifactMetadataResult, metadata.Result{Facts: []metadata.Fact{{CameraModel: "Camera A"}}})
	stage := NewStage(NewService())

	result, err := stage.Run(context.Background(), runContext)
	if err != nil {
		t.Fatalf("expected aggregate stage success: %v", err)
	}
	if result.Status != app.StageStatusSuccess {
		t.Fatalf("unexpected stage result: %+v", result)
	}
	artifact, ok := runContext.Artifact(app.ArtifactAggregateResult)
	if !ok {
		t.Fatal("expected aggregate artifact")
	}
	aggregateResult, ok := artifact.(Result)
	if !ok {
		t.Fatalf("unexpected aggregate artifact type: %T", artifact)
	}
	if aggregateResult.Totals.Facts != 1 {
		t.Fatalf("unexpected aggregate totals: %+v", aggregateResult.Totals)
	}
}

func floatPointer(value float64) *float64 {
	return &value
}

func intPointer(value int) *int {
	return &value
}

func timePointer(value time.Time) *time.Time {
	return &value
}
