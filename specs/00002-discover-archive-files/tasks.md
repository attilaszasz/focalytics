# Tasks: Discover Archive Files

**Input**: Design documents from `specs/00002-discover-archive-files/`  
**Prerequisites**: `plan.md` (required), `spec.md` (required)

**Tests**: Include unit, integration, and coverage tasks because the spec requires deterministic traversal, progress reporting, and warning handling.

**Organization**: Tasks are grouped by user story (`US#`). Shared discovery scaffolding stays in Setup or Foundational only when it blocks multiple stories.

## Project Mode

`Greenfield extension`

## Epic / Capability Map

- `[US1]` → Traverse nested archives deterministically and emit supported candidates.
- `[US2]` → Surface non-interactive progress updates during traversal.
- `[US3]` → Fail fast on invalid roots while continuing past unreadable child entries.

## Phase 1: Setup (Repository / Workspace Delta)

- [X] T001 [P] [US1] {FR-001,FR-002} Create the discovery package skeleton in src/internal/discovery/candidate.go and src/internal/discovery/stage.go
- [X] T002 [P] [US2] {FR-004} Extend the progress event model and add a text sink in src/internal/progress/event.go and src/internal/progress/text.go

---

## Phase 2: Foundational (Cross-Story Blockers)

- [X] T003 {FR-001,FR-002,FR-003,FR-006} Implement deterministic recursive discovery, filtering, and warning capture in src/internal/discovery/service.go
- [X] T004 {FR-001,FR-004} Wire the discovery stage and text progress sink into src/cmd/root.go

---

## Phase 3: User Story 1 - Scan a Real Archive (Priority: P1) 🎯 MVP

- [X] T005 [US1] {FR-001,FR-002,FR-003} Add deterministic traversal and candidate emission tests in src/internal/discovery/service_test.go and src/cmd/run_integration_test.go

---

## Phase 4: User Story 2 - Observe Scan Progress (Priority: P1) 🎯 MVP

- [X] T006 [US2] {FR-004} Add text progress sink and command progress tests in src/internal/progress/text_test.go and src/cmd/run_test.go

---

## Phase 5: User Story 3 - Continue Past Child Read Errors (Priority: P1) 🎯 MVP

- [X] T007 [US3] {FR-005,FR-006} Add unreadable-child warning coverage in src/internal/discovery/service_test.go and preserve invalid-root command coverage in src/cmd/run_test.go

## Dependencies

Setup → Foundational → User Story 1 → User Story 2 → User Story 3

- Tasks marked `[P]` can run in parallel within their phase.
- User Story 1 depends on the foundational discovery service and command wiring.
- User Story 2 depends on the discovery service so progress can report real traversal metrics.
- User Story 3 depends on the same discovery service and command validation paths used by the earlier stories.