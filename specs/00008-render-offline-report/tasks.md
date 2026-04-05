# Tasks: Render Offline Report

**Input**: Design documents from `specs/00008-render-offline-report/`  
**Prerequisites**: `plan.md` (required), `spec.md` (required)

**Tests**: Include unit, integration, golden-file, and coverage tasks because the spec requires one portable HTML artifact, deterministic output naming, and visible exclusion notes.

**Organization**: Tasks are grouped by user story (`US#`). Shared renderer scaffolding and stage wiring stay in Setup or Foundational only when they block multiple stories.

## Project Mode

`Brownfield extension`

## Epic / Capability Map

- `[US1]` → Write one self-contained offline report artifact.
- `[US2]` → Render overview, timeline, gear, and technical sections from aggregate data.
- `[US3]` → Surface exclusion notes and clear render-failure behavior.

## Phase 1: Setup (Repository / Workspace Delta)

- [X] T001 [P] [US1] {FR-001,FR-004} Create the render package skeleton and embedded asset files in src/internal/render/assets.go, src/internal/render/templates/report.html.tmpl, and src/internal/render/templates/report.css
- [X] T002 [P] [US2] {FR-003,FR-005} Add the report model and output-path scaffolding in src/internal/render/model.go and src/internal/render/path.go

---

## Phase 2: Foundational (Cross-Story Blockers)

- [X] T003 {FR-001,FR-006} Implement render service file output and timestamped naming in src/internal/render/service.go and src/internal/render/path.go
- [X] T004 {FR-001,FR-006} Add the render stage and wire it into src/internal/app/context.go, src/internal/render/stage.go, src/cmd/root.go, and src/cmd/run.go

---

## Phase 3: User Story 1 - Generate One Offline Report (Priority: P1) 🎯 MVP

- [X] T005 [US1] {FR-001,FR-002,FR-006} Add renderer output and CLI artifact-path tests in src/internal/render/service_test.go, src/cmd/run_test.go, and src/cmd/run_integration_test.go

---

## Phase 4: User Story 2 - Show Archive Insights Clearly (Priority: P1) 🎯 MVP

- [X] T006 [US2] {FR-003,FR-004} Implement overview, timeline, gear, and technical section rendering in src/internal/render/service.go and src/internal/render/templates/report.html.tmpl
- [X] T007 [US2] {FR-003,FR-004} Add golden and structural render coverage in src/internal/render/service_test.go and testdata under src/internal/render/testdata/

---

## Phase 5: User Story 3 - Explain Missing Data Honestly (Priority: P1) 🎯 MVP

- [X] T008 [US3] {FR-005,FR-006} Add exclusion-note and render-failure tests in src/internal/render/service_test.go and src/cmd/run_test.go

## Dependencies

Setup → Foundational → User Story 1 → User Story 2 → User Story 3

- Tasks marked `[P]` can run in parallel within their phase.
- User Story 1 depends on shared renderer scaffolding and stage wiring.
- User Story 2 depends on the render service and model introduced in the foundational phase.
- User Story 3 depends on the fully rendered section model so exclusion notes and failure behavior reflect the final artifact path.