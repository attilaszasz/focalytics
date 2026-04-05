package metadata

import "time"

type Metric string

const (
	MetricCapturedAt            Metric = "captured_at"
	MetricCameraModel           Metric = "camera_model"
	MetricLensModel             Metric = "lens_model"
	MetricFocalLengthMM         Metric = "focal_length_mm"
	MetricNormalizedFocalLength Metric = "normalized_focal_length_mm"
	MetricApertureF             Metric = "aperture_f"
	MetricShutterSeconds        Metric = "shutter_seconds"
	MetricISO                   Metric = "iso"
)

type ProvenanceSource string

const (
	ProvenanceEmbedded                 ProvenanceSource = "embedded"
	ProvenanceSidecar                  ProvenanceSource = "sidecar"
	ProvenanceFileTimestamp            ProvenanceSource = "file_timestamp"
	ProvenanceDirectoryHint            ProvenanceSource = "directory_hint"
	ProvenanceDerivedCropFactor        ProvenanceSource = "derived_crop_factor"
	ProvenanceDerivedActualFocalLength ProvenanceSource = "derived_actual_focal_length"
)

type Exclusion struct {
	Metric Metric
	Reason string
}

type Fact struct {
	Path                    string
	RelativePath            string
	SidecarPath             string
	CapturedAt              *time.Time
	CameraModel             string
	LensModel               string
	FocalLengthMM           *float64
	NormalizedFocalLengthMM *float64
	ApertureF               *float64
	ShutterSeconds          *float64
	ISO                     *int
	Provenance              map[Metric]ProvenanceSource
	Exclusions              []Exclusion
}

type Warning struct {
	Path    string
	Message string
}

type Result struct {
	Facts    []Fact
	Warnings []Warning
}
