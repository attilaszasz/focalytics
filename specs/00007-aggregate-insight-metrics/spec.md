---
feature_branch: "00007-aggregate-insight-metrics"
created: "2026-04-05"
input: "E004 Aggregate Insight Metrics"
spec_type: "product"
spec_maturity: "draft"
epic_id: "E004"
epic_sources: "{PRD:CAP-004,CAP-005,CAP-006}{SAD:ADR-002,ADR-003}"
---

# Feature Specification: Aggregate Insight Metrics

**Feature Branch**: `00007-aggregate-insight-metrics`  
**Created**: 2026-04-05  
**Status**: Draft  
**Spec Type**: product  
**Spec Maturity**: draft  
**Epic ID**: E004  
**Epic Sources**: {PRD:CAP-004,CAP-005,CAP-006}{SAD:ADR-002,ADR-003}  
**Product Document**: specs/prd.md

## Problem Statement

focalytics can now recover per-photo metadata facts, but it still lacks the in-memory summaries that turn those facts into report-ready insight. Without deterministic archive-level aggregation for timeline, gear, technical settings, and data-quality visibility, the report layer would either need to reprocess every photo or risk presenting inconsistent counts and misleading omissions.

## Scope

### Included

- Aggregate normalized metadata facts into archive-level timeline summaries.
- Build deterministic camera, lens, focal-length, aperture, shutter-speed, and ISO summaries from recovered facts.
- Preserve exclusion and warning counts alongside every affected metric family.
- Produce an aggregate artifact that later rendering can consume without re-reading per-photo metadata.

### Excluded

- HTML template rendering and chart presentation details — those belong to E005.
- Additional metadata recovery or parser changes — those were delivered in E003.
- Persistent caching or on-disk summary storage — this wave must remain stateless per run.

### Edge Cases & Boundaries

- Empty or metadata-poor archives must still produce structurally valid summaries with explicit zero or exclusion counts.
- Tied bucket counts must resolve in a deterministic order so repeated runs over the same archive produce stable outputs.
- Metrics missing on some facts must exclude only those dimensions, not remove valid data from unrelated summaries.

## User Scenarios & Testing

### User Story 1 - Summarize Shooting Timeline (Priority: P1)

As a photographer, I want focalytics to turn capture dates into archive-level timeline summaries, so I can understand when I was most active without scanning individual files again.

**Why this priority**: Timeline insight is a core report pillar and depends directly on this aggregation layer.

**Independent Test**: Aggregate facts spanning multiple years and days, then confirm the resulting timeline buckets and archive date span are deterministic.

**Acceptance Scenarios**:
1. **Given** metadata facts with capture dates across multiple years, **When** aggregation runs, **Then** the system emits stable year-level counts and the overall archive date span.
2. **Given** metadata facts with repeated capture days, **When** aggregation runs, **Then** the system emits day-level activity buckets that are sorted deterministically for later heatmap rendering.

### User Story 2 - Summarize Gear And Technical Patterns (Priority: P1)

As a photographer, I want focalytics to aggregate camera, lens, focal-length, aperture, shutter-speed, and ISO usage, so the report can show how I actually shoot instead of listing raw metadata.

**Why this priority**: Gear and technical analytics are the main value of the product after timeline insight and cannot be rendered without stable aggregate buckets.

**Independent Test**: Aggregate a mixed metadata fixture and confirm camera, lens, focal-length, aperture, shutter-speed, and ISO summaries return the expected counts and sort order.

**Acceptance Scenarios**:
1. **Given** metadata facts with repeated camera and lens values, **When** aggregation runs, **Then** the system emits ranked gear summaries with deterministic tie-breaking.
2. **Given** metadata facts with normalized focal length, aperture, shutter speed, and ISO values, **When** aggregation runs, **Then** the system emits stable technical buckets suitable for direct report rendering.

### User Story 3 - Preserve Data Quality Context (Priority: P1)

As a photographer, I want focalytics to carry exclusion and warning totals into the aggregate output, so later charts can explain what was omitted and why.

**Why this priority**: Data quality transparency is part of the product promise and must survive aggregation instead of being reduced to raw logs.

**Independent Test**: Aggregate facts and warnings with missing metrics, then confirm the aggregate output records per-metric exclusion totals and warning counts without dropping valid summaries.

**Acceptance Scenarios**:
1. **Given** facts with exclusions for specific metrics, **When** aggregation runs, **Then** the system records per-metric exclusion totals and grouped reasons in the aggregate result.
2. **Given** metadata recovery warnings, **When** aggregation runs, **Then** the system preserves warning totals alongside summary data without failing the run.

## Requirements

### Functional Requirements

- **FR-001**: System MUST aggregate recovered metadata facts into archive-level timeline summaries that include a date span and deterministic year/day bucket counts.
- **FR-002**: System MUST aggregate camera and lens usage into ranked summaries with deterministic tie-breaking.
- **FR-003**: System MUST aggregate normalized focal length, aperture, shutter speed, and ISO values into stable technical buckets suitable for direct report rendering.
- **FR-004**: System MUST preserve warning totals and explicit per-metric exclusion counts in the aggregate output.
- **FR-005**: System MUST store the aggregate result as a shared stage artifact so later rendering can consume it without re-reading file-level facts.
- **FR-006**: System MUST avoid carrying per-photo frontend payloads in the aggregate output so memory use remains bound to aggregate cardinality rather than file count.

### Key Entities

- **Archive Summary**: The top-level aggregate output containing date span, summary sections, warning totals, and exclusion totals.
- **Timeline Bucket**: A deterministic year-level or day-level count derived from normalized capture dates.
- **Gear Summary**: Ranked camera or lens counts derived from metadata facts.
- **Technical Bucket**: A grouped count for focal length, aperture, shutter speed, or ISO.
- **Exclusion Summary**: Per-metric totals and reasons describing which data could not be included in a summary.

## Assumptions & Risks

### Assumptions

- E003 metadata facts remain the canonical input model for aggregation.
- The report layer in E005 will consume aggregate-only structures rather than raw file records.
- Deterministic ordering rules can be validated fully through unit and integration tests without a browser runtime.

### Risks

- **Bucket drift** *(likelihood: medium, impact: high)*: If boundary rules for time or technical metrics are inconsistent, repeated runs will produce unstable report sections.
- **Exclusion dilution** *(likelihood: medium, impact: high)*: If exclusions are merged too aggressively, later charts may understate missing data.
- **Memory creep** *(likelihood: low, impact: medium)*: If aggregate structures retain file-level detail, large archives may violate the stateless aggregate-only design.

## Implementation Signals

- `NEW-MODULE` — Add a dedicated aggregation package under `/src/internal/aggregate` for typed summary construction.
- `NEW-API` — Introduce a pipeline stage and shared run-context artifact for aggregate results.
- `NEW-ENTITY` — Define archive summary, timeline bucket, technical bucket, gear summary, and exclusion summary models.

## Success Criteria

### Measurable Outcomes

- **SC-001** [US1]: Aggregate output provides deterministic timeline summaries and archive date span for the same input facts across repeated runs.
- **SC-002** [US2]: Aggregate output provides stable gear and technical summaries from normalized metadata facts without requiring access to per-photo records.
- **SC-003** [US3]: Aggregate output preserves warning totals and exclusion counts for every affected metric family.

## Glossary

| Term | Definition |
|------|------------|
| Aggregate-only output | Summary data that contains counts and rollups without embedding per-photo records for rendering. |
| Timeline bucket | A count representing archive activity for a canonical time unit such as year or day. |
| Technical bucket | A grouped count for a numeric shooting metric such as focal length, aperture, shutter speed, or ISO. |

## Compliance Check

- **Status**: PASS
- **Policy**: Aligns with stateless-per-run processing, modular `/src` boundaries, and explicit data-quality reporting from project-instructions.md.
- **Notes**: The feature remains offline, aggregate-only, and preserves exclusion transparency for later rendering.