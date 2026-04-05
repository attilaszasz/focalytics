# Tasks: Interactive Progress Display

**Input**: Design documents from `specs/00009-interactive-progress-display/`
**Prerequisites**: `plan.md` (required), `spec.md` (required), `research.md`, `contracts/`

**Tests**: Regression and coverage updates are included because the feature changes existing stdout/stderr contracts and progress behavior.

**Organization**: Tasks are grouped by user story, with a shared foundational phase for the cross-story progress contract changes.

## Project Mode

`Brownfield`

## Epic / Capability Map

- `[US1]` → Interactive terminal progress for long archive runs
- `[US2]` → Quiet non-interactive output with durable report-path stdout
- `[US3]` → Persistent warning and failure visibility in the new progress UX

## Brownfield Notes

- Existing flows touched: `src/cmd/root.go`, `src/internal/progress/`, `src/internal/pipeline/runner.go`, `src/internal/discovery/service.go`, `src/internal/render/service.go`
- Compatibility or migration concerns: preserve `progress.Sink`, keep stdout script-friendly, and avoid leaking TUI control codes into redirected runs
- Regression focus: interactive progress behavior, non-interactive output contracts, warning visibility, report-path success output

## Phase 1: Foundational (Cross-Work-Item Blockers)

- [X] T001 {FR-001,FR-002,FR-003,FR-004} Extend the progress event contract and add TTY dependencies in src/internal/progress/event.go and src/go.mod
- [X] T002 {FR-002,FR-003,FR-007} Publish stage lifecycle and aggregate progress signals in src/internal/pipeline/runner.go and src/internal/metadata/service.go

---

## Phase 2: Work Item 1 - Track Long Runs Clearly (Priority: P1) 🎯 MVP

- [X] T003 [US1] {FR-001,FR-002,FR-003} Implement the Bubble Tea progress model and sink in src/internal/progress/tui.go
- [X] T004 [US1] {FR-001,FR-002,FR-003,FR-004} Wire the interactive TTY execution path in src/cmd/root.go
- [X] T005 [US1] {FR-001,FR-002,FR-003} Add progress model and runner regression coverage in src/internal/progress/tui_test.go and src/internal/pipeline/runner_test.go

---

## Phase 3: Work Item 2 - Keep Scripted Runs Clean (Priority: P1) 🎯 MVP

- [X] T006 [US2] {FR-004,FR-005,FR-006,FR-008} Remove routine candidate and report chatter from src/internal/discovery/service.go and src/internal/render/service.go
- [X] T007 [US2] {FR-004,FR-005,FR-006} Update non-interactive sink behavior in src/internal/progress/text.go and src/internal/progress/text_test.go
- [X] T008 [US2] {FR-004,FR-005,FR-006,FR-008} Update CLI and integration regressions in src/cmd/run_test.go and src/cmd/run_integration_test.go

---

## Phase 4: Work Item 3 - Preserve Honest Runtime Signals (Priority: P2)

- [X] T009 [US3] {FR-007,FR-009} Implement persistent warning presentation and final-status handling in src/internal/progress/tui.go
- [X] T010 [US3] {FR-007,FR-009} Add warning and failure visibility coverage in src/internal/discovery/service_test.go, src/internal/metadata/service_test.go, and src/internal/render/service_test.go

---

## Phase 5: Polish & Cross-Cutting Concerns

- [X] T011 [P] Update progress feature documentation and release notes in README.md and specs/00009-interactive-progress-display/contracts/runtime-progress-contracts.md

## Dependencies

Foundational → US1 → US2 → US3 → Polish

- Phase 1 establishes the shared event model and lifecycle signals required by every later phase.
- Phase 2 depends on Phase 1 because the TUI consumes the extended event contract.
- Phase 3 depends on Phase 2 so the non-interactive path can coexist with the interactive runner wiring.
- Phase 4 depends on Phases 2 and 3 because warning persistence is implemented inside the new TUI and must honor the quiet-output contract.
- T011 can run after Phases 2-4 are complete.