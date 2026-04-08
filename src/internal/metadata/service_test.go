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
	if fact.ISO == nil || *fact.ISO != 200 {
		t.Fatalf("unexpected ISO: %+v", fact.ISO)
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
	foundMetric := false
	for _, event := range sink.events {
		if event.Kind == progress.EventKindMetric && event.Stage == "metadata" && event.ProcessedCount == 1 && event.TotalCount == 1 {
			foundMetric = true
		}
	}
	if !foundMetric {
		t.Fatalf("expected metadata progress metric, got %+v", sink.events)
	}
}

func TestRecoverUsesSidecarISOFromSequence(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "gallery", "legacy.crw")
	sidecarPath := filepath.Join(root, "gallery", "legacy.xmp")
	sink := &recordingSink{}
	writeFixtureFile(t, imagePath, []byte("not-a-real-raw"))
	writeFixtureFile(t, sidecarPath, []byte(`<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" xmlns:xmp="http://ns.adobe.com/xap/1.0/"><tiff:Model>Canon PowerShot G6</tiff:Model><exif:FocalLength>28.8</exif:FocalLength><exif:ExposureTime>1/250</exif:ExposureTime><exif:FNumber>3</exif:FNumber><exif:ISOSpeedRatings><rdf:Seq><rdf:li>50</rdf:li></rdf:Seq></exif:ISOSpeedRatings><xmp:CreateDate>2005-06-14T10:03:26+03:00</xmp:CreateDate></rdf:Description></rdf:RDF></x:xmpmeta>`))

	result, err := NewService().Recover(discovery.Result{Candidates: []discovery.Candidate{
		{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "gallery/legacy.crw"},
		{Kind: discovery.CandidateKindSidecar, Path: sidecarPath, RelativePath: "gallery/legacy.xmp"},
	}}, sink)
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}

	fact := result.Facts[0]
	if fact.ISO == nil || *fact.ISO != 50 {
		t.Fatalf("expected ISO from sidecar sequence, got %+v", fact.ISO)
	}
	if fact.Provenance[MetricISO] != ProvenanceSidecar {
		t.Fatalf("expected sidecar ISO provenance, got %s", fact.Provenance[MetricISO])
	}
	if len(result.Warnings) != 0 {
		t.Fatalf("expected CRW embedded warning to be suppressed after sidecar recovery, got %+v", result.Warnings)
	}
	for _, event := range sink.events {
		if event.Kind == progress.EventKindWarning {
			t.Fatalf("expected no warning event after CRW fallback recovery, got %+v", sink.events)
		}
	}
}

func TestRecoverDerivesNormalizedFocalLengthFromCropFactor(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "gallery", "aps-c.jpg")
	sidecarPath := filepath.Join(root, "gallery", "aps-c.xmp")
	writeFixtureFile(t, imagePath, []byte("not-a-real-jpeg"))
	writeFixtureFile(t, sidecarPath, []byte(`<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" exif:FocalLength="85" tiff:Model="Canon EOS DIGITAL REBEL XT"/></rdf:RDF></x:xmpmeta>`))

	service := NewService()
	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{
		{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "gallery/aps-c.jpg"},
		{Kind: discovery.CandidateKindSidecar, Path: sidecarPath, RelativePath: "gallery/aps-c.xmp"},
	}}, &recordingSink{})
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}
	fact := result.Facts[0]
	if fact.NormalizedFocalLengthMM == nil || *fact.NormalizedFocalLengthMM != 136 {
		t.Fatalf("unexpected crop-derived normalized focal length: %+v", fact.NormalizedFocalLengthMM)
	}
	if fact.Provenance[MetricNormalizedFocalLength] != ProvenanceDerivedCropFactor {
		t.Fatalf("expected crop-factor provenance, got %s", fact.Provenance[MetricNormalizedFocalLength])
	}
}

func TestRecoverDerivesNormalizedFocalLengthForFullFrameBody(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "gallery", "full-frame.jpg")
	sidecarPath := filepath.Join(root, "gallery", "full-frame.xmp")
	writeFixtureFile(t, imagePath, []byte("not-a-real-jpeg"))
	writeFixtureFile(t, sidecarPath, []byte(`<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" exif:FocalLength="50" tiff:Model="Canon EOS 5D Mark IV"/></rdf:RDF></x:xmpmeta>`))

	service := NewService()
	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{
		{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "gallery/full-frame.jpg"},
		{Kind: discovery.CandidateKindSidecar, Path: sidecarPath, RelativePath: "gallery/full-frame.xmp"},
	}}, &recordingSink{})
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}
	fact := result.Facts[0]
	if fact.NormalizedFocalLengthMM == nil || *fact.NormalizedFocalLengthMM != 50 {
		t.Fatalf("unexpected full-frame normalized focal length: %+v", fact.NormalizedFocalLengthMM)
	}
	if fact.Provenance[MetricNormalizedFocalLength] != ProvenanceDerivedCropFactor {
		t.Fatalf("expected crop-factor provenance, got %s", fact.Provenance[MetricNormalizedFocalLength])
	}
}

func TestRecoverSkipsActualFocalFallbackForPhoneWithoutEquivalent(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "gallery", "phone.jpg")
	sidecarPath := filepath.Join(root, "gallery", "phone.xmp")
	writeFixtureFile(t, imagePath, []byte("not-a-real-jpeg"))
	writeFixtureFile(t, sidecarPath, []byte(`<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" exif:FocalLength="6.9" tiff:Model="iPhone 15 Pro"/></rdf:RDF></x:xmpmeta>`))

	service := NewService()
	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{
		{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "gallery/phone.jpg"},
		{Kind: discovery.CandidateKindSidecar, Path: sidecarPath, RelativePath: "gallery/phone.xmp"},
	}}, &recordingSink{})
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}
	fact := result.Facts[0]
	if fact.NormalizedFocalLengthMM != nil {
		t.Fatalf("expected no normalized focal fallback for phone, got %+v", fact.NormalizedFocalLengthMM)
	}
	if !hasExclusion(fact, MetricNormalizedFocalLength) {
		t.Fatalf("expected normalized focal exclusion, got %+v", fact.Exclusions)
	}
}

func TestRecoverFallsBackToActualFocalLengthForUnknownCamera(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "gallery", "unknown.jpg")
	sidecarPath := filepath.Join(root, "gallery", "unknown.xmp")
	writeFixtureFile(t, imagePath, []byte("not-a-real-jpeg"))
	writeFixtureFile(t, sidecarPath, []byte(`<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" exif:FocalLength="42" tiff:Model="Unknown Prototype Camera"/></rdf:RDF></x:xmpmeta>`))

	service := NewService()
	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{
		{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "gallery/unknown.jpg"},
		{Kind: discovery.CandidateKindSidecar, Path: sidecarPath, RelativePath: "gallery/unknown.xmp"},
	}}, &recordingSink{})
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}
	fact := result.Facts[0]
	if fact.NormalizedFocalLengthMM == nil || *fact.NormalizedFocalLengthMM != 42 {
		t.Fatalf("expected raw focal fallback for unknown camera, got %+v", fact.NormalizedFocalLengthMM)
	}
	if fact.Provenance[MetricNormalizedFocalLength] != ProvenanceDerivedActualFocalLength {
		t.Fatalf("expected raw-focal provenance, got %s", fact.Provenance[MetricNormalizedFocalLength])
	}
}

func TestRecoverMixedCropFactorsProduceDifferentNormalizedFocalLengths(t *testing.T) {
	root := t.TempDir()
	firstImagePath := filepath.Join(root, "gallery", "aps-c.jpg")
	firstSidecarPath := filepath.Join(root, "gallery", "aps-c.xmp")
	secondImagePath := filepath.Join(root, "gallery", "full-frame.jpg")
	secondSidecarPath := filepath.Join(root, "gallery", "full-frame.xmp")
	writeFixtureFile(t, firstImagePath, []byte("not-a-real-jpeg"))
	writeFixtureFile(t, secondImagePath, []byte("not-a-real-jpeg"))
	writeFixtureFile(t, firstSidecarPath, []byte(`<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" exif:FocalLength="50" tiff:Model="Canon EOS 450D"/></rdf:RDF></x:xmpmeta>`))
	writeFixtureFile(t, secondSidecarPath, []byte(`<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" exif:FocalLength="50" tiff:Model="Canon EOS 5D Mark IV"/></rdf:RDF></x:xmpmeta>`))

	service := NewService()
	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{
		{Kind: discovery.CandidateKindImage, Path: firstImagePath, RelativePath: "gallery/aps-c.jpg"},
		{Kind: discovery.CandidateKindSidecar, Path: firstSidecarPath, RelativePath: "gallery/aps-c.xmp"},
		{Kind: discovery.CandidateKindImage, Path: secondImagePath, RelativePath: "gallery/full-frame.jpg"},
		{Kind: discovery.CandidateKindSidecar, Path: secondSidecarPath, RelativePath: "gallery/full-frame.xmp"},
	}}, &recordingSink{})
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}
	if len(result.Facts) != 2 {
		t.Fatalf("expected two facts, got %d", len(result.Facts))
	}
	if *result.Facts[0].NormalizedFocalLengthMM == *result.Facts[1].NormalizedFocalLengthMM {
		t.Fatalf("expected different normalized focal lengths, got %+v and %+v", result.Facts[0].NormalizedFocalLengthMM, result.Facts[1].NormalizedFocalLengthMM)
	}
}

func TestRecoverCanonicalizesHyphenatedCameraProfiles(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "gallery", "mft.orf")
	sidecarPath := filepath.Join(root, "gallery", "mft.xmp")
	writeFixtureFile(t, imagePath, []byte("not-a-real-raw"))
	writeFixtureFile(t, sidecarPath, []byte(`<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" exif:FocalLength="25" tiff:Model="OM-5"/></rdf:RDF></x:xmpmeta>`))

	result, err := NewService().Recover(discovery.Result{Candidates: []discovery.Candidate{
		{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "gallery/mft.orf"},
		{Kind: discovery.CandidateKindSidecar, Path: sidecarPath, RelativePath: "gallery/mft.xmp"},
	}}, &recordingSink{})
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}

	fact := result.Facts[0]
	if fact.NormalizedFocalLengthMM == nil || *fact.NormalizedFocalLengthMM != 50 {
		t.Fatalf("expected crop-factor normalization for OM-5, got %+v", fact.NormalizedFocalLengthMM)
	}
	if fact.Provenance[MetricNormalizedFocalLength] != ProvenanceDerivedCropFactor {
		t.Fatalf("expected crop-factor provenance, got %s", fact.Provenance[MetricNormalizedFocalLength])
	}
}

func TestRecoverUsesPlatformMetadataAfterEmbeddedFailure(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "gallery", "unsupported.rw2")
	writeFixtureFile(t, imagePath, []byte("not-a-real-raw"))
	capturedAt := time.Date(2022, time.August, 13, 8, 51, 55, 0, time.UTC)
	focalLength := 15.0
	aperture := 5.6
	shutter := 0.0025
	iso := 200

	service := NewService()
	sink := &recordingSink{}
	service.platformMetadata = func(string) embeddedValues {
		return embeddedValues{
			CapturedAt:     &capturedAt,
			CameraModel:    "DC-G100",
			LensModel:      "LEICA DG 12-60/F2.8-4.0",
			FocalLengthMM:  &focalLength,
			ApertureF:      &aperture,
			ShutterSeconds: &shutter,
			ISO:            &iso,
		}
	}

	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "gallery/unsupported.rw2"}}}, sink)
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}

	fact := result.Facts[0]
	if fact.CameraModel != "DC-G100" || fact.LensModel != "LEICA DG 12-60/F2.8-4.0" {
		t.Fatalf("expected platform metadata strings, got %+v", fact)
	}
	if fact.CapturedAt == nil || !fact.CapturedAt.Equal(capturedAt) {
		t.Fatalf("expected platform capture time, got %+v", fact.CapturedAt)
	}
	if fact.ISO == nil || *fact.ISO != 200 {
		t.Fatalf("expected platform ISO, got %+v", fact.ISO)
	}
	if fact.Provenance[MetricCapturedAt] != ProvenancePlatformMetadata {
		t.Fatalf("expected platform capture provenance, got %s", fact.Provenance[MetricCapturedAt])
	}
	if fact.Provenance[MetricISO] != ProvenancePlatformMetadata {
		t.Fatalf("expected platform ISO provenance, got %s", fact.Provenance[MetricISO])
	}
	if fact.NormalizedFocalLengthMM == nil || *fact.NormalizedFocalLengthMM != 30 {
		t.Fatalf("expected derived normalized focal length from platform metadata, got %+v", fact.NormalizedFocalLengthMM)
	}
	if fact.Provenance[MetricNormalizedFocalLength] != ProvenanceDerivedCropFactor {
		t.Fatalf("expected derived crop-factor provenance, got %s", fact.Provenance[MetricNormalizedFocalLength])
	}
	if len(result.Warnings) != 0 {
		t.Fatalf("expected RW2 embedded warning to be suppressed after platform fallback, got %+v", result.Warnings)
	}
	for _, event := range sink.events {
		if event.Kind == progress.EventKindWarning {
			t.Fatalf("expected no warning event after RW2 platform fallback, got %+v", sink.events)
		}
	}
	if sink.events[len(sink.events)-1].Kind != progress.EventKindMetric {
		t.Fatalf("expected metadata metric event, got %+v", sink.events)
	}
}

func TestRecoverKeepsRW2WarningWhenFallbackDoesNotRecoverMetadata(t *testing.T) {
	root := t.TempDir()
	imagePath := filepath.Join(root, "gallery", "unsupported.rw2")
	writeFixtureFile(t, imagePath, []byte("not-a-real-raw"))

	service := NewService()
	service.platformMetadata = func(string) embeddedValues { return embeddedValues{} }
	sink := &recordingSink{}
	result, err := service.Recover(discovery.Result{Candidates: []discovery.Candidate{{Kind: discovery.CandidateKindImage, Path: imagePath, RelativePath: "gallery/unsupported.rw2"}}}, sink)
	if err != nil {
		t.Fatalf("expected metadata recovery success: %v", err)
	}
	if len(result.Warnings) != 1 {
		t.Fatalf("expected unresolved RW2 warning to remain visible, got %+v", result.Warnings)
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
	service.platformMetadata = nil
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
