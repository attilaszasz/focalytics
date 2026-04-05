# Tasks: Publish Installable Releases

**Input**: Design documents from `specs/00006-publish-installable-releases/`  
**Prerequisites**: `plan.md` (required), `spec.md` (required)

**Tests**: Include release contract tests and run the repository QC commands because the epic gates publication quality rather than a runtime feature.

**Organization**: Tasks are grouped by technical objective (`OBJ#`). Shared release scaffolding stays in Setup or Foundational only when it blocks multiple objectives.

## Project Mode

`Greenfield extension`

## Epic / Capability Map

- `[OBJ1]` → Publish one canonical cross-platform artifact set and checksum manifest.
- `[OBJ2]` → Generate Homebrew and WinGet update inputs from published release assets.
- `[OBJ3]` → Enforce drift detection before release metadata can be published.

## Phase 1: Setup (Repository / Workspace Delta)

- [X] T001 [P] {TR-001} Create the release contract package skeleton in src/internal/release/contract.go and src/internal/release/checksums.go
- [X] T002 [P] {TR-002} Add the release metadata generator entrypoint in src/tools/releasegen/main.go
- [X] T003 [P] {TR-001} Add the reusable archive packaging script in scripts/release/build-archive.sh

---

## Phase 2: Foundational (Cross-Objective Blockers)

- [X] T004 {TR-001,TR-004} Implement the canonical target matrix, asset naming, checksum parsing, and drift verification in src/internal/release/contract.go and src/internal/release/checksums.go
- [X] T005 [P] {TR-002} Implement Homebrew and WinGet metadata generation in src/internal/release/metadata.go and src/tools/releasegen/main.go
- [X] T006 [P] {TR-004} Add release contract regression tests in src/internal/release/release_test.go

---

## Phase 3: Objective 1 - Publish Canonical Release Assets (Priority: P1) 🎯 MVP

- [X] T007 [OBJ1] {TR-001,TR-003} Implement the GitHub Actions release workflow in .github/workflows/release.yml

---

## Phase 4: Objective 2 - Generate Package-Manager Update Inputs (Priority: P1) 🎯 MVP

- [X] T008 [OBJ2] {TR-002,TR-004} Wire metadata generation and upload from release assets in .github/workflows/release.yml and src/tools/releasegen/main.go

---

## Phase 5: Objective 3 - Enforce Release Drift Detection (Priority: P1) 🎯 MVP

- [X] T009 [OBJ3] {TR-003,TR-004} Add workflow verification steps that fail on validation or canonical artifact drift in .github/workflows/release.yml

## Dependencies

Setup → Foundational → Objective 1 → Objective 2 → Objective 3

- Tasks marked `[P]` can run in parallel within their phase.
- Objective 1 depends on the shared release contract and packaging utilities.
- Objective 2 depends on Objective 1 because package-manager metadata must point to the published release asset set.
- Objective 3 depends on Objectives 1 and 2 so drift checks validate the real publishing path.