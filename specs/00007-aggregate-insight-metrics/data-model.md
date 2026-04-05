# Data Model: Aggregate Insight Metrics

## Entities

### ArchiveSummary
- **Fields**: `date_span`, `totals`, `timeline`, `gear`, `technical`, `exclusions`, `warnings_total`
- **Role**: top-level in-memory artifact passed from aggregation to rendering
- **Relationships**: owns `TimelineBucket`, `RankedBucket`, and `ExclusionSummary` collections

### DateSpan
- **Fields**: `first_captured_at`, `last_captured_at`
- **Role**: archive-level capture range derived from aggregated facts

### TimelineBucket
- **Fields**: `key`, `label`, `count`
- **Role**: deterministic count for a canonical year or day key

### RankedBucket
- **Fields**: `key`, `label`, `count`
- **Role**: deterministic ranked usage bucket for camera, lens, or technical metrics

### ExclusionSummary
- **Fields**: `metric`, `reason`, `count`
- **Role**: grouped omission count for a specific metric family and exclusion reason

## Validation Rules

- Bucket keys must be canonical and sortable without locale-specific formatting.
- Counts must be non-negative integers.
- Empty archives must still yield a valid `ArchiveSummary` with zero-value sections.
- Aggregate models must not retain file-level payloads beyond derived counts and archive span.

## State Notes

- All entities are transient and exist only for a single run.
- Ordering is applied before artifact publication, not inferred from map iteration.