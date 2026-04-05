package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/attila/focalytics/internal/aggregate"
	"github.com/attila/focalytics/internal/metadata"
)

type Service struct {
	now       func() time.Time
	getwd     func() (string, error)
	writeFile func(string, []byte, os.FileMode) error
}

func NewService() Service {
	return Service{
		now:       time.Now,
		getwd:     os.Getwd,
		writeFile: os.WriteFile,
	}
}

func (s Service) Generate(summary aggregate.Result, archiveRoot string, stdout io.Writer) (Result, error) {
	generatedAt := s.now().UTC()
	path, err := s.reportPath(generatedAt)
	if err != nil {
		return Result{}, err
	}

	model := buildModel(summary, archiveRoot, generatedAt)
	htmlDocument, err := renderDocument(model)
	if err != nil {
		return Result{}, err
	}
	if err := s.writeFile(path, htmlDocument, 0o644); err != nil {
		return Result{}, fmt.Errorf("write report %q: %w", path, err)
	}
	if stdout != nil {
		_, _ = fmt.Fprintf(stdout, "report\t%s\n", path)
	}
	return Result{Path: path, GeneratedAt: generatedAt}, nil
}

func renderDocument(model Model) ([]byte, error) {
	styles, err := embeddedAssets.ReadFile("templates/report.css")
	if err != nil {
		return nil, fmt.Errorf("read embedded css: %w", err)
	}
	tmpl, err := template.New("report.html.tmpl").ParseFS(embeddedAssets, "templates/report.html.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parse report template: %w", err)
	}
	viewData := struct {
		Styles template.CSS
		Report Model
	}{
		Styles: template.CSS(string(styles)),
		Report: model,
	}
	var buffer bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buffer, "report.html.tmpl", viewData); err != nil {
		return nil, fmt.Errorf("execute report template: %w", err)
	}
	return buffer.Bytes(), nil
}

func buildModel(summary aggregate.Result, archiveRoot string, generatedAt time.Time) Model {
	archiveName := filepath.Base(archiveRoot)
	if archiveName == "." || archiveName == string(filepath.Separator) || archiveName == "" {
		archiveName = archiveRoot
	}

	return Model{
		ReportTitle: "focalytics archive report",
		ArchiveName: archiveName,
		GeneratedAt: generatedAt.Format("2006-01-02 15:04 UTC"),
		Overview: OverviewSection{
			TotalPhotos: summary.Totals.Facts,
			DateSpan:    formatDateSpan(summary.DateSpan),
			WarningText: formatWarnings(summary.WarningsTotal),
			TopCamera:   topLabel(summary.Gear.Cameras),
			TopLens:     topLabel(summary.Gear.Lenses),
			TopFocal:    topLabel(summary.Technical.FocalLengths),
		},
		Timeline: TimelineSection{
			YearBars:      barRows(summary.Timeline.Years),
			HeatmapCells:  heatmapCells(summary.Timeline.Days, summary.DateSpan),
			HeatmapWidth:  heatmapWidth(summary.Timeline.Days, summary.DateSpan),
			HeatmapHeight: 112,
			Note:          sectionNote(summary.Exclusions, metadata.MetricCapturedAt),
		},
		Cameras:  metricSection("Camera bodies", "Most-used cameras in the archive", rankedRows(summary.Gear.Cameras), sectionNote(summary.Exclusions, metadata.MetricCameraModel)),
		Lenses:   metricSection("Lenses", "Most-used lenses in the archive", rankedRows(summary.Gear.Lenses), sectionNote(summary.Exclusions, metadata.MetricLensModel)),
		Focal:    metricSection("Focal length", "Normalized focal length usage", rankedRows(summary.Technical.FocalLengths), sectionNote(summary.Exclusions, metadata.MetricNormalizedFocalLength)),
		Aperture: metricSection("Aperture", "How often apertures were used", rankedRows(summary.Technical.Apertures), sectionNote(summary.Exclusions, metadata.MetricApertureF)),
		Shutter:  metricSection("Shutter speed", "Exposure duration distribution", rankedRows(summary.Technical.ShutterSpeeds), sectionNote(summary.Exclusions, metadata.MetricShutterSeconds)),
		ISO:      metricSection("ISO", "Sensitivity distribution", rankedRows(summary.Technical.ISOs), sectionNote(summary.Exclusions, metadata.MetricISO)),
	}
}

func metricSection(title, subtitle string, rows []BarRow, note *SectionNote) MetricSection {
	return MetricSection{Title: title, Subtitle: subtitle, Rows: rows, Note: note}
}

func formatDateSpan(span aggregate.DateSpan) string {
	if span.FirstCapturedAt == nil || span.LastCapturedAt == nil {
		return "No capture dates were recovered for this archive."
	}
	first := span.FirstCapturedAt.Format("2006-01-02")
	last := span.LastCapturedAt.Format("2006-01-02")
	if first == last {
		return fmt.Sprintf("Photos taken on %s.", first)
	}
	return fmt.Sprintf("Photos taken between %s and %s.", first, last)
}

func formatWarnings(count int) string {
	if count == 0 {
		return "No metadata warnings were recorded."
	}
	if count == 1 {
		return "1 metadata warning was recorded during analysis."
	}
	return fmt.Sprintf("%d metadata warnings were recorded during analysis.", count)
}

func topLabel(rows []aggregate.RankedBucket) string {
	if len(rows) == 0 {
		return "Unavailable"
	}
	top := rows[0]
	for _, row := range rows[1:] {
		if row.Count > top.Count {
			top = row
		}
	}
	return top.Label
}

func barRows(rows []aggregate.TimelineBucket) []BarRow {
	converted := make([]BarRow, 0, len(rows))
	maxCount := 0
	for _, row := range rows {
		if row.Count > maxCount {
			maxCount = row.Count
		}
	}
	for _, row := range rows {
		converted = append(converted, BarRow{Label: row.Label, DisplayValue: fmt.Sprintf("%d", row.Count), Count: row.Count, WidthPercent: widthPercent(row.Count, maxCount)})
	}
	return converted
}

func rankedRows(rows []aggregate.RankedBucket) []BarRow {
	converted := make([]BarRow, 0, len(rows))
	maxCount := 0
	for _, row := range rows {
		if row.Count > maxCount {
			maxCount = row.Count
		}
	}
	for _, row := range rows {
		converted = append(converted, BarRow{Label: row.Label, DisplayValue: fmt.Sprintf("%d", row.Count), Count: row.Count, WidthPercent: widthPercent(row.Count, maxCount)})
	}
	return converted
}

func widthPercent(count, maxCount int) float64 {
	if maxCount == 0 {
		return 0
	}
	return (float64(count) / float64(maxCount)) * 100
}

func heatmapCells(days []aggregate.TimelineBucket, span aggregate.DateSpan) []HeatmapCell {
	if span.FirstCapturedAt == nil || span.LastCapturedAt == nil {
		return nil
	}
	counts := map[string]int{}
	maxCount := 0
	for _, bucket := range days {
		counts[bucket.Key] = bucket.Count
		if bucket.Count > maxCount {
			maxCount = bucket.Count
		}
	}
	start := span.FirstCapturedAt.UTC()
	start = start.AddDate(0, 0, -int(start.Weekday()))
	end := span.LastCapturedAt.UTC()
	end = end.AddDate(0, 0, 6-int(end.Weekday()))

	cellWidth := 12
	cellHeight := 12
	gap := 4
	week := 0
	cells := make([]HeatmapCell, 0)
	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		count := counts[current.Format("2006-01-02")]
		cells = append(cells, HeatmapCell{
			X:      week * (cellWidth + gap),
			Y:      int(current.Weekday()) * (cellHeight + gap),
			Width:  cellWidth,
			Height: cellHeight,
			Fill:   heatColor(count, maxCount),
			Label:  fmt.Sprintf("%s: %d photos", current.Format("2006-01-02"), count),
		})
		if current.Weekday() == time.Saturday {
			week++
		}
	}
	return cells
}

func heatmapWidth(days []aggregate.TimelineBucket, span aggregate.DateSpan) int {
	if span.FirstCapturedAt == nil || span.LastCapturedAt == nil {
		return 0
	}
	start := span.FirstCapturedAt.UTC().AddDate(0, 0, -int(span.FirstCapturedAt.UTC().Weekday()))
	end := span.LastCapturedAt.UTC().AddDate(0, 0, 6-int(span.LastCapturedAt.UTC().Weekday()))
	weeks := int(end.Sub(start).Hours()/24)/7 + 1
	return weeks * 16
}

func heatColor(count, maxCount int) string {
	if count == 0 || maxCount == 0 {
		return "#222b31"
	}
	ratio := float64(count) / float64(maxCount)
	switch {
	case ratio >= 0.75:
		return "#f3a64f"
	case ratio >= 0.5:
		return "#d68137"
	case ratio >= 0.25:
		return "#8db7a2"
	default:
		return "#436169"
	}
}

func sectionNote(exclusions []aggregate.ExclusionSummary, metrics ...metadata.Metric) *SectionNote {
	metricSet := map[metadata.Metric]struct{}{}
	for _, metric := range metrics {
		metricSet[metric] = struct{}{}
	}
	total := 0
	details := make([]string, 0)
	for _, exclusion := range exclusions {
		if _, ok := metricSet[exclusion.Metric]; !ok {
			continue
		}
		total += exclusion.Count
		details = append(details, fmt.Sprintf("%s (%d)", exclusion.Reason, exclusion.Count))
	}
	if total == 0 {
		return nil
	}
	sort.Strings(details)
	summary := fmt.Sprintf("Note: %d missing-data exclusions affected this section.", total)
	if total == 1 {
		summary = "Note: 1 missing-data exclusion affected this section."
	}
	return &SectionNote{Summary: summary, Details: details}
}
