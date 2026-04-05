package aggregate

import (
	"sort"
	"strings"
	"time"

	"github.com/attila/focalytics/internal/metadata"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) Aggregate(metadataResult metadata.Result) Result {
	result := Result{}
	yearCounts := map[string]int{}
	dayCounts := map[string]int{}
	cameraCounts := map[string]int{}
	lensCounts := map[string]int{}
	focalBuckets := map[string]RankedBucket{}
	apertureBuckets := map[string]RankedBucket{}
	shutterBuckets := map[string]RankedBucket{}
	isoBuckets := map[string]RankedBucket{}
	exclusionBuckets := map[string]ExclusionSummary{}

	result.WarningsTotal = len(metadataResult.Warnings)

	for _, fact := range metadataResult.Facts {
		result.Totals.Facts++

		if fact.CapturedAt != nil {
			result.Totals.FactsWithCapturedAt++
			s.updateDateSpan(&result.DateSpan, *fact.CapturedAt)
			yearCounts[yearBucketKey(*fact.CapturedAt)]++
			dayCounts[dayBucketKey(*fact.CapturedAt)]++
		}

		if label := strings.TrimSpace(fact.CameraModel); label != "" {
			cameraCounts[label]++
		}
		if label := strings.TrimSpace(fact.LensModel); label != "" {
			lensCounts[label]++
		}
		if fact.NormalizedFocalLengthMM != nil {
			key, label := focalLengthBucket(*fact.NormalizedFocalLengthMM)
			incrementBucket(focalBuckets, key, label)
		}
		if fact.ApertureF != nil {
			key, label := apertureBucket(*fact.ApertureF)
			incrementBucket(apertureBuckets, key, label)
		}
		if fact.ShutterSeconds != nil {
			key, label := shutterSpeedBucket(*fact.ShutterSeconds)
			incrementBucket(shutterBuckets, key, label)
		}
		if fact.ISO != nil {
			key, label := isoBucket(*fact.ISO)
			incrementBucket(isoBuckets, key, label)
		}

		for _, exclusion := range fact.Exclusions {
			result.Totals.ExcludedMetrics++
			bucketKey := string(exclusion.Metric) + "|" + exclusion.Reason
			bucket := exclusionBuckets[bucketKey]
			bucket.Metric = exclusion.Metric
			bucket.Reason = exclusion.Reason
			bucket.Count++
			exclusionBuckets[bucketKey] = bucket
		}
	}

	result.Timeline.Years = timelineBuckets(yearCounts)
	result.Timeline.Days = timelineBuckets(dayCounts)
	result.Gear.Cameras = rankedCounts(cameraCounts)
	result.Gear.Lenses = rankedCounts(lensCounts)
	result.Technical.FocalLengths = rankedBucketValues(focalBuckets)
	result.Technical.Apertures = rankedBucketValues(apertureBuckets)
	result.Technical.ShutterSpeeds = rankedBucketValues(shutterBuckets)
	result.Technical.ISOs = rankedBucketValues(isoBuckets)
	result.Exclusions = exclusionValues(exclusionBuckets)

	return result
}

func (s Service) updateDateSpan(span *DateSpan, value time.Time) {
	utc := value.UTC()
	if span.FirstCapturedAt == nil || utc.Before(*span.FirstCapturedAt) {
		copied := utc
		span.FirstCapturedAt = &copied
	}
	if span.LastCapturedAt == nil || utc.After(*span.LastCapturedAt) {
		copied := utc
		span.LastCapturedAt = &copied
	}
}

func incrementBucket(buckets map[string]RankedBucket, key, label string) {
	bucket := buckets[key]
	bucket.Key = key
	bucket.Label = label
	bucket.Count++
	buckets[key] = bucket
}

func timelineBuckets(counts map[string]int) []TimelineBucket {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	buckets := make([]TimelineBucket, 0, len(keys))
	for _, key := range keys {
		buckets = append(buckets, TimelineBucket{Key: key, Label: key, Count: counts[key]})
	}
	return buckets
}

func rankedCounts(counts map[string]int) []RankedBucket {
	buckets := make([]RankedBucket, 0, len(counts))
	for key, count := range counts {
		buckets = append(buckets, RankedBucket{Key: key, Label: key, Count: count})
	}
	sort.Slice(buckets, func(left, right int) bool {
		if buckets[left].Count != buckets[right].Count {
			return buckets[left].Count > buckets[right].Count
		}
		if buckets[left].Label != buckets[right].Label {
			return buckets[left].Label < buckets[right].Label
		}
		return buckets[left].Key < buckets[right].Key
	})
	return buckets
}

func rankedBucketValues(values map[string]RankedBucket) []RankedBucket {
	buckets := make([]RankedBucket, 0, len(values))
	for _, bucket := range values {
		buckets = append(buckets, bucket)
	}
	sort.Slice(buckets, func(left, right int) bool {
		return buckets[left].Key < buckets[right].Key
	})
	return buckets
}

func exclusionValues(values map[string]ExclusionSummary) []ExclusionSummary {
	buckets := make([]ExclusionSummary, 0, len(values))
	for _, bucket := range values {
		buckets = append(buckets, bucket)
	}
	sort.Slice(buckets, func(left, right int) bool {
		if buckets[left].Metric != buckets[right].Metric {
			return buckets[left].Metric < buckets[right].Metric
		}
		return buckets[left].Reason < buckets[right].Reason
	})
	return buckets
}
