# Tasks: Recover Metadata Facts

**Input**: Design documents from `specs/00003-recover-metadata-facts/`  
**Prerequisites**: `plan.md` (required), `spec.md` (required)

**Tests**: Include unit, integration, and coverage tasks because the spec requires layered metadata recovery, provenance tracking, and warning-first degradation.

**Organization**: Tasks are grouped by user story (`US#`). Shared plumbing stays in Setup or Foundational only when it blocks multiple stories.

## Project Mode

`Greenfield extension`

## Epic / Capability Map

- `[US1]` → Recover embedded EXIF metadata into normalized facts.
- `[US2]` → Fill missing metrics from XMP sidecars and explicit fallback layers.
- `[US3]` → Preserve warnings and exclusions for corrupt or incomplete metadata.

## Phase 1: Setup (Repository / Workspace Delta)

- [X] T001 [P] [US1] {FR-001,FR-005} Create the metadata package skeleton in src/internal/metadata/fact.go and src/internal/metadata/stage.go
- [X] T002 [P] [US1] {FR-001} Add the EXIF dependency and metadata parser scaffolding in src/go.mod and src/internal/metadata/service.go

---

## Phase 2: Foundational (Cross-Story Blockers)

- [X] T003 {FR-001,FR-005} Add shared run-context artifact storage in src/internal/app/context.go and wire discovery stage output storage in src/internal/discovery/stage.go
- [X] T004 {FR-001,FR-002,FR-003,FR-004,FR-005,FR-006} Implement layered metadata recovery, provenance, exclusions, and warning publication in src/internal/metadata/service.go and src/internal/metadata/xmp.go
- [X] T005 {FR-001,FR-005} Wire the metadata stage into src/cmd/root.go

---

## Phase 3: User Story 1 - Recover Embedded Metadata (Priority: P1) 🎯 MVP

- [X] T006 [US1] {FR-001,FR-005} Add embedded EXIF recovery tests in src/internal/metadata/service_test.go using the gallery fixture

---

## Phase 4: User Story 2 - Use Sidecars and Fallbacks Honestly (Priority: P1) 🎯 MVP

- [X] T007 [US2] {FR-002,FR-003,FR-004,FR-005} Add sidecar and fallback recovery tests in src/internal/metadata/service_test.go

---

## Phase 5: User Story 3 - Surface Exclusions and Corrupt Metadata (Priority: P1) 🎯 MVP

- [X] T008 [US3] {FR-005,FR-006} Add corrupt-metadata warning and exclusion tests in src/internal/metadata/service_test.go and CLI integration coverage in src/cmd/run_integration_test.go

## Dependencies

Setup → Foundational → User Story 1 → User Story 2 → User Story 3

- Tasks marked `[P]` can run in parallel within their phase.
- User Story 1 depends on the shared run-context artifacts and metadata package scaffolding.
- User Story 2 depends on the metadata service helpers introduced for User Story 1.
- User Story 3 depends on the full layered-recovery pipeline so warnings and exclusions reflect real parsing paths.