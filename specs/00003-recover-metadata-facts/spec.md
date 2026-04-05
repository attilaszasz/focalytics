---
feature_branch: "00003-recover-metadata-facts"
created: "2026-04-05"
input: "E003 Recover Metadata Facts"
spec_type: "product"
spec_maturity: "draft"
epic_id: "E003"
epic_sources: "{PRD:CAP-002,CAP-006}{SAD:ADR-003}"
---

# Feature Specification: Recover Metadata Facts

**Feature Branch**: `00003-recover-metadata-facts`  
**Created**: 2026-04-05  
**Status**: Draft  
**Spec Type**: product  
**Spec Maturity**: draft  
**Epic ID**: E003  
**Epic Sources**: {PRD:CAP-002,CAP-006}{SAD:ADR-003}  
**Product Document**: specs/prd.md

## Problem Statement

Discovery now finds candidate image and sidecar files, but focalytics still cannot turn them into trustworthy photo facts. Without layered metadata recovery that prioritizes embedded tags, supplements gaps with matching XMP sidecars, and falls back only where explicitly allowed, later aggregation work would either misstate archive behavior or drop too many files to be useful.

## Scope

### Included

- Parse embedded EXIF metadata from discovered image candidates.
- Parse matching `.xmp` sidecar files when present and use them to fill missing fields.
- Apply date fallback rules using filesystem timestamps and directory-derived hints when canonical capture time is unavailable.
- Normalize focal-length reporting into a usable metric while preserving provenance.
- Record explicit provenance and exclusion details per metric.
- Continue past corrupt metadata and unsupported per-file parse failures while publishing warnings.

### Excluded

- Archive-level aggregation and chart-ready bucketing.
- Report rendering.
- Metadata write-back or sidecar mutation.

### Edge Cases & Boundaries

- A broken EXIF block must not prevent XMP or fallback recovery for the same file.
- Missing lens, exposure, or ISO data must exclude only those metrics, not the entire file.
- If no XMP sidecars exist in a given archive, embedded and fallback recovery must still proceed.

## User Scenarios & Testing

### User Story 1 - Recover Embedded Metadata (Priority: P1)

As a photographer, I want focalytics to read useful facts directly from image metadata, so my archive can be analyzed without any manual cataloging step.

**Why this priority**: Embedded metadata is the primary source and the baseline dependency for all later insight work.

**Independent Test**: Parse a real gallery image with EXIF metadata and confirm capture date plus selected technical fields are recovered with embedded provenance.

**Acceptance Scenarios**:
1. **Given** a discovered image with embedded EXIF metadata, **When** metadata recovery runs, **Then** the system extracts available capture, camera, lens, focal-length, aperture, shutter, and ISO facts.
2. **Given** embedded metadata is present for a metric, **When** the fact is emitted, **Then** its provenance records the embedded source.

### User Story 2 - Use Sidecars and Fallbacks Honestly (Priority: P1)

As a photographer, I want focalytics to preserve partial value when embedded tags are missing, so old or edited files still contribute to the analysis without pretending their metadata is complete.

**Why this priority**: Real archives are inconsistent, and the product promise depends on graceful degradation.

**Independent Test**: Feed an image with missing embedded metadata plus a matching XMP sidecar and confirm the sidecar fills missing metrics while fallback dates are applied when both embedded and sidecar dates are absent.

**Acceptance Scenarios**:
1. **Given** a matching `.xmp` sidecar supplies missing metadata, **When** recovery runs, **Then** the missing metrics are filled from the sidecar and tagged with sidecar provenance.
2. **Given** both embedded and sidecar capture time are unavailable, **When** recovery runs, **Then** the system falls back first to file timestamps and then to directory-derived date hints while marking those facts as derived.

### User Story 3 - Surface Exclusions and Corrupt Metadata (Priority: P1)

As a photographer, I want focalytics to warn me about corrupt files and explicitly mark missing metrics as excluded, so I can trust what later charts include and omit.

**Why this priority**: Data quality transparency is a core product requirement, not a polish item.

**Independent Test**: Process a corrupt or metadata-empty image and confirm the run continues, warnings are emitted, and excluded metrics are recorded explicitly.

**Acceptance Scenarios**:
1. **Given** an image has corrupt embedded metadata, **When** recovery runs, **Then** the system emits a warning and still attempts sidecar or fallback recovery.
2. **Given** a metric remains unavailable after all recovery layers, **When** the fact is emitted, **Then** the metric is recorded as excluded with a reason rather than hidden silently.

## Requirements

### Functional Requirements

- **FR-001**: System MUST parse embedded EXIF metadata from discovered image candidates when available.
- **FR-002**: System MUST consult a matching `.xmp` sidecar to fill missing metrics when embedded metadata is absent or incomplete.
- **FR-003**: System MUST apply capture-date fallbacks using file timestamps and directory-derived hints when canonical metadata is missing.
- **FR-004**: System MUST normalize focal-length reporting into a reusable metric while preserving the provenance of that value.
- **FR-005**: System MUST record provenance for every recovered metric and explicit exclusions for every unrecovered metric.
- **FR-006**: System MUST emit warnings for corrupt or unsupported metadata without aborting the run.

### Key Entities

- **Metadata Fact**: The per-image normalized record passed to later aggregation stages.
- **Provenance Record**: The source label attached to a recovered metric, such as embedded, sidecar, file timestamp, or directory hint.
- **Exclusion Record**: The reason a metric could not be recovered after all fallback layers.
- **Normalization Input**: Raw metadata values used to derive normalized fields such as display focal length.

## Assumptions & Risks

### Assumptions

- Discovered image and sidecar candidates from E002 remain deterministic and share stable relative paths.
- The gallery fixture contains enough embedded metadata to validate the primary recovery path.
- XMP sidecar coverage can be validated with synthetic fixtures even if the current gallery lacks sidecars.

### Risks

- **Parser fragility** *(likelihood: medium, impact: high)*: Real-world metadata inconsistencies may break naive EXIF or XMP parsing.
- **Provenance drift** *(likelihood: medium, impact: high)*: If fallback layers are not tracked precisely, later charts will overstate data confidence.
- **Over-normalization** *(likelihood: low, impact: medium)*: Aggressive focal-length normalization could hide uncertainty if derived values are not labeled clearly.

## Implementation Signals

- `NEW-MODULE` — Add a dedicated metadata package under `/src/internal/metadata`.
- `NEW-API` — Extend run context with shared stage artifacts so discovery feeds metadata recovery directly.
- `NEW-ENTITY` — Introduce normalized metadata facts, provenance records, and exclusion records for later aggregation.

## Success Criteria

### Measurable Outcomes

- **SC-001** [US1]: Embedded EXIF metadata is recovered into normalized facts for real gallery images with available metadata.
- **SC-002** [US2]: Sidecar and fallback recovery fill partial facts while preserving explicit provenance for each recovered metric.
- **SC-003** [US3]: Corrupt or metadata-empty files emit warnings and per-metric exclusions without aborting the run.

## Glossary

| Term | Definition |
|------|------------|
| Embedded metadata | Metadata stored inside the image file, typically EXIF. |
| Sidecar | A separate metadata file, usually `.xmp`, associated with an image candidate. |
| Exclusion | An explicit declaration that a specific metric could not be recovered. |

## Compliance Check

- **Status**: PASS
- **Policy**: Aligns with local-first safety, modular `/src` boundaries, and explicit data-quality reporting from project-instructions.md.
- **Notes**: The feature is read-only, offline, and degrades per metric rather than per run.