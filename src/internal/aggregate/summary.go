package aggregate

import (
	"time"

	"github.com/attila/focalytics/internal/metadata"
)

type Result struct {
	Totals        Totals
	DateSpan      DateSpan
	Timeline      TimelineSummary
	Gear          GearSummary
	Technical     TechnicalSummary
	Exclusions    []ExclusionSummary
	WarningsTotal int
}

type Totals struct {
	Facts               int
	FactsWithCapturedAt int
	ExcludedMetrics     int
}

type DateSpan struct {
	FirstCapturedAt *time.Time
	LastCapturedAt  *time.Time
}

type TimelineSummary struct {
	Years []TimelineBucket
	Days  []TimelineBucket
}

type TimelineBucket struct {
	Key   string
	Label string
	Count int
}

type GearSummary struct {
	Cameras []RankedBucket
	Lenses  []RankedBucket
}

type TechnicalSummary struct {
	FocalLengths      []RankedBucket
	FocalLengthLenses []BucketLensSummary
	Apertures         []RankedBucket
	ApertureLenses    []BucketLensSummary
	ShutterSpeeds     []RankedBucket
	ISOs              []RankedBucket
}

type RankedBucket struct {
	Key   string
	Label string
	Count int
}

type BucketLensSummary struct {
	BucketKey string
	Lenses    []RankedBucket
}

type ExclusionSummary struct {
	Metric metadata.Metric
	Reason string
	Count  int
}
