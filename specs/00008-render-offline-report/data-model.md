# Data Model: Render Offline Report

## Entities

### ReportModel
- **Fields**: `title`, `generated_at`, `archive_overview`, `timeline`, `gear`, `technical`, `exclusion_notes`, `report_path`
- **Role**: render-ready view derived from aggregate summaries for one HTML document
- **Relationships**: owns section-specific view models and exclusion note collections

### OverviewCard
- **Fields**: `label`, `value`, `supporting_text`
- **Role**: compact hero metric used in the overview section

### ChartBlock
- **Fields**: `title`, `subtitle`, `svg_markup`, `table_rows`, `exclusion_note`
- **Role**: reusable report block for timeline, gear, or technical summaries

### ExclusionNote
- **Fields**: `section_key`, `text`, `details`
- **Role**: visible explanation of omitted data for one report section

### RenderResult
- **Fields**: `path`, `generated_at`
- **Role**: final renderer output recorded after the HTML file is written

## Validation Rules

- ReportModel must be renderable from aggregate-only data without per-photo payloads.
- Every section affected by exclusions must have either an explicit note or a verified zero-exclusion state.
- Report artifact paths must end in `.html` and resolve inside the working directory for default output.
- Empty-state sections must remain structurally valid and readable when aggregate buckets are absent.

## State Notes

- All models are transient except the final HTML artifact written to disk.
- Timestamp handling for naming and display must be injectable in tests to keep golden output deterministic.