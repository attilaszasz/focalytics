# Data Model: Ignore Phone Photos

## Entities

### AnalysisFilter
- **Fields**: `ignore_phone_photos`
- **Role**: captures the user-selected runtime option that narrows affected analytics
- **Relationships**: attached to `ScanRequest`; controls aggregate and report filtering behavior

### DeviceClassification
- **Fields**: `kind`, `evidence_source`, `camera_model`
- **Role**: records whether a fact is classified as `phone`, `non_phone`, or `unknown`
- **Relationships**: attached to each metadata fact and consumed by aggregation

### FilteredScopeSummary
- **Fields**: `filter_active`, `filtered_photos`, `affected_sections`, `full_archive_totals_preserved`
- **Role**: carries the user-visible summary of how filtering changed the report scope
- **Relationships**: derived during aggregation and consumed by rendering and terminal completion feedback

### FilteredMetricSection
- **Fields**: `section_key`, `eligible_count`, `empty_state`, `scope_note`
- **Role**: describes one gear or technical section after the phone filter is applied
- **Relationships**: nested inside the render model and tied to `FilteredScopeSummary`

## Validation Rules

- `AnalysisFilter.ignore_phone_photos` defaults to `false` and must not change behavior when unset.
- `DeviceClassification.kind` is `unknown` whenever trusted device identity evidence is missing or conflicting.
- `FilteredScopeSummary.filtered_photos` counts only confidently classified phone-made photos removed from affected analytics.
- `FilteredMetricSection.empty_state` must remain renderable when no eligible non-phone photos remain.
- Timeline, archive totals, and date span remain derived from the full archive in this increment.

## State Notes

- All entities are transient runtime models; no persistent storage is introduced.
- Classification and filtered counts must be computed once and propagated, not re-derived independently in rendering.