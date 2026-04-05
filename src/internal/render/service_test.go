package render

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/attila/focalytics/internal/aggregate"
	"github.com/attila/focalytics/internal/app"
	"github.com/attila/focalytics/internal/metadata"
)

func TestGenerateWritesReportAndPrintsPath(t *testing.T) {
	workingDir := t.TempDir()
	fixedTime := time.Date(2026, time.April, 5, 11, 5, 0, 0, time.UTC)
	service := NewService()
	service.now = func() time.Time { return fixedTime }
	service.getwd = func() (string, error) { return workingDir, nil }
	stdout := &bytes.Buffer{}

	result, err := service.Generate(renderFixtureSummary(), "/archives/gallery", stdout)
	if err != nil {
		t.Fatalf("expected render success: %v", err)
	}
	expectedPath := filepath.Join(workingDir, "focalytics_report_20260405_1105.html")
	if result.Path != expectedPath {
		t.Fatalf("unexpected report path: got %q want %q", result.Path, expectedPath)
	}
	if stdout.String() != expectedPath+"\n" {
		t.Fatalf("expected stdout to equal report path, got %q", stdout.String())
	}
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	html := string(content)
	for _, fragment := range []string{"<h2>Timeline</h2>", "Camera bodies", "focalytics report", "Note: 3 missing-data exclusions affected this section."} {
		if !strings.Contains(html, fragment) {
			t.Fatalf("expected report to contain %q", fragment)
		}
	}
	for _, fragment := range []string{"<div class=\"card-value-item\">Camera A</div>", "<div class=\"card-value-item\">Camera B</div>", "<div class=\"card-value-item\">Lens Prime</div>", "<div class=\"card-value-item\">Lens Zoom</div>", "<div class=\"card-value-item\">50mm</div>", "<div class=\"card-value-item\">85mm</div>"} {
		if !strings.Contains(html, fragment) {
			t.Fatalf("expected report to contain %q", fragment)
		}
	}
	if !strings.Contains(html, "data-hint=\"Lenses:\nLens B (2)\"") {
		t.Fatalf("expected report to contain focal/aperture lens tooltip, got %q", html)
	}
	if !strings.Contains(html, "<svg class=\"heatmap\"") {
		t.Fatalf("expected inline SVG heatmap, got %q", html)
	}
	golden, err := os.ReadFile(filepath.Join("testdata", "report.normalized.golden"))
	if err != nil {
		t.Fatalf("read golden report: %v", err)
	}
	if normalized := normalizeHTML(html); normalized != strings.TrimSpace(string(golden)) {
		t.Fatalf("normalized report mismatch\nwant: %s\n got: %s", strings.TrimSpace(string(golden)), normalized)
	}
}

func TestGenerateFailsWhenWriteFails(t *testing.T) {
	service := NewService()
	service.getwd = func() (string, error) { return t.TempDir(), nil }
	service.writeFile = func(string, []byte, os.FileMode) error { return errors.New("disk full") }
	stdout := &bytes.Buffer{}

	_, err := service.Generate(renderFixtureSummary(), "/archives/gallery", stdout)
	if err == nil || !strings.Contains(err.Error(), "disk full") {
		t.Fatalf("expected write failure, got %v", err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout on write failure, got %q", stdout.String())
	}
}

func TestRenderStageStoresArtifact(t *testing.T) {
	workingDir := t.TempDir()
	service := NewService()
	service.now = func() time.Time { return time.Date(2026, time.April, 5, 11, 5, 0, 0, time.UTC) }
	service.getwd = func() (string, error) { return workingDir, nil }
	runContext := app.NewRunContext(app.ScanRequest{ArchiveRoot: "/archives/gallery", Stdout: &bytes.Buffer{}}, app.DefaultExitPolicy(), nil, nil)
	runContext.SetArtifact(app.ArtifactAggregateResult, renderFixtureSummary())
	stage := NewStage(service)

	result, err := stage.Run(context.Background(), runContext)
	if err != nil {
		t.Fatalf("expected render stage success: %v", err)
	}
	if result.Status != app.StageStatusSuccess {
		t.Fatalf("unexpected stage status: %+v", result)
	}
	artifact, ok := runContext.Artifact(app.ArtifactRenderResult)
	if !ok {
		t.Fatal("expected render artifact")
	}
	if _, ok := artifact.(Result); !ok {
		t.Fatalf("unexpected render artifact type: %T", artifact)
	}
}

func TestBuildModelSelectsTopFocalByHighestCount(t *testing.T) {
	summary := renderFixtureSummary()
	summary.Technical.FocalLengths = []aggregate.RankedBucket{
		{Key: "000176", Label: "17.6mm", Count: 4},
		{Key: "000850", Label: "85mm", Count: 127},
		{Key: "000500", Label: "50mm", Count: 110},
		{Key: "000240", Label: "24mm", Count: 95},
	}

	model := buildModel(summary, "/archives/gallery", time.Date(2026, time.April, 5, 11, 31, 0, 0, time.UTC))

	want := []string{"85mm", "50mm", "24mm"}
	if !reflect.DeepEqual(model.Overview.TopFocal, want) {
		t.Fatalf("unexpected top focal labels: got %v want %v", model.Overview.TopFocal, want)
	}
}

func TestBuildModelLimitsTopHeroListsToAvailableLabels(t *testing.T) {
	summary := renderFixtureSummary()
	summary.Gear.Cameras = []aggregate.RankedBucket{{Key: "Camera A", Label: "Camera A", Count: 2}, {Key: "Camera B", Label: "Camera B", Count: 1}}

	model := buildModel(summary, "/archives/gallery", time.Date(2026, time.April, 5, 11, 31, 0, 0, time.UTC))

	want := []string{"Camera A", "Camera B"}
	if !reflect.DeepEqual(model.Overview.TopCamera, want) {
		t.Fatalf("unexpected top camera labels: got %v want %v", model.Overview.TopCamera, want)
	}
}

func renderFixtureSummary() aggregate.Result {
	first := time.Date(2020, time.January, 2, 10, 0, 0, 0, time.UTC)
	last := time.Date(2021, time.May, 3, 12, 0, 0, 0, time.UTC)
	return aggregate.Result{
		Totals:        aggregate.Totals{Facts: 3, FactsWithCapturedAt: 3, ExcludedMetrics: 4},
		DateSpan:      aggregate.DateSpan{FirstCapturedAt: &first, LastCapturedAt: &last},
		WarningsTotal: 2,
		Timeline: aggregate.TimelineSummary{
			Years: []aggregate.TimelineBucket{{Key: "2020", Label: "2020", Count: 2}, {Key: "2021", Label: "2021", Count: 1}},
			Days:  []aggregate.TimelineBucket{{Key: "2020-01-02", Label: "2020-01-02", Count: 2}, {Key: "2021-05-03", Label: "2021-05-03", Count: 1}},
		},
		Gear: aggregate.GearSummary{
			Cameras: []aggregate.RankedBucket{{Key: "Camera A", Label: "Camera A", Count: 2}, {Key: "Camera B", Label: "Camera B", Count: 1}},
			Lenses:  []aggregate.RankedBucket{{Key: "Lens Prime", Label: "Lens Prime", Count: 2}, {Key: "Lens Zoom", Label: "Lens Zoom", Count: 1}},
		},
		Technical: aggregate.TechnicalSummary{
			FocalLengths:      []aggregate.RankedBucket{{Key: "000500", Label: "50mm", Count: 2}, {Key: "000850", Label: "85mm", Count: 1}},
			FocalLengthLenses: []aggregate.BucketLensSummary{{BucketKey: "000500", Lenses: []aggregate.RankedBucket{{Key: "Lens B", Label: "Lens B", Count: 2}}}, {BucketKey: "000850", Lenses: []aggregate.RankedBucket{{Key: "Lens A", Label: "Lens A", Count: 1}}}},
			Apertures:         []aggregate.RankedBucket{{Key: "000028", Label: "f/2.8", Count: 2}, {Key: "000040", Label: "f/4", Count: 1}},
			ApertureLenses:    []aggregate.BucketLensSummary{{BucketKey: "000028", Lenses: []aggregate.RankedBucket{{Key: "Lens B", Label: "Lens B", Count: 2}}}, {BucketKey: "000040", Lenses: []aggregate.RankedBucket{{Key: "Lens A", Label: "Lens A", Count: 1}}}},
			ShutterSpeeds:     []aggregate.RankedBucket{{Key: "000000008000", Label: "1/125s", Count: 2}, {Key: "000000016667", Label: "1/60s", Count: 1}},
			ISOs:              []aggregate.RankedBucket{{Key: "000100", Label: "ISO 100", Count: 1}, {Key: "000200", Label: "ISO 200", Count: 2}},
		},
		Exclusions: []aggregate.ExclusionSummary{{Metric: metadata.MetricCapturedAt, Reason: "capture time unavailable after embedded, sidecar, and fallback recovery", Count: 3}, {Metric: metadata.MetricCameraModel, Reason: "camera model unavailable", Count: 1}},
	}
}

func normalizeHTML(content string) string {
	stylePattern := regexp.MustCompile(`(?s)<style>.*?</style>`)
	svgPattern := regexp.MustCompile(`(?s)<svg class="heatmap".*?</svg>`)
	spacePattern := regexp.MustCompile(`\s+`)
	normalized := stylePattern.ReplaceAllString(content, `<style>[embedded-css]</style>`)
	normalized = svgPattern.ReplaceAllString(normalized, `<svg class="heatmap">[heatmap]</svg>`)
	normalized = spacePattern.ReplaceAllString(normalized, " ")
	return strings.TrimSpace(normalized)
}
