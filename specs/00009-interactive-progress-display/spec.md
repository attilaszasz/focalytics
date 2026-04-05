---
feature_branch: "00009-interactive-progress-display"
created: "2026-04-05"
input: "E007 Interactive Progress Display"
spec_type: "product"
spec_maturity: "clarified"
epic_id: "E007"
epic_sources: "{PRD:CAP-001}{SAD:ADR-001}"
---

# Feature Specification: Interactive Progress Display

**Feature Branch**: `00009-interactive-progress-display`  
**Created**: 2026-04-05  
**Status**: Draft  
**Spec Type**: product  
**Spec Maturity**: clarified  
**Epic ID**: E007  
**Epic Sources**: {PRD:CAP-001}{SAD:ADR-001}  
**Product Document**: specs/prd.md

## Problem Statement

focalytics can complete a full archive scan, but large runs currently flood the terminal with per-file output instead of showing clear progress. That makes the core one-command workflow feel noisy and uncertain during the exact long-running scans the product is meant to handle. If the tool keeps reporting progress this way, users lose trust in whether the run is advancing and scripted use remains awkward because terminal chatter and durable output are mixed together.

## Scope

### Included

- Replace per-file progress chatter with an interactive terminal progress display for TTY runs.
- Show stage-level progress states for discovery, metadata recovery, aggregation, and report rendering.
- Surface live counters for discovery and metadata work without requiring per-file logging.
- Preserve a non-interactive fallback that avoids TUI control codes and keeps stdout script-friendly.
- Keep non-interactive success output quiet apart from warnings, errors, and the final report path.
- Keep warnings and final report-path visibility explicit during successful and failed runs.

### Excluded

- Changing archive discovery rules, metadata extraction logic, or aggregate calculations — those behaviors belong to earlier completed epics.
- Adding browser-side interactivity or report-view filtering — this feature affects terminal execution feedback only.
- Introducing persistent indexing, background daemons, or remote progress streaming — the runtime remains a local one-shot CLI.

### Edge Cases & Boundaries

- Piped or redirected runs must not emit terminal control sequences or transient redraw frames.
- Non-interactive runs must avoid routine progress chatter while still surfacing warnings and fatal failures.
- Very large archives must continue to feel responsive without rendering one UI update per discovered file.
- Warnings and fatal errors must remain visible even when the interactive display is active.
- Interactive warnings must remain readable after redraws rather than disappearing as transient UI state.
- Fast runs on small archives may show the progress UI only briefly, but they must still exit cleanly and report the generated report path.

## User Scenarios & Testing

### User Story 1 - Track Long Runs Clearly (Priority: P1)

As a photographer, I want focalytics to show stage-level progress and live counters while it scans a large archive, so I can tell the command is still working without reading thousands of file-level messages.

**Why this priority**: Core value proposition for large-archive use — without understandable progress feedback, the primary workflow feels unreliable.

**Independent Test**: Run the CLI in a terminal against a large fixture archive and confirm the display shows stage transitions plus live counters instead of per-file output.

**Acceptance Scenarios**:

1. **Given** an interactive terminal session and a large archive, **When** focalytics runs, **Then** it shows stage-level progress states and live counters rather than printing one line per discovered file.
2. **Given** discovery and metadata recovery are in progress, **When** the command advances through the pipeline, **Then** the visible progress updates remain concise and clearly indicate forward movement.

### User Story 2 - Keep Scripted Runs Clean (Priority: P1)

As a photographer or shell user, I want focalytics to keep non-interactive runs free of TUI redraw noise, so I can pipe, redirect, or capture the report path reliably.

**Why this priority**: Scriptability is part of the CLI contract and must survive the new progress experience.

**Independent Test**: Run the CLI with stdout piped or captured and confirm only the report path is written to stdout on success while progress UI control sequences are absent.

**Acceptance Scenarios**:

1. **Given** stdout or stderr is redirected, **When** focalytics runs successfully, **Then** it does not emit interactive progress control sequences into the redirected output.
2. **Given** a successful non-interactive run, **When** the command finishes, **Then** the report path remains the durable success output that scripts can capture and routine progress chatter is absent.

### User Story 3 - Preserve Honest Runtime Signals (Priority: P2)

As a photographer, I want warnings and completion status to remain explicit even after the UI becomes quieter, so I still understand skipped files, parse issues, and whether the final report was produced.

**Why this priority**: Honest reporting remains important, but the MVP value is already achieved once progress is clear and script-friendly.

**Independent Test**: Exercise warning and failure paths with the interactive display enabled and confirm warnings or fatal errors remain visible and do not get hidden by the TUI.

**Acceptance Scenarios**:

1. **Given** recoverable warnings occur during a run, **When** focalytics reports progress, **Then** those warnings remain visible and understandable alongside the quieter progress UI instead of disappearing in redraw noise.
2. **Given** the run completes or fails, **When** the command exits, **Then** the terminal output still makes the final status and report-path outcome unambiguous.

## Requirements

### Functional Requirements

- **FR-001**: System MUST show an interactive progress display during TTY runs instead of emitting per-file progress lines.
- **FR-002**: System MUST present stage-level status for discovery, metadata recovery, aggregation, and report rendering.
- **FR-003**: System MUST surface live aggregate progress metrics for long-running work, including discovery counts and metadata processing progress.
- **FR-004**: System MUST detect non-interactive execution and disable interactive terminal rendering for those runs.
- **FR-005**: System MUST avoid routine progress output in non-interactive success paths while still surfacing warnings and fatal errors.
- **FR-006**: System MUST preserve the generated report path as the sole stdout success output.
- **FR-007**: System MUST keep warnings and fatal errors visible when the interactive progress display is active.
- **FR-008**: System MUST stop emitting per-file candidate listings as the default user-facing progress feedback.
- **FR-009**: System MUST render interactive warnings in a way that remains readable after screen redraws.

### Key Entities

- **Progress Event**: A stage or metric update emitted by the pipeline for UI and fallback rendering.
- **Interactive Progress Session**: The transient terminal view shown during a TTY run.
- **Output Mode**: The runtime decision between interactive terminal rendering and non-interactive fallback behavior.
- **Report Path Output**: The durable success result that scripts or users consume after the pipeline finishes.

## Assumptions & Risks

### Assumptions

- The existing progress sink abstraction remains the stable integration point for a richer terminal UI.
- Bubble Tea is an acceptable dependency for the progress layer because it is already part of the architectural context.
- The report path should remain the only stdout success output even after the progress UX changes.

### Risks

- **TTY detection drift** *(likelihood: medium, impact: high)*: Incorrect interactive-mode detection could send TUI control codes into redirected output or suppress useful feedback in real terminals.
- **Stage visibility gaps** *(likelihood: medium, impact: medium)*: If stages do not publish enough lifecycle information, the new UI may look stalled even while work continues.
- **Regression in output contracts** *(likelihood: medium, impact: high)*: Removing legacy per-file output may break tests or downstream assumptions unless the report-path contract remains explicit.

## Implementation Signals

- `NEW-UI` — Add a terminal progress UI for interactive runs using the existing pipeline event flow.
- `NEW-API` — Extend progress events to represent stage lifecycle and aggregate progress metrics.
- `BREAKING-CHANGE` — Remove default per-file stdout/stderr progress chatter and preserve only the durable report-path output contract.

## Success Criteria

### Measurable Outcomes

- **SC-001** [US1]: In an interactive terminal, large archive runs show concise stage-level progress and live counters instead of thousands of file-level lines.
- **SC-002** [US2]: In a piped or redirected run, the command emits no interactive control noise and preserves the report path as the durable success output.
- **SC-003** [US3]: Warning and failure conditions remain explicit enough that users can tell whether the run completed successfully and what issues were encountered.

## Clarifications

### Session 2026-04-05

- Q: Should non-interactive runs still emit routine text progress updates? -> A: No. Non-interactive runs stay quiet except for warnings, fatal errors, and the final report path.
- Q: How must warnings behave while the interactive UI is active? -> A: Warnings must remain persistently readable and not disappear as transient redraw-only state.

## Glossary

| Term | Definition |
|------|------------|
| TTY | A terminal-attached session where interactive redraw-based output is appropriate. |
| Progress sink | The runtime interface that receives progress events from pipeline stages. |
| Durable output | Stable command output that users or scripts can rely on after the run completes. |

## Compliance Check

### Instructions Check Report
**Target**: spec.md
**Status**: PASS

| Principle | Verdict | Notes |
|-----------|---------|-------|
| Narrow Product Scope | PASS | Limits scope to terminal progress UX and excludes broader archive-management or browser interactivity changes. |
| Local-First Safety | PASS | Keeps execution offline and read-only; no source archive mutation or external service dependency is introduced. |
| Modular Pipeline Design | PASS | Builds on the existing progress sink and stage pipeline boundaries rather than collapsing logic into the CLI entry point. |
| Honest Data Reporting | PASS | Requires warnings and final status to remain explicit despite quieter progress output. |
| Cross-Platform Release Quality | PASS | Preserves CLI behavior across interactive and non-interactive runs without narrowing supported operating systems. |
| Agent Output Style | N/A | Repository governance for agent responses rather than a feature-spec behavior requirement. |

**Violations**:
- None.