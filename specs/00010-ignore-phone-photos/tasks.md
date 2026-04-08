# Tasks: Ignore Phone Photos

**Input**: Design documents from `specs/00010-ignore-phone-photos/`
**Prerequisites**: `plan.md` (required), `spec.md` (required), `research.md`, `data-model.md`, `contracts/`

**Tests**: Regression and coverage updates are included because the feature changes runtime filtering, report interpretation, and the stdout/stderr success contract.

**Organization**: Tasks are grouped by user story, with a shared foundational phase for runtime contract and classification changes.

## Project Mode

`Brownfield`

## Epic / Capability Map

- `[US1]` → Filter gear and technical insights to compare dedicated-camera work
- `[US2]` → Preserve unchanged default behavior and conservative classification
- `[US3]` → Disclose filtered scope clearly in reports and terminal completion output

## Brownfield Notes

- Existing flows touched: `src/cmd/run.go`, `src/internal/app/request.go`, `src/internal/metadata/`, `src/internal/aggregate/`, `src/internal/render/`, `src/internal/progress/tui.go`
- Compatibility or migration concerns: preserve report-path-only stdout, keep timeline and totals whole-archive, and avoid speculative phone classification
- Regression focus: filtered versus unfiltered CLI behavior, classification false positives, shared filtered counts, and empty-state report rendering

## Phase 1: Foundational (Cross-Work-Item Blockers)

- [X] T001 {FR-001,FR-003,FR-004} Extend runtime filter and classification contracts in src/internal/app/request.go and src/internal/metadata/fact.go
- [X] T002 {FR-003,FR-004,FR-010} Implement conservative phone classification rules in src/internal/metadata/camera_profiles.go and src/internal/metadata/service.go

---

## Phase 2: Work Item 1 - Compare Camera-Only Insights (Priority: P1) 🎯 MVP

- [X] T003 [US1] {FR-005,FR-006} Implement filtered gear and technical aggregation with preserved whole-archive totals in src/internal/aggregate/summary.go and src/internal/aggregate/service.go
- [X] T004 [US1] {FR-005,FR-006} Add filtered aggregation regression coverage in src/internal/aggregate/service_test.go

---

## Phase 3: Work Item 2 - Preserve Trustworthy Defaults (Priority: P1) 🎯 MVP

- [X] T005 [US2] {FR-001,FR-002,FR-008} Wire the ignore-phone-photos flag and request propagation in src/cmd/run.go and src/cmd/run_test.go
- [X] T006 [US2] {FR-002,FR-003,FR-004,FR-010} Add metadata classification guardrail coverage in src/internal/metadata/service_test.go
- [X] T007 [US2] {FR-001,FR-002,FR-008} Add filtered and unfiltered CLI integration coverage in src/cmd/run_integration_test.go

---

## Phase 4: Work Item 3 - Understand Filtered Scope (Priority: P2)

- [X] T008 [US3] {FR-007,FR-009} Extend render scope-note models and template handling in src/internal/render/model.go and src/internal/render/templates/report.html.tmpl
- [X] T009 [US3] {FR-007,FR-008,FR-009} Implement filtered report and completion disclosure in src/internal/render/service.go, src/internal/render/service_test.go, and src/internal/progress/tui.go
- [X] T010 [US3] {FR-007,FR-008,FR-009} Add filtered report disclosure regression coverage in src/internal/render/service_test.go and src/cmd/run_integration_test.go

## Dependencies

Foundational → US1 → US2 → US3

- Phase 1 establishes the request and fact contracts required by every later task.
- Phase 2 depends on Phase 1 because filtered aggregation consumes device classification state.
- Phase 3 depends on Phases 1 and 2 so CLI filtering can drive the aggregate path and preserve the default behavior contract.
- Phase 4 depends on Phases 2 and 3 because disclosure consumes shared filtered counts and must honor the final stdout/stderr contract.
