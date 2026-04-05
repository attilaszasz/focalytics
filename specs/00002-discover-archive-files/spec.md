---
feature_branch: "00002-discover-archive-files"
created: "2026-04-05"
input: "E002 Discover Archive Files"
spec_type: "product"
spec_maturity: "draft"
epic_id: "E002"
epic_sources: "{PRD:CAP-001}{SAD:ADR-001,ADR-002}"
---

# Feature Specification: Discover Archive Files

**Feature Branch**: `00002-discover-archive-files`  
**Created**: 2026-04-05  
**Status**: Draft  
**Spec Type**: product  
**Spec Maturity**: draft  
**Epic ID**: E002  
**Epic Sources**: {PRD:CAP-001}{SAD:ADR-001,ADR-002}  
**Product Document**: specs/prd.md

## Problem Statement

focalytics has a runnable CLI shell, but it still does not inspect a real archive. Until the binary can traverse nested folders deterministically, identify supported image and sidecar candidates, and report live scan progress without crashing on routine filesystem issues, later metadata and reporting epics have no trustworthy input stream.

## Scope

### Included

- Traverse the archive root recursively with deterministic ordering.
- Detect supported image candidates and XMP sidecar candidates.
- Skip symlink surprises instead of following them implicitly.
- Surface live progress updates and warnings during traversal.
- Continue past unreadable child entries while failing fast for an invalid root path.

### Excluded

- Parsing EXIF, XMP, or filesystem-derived metadata values.
- Aggregating archive insights or rendering reports.
- Persisting an index or cache of discovered files.

### Edge Cases & Boundaries

- Unsupported file types must be ignored without aborting traversal.
- Unreadable child directories or files must produce warnings and allow the scan to continue.
- The archive root itself must still fail fast if it cannot be opened or is not a directory.
- Symlinked files and directories must not be traversed automatically.

## User Scenarios & Testing

### User Story 1 - Scan a Real Archive (Priority: P1)

As a photographer, I want focalytics to walk my archive folders and identify candidate image and sidecar files, so later analysis stages operate on a deterministic file set.

**Why this priority**: This is the first product increment and the dependency root for metadata recovery.

**Independent Test**: Run the binary against a nested temp archive and confirm that supported candidates are emitted in deterministic order.

**Acceptance Scenarios**:
1. **Given** a valid archive root with nested folders and supported files, **When** focalytics runs, **Then** it traverses the tree deterministically and emits supported image and sidecar candidates in a stable order.
2. **Given** unsupported files in the same tree, **When** focalytics runs, **Then** those files are skipped without affecting the emitted candidate order.

### User Story 2 - Observe Scan Progress (Priority: P1)

As a photographer, I want to see what path focalytics is scanning and how fast it is progressing, so I know the tool is still working during large archive walks.

**Why this priority**: Responsiveness is required for trust during long-running scans.

**Independent Test**: Run the binary against a tree with multiple files and confirm stderr progress lines include current path, throughput, and running counts.

**Acceptance Scenarios**:
1. **Given** a traversal in progress, **When** focalytics visits files and directories, **Then** it surfaces progress updates with current path, files seen, candidate count, and throughput.
2. **Given** a non-interactive environment, **When** focalytics runs, **Then** the progress output remains text-based and does not require a separate UI process.

### User Story 3 - Continue Past Child Read Errors (Priority: P1)

As a photographer, I want focalytics to warn me about unreadable child entries without aborting the full scan, so one bad folder does not invalidate the rest of the archive.

**Why this priority**: Real archives often contain stale mounts, partial copies, or permission issues.

**Independent Test**: Exercise the discovery service with a simulated unreadable child directory and confirm it records a warning while still returning other candidates.

**Acceptance Scenarios**:
1. **Given** an invalid archive root, **When** focalytics starts, **Then** it fails immediately with a clear invalid-input error.
2. **Given** an unreadable child entry below a valid root, **When** focalytics scans, **Then** it emits a warning and continues traversing readable siblings.

## Requirements

### Functional Requirements

- **FR-001**: System MUST traverse the provided archive root recursively in deterministic path order.
- **FR-002**: System MUST identify supported image files and `.xmp` sidecar files as discovery candidates.
- **FR-003**: System MUST ignore unsupported files and skip symlinked entries without following them automatically.
- **FR-004**: System MUST emit scan progress that includes current path, files seen, candidate count, warning count, and throughput.
- **FR-005**: System MUST fail fast when the archive root is invalid, unreadable, or not a directory.
- **FR-006**: System MUST record warnings for unreadable child entries and continue the scan when the root remains valid.

### Key Entities

- **Archive Root**: The top-level directory the user supplied for traversal.
- **File Candidate**: A supported image or sidecar file selected for later metadata recovery.
- **Sidecar Candidate**: An `.xmp` file eligible for later pairing with an image candidate.
- **Progress Snapshot**: A point-in-time traversal update including path, counts, warnings, and throughput.

## Assumptions & Risks

### Assumptions

- The existing E001 CLI contract and stage runner remain the execution boundary for discovery.
- Text progress output on stderr is acceptable for the first user-visible increment.
- Supported file extensions can start with a practical baseline and be extended in later epics if needed.

### Risks

- **Traversal drift** *(likelihood: medium, impact: high)*: Non-deterministic directory handling would destabilize later metadata and report tests.
- **Progress spam** *(likelihood: medium, impact: medium)*: Excessive progress events could make stderr noisy if not kept concise.
- **Filesystem surprises** *(likelihood: medium, impact: medium)*: Symlinks or child read errors could cause misleading results if not surfaced explicitly.

## Implementation Signals

- `NEW-MODULE` — Add a dedicated discovery package under `/src/internal/discovery`.
- `NEW-API` — Extend progress events and add a text sink for non-interactive progress reporting.
- `NEW-ENTITY` — Introduce file candidate and discovery result models for later metadata work.

## Success Criteria

### Measurable Outcomes

- **SC-001** [US1]: Running focalytics against a nested archive emits supported candidates in deterministic order.
- **SC-002** [US2]: Progress output includes current path, counts, and throughput without requiring a second process.
- **SC-003** [US3]: Invalid roots fail fast, while unreadable child entries are reported as warnings and do not abort sibling traversal.

## Glossary

| Term | Definition |
|------|------------|
| Candidate | A discovered file that later stages may inspect for metadata. |
| Deterministic traversal | A stable walk order that does not vary between runs on the same tree. |
| Sidecar | A companion file, such as `.xmp`, that can supplement image metadata later. |

## Compliance Check

- **Status**: PASS
- **Policy**: Aligns with local-first safety, modular `/src` package boundaries, and honest reporting requirements from project-instructions.md.
- **Notes**: The feature stays offline, read-only, and explicit about skipped or unreadable inputs.