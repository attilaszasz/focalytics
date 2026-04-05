package metadata

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/attila/focalytics/internal/app"
	"github.com/attila/focalytics/internal/discovery"
	"github.com/attila/focalytics/internal/progress"
)

type recordingSink struct {
	events []progress.Event
}

func (r *recordingSink) Publish(event progress.Event) error {
	r.events = append(r.events, event)
	return nil
}

func TestRecoverEmbeddedMetadataFromGalleryFixture(t *testing.T) {
	imagePath := filepath.Join("..", "..", "..", "gallery", "2009_10_06_iCory", "IMG_7554.jpg")
	if _, err := os.Stat(imagePath); err != nil {
		t.Fatalf("gallery fixture missing: %v", err)
	}

	service := NewService()
	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "2009_10_06_iCory/IMG_7554.jpg"}}}, &recordingSink{})
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}
	if len(result.Facts) != 1 {
		t.Fatalf("expected one fact, got %d", len(result.Facts))
	}
	fact := result.Facts[0]
	if fact.CapturedAt == nil {
		t.Fatal("expected embedded capture time")
	}
	if fact.Provenance[MetricCapturedAt] != ProvenanceEmbedded {
		t.Fatalf("expected embedded date provenance, got %s", fact.Provenance[MetricCapturedAt])
	}
	if fact.CameraModel == "" {
		t.Fatal("expected embedded camera model")
	}
	if fact.Provenance[MetricCameraModel] != ProvenanceEmbedded {
		t.Fatalf("expected embedded camera provenance, got %s", fact.Provenance[MetricCameraModel])
	}
	if fact.FocalLengthMM == nil {
		t.Fatal("expected embedded focal length")
	}
}

func TestRecoverUsesSidecarFallbacks(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "2022_04_30_Flowers", "image.jpg")
	sidecarPath := filepath.Join(root, "2022_04_30_Flowers", "image.xmp")
	writeFixtureFile(t, imagePath, []byte("not-a-real-jpeg"))
	writeFixtureFile(t, sidecarPath, []byte(`<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" xmlns:xmp="http://ns.adobe.com/xap/1.0/" exif:FocalLength="85" exif:FocalLengthIn35mmFormat="127" exif:ExposureTime="1/125" exif:FNumber="2.8" exif:ISOSpeedRatings="200" tiff:Model="Test Camera" aux:Lens="Test Lens" xmp:CreateDate="2022-04-30T10:11:12" xmlns:aux="http://ns.adobe.com/exif/1.0/aux/"/></rdf:RDF></x:xmpmeta>`))

	service := NewService()
	sink := &recordingSink{}
	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{
		{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "2022_04_30_Flowers/image.jpg"},
		{Kind: discovery.CandidateKindSidecar, Path: sidecarPath, RelativePath: "2022_04_30_Flowers/image.xmp"},
	}}, sink)
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}
	fact := result.Facts[0]
	if fact.CapturedAt == nil || fact.Provenance[MetricCapturedAt] != ProvenanceSidecar {
		t.Fatalf("expected sidecar capture time provenance, got %+v", fact.Provenance)
	}
	if fact.CameraModel != "Test Camera" {
		t.Fatalf("unexpected camera model: %s", fact.CameraModel)
	}
	if fact.LensModel != "Test Lens" {
		t.Fatalf("unexpected lens model: %s", fact.LensModel)
	}
	if fact.NormalizedFocalLengthMM == nil || *fact.NormalizedFocalLengthMM != 127 {
		t.Fatalf("unexpected normalized focal length: %+v", fact.NormalizedFocalLengthMM)
	}
	if len(result.Warnings) == 0 {
		t.Fatal("expected embedded metadata warning for synthetic image")
	}
	if len(sink.events) == 0 || sink.events[0].Kind != progress.EventKindWarning {
		t.Fatalf("expected warning event, got %+v", sink.events)
	}
}

func TestRecoverUsesFileTimestampFallbackAndExclusions(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "gallery", "no-metadata.jpg")
	writeFixtureFile(t, imagePath, []byte("not-a-real-jpeg"))
	fallbackTime := time.Date(2020, time.March, 14, 12, 0, 0, 0, time.UTC)
	if err := os.Chtimes(imagePath, fallbackTime, fallbackTime); err != nil {
		t.Fatalf("set file times: %v", err)
	}

	service := NewService()
	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "gallery/no-metadata.jpg"}}}, &recordingSink{})
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}
	fact := result.Facts[0]
	if fact.CapturedAt == nil || !fact.CapturedAt.Equal(fallbackTime) {
		t.Fatalf("expected file timestamp fallback, got %+v", fact.CapturedAt)
	}
	if fact.Provenance[MetricCapturedAt] != ProvenanceFileTimestamp {
		t.Fatalf("expected file timestamp provenance, got %s", fact.Provenance[MetricCapturedAt])
	}
	if !hasExclusion(fact, MetricCameraModel) || !hasExclusion(fact, MetricISO) {
		t.Fatalf("expected exclusions for missing metrics, got %+v", fact.Exclusions)
	}
}

func TestRecoverUsesDirectoryHintWhenStatUnavailable(t *testing.T) {
	service := NewService()
	service.stat = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{{Kind: discovery.CandidateKindImage, Path: "/archive/2018_trip/image.jpg", RelativePath: "2018_trip/image.jpg"}}}, &recordingSink{})
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}
	fact := result.Facts[0]
	if fact.CapturedAt == nil || fact.CapturedAt.Year() != 2018 {
		t.Fatalf("expected directory year fallback, got %+v", fact.CapturedAt)
	}
	if fact.Provenance[MetricCapturedAt] != ProvenanceDirectoryHint {
		t.Fatalf("expected directory hint provenance, got %s", fact.Provenance[MetricCapturedAt])
	}
}

func TestMetadataStageStoresArtifacts(t *testing.T) {
	imagePath := filepath.Join("..", "..", "..", "gallery", "2009_10_06_iCory", "IMG_7554.jpg")
	if _, err := os.Stat(imagePath); err != nil {
		t.Fatalf("gallery fixture missing: %v", err)
	}

	runContext := appContextForTest(discovery.Result{Candidates: []discovery.Candidate{{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "2009_10_06_iCory/IMG_7554.jpg"}}})
	stage := NewStage(NewService())
	result, err := stage.Run(context.Background(), runContext)
	if err != nil {
		t.Fatalf("expected metadata stage success: %v", err)
	}
	if result.Status != app.StageStatusSuccess {
		t.Fatalf("unexpected stage status: %+v", result)
	}
	if _, ok := runContext.Artifact(app.ArtifactMetadataResult); !ok {
		t.Fatal("expected metadata artifact to be stored")
	}
}

func hasExclusion(fact Fact, metric Metric) bool {
	for _, exclusion := range fact.Exclusions {
		if exclusion.Metric == metric {
			return true
		}
	}
	return false
}

func writeFixtureFile(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create dir: %v", err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func appContextForTest(discoveryResult discovery.Result) app.RunContext {
	runContext := app.NewRunContext(app.ScanRequest{}, app.DefaultExitPolicy(), &recordingSink{}, nil)
	runContext.SetArtifact(app.ArtifactDiscoveryResult, discoveryResult)
	return runContext
}
