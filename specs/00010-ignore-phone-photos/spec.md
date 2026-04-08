---
feature_branch: "00010-ignore-phone-photos"
created: "2026-04-08"
input: "E008 Ignore Phone Photos"
spec_type: "product"
spec_maturity: "clarified"
epic_id: "E008"
epic_sources: "{PRD:CAP-005,CAP-006}{SAD:ADR-003}"
---

# Feature Specification: Ignore Phone Photos

**Feature Branch**: `00010-ignore-phone-photos`  
**Created**: 2026-04-08  
**Status**: Draft  
**Spec Type**: product  
**Spec Maturity**: clarified  
**Epic ID**: E008  
**Epic Sources**: {PRD:CAP-005,CAP-006}{SAD:ADR-003}  
**Product Document**: specs/prd.md

## Problem Statement

focalytics can already recover camera-model metadata and even treats phone models specially for focal-length fallback, but it does not give photographers a user-facing way to exclude phone-made images from analysis. In mixed archives, phone captures can dominate gear and technical charts, which makes it harder to understand dedicated-camera shooting habits. If filtered analysis is not available, users have to accept skewed insights or manually curate their archive outside the product.

## Scope

### Included

- Add an optional CLI parameter that excludes confidently identified phone-made images from affected analytics.
- Filter gear and technical insights, including hero rankings derived from those sections, when the phone filter is enabled.
- Classify phone-origin captures conservatively from available embedded, sidecar, or platform camera identity metadata.
- Disclose filtered counts and affected sections in an always-visible report scope note, per-section notes, and terminal feedback.
- Preserve current behavior when the filter parameter is not supplied.

### Excluded

- Changing the default one-command behavior to exclude phone photos automatically — this feature is opt-in only.
- Introducing broader device-taxonomy filtering such as tablets, drones, or custom include and exclude lists — that would expand scope beyond the current epic.
- Rewriting timeline, archive-span, or total-photo counts to be filter-specific in this increment — the initial filter targets gear and technical interpretation only.

### Edge Cases & Boundaries

- Files with missing, conflicting, or weak device metadata must remain included rather than being excluded speculatively.
- Edited derivatives or exported copies are filtered only when original device provenance remains unambiguous in trusted metadata.
- A filtered run may leave an affected section with zero eligible photos; the report must remain valid and explain the empty result.
- The final report path must remain scriptable even when filter-specific terminal feedback is shown.
- Discovery, metadata recovery, and warning counts still reflect the full archive scan even when affected analytics are filtered.
- Terminal completion feedback for filtered runs must stay off stdout so shell capture of the report path remains stable.

## User Scenarios & Testing

### User Story 1 - Compare Camera-Only Insights (Priority: P1)

As a photographer, I want to run focalytics with a phone-photo filter so that gear and technical analytics reflect my dedicated-camera work instead of being dominated by phone captures.

**Why this priority**: This is the core value of the feature; without filtered analytics, the user still cannot separate phone-heavy archives from dedicated-camera shooting patterns.

**Independent Test**: Run focalytics against a mixed archive with the phone filter enabled and confirm that affected gear and technical sections exclude confidently identified phone-made images.

**Acceptance Scenarios**:

1. **Given** an archive containing both dedicated-camera and phone-made images, **When** the user enables the phone filter, **Then** gear and technical analytics exclude confidently identified phone-made images while timeline and archive-total reporting still represent the full scanned archive.
2. **Given** the phone filter is enabled, **When** focalytics renders hero rankings derived from gear or technical analytics, **Then** those rankings are based on the filtered set rather than the full archive.

### User Story 2 - Preserve Trustworthy Defaults (Priority: P1)

As a photographer, I want the default run to stay unchanged and ambiguous device classifications to stay included, so I can trust that the new option narrows scope only when I ask for it.

**Why this priority**: Backward-compatible behavior and conservative classification are necessary to keep the CLI trustworthy and avoid silent over-filtering.

**Independent Test**: Run focalytics on the same archive with and without the phone filter and confirm that the default run matches current behavior while uncertain files remain included during filtered runs.

**Acceptance Scenarios**:

1. **Given** the user runs focalytics without the phone filter, **When** the scan completes, **Then** the resulting analytics and output behavior match the existing unfiltered experience.
2. **Given** a file with missing or conflicting device identity metadata, **When** the user enables the phone filter, **Then** the file remains included in affected analytics rather than being excluded by heuristic guesswork.

### User Story 3 - Understand Filtered Scope (Priority: P2)

As a photographer, I want the report and terminal feedback to state how many photos were filtered and which sections changed, so I can interpret the narrowed analysis correctly.

**Why this priority**: Scope disclosure strengthens trust, but the feature already provides its main utility once filtered analytics and conservative defaults are in place.

**Independent Test**: Run a filtered analysis and confirm that the generated report and terminal completion feedback state the active filter, filtered counts, and affected sections.

**Acceptance Scenarios**:

1. **Given** the phone filter is enabled, **When** focalytics finishes the run, **Then** terminal feedback on stderr or the interactive UI states that phone-made images were excluded from affected analytics and reports how many photos were filtered.
2. **Given** the generated HTML report contains filtered analytics, **When** the user opens it, **Then** the report shows an always-visible scope note near the overview plus per-section notes for affected analytics explaining the active filter and filtered counts.

## Requirements

### Functional Requirements

- **FR-001**: System MUST expose a single opt-in run parameter that enables exclusion of confidently identified phone-made images from affected analytics.
- **FR-002**: System MUST preserve current analytics behavior when the phone filter parameter is absent.
- **FR-003**: System MUST classify a file as phone-made only from high-confidence device identity metadata such as trusted make or model fields recovered from embedded, sidecar, or platform sources.
- **FR-004**: System MUST keep files included when device metadata is missing, conflicting, or insufficient to classify confidently as phone-made.
- **FR-005**: System MUST exclude filtered phone-made images from gear and technical analytics, including hero rankings derived from those analytics.
- **FR-006**: System MUST leave archive totals, timeline metrics, and date-span reporting based on the full scanned archive for this increment.
- **FR-007**: System MUST disclose the active phone filter and filtered-photo counts in an always-visible report scope note and in every affected analytics section.
- **FR-008**: System MUST preserve the final report path as the sole stdout success output while surfacing filter-specific completion feedback through stderr or the interactive UI channel.
- **FR-009**: System MUST render valid empty states for affected analytics when filtering removes all eligible photos from a section.
- **FR-010**: System MUST NOT classify files as phone-made from editing software tags, focal-length heuristics, or provenance-only metadata without trusted device identity.

### Key Entities

- **Analysis Filter**: The user-selected run option that narrows which photos contribute to affected analytics.
- **Device Classification**: The trusted determination of whether a photo is phone-made, not phone-made, or unknown based on recovered metadata.
- **Filtered Insight Scope**: The subset of analytics that are recalculated from non-phone photos when the filter is active.
- **Scope Disclosure Note**: Report or terminal text that explains the active filter and the number of filtered photos affecting interpretation.

## Assumptions & Risks

### Assumptions

- Existing camera-model recovery already provides enough metadata coverage for conservative phone classification to be useful on real archives.
- The current report model and terminal completion path can be extended to disclose filtered scope without changing the offline-only product posture.
- Users want phone-photo exclusion as an explicit comparison mode rather than a new default.

### Risks

- **False-positive filtering** *(likelihood: medium, impact: high)*: Overly broad phone classification could exclude legitimate camera photos and reduce trust in the report.
- **Scope mismatch across layers** *(likelihood: medium, impact: high)*: If aggregation, report notes, and terminal summaries disagree on filtered counts, users may distrust the output.
- **Sparse filtered results** *(likelihood: medium, impact: medium)*: Archives dominated by phone captures may leave some sections empty, which could feel like a failure unless the UI explains why.

## Implementation Signals

- `NEW-CONFIG` — Add an opt-in CLI parameter for excluding phone-made images from affected analytics.
- `NEW-ENTITY` — Carry device-classification state and filtered-scope counts through aggregation and rendering inputs.
- `NEW-API` — Extend request, aggregate, and report contracts to express active filters and filtered-scope disclosures.

## Success Criteria

### Measurable Outcomes

- **SC-001** [US1]: Users can generate gear and technical insights from a mixed archive with phone-made images excluded when they enable the filter.
- **SC-002** [US2]: Users running the default command receive the same inclusion behavior as before, and filtered runs keep ambiguous-device files included.
- **SC-003** [US3]: Users opening a filtered report can identify the active filter and the filtered-photo counts for affected analytics without inspecting source files.

## Glossary

| Term | Definition |
|------|------------|
| Phone-made image | A photo whose trusted device identity metadata indicates capture on a smartphone-class device. |
| Unknown device | A photo whose recovered metadata is too weak or conflicting to classify confidently as phone-made. |
| Filtered insight scope | The report sections whose values are recalculated after applying the optional phone-photo filter. |

## Clarifications

### Session 2026-04-08

- Q: Which report sections should the phone filter change? -> A: Gear and technical analytics, plus hero rankings derived from them; timeline, archive totals, and date span remain whole-archive in this increment.
- Q: Where must filtered-scope disclosure appear? -> A: In an always-visible report scope note, in every affected section, and in terminal completion feedback.
- Q: How should filtered completion feedback preserve scriptability? -> A: Keep the report path as the sole stdout success output and send filter-specific feedback to stderr or the interactive UI.
- Q: What metadata evidence is sufficient for phone classification? -> A: Trusted device identity fields such as make or model from embedded, sidecar, or platform metadata; not software tags, focal heuristics, or provenance-only hints.

## Compliance Check

### Instructions Check Report
**Target**: spec.md
**Status**: PASS

| Principle | Verdict | Notes |
|-----------|---------|-------|
| Narrow Product Scope | PASS | Limits scope to an opt-in analysis filter and excludes broader device-taxonomy, editing, or archive-management expansion. |
| Local-First Safety | PASS | Keeps classification, filtering, and reporting entirely local and read-only, with no archive mutation or network dependency. |
| Modular Pipeline Design | PASS | Describes contract extensions across request, metadata, aggregation, and report layers instead of collapsing behavior into one command handler. |
| Honest Data Reporting | PASS | Requires conservative classification and explicit filtered-scope disclosure so users can understand what changed. |
| Cross-Platform Release Quality | PASS | Preserves the existing default CLI contract and durable output expectations across supported desktop platforms. |
| Agent Output Style | N/A | Repository governance for agent responses rather than a feature behavior requirement. |

**Violations**:
- None.
