---
feature_branch: "00008-render-offline-report"
created: "2026-04-05"
input: "E005 Render Offline Report"
spec_type: "product"
spec_maturity: "draft"
epic_id: "E005"
epic_sources: "{PRD:CAP-003,CAP-004,CAP-005,CAP-006}{SAD:ADR-004}"
---

# Feature Specification: Render Offline Report

**Feature Branch**: `00008-render-offline-report`  
**Created**: 2026-04-05  
**Status**: Draft  
**Spec Type**: product  
**Spec Maturity**: draft  
**Epic ID**: E005  
**Epic Sources**: {PRD:CAP-003,CAP-004,CAP-005,CAP-006}{SAD:ADR-004}  
**Product Document**: specs/prd.md

## Problem Statement

focalytics now discovers files, recovers metadata, and aggregates archive insights, but users still do not receive the finished artifact they came to the product for. Without a self-contained offline report that turns aggregate data into a polished HTML dashboard with clear exclusion notes, the CLI stops short of delivering any durable end-user value.

## Scope

### Included

- Render one self-contained HTML report from the aggregate summary artifact.
- Write the report to the current working directory with a timestamped `.html` filename.
- Present overview, timeline, gear, technical analytics, and missing-data notes using aggregate-only data.
- Surface the generated report path to the user at the end of a successful run.

### Excluded

- Interactive filtering, search, or browser-side data exploration — the MVP report stays static and offline.
- Changes to discovery, metadata recovery, or aggregate calculations — those belong to earlier epics.
- Multi-file asset bundles or CDN-backed resources — the output must remain one portable HTML file.

### Edge Cases & Boundaries

- Empty or sparse archives must still render a valid report with zero-state messaging instead of broken sections.
- Missing metric families must produce visible exclusion notes rather than blank charts that imply complete data.
- Report generation failures must stop the run clearly and avoid leaving a misleading partial success message.

## User Scenarios & Testing

### User Story 1 - Generate One Offline Report (Priority: P1)

As a photographer, I want focalytics to write one HTML report I can open locally, so the command produces a durable artifact instead of only terminal output.

**Why this priority**: The product promise is not fulfilled until the CLI emits the final report artifact.

**Independent Test**: Run the CLI against a fixture archive and confirm it writes exactly one `.html` report in the working directory and prints the report path.

**Acceptance Scenarios**:
1. **Given** a successful archive run, **When** rendering completes, **Then** the system writes one self-contained HTML file to the working directory with a timestamped name.
2. **Given** the report file was written successfully, **When** the command exits, **Then** the user sees the generated report path in command output.

### User Story 2 - Show Archive Insights Clearly (Priority: P1)

As a photographer, I want the report to present overview, timeline, gear, and technical sections in a readable offline layout, so I can understand the archive without inspecting raw counts.

**Why this priority**: The report must translate aggregate data into a coherent artifact, not just dump serialized structures.

**Independent Test**: Render a report from a representative aggregate fixture and confirm it contains the expected overview, timeline, gear, and technical sections with aggregate-backed content.

**Acceptance Scenarios**:
1. **Given** an aggregate summary with overview, timeline, gear, and technical data, **When** rendering runs, **Then** the report includes dedicated sections for each summary family.
2. **Given** aggregate-only chart inputs, **When** the report is generated, **Then** the rendered visuals rely only on embedded HTML, CSS, and inline SVG rather than external libraries.

### User Story 3 - Explain Missing Data Honestly (Priority: P1)

As a photographer, I want the report to explain which charts exclude incomplete data, so I can trust the visuals instead of assuming the archive was fully represented.

**Why this priority**: Honest data reporting is a core product principle and must remain visible in the final artifact.

**Independent Test**: Render a report from an aggregate fixture with exclusion counts and confirm every affected section includes a visible missing-data note.

**Acceptance Scenarios**:
1. **Given** an aggregate summary with exclusion counts for a metric family, **When** the corresponding section renders, **Then** the report includes a visible note describing the omitted data.
2. **Given** an aggregate summary with warnings and sparse sections, **When** rendering runs, **Then** the report stays valid and explains missing or excluded data without hiding the affected section entirely.

## Requirements

### Functional Requirements

- **FR-001**: System MUST render one self-contained HTML report from the aggregate summary artifact produced by E004.
- **FR-002**: System MUST write the report to the current working directory using a timestamped `.html` filename.
- **FR-003**: System MUST present overview, timeline, gear, and technical analytics sections backed only by aggregate data structures.
- **FR-004**: System MUST embed required presentation assets in the generated HTML so the report opens offline without auxiliary files or network requests.
- **FR-005**: System MUST display visible exclusion notes for each rendered section affected by missing metric data.
- **FR-006**: System MUST print the generated report path after a successful run and fail clearly if the report cannot be written.

### Key Entities

- **Report Model**: The render-ready view of aggregate summary data used to populate the HTML template.
- **Report Artifact**: The generated standalone `.html` file written to disk.
- **Section View**: A report subsection such as overview, timeline, gear, or technical analytics.
- **Exclusion Note**: Visible explanatory text attached to a section when missing data affected its inputs.

## Assumptions & Risks

### Assumptions

- E004 aggregate summaries remain the only input required for the first report implementation.
- Modern desktop browsers can render semantic HTML, CSS, and inline SVG without extra runtime dependencies.
- The current working directory is a suitable default output location for the MVP command flow.

### Risks

- **Template drift** *(likelihood: medium, impact: high)*: If the report model and template expectations diverge, rendering may fail late in the pipeline.
- **Visual ambiguity** *(likelihood: medium, impact: medium)*: Dense sections or unclear labels may reduce report usefulness even when data is correct.
- **Artifact fragility** *(likelihood: low, impact: high)*: If the output depends on external assets or unstable naming, the report will not remain portable offline.

## Implementation Signals

- `NEW-MODULE` — Add a dedicated rendering package under `/src/internal/render` for report modeling, templating, and file output.
- `NEW-UI` — Create a self-contained HTML/CSS presentation with inline SVG-based visual summaries.
- `NEW-ENTITY` — Define a render-ready report model and result type derived from aggregate summaries.

## Success Criteria

### Measurable Outcomes

- **SC-001** [US1]: A successful run writes exactly one self-contained HTML report artifact and reports its path to the user.
- **SC-002** [US2]: The generated report includes overview, timeline, gear, and technical sections populated from aggregate-only data.
- **SC-003** [US3]: Every affected report section surfaces exclusion notes when input data was missing or omitted.

## Glossary

| Term | Definition |
|------|------------|
| Self-contained report | A single HTML artifact that contains the assets it needs and does not require network-hosted resources. |
| Inline SVG | Vector markup embedded directly in the HTML document for charts or heatmaps. |
| Exclusion note | Visible explanation that a report section omitted part of the archive because required data was missing. |

## Compliance Check

- **Status**: PASS
- **Policy**: Aligns with offline-only runtime, modular `/src` boundaries, and explicit data-quality reporting from project-instructions.md.
- **Notes**: The feature stays within the product’s one-command local workflow and preserves the self-contained report constraint.