package aggregate

import (
	"fmt"
	"strings"
	"time"

	"github.com/attila/focalytics/internal/metadata"
)

type Result struct {
	Totals        Totals
	DateSpan      DateSpan
	Timeline      TimelineSummary
	Gear          GearSummary
	Technical     TechnicalSummary
	Filter        FilteredScopeSummary
	Exclusions    []ExclusionSummary
	WarningsTotal int
}

type FilteredScopeSummary struct {
	Active           bool
	FilteredPhotos   int
	AffectedSections []string
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

func (s FilteredScopeSummary) OverviewNote() string {
	if !s.Active {
		return ""
	}
	return fmt.Sprintf("Phone filter active for gear and technical insights: excluded %d %s. Timeline and total-photo figures still reflect the full archive.", s.FilteredPhotos, pluralizePhotos(s.FilteredPhotos))
}

func (s FilteredScopeSummary) SectionNote() string {
	if !s.Active {
		return ""
	}
	return fmt.Sprintf("Phone filter active: excluded %d %s from this section.", s.FilteredPhotos, pluralizePhotos(s.FilteredPhotos))
}

func (s FilteredScopeSummary) CompletionNote() string {
	if !s.Active {
		return ""
	}
	return fmt.Sprintf("Phone filter active: excluded %d %s from gear and technical insights; timeline and total photos still reflect the full archive.", s.FilteredPhotos, pluralizePhotos(s.FilteredPhotos))
}

func (s FilteredScopeSummary) EmptyMessage() string {
	if !s.Active {
		return "No aggregated data was available for this section."
	}
	return "No non-phone photos remained for this section after filtering."
}

func pluralizePhotos(count int) string {
	if count == 1 {
		return "phone-made photo"
	}
	return "phone-made photos"
}

func (s FilteredScopeSummary) AffectedSectionList() string {
	if len(s.AffectedSections) == 0 {
		return ""
	}
	return strings.Join(s.AffectedSections, ", ")
}
