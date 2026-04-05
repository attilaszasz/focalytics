# Tasks: Aggregate Insight Metrics

**Input**: Design documents from `specs/00007-aggregate-insight-metrics/`  
**Prerequisites**: `plan.md` (required), `spec.md` (required)

**Tests**: Include unit, integration, and coverage tasks because the spec requires deterministic aggregation, artifact handoff, and explicit exclusion visibility.

**Organization**: Tasks are grouped by user story (`US#`). Shared summary models and pipeline plumbing stay in Setup or Foundational only when they block multiple stories.

## Project Mode

`Brownfield extension`

## Epic / Capability Map

- `[US1]` → Aggregate timeline summaries and archive date span.
- `[US2]` → Aggregate gear and technical usage patterns deterministically.
- `[US3]` → Preserve warnings and exclusions in aggregate output.

## Phase 1: Setup (Repository / Workspace Delta)

- [X] T001 [P] [US1] {FR-001,FR-004,FR-006} Create the aggregate package skeleton and summary models in src/internal/aggregate/summary.go and src/internal/aggregate/stage.go
- [X] T002 [P] [US2] {FR-003} Add bucket helper scaffolding in src/internal/aggregate/buckets.go and src/internal/aggregate/service.go

---

## Phase 2: Foundational (Cross-Story Blockers)

- [X] T003 {FR-005} Add the aggregate artifact key in src/internal/app/context.go and wire the aggregate stage into src/cmd/root.go
- [X] T004 {FR-001,FR-002,FR-003,FR-004,FR-006} Implement aggregate service logic for timeline, gear, technical, warning, and exclusion summaries in src/internal/aggregate/service.go and src/internal/aggregate/buckets.go

---

## Phase 3: User Story 1 - Summarize Shooting Timeline (Priority: P1) 🎯 MVP

- [X] T005 [US1] {FR-001,FR-005} Add deterministic timeline and archive-span tests in src/internal/aggregate/service_test.go and stage coverage in src/internal/aggregate/stage.go

---

## Phase 4: User Story 2 - Summarize Gear And Technical Patterns (Priority: P1) 🎯 MVP

- [X] T006 [US2] {FR-002,FR-003,FR-006} Add gear ranking and technical bucket tests in src/internal/aggregate/service_test.go

---

## Phase 5: User Story 3 - Preserve Data Quality Context (Priority: P1) 🎯 MVP

- [X] T007 [US3] {FR-004,FR-005} Add warning and exclusion summary tests in src/internal/aggregate/service_test.go
- [X] T008 [US3] {FR-005} Extend CLI integration coverage for aggregate stage wiring in src/cmd/run_integration_test.go

## Dependencies

Setup → Foundational → User Story 1 → User Story 2 → User Story 3

- Tasks marked `[P]` can run in parallel within their phase.
- User Story 1 depends on the shared aggregate models and pipeline artifact wiring.
- User Story 2 depends on the bucket helpers introduced in the foundational aggregate service.
- User Story 3 depends on the full summary model so warnings and exclusions reflect the same aggregate output used by later rendering.