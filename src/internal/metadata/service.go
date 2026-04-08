package metadata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"

	"github.com/attila/focalytics/internal/discovery"
	"github.com/attila/focalytics/internal/progress"
)

var (
	fullDateHintPattern = regexp.MustCompile(`(?i)(\d{4})[_-](\d{2})[_-](\d{2})`)
	yearHintPattern     = regexp.MustCompile(`(?i)(?:^|[^0-9])(19\d{2}|20\d{2})(?:[^0-9]|$)`)
)

type Service struct {
	now              func() time.Time
	stat             func(string) (os.FileInfo, error)
	platformMetadata func(string) embeddedValues
}

func NewService() Service {
	return Service{
		now:              time.Now,
		stat:             os.Stat,
		platformMetadata: readPlatformMetadata,
	}
}

func (s Service) Recover(discoveryResult discovery.Result, sink progress.Sink) (Result, error) {
	if sink == nil {
		sink = progress.NoopSink{}
	}

	totalImages := 0
	for _, candidate := range discoveryResult.Candidates {
		if candidate.Kind == discovery.CandidateKindImage {
			totalImages++
		}
	}

	sidecars := make(map[string]discovery.Candidate)
	for _, candidate := range discoveryResult.Candidates {
		if candidate.Kind != discovery.CandidateKindSidecar {
			continue
		}
		sidecars[trimExtension(candidate.RelativePath)] = candidate
	}

	result := Result{Facts: make([]Fact, 0, len(discoveryResult.Candidates))}
	processedImages := 0
	for _, candidate := range discoveryResult.Candidates {
		if candidate.Kind != discovery.CandidateKindImage {
			continue
		}

		sidecar := sidecars[trimExtension(candidate.RelativePath)]
		fact := Fact{
			Path:         candidate.Path,
			RelativePath: candidate.RelativePath,
			Provenance:   map[Metric]ProvenanceSource{},
		}
		if sidecar.Path != "" {
			fact.SidecarPath = sidecar.Path
		}

		embedded, embeddedWarning := s.readEmbeddedMetadata(candidate.Path)
		platformValues := embeddedValues{}
		if embeddedWarning != "" {
			if s.platformMetadata != nil {
				platformValues = s.platformMetadata(candidate.Path)
			}
		}

		xmpValues := xmpValues{}
		if sidecar.Path != "" {
			parsed, err := parseXMP(sidecar.Path)
			if err != nil {
				message := fmt.Sprintf("unable to parse XMP sidecar: %v", err)
				result.Warnings = append(result.Warnings, Warning{Path: sidecar.Path, Message: message})
				publishWarning(sink, sidecar.Path, message)
			} else {
				xmpValues = parsed
			}
		}

		s.applyMetricTime(&fact, MetricCapturedAt, embedded.CapturedAt, ProvenanceEmbedded)
		if fact.CapturedAt == nil {
			s.applyMetricTime(&fact, MetricCapturedAt, xmpValues.DateCreated, ProvenanceSidecar)
		}
		if fact.CapturedAt == nil {
			s.applyMetricTime(&fact, MetricCapturedAt, platformValues.CapturedAt, ProvenancePlatformMetadata)
		}
		if fact.CapturedAt == nil {
			if fileTime := s.fileTimestamp(candidate.Path); fileTime != nil {
				s.applyMetricTime(&fact, MetricCapturedAt, fileTime, ProvenanceFileTimestamp)
			}
		}
		if fact.CapturedAt == nil {
			if hint := directoryDateHint(candidate.Path); hint != nil {
				s.applyMetricTime(&fact, MetricCapturedAt, hint, ProvenanceDirectoryHint)
			}
		}
		if fact.CapturedAt == nil {
			fact.Exclusions = append(fact.Exclusions, Exclusion{Metric: MetricCapturedAt, Reason: "capture time unavailable after embedded, sidecar, and fallback recovery"})
		}

		s.applyMetricString(&fact, MetricCameraModel, embedded.CameraModel, ProvenanceEmbedded)
		if fact.CameraModel == "" {
			s.applyMetricString(&fact, MetricCameraModel, xmpValues.CameraModel, ProvenanceSidecar)
		}
		if fact.CameraModel == "" {
			s.applyMetricString(&fact, MetricCameraModel, platformValues.CameraModel, ProvenancePlatformMetadata)
		}
		if fact.CameraModel == "" {
			fact.Exclusions = append(fact.Exclusions, Exclusion{Metric: MetricCameraModel, Reason: "camera model unavailable"})
		}

		s.applyMetricString(&fact, MetricLensModel, embedded.LensModel, ProvenanceEmbedded)
		if fact.LensModel == "" {
			s.applyMetricString(&fact, MetricLensModel, xmpValues.LensModel, ProvenanceSidecar)
		}
		if fact.LensModel == "" {
			s.applyMetricString(&fact, MetricLensModel, platformValues.LensModel, ProvenancePlatformMetadata)
		}
		if fact.LensModel == "" {
			fact.Exclusions = append(fact.Exclusions, Exclusion{Metric: MetricLensModel, Reason: "lens model unavailable"})
		}

		s.applyMetricFloat(&fact, MetricFocalLengthMM, embedded.FocalLengthMM, ProvenanceEmbedded)
		if fact.FocalLengthMM == nil {
			s.applyMetricFloat(&fact, MetricFocalLengthMM, xmpValues.FocalLengthMM, ProvenanceSidecar)
		}
		if fact.FocalLengthMM == nil {
			s.applyMetricFloat(&fact, MetricFocalLengthMM, platformValues.FocalLengthMM, ProvenancePlatformMetadata)
		}
		if fact.FocalLengthMM == nil {
			fact.Exclusions = append(fact.Exclusions, Exclusion{Metric: MetricFocalLengthMM, Reason: "focal length unavailable"})
		}

		s.applyMetricFloat(&fact, MetricNormalizedFocalLength, embedded.NormalizedFocalLengthMM, ProvenanceEmbedded)
		if fact.NormalizedFocalLengthMM == nil {
			s.applyMetricFloat(&fact, MetricNormalizedFocalLength, xmpValues.NormalizedFocalLengthMM, ProvenanceSidecar)
		}
		if fact.NormalizedFocalLengthMM == nil {
			if value, ok := deriveNormalizedFocalLength(fact.CameraModel, fact.FocalLengthMM); ok {
				s.applyMetricFloat(&fact, MetricNormalizedFocalLength, value, ProvenanceDerivedCropFactor)
			}
		}
		if fact.NormalizedFocalLengthMM == nil && shouldAllowActualFocalFallback(fact.CameraModel) && fact.FocalLengthMM != nil {
			value := *fact.FocalLengthMM
			s.applyMetricFloat(&fact, MetricNormalizedFocalLength, &value, ProvenanceDerivedActualFocalLength)
		}
		if fact.NormalizedFocalLengthMM == nil {
			fact.Exclusions = append(fact.Exclusions, Exclusion{Metric: MetricNormalizedFocalLength, Reason: "normalized focal length unavailable"})
		}

		s.applyMetricFloat(&fact, MetricApertureF, embedded.ApertureF, ProvenanceEmbedded)
		if fact.ApertureF == nil {
			s.applyMetricFloat(&fact, MetricApertureF, xmpValues.ApertureF, ProvenanceSidecar)
		}
		if fact.ApertureF == nil {
			s.applyMetricFloat(&fact, MetricApertureF, platformValues.ApertureF, ProvenancePlatformMetadata)
		}
		if fact.ApertureF == nil {
			fact.Exclusions = append(fact.Exclusions, Exclusion{Metric: MetricApertureF, Reason: "aperture unavailable"})
		}

		s.applyMetricFloat(&fact, MetricShutterSeconds, embedded.ShutterSeconds, ProvenanceEmbedded)
		if fact.ShutterSeconds == nil {
			s.applyMetricFloat(&fact, MetricShutterSeconds, xmpValues.ShutterSeconds, ProvenanceSidecar)
		}
		if fact.ShutterSeconds == nil {
			s.applyMetricFloat(&fact, MetricShutterSeconds, platformValues.ShutterSeconds, ProvenancePlatformMetadata)
		}
		if fact.ShutterSeconds == nil {
			fact.Exclusions = append(fact.Exclusions, Exclusion{Metric: MetricShutterSeconds, Reason: "shutter speed unavailable"})
		}

		s.applyMetricInt(&fact, MetricISO, embedded.ISO, ProvenanceEmbedded)
		if fact.ISO == nil {
			s.applyMetricInt(&fact, MetricISO, xmpValues.ISO, ProvenanceSidecar)
		}
		if fact.ISO == nil {
			s.applyMetricInt(&fact, MetricISO, platformValues.ISO, ProvenancePlatformMetadata)
		}
		if fact.ISO == nil {
			fact.Exclusions = append(fact.Exclusions, Exclusion{Metric: MetricISO, Reason: "ISO unavailable"})
		}

		if embeddedWarning != "" && shouldPublishEmbeddedWarning(candidate.Path, embeddedWarning, fact) {
			result.Warnings = append(result.Warnings, Warning{Path: candidate.Path, Message: embeddedWarning})
			publishWarning(sink, candidate.Path, embeddedWarning)
		}

		result.Facts = append(result.Facts, fact)
		processedImages++
		if shouldPublishMetric(processedImages, totalImages) {
			_ = sink.Publish(progress.Event{Kind: progress.EventKindMetric, Stage: "metadata", Message: "metadata progress", CurrentPath: candidate.Path, ProcessedCount: processedImages, TotalCount: totalImages, Warnings: len(result.Warnings)})
		}
	}

	return result, nil
}

func shouldPublishMetric(processed, total int) bool {
	if total == 0 {
		return false
	}
	return processed == total || processed%25 == 0
}

type embeddedValues struct {
	CapturedAt              *time.Time
	CameraModel             string
	LensModel               string
	FocalLengthMM           *float64
	NormalizedFocalLengthMM *float64
	ApertureF               *float64
	ShutterSeconds          *float64
	ISO                     *int
}

func (s Service) readEmbeddedMetadata(path string) (embeddedValues, string) {
	file, err := os.Open(path)
	if err != nil {
		return embeddedValues{}, fmt.Sprintf("unable to open file for metadata recovery: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	decoded, err := decodeEmbeddedEXIF(file, path)
	if err != nil {
		return embeddedValues{}, fmt.Sprintf("embedded metadata unavailable: %v", err)
	}

	values := embeddedValues{}
	if capturedAt, err := decoded.DateTime(); err == nil {
		capturedAt = capturedAt.UTC()
		values.CapturedAt = &capturedAt
	}
	values.CameraModel = exifString(decoded, exif.Model)
	values.LensModel = firstNonEmpty(exifString(decoded, exif.LensModel), exifString(decoded, exif.LensMake))
	values.FocalLengthMM = exifFloat(decoded, exif.FocalLength)
	values.NormalizedFocalLengthMM = exifFloat(decoded, exif.FocalLengthIn35mmFilm)
	values.ApertureF = exifFloat(decoded, exif.FNumber)
	values.ShutterSeconds = exifFloat(decoded, exif.ExposureTime)
	values.ISO = exifInt(decoded, exif.ISOSpeedRatings)

	return values, ""
}

func decodeEmbeddedEXIF(file *os.File, path string) (*exif.Exif, error) {
	decoded, err := exif.Decode(file)
	if err == nil {
		return decoded, nil
	}

	header, headerErr := readHeader(file, 4)
	if headerErr != nil {
		return nil, err
	}
	if !shouldRetryWithPatchedTIFF(path, header) {
		return nil, err
	}

	patched, patchedErr := readPatchedTIFFBuffer(file)
	if patchedErr != nil {
		return nil, fmt.Errorf("%v; patched TIFF retry unavailable: %w", err, patchedErr)
	}

	decoded, patchedErr = exif.Decode(bytes.NewReader(patched))
	if patchedErr != nil {
		return nil, fmt.Errorf("%v; patched TIFF retry failed: %w", err, patchedErr)
	}

	return decoded, nil
}

func readHeader(file *os.File, size int) ([]byte, error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	header := make([]byte, size)
	n, err := io.ReadFull(file, header)
	if err != nil {
		return nil, err
	}
	return header[:n], nil
}

func shouldRetryWithPatchedTIFF(path string, header []byte) bool {
	if len(header) < 4 {
		return false
	}
	extension := strings.ToLower(filepath.Ext(path))
	var magic uint16
	switch {
	case header[0] == 'I' && header[1] == 'I':
		magic = binary.LittleEndian.Uint16(header[2:4])
	case header[0] == 'M' && header[1] == 'M':
		magic = binary.BigEndian.Uint16(header[2:4])
	default:
		return false
	}

	switch extension {
	case ".orf":
		return magic == 0x4f52 || magic == 0x5352
	case ".rw2":
		return magic == 0x0055
	default:
		return false
	}
}

func readPatchedTIFFBuffer(file *os.File) ([]byte, error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(content) < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	patched := make([]byte, len(content))
	copy(patched, content)
	if patched[0] == 'I' && patched[1] == 'I' {
		binary.LittleEndian.PutUint16(patched[2:4], 0x002a)
		return patched, nil
	}
	if patched[0] == 'M' && patched[1] == 'M' {
		binary.BigEndian.PutUint16(patched[2:4], 0x002a)
		return patched, nil
	}
	return nil, fmt.Errorf("unsupported TIFF byte order")
}

func readPlatformMetadata(path string) embeddedValues {
	if runtime.GOOS != "darwin" {
		return embeddedValues{}
	}

	output, err := exec.Command(
		"mdls",
		"-name", "kMDItemAcquisitionModel",
		"-name", "kMDItemLensModel",
		"-name", "kMDItemFocalLength",
		"-name", "kMDItemFNumber",
		"-name", "kMDItemExposureTimeSeconds",
		"-name", "kMDItemISOSpeed",
		"-name", "kMDItemContentCreationDate",
		path,
	).Output()
	if err != nil {
		return embeddedValues{}
	}

	values := parsePlatformMetadataValues(string(output))
	return embeddedValues{
		CapturedAt:     parsePlatformTime(values["kMDItemContentCreationDate"]),
		CameraModel:    values["kMDItemAcquisitionModel"],
		LensModel:      values["kMDItemLensModel"],
		FocalLengthMM:  parseFloatPointer(values["kMDItemFocalLength"]),
		ApertureF:      parseFloatPointer(values["kMDItemFNumber"]),
		ShutterSeconds: parseFloatPointer(values["kMDItemExposureTimeSeconds"]),
		ISO:            parseIntPointer(values["kMDItemISOSpeed"]),
	}
}

func parsePlatformMetadataValues(output string) map[string]string {
	values := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" || value == "" || value == "(null)" {
			continue
		}
		values[key] = strings.Trim(strings.TrimSpace(value), "\"")
	}
	return values
}

func parsePlatformTime(value string) *time.Time {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02 15:04:05 -0700", trimmed)
	if err != nil {
		return nil
	}
	parsed = parsed.UTC()
	return &parsed
}

func shouldPublishEmbeddedWarning(path, message string, fact Fact) bool {
	if !isSuppressibleEmbeddedWarning(path, message) {
		return true
	}
	return !hasRecoveredMetadataFallback(fact)
}

func isSuppressibleEmbeddedWarning(path, message string) bool {
	if !strings.HasPrefix(message, "embedded metadata unavailable:") {
		return false
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".crw", ".rw2":
		return true
	default:
		return false
	}
}

func hasRecoveredMetadataFallback(fact Fact) bool {
	for _, metric := range []Metric{MetricCameraModel, MetricLensModel, MetricFocalLengthMM, MetricApertureF, MetricShutterSeconds, MetricISO} {
		source, ok := fact.Provenance[metric]
		if !ok || source == ProvenanceEmbedded {
			continue
		}
		return true
	}
	return false
}

func exifString(decoded *exif.Exif, field exif.FieldName) string {
	tag, err := decoded.Get(field)
	if err != nil {
		return ""
	}
	if value, err := tag.StringVal(); err == nil {
		return strings.TrimSpace(value)
	}
	return strings.Trim(strings.TrimSpace(tag.String()), "\"")
}

func exifFloat(decoded *exif.Exif, field exif.FieldName) *float64 {
	tag, err := decoded.Get(field)
	if err != nil {
		return nil
	}
	if rat, err := tag.Rat(0); err == nil {
		value, _ := rat.Float64()
		return &value
	}
	if value, err := tag.Float(0); err == nil {
		return &value
	}
	if integer, err := tag.Int(0); err == nil {
		value := float64(integer)
		return &value
	}
	return parseFloatPointer(strings.Trim(strings.TrimSpace(tag.String()), "\""))
}

func deriveNormalizedFocalLength(cameraModel string, focalLengthMM *float64) (*float64, bool) {
	if focalLengthMM == nil {
		return nil, false
	}
	profile, ok := lookupCameraProfile(cameraModel)
	if !ok || profile.CropFactor == nil {
		return nil, false
	}
	value := math.Round((*focalLengthMM)*(*profile.CropFactor)*10) / 10
	return &value, true
}

func shouldAllowActualFocalFallback(cameraModel string) bool {
	profile, ok := lookupCameraProfile(cameraModel)
	if !ok {
		return true
	}
	return profile.Fallback == focalFallbackAllowActual
}

func exifInt(decoded *exif.Exif, field exif.FieldName) *int {
	tag, err := decoded.Get(field)
	if err != nil {
		return nil
	}
	if value, err := tag.Int(0); err == nil {
		return &value
	}
	return parseIntPointer(strings.Trim(strings.TrimSpace(tag.String()), "\""))
}

func (s Service) applyMetricTime(fact *Fact, metric Metric, value *time.Time, source ProvenanceSource) {
	if value == nil || fact.CapturedAt != nil {
		return
	}
	fact.CapturedAt = value
	fact.Provenance[metric] = source
}

func (s Service) applyMetricString(fact *Fact, metric Metric, value string, source ProvenanceSource) {
	if strings.TrimSpace(value) == "" {
		return
	}
	switch metric {
	case MetricCameraModel:
		if fact.CameraModel != "" {
			return
		}
		fact.CameraModel = strings.TrimSpace(value)
	case MetricLensModel:
		if fact.LensModel != "" {
			return
		}
		fact.LensModel = strings.TrimSpace(value)
	default:
		return
	}
	fact.Provenance[metric] = source
}

func (s Service) applyMetricFloat(fact *Fact, metric Metric, value *float64, source ProvenanceSource) {
	if value == nil {
		return
	}
	switch metric {
	case MetricFocalLengthMM:
		if fact.FocalLengthMM != nil {
			return
		}
		copied := *value
		fact.FocalLengthMM = &copied
	case MetricNormalizedFocalLength:
		if fact.NormalizedFocalLengthMM != nil {
			return
		}
		copied := *value
		fact.NormalizedFocalLengthMM = &copied
	case MetricApertureF:
		if fact.ApertureF != nil {
			return
		}
		copied := *value
		fact.ApertureF = &copied
	case MetricShutterSeconds:
		if fact.ShutterSeconds != nil {
			return
		}
		copied := *value
		fact.ShutterSeconds = &copied
	default:
		return
	}
	fact.Provenance[metric] = source
}

func (s Service) applyMetricInt(fact *Fact, metric Metric, value *int, source ProvenanceSource) {
	if value == nil || metric != MetricISO || fact.ISO != nil {
		return
	}
	copied := *value
	fact.ISO = &copied
	fact.Provenance[metric] = source
}

func (s Service) fileTimestamp(path string) *time.Time {
	info, err := s.stat(path)
	if err != nil {
		return nil
	}
	modTime := info.ModTime().UTC()
	return &modTime
}

func directoryDateHint(path string) *time.Time {
	current := filepath.Dir(path)
	for {
		base := filepath.Base(current)
		if matches := fullDateHintPattern.FindStringSubmatch(base); len(matches) == 4 {
			year, _ := strconv.Atoi(matches[1])
			month, _ := strconv.Atoi(matches[2])
			day, _ := strconv.Atoi(matches[3])
			value := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			return &value
		}
		if matches := yearHintPattern.FindStringSubmatch(base); len(matches) == 2 {
			year, _ := strconv.Atoi(matches[1])
			value := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
			return &value
		}
		next := filepath.Dir(current)
		if next == current {
			break
		}
		current = next
	}
	return nil
}

func parseFloatPointer(value string) *float64 {
	trimmed := strings.TrimSpace(strings.Trim(value, "\""))
	trimmed = strings.TrimSuffix(trimmed, "mm")
	if trimmed == "" {
		return nil
	}
	if strings.Contains(trimmed, "/") {
		parts := strings.SplitN(trimmed, "/", 2)
		if len(parts) == 2 {
			numerator, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			denominator, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil && denominator != 0 {
				value := numerator / denominator
				return &value
			}
		}
	}
	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseIntPointer(value string) *int {
	trimmed := strings.TrimSpace(strings.Trim(value, "\""))
	if trimmed == "" {
		return nil
	}
	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseExposurePointer(value string) *float64 {
	if parsed := parseFloatPointer(value); parsed != nil {
		return parsed
	}
	trimmed := strings.TrimSpace(strings.Trim(value, "\""))
	if trimmed == "" {
		return nil
	}
	rat, ok := new(big.Rat).SetString(trimmed)
	if !ok {
		return nil
	}
	parsed, _ := rat.Float64()
	return &parsed
}

func firstParsedTime(values ...string) *time.Time {
	formats := []string{time.RFC3339, "2006:01:02 15:04:05", "2006-01-02T15:04:05", "2006-01-02 15:04:05"}
	for _, candidate := range values {
		trimmed := strings.TrimSpace(candidate)
		if trimmed == "" {
			continue
		}
		for _, format := range formats {
			parsed, err := time.Parse(format, trimmed)
			if err == nil {
				value := parsed.UTC()
				return &value
			}
		}
	}
	return nil
}

func trimExtension(path string) string {
	extension := filepath.Ext(path)
	return strings.TrimSuffix(path, extension)
}

func publishWarning(sink progress.Sink, path, message string) {
	_ = sink.Publish(progress.Event{Kind: progress.EventKindWarning, Stage: "metadata", Message: message, CurrentPath: path})
}
