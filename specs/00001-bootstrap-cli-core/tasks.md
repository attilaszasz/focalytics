# Tasks: Bootstrap CLI Core

**Input**: Design documents from `specs/00001-bootstrap-cli-core/`  
**Prerequisites**: `plan.md` (required), `spec.md` (required), `research.md`, `data-model.md`, `contracts/`

**Tests**: Include test tasks because the spec explicitly requires a reusable baseline test command and foundational automated coverage.

**Organization**: Tasks are grouped by technical objective (`OBJ#`). Shared repo scaffolding stays in Setup or Foundational only when it blocks multiple objectives.

## Project Mode

`Greenfield`

## Epic / Capability Map

- `[OBJ1]` → Establish the executable entrypoint, command constructors, and deterministic exit-code skeleton.
- `[OBJ2]` → Define runtime stage and progress contracts that later epics can implement safely.
- `[OBJ3]` → Provide the foundational automated test harness and QC entrypoints.

## Phase 1: Setup (Repository / Workspace Delta)

- [X] T001 Initialize the Go module and Cobra dependency manifest in src/go.mod
- [X] T002 [P] Add the shared lint configuration in .golangci.yml
- [X] T003 [P] Create the executable bootstrap entrypoint in src/main.go

---

## Phase 2: Foundational (Cross-Objective Blockers)

- [X] T004 {TR-003} Create the scan request model in src/internal/app/request.go
- [X] T005 [P] {TR-003,TR-004} Create the run context and exit policy types in src/internal/app/context.go and src/internal/app/exitcodes.go
- [X] T006 [P] {TR-003} Define the pipeline stage contract in src/internal/pipeline/stage.go
- [X] T007 [P] {TR-005} Define the progress event and sink contracts in src/internal/progress/event.go and src/internal/progress/sink.go

---

## Phase 3: Objective 1 - Establish Command and Runtime Skeleton (Priority: P1) 🎯 MVP

- [X] T008 [OBJ1] {TR-001,TR-002} Implement the root command constructor in src/cmd/root.go
- [X] T009 [OBJ1] {TR-001,TR-002,TR-004} Implement run command validation and exit mapping in src/cmd/run.go
- [X] T010 [OBJ1] {TR-001,TR-004} Wire binary startup to command execution in src/main.go

---

## Phase 4: Objective 2 - Define Progress and Pipeline Boundaries (Priority: P1) 🎯 MVP

- [X] T011 [OBJ2] {TR-003,TR-005} Implement the ordered runtime runner in src/internal/pipeline/runner.go
- [X] T012 [P] [OBJ2] {TR-005} Implement the no-op progress sink in src/internal/progress/noop.go
- [X] T013 [OBJ2] {TR-003,TR-005} Wire default stage and progress dependencies into the run command in src/cmd/run.go

---

## Phase 5: Objective 3 - Provide a Baseline Test Harness (Priority: P1) 🎯 MVP

- [X] T014 [P] [OBJ3] {TR-006} Add root command and validation tests in src/cmd/root_test.go
- [X] T015 [P] [OBJ3] {TR-006} Add exit-code policy tests in src/internal/app/exitcodes_test.go
- [X] T016 [OBJ3] {TR-006} Add run command helpers and integration-tagged tests in src/cmd/run_test.go

## Dependencies

Setup → Foundational → Objective 1 → Objective 2 → Objective 3

- Tasks marked `[P]` can run in parallel within their phase.
- Objective 1 depends on the shared runtime models and contracts from Foundational.
- Objective 2 depends on Objective 1's command skeleton so runtime wiring has a stable entrypoint.
- Objective 3 depends on Objectives 1 and 2 so tests exercise real command and runtime seams.
