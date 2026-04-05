package render

import "time"

type Result struct {
	Path        string
	GeneratedAt time.Time
}

type Model struct {
	ReportTitle string
	ArchiveName string
	GeneratedAt string
	Overview    OverviewSection
	Timeline    TimelineSection
	Cameras     MetricSection
	Lenses      MetricSection
	Focal       MetricSection
	Aperture    MetricSection
	Shutter     MetricSection
	ISO         MetricSection
}

type OverviewSection struct {
	TotalPhotos int
	DateSpan    string
	WarningText string
	TopCamera   string
	TopLens     string
	TopFocal    string
}

type TimelineSection struct {
	YearBars      []BarRow
	HeatmapCells  []HeatmapCell
	HeatmapWidth  int
	HeatmapHeight int
	Note          *SectionNote
}

type MetricSection struct {
	Title    string
	Subtitle string
	Rows     []BarRow
	Note     *SectionNote
}

type BarRow struct {
	Label        string
	DisplayValue string
	Count        int
	WidthPercent float64
}

type HeatmapCell struct {
	X      int
	Y      int
	Width  int
	Height int
	Fill   string
	Label  string
}

type SectionNote struct {
	Summary string
	Details []string
}
