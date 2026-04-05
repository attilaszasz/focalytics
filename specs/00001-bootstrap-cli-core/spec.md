---
feature_branch: "00001-bootstrap-cli-core"
created: "2026-04-05"
input: "E001 Bootstrap CLI Core"
spec_type: "technical"
spec_maturity: "draft"
epic_id: "E001"
epic_sources: "{SAD:ADR-001,ADR-002}"
---

# Feature Specification: Bootstrap CLI Core

**Feature Branch**: `00001-bootstrap-cli-core`  
**Created**: 2026-04-05  
**Status**: Draft  
**Spec Type**: technical  
**Spec Maturity**: draft  
**Epic ID**: E001  
**Epic Sources**: {SAD:ADR-001,ADR-002}  
**Product Document**: specs/prd.md

## Problem Statement

focalytics has a defined product and architecture, but no executable foundation yet. Until the repository has a stable CLI entrypoint, shared runtime contracts, and a repeatable test command under /src, every later epic depends on assumptions instead of real interfaces. This blocks delivery of discovery, metadata recovery, aggregation, rendering, and release automation.

## Scope

### Included

- Establish the initial Go module layout and executable entrypoint under /src for a single-binary CLI.
- Define shared runtime contracts for command execution, scan requests, progress events, and pipeline stage boundaries.
- Provide deterministic exit-code handling for success, invalid input, and fatal runtime failure.
- Add a baseline automated test path suitable for local development and CI reuse.

### Excluded

- Recursive archive traversal and file filtering logic — delivered by E002 once the runtime contract exists.
- Metadata parsing, normalization, aggregation, and report rendering behavior — those remain in later epics even if package placeholders exist.
- Release automation and package-manager metadata — delivered by E006 after the build contract is real.

### Edge Cases & Boundaries

- Invalid or missing archive-root arguments must fail before any pipeline stage executes.
- The scaffold must not require a TTY; progress integration needs a non-interactive boundary for later automation and tests.
- Shared interfaces may define future module seams, but they must not force fake business behavior or placeholder outputs that later epics have to unwind.

## Technical Objectives

### Objective 1 - Establish Command and Runtime Skeleton (Priority: P1)

Create the initial executable structure for focalytics so later epics inherit a stable command contract instead of inventing their own entrypoints. The foundation must keep Cobra focused on CLI wiring while domain packages own discovery, metadata recovery, aggregation, rendering, and progress contracts.

**Why this priority**: Core delivery blocker — every remaining epic depends on a stable runtime surface.

**Rationale**: A thin command layer and explicit package seams reduce coupling early, which is critical in a modular monolith without process boundaries.

**Deliverables**:
- Go module and executable entrypoint rooted under /src
- Root command and run command constructors
- Shared runtime types for scan requests, run context, and stage interfaces
- Exit-code contract documented in code and tests

**Validation Criteria**:
1. **Given** a clean checkout, **When** a developer runs the build command for the scaffold, **Then** the repository produces a buildable focalytics binary from source rooted under /src.
2. **Given** a valid archive path or an invalid path, **When** the CLI is invoked, **Then** it returns the documented exit-code class for success or invalid input without depending on unimplemented downstream stages.

### Objective 2 - Define Progress and Pipeline Boundaries (Priority: P1)

Introduce runtime contracts that let future epics publish scan progress and consume pipeline stages without binding orchestration logic to Bubble Tea or any specific parser implementation.

**Why this priority**: Core architectural seam — later product work depends on these interfaces staying stable.

**Rationale**: Progress reporting and pipeline coordination will be used by nearly every later epic, so the boundary needs to exist before the first feature logic lands.

**Deliverables**:
- Typed progress event model and observer interface
- Pipeline stage interfaces for discovery, metadata recovery, aggregation, and rendering
- Default no-op or not-yet-configured wiring that keeps the runtime executable

**Validation Criteria**:
1. **Given** the scaffolded runtime, **When** future modules are linked through the shared interfaces, **Then** command orchestration depends on interfaces rather than concrete package globals.
2. **Given** a non-interactive environment, **When** progress reporting is disabled or absent, **Then** the command still executes deterministically without terminal-control requirements.

### Objective 3 - Provide a Baseline Test Harness (Priority: P1)

Create the first automated tests that prove command construction, argument validation, and exit-code handling can be verified without running the full product pipeline.

**Why this priority**: CI gate blocker — later epics need a reusable test command and regression baseline from the start.

**Rationale**: Early tests prevent the runtime contract from drifting while the rest of the system is built on top of it.

**Deliverables**:
- Package-level tests for command creation and input validation
- Test helpers for temporary archive paths and command execution
- A documented baseline `go test` target for foundational packages

**Validation Criteria**:
1. **Given** the scaffolded repository, **When** the baseline test command is run, **Then** foundational package tests pass without network access or real archive fixtures.
2. **Given** a regression in command validation or exit-code mapping, **When** tests are re-run, **Then** the failure is isolated to the foundational packages rather than hidden in later epics.

### Technical Constraints

- All implementation code for this epic must live under /src and preserve the single-binary architecture.
- The runtime must remain offline-only, stateless per run, and free of source-archive mutation.
- The scaffold must compile on macOS, Windows, and Linux without introducing platform-specific command behavior.
- Bubble Tea may be prepared as an adapter boundary, but the runtime contract must not require an interactive terminal to function.

## Integration Points

- **IP-001**: E002 depends on the command entrypoint, scan request model, and progress event boundary introduced by this epic.
- **IP-002**: E003 and E004 depend on pipeline-stage interfaces and shared run context defined by this epic.
- **IP-003**: E005 depends on the renderer boundary and exit-code contract remaining stable through later waves.
- **IP-004**: E006 depends on the binary name, build target, and reusable test command established by this epic.

## Requirements

### Technical Requirements

- **TR-001**: System MUST provide a buildable focalytics executable with its source rooted under /src and a thin CLI entrypoint that delegates work to internal runtime packages.
- **TR-002**: System MUST expose constructor-based command wiring for the root command and primary run command rather than relying on package-global initialization as the orchestration mechanism.
- **TR-003**: System MUST define shared runtime types for scan requests, run context, and pipeline stage contracts used by discovery, metadata recovery, aggregation, and rendering.
- **TR-004**: System MUST define deterministic exit-code handling for successful completion, invalid input, and fatal runtime failure.
- **TR-005**: System MUST define a progress reporting contract that supports both interactive and non-interactive execution without making terminal control mandatory.
- **TR-006**: System MUST include foundational automated tests that validate command construction, argument validation, and exit-code behavior through a repeatable `go test` command.

### Key Entities

- **Scan Request**: The user-supplied invocation intent, including archive root and command-level options needed to start one run.
- **Run Context**: Shared execution metadata passed across pipeline stages, including configuration, logging hooks, and progress sinks.
- **Progress Event**: A typed runtime signal describing scan state, counters, current path, or warnings without embedding UI concerns.
- **Pipeline Stage Contract**: The interface boundary a module implements so the command runtime can coordinate discovery, metadata recovery, aggregation, or rendering without concrete coupling.

## Assumptions & Risks

### Assumptions

- Go 1.24 is available in local development and CI as the project baseline.
- The initial scaffold can defer real business logic while still exposing stable runtime interfaces.
- Later epics will consume shared contracts through package imports rather than replacing the entrypoint structure.
- A baseline `go test` command is sufficient for this epic before golden files or large fixtures exist.

### Risks

- **Boundary overdesign** *(likelihood: medium, impact: medium)*: Overly abstract interfaces could slow later epics if the scaffold predicts the pipeline incorrectly.
- **Runtime drift** *(likelihood: medium, impact: high)*: If command and exit-code behavior are not pinned down now, later epics may diverge and break CI or release assumptions.
- **TTY coupling** *(likelihood: low, impact: medium)*: If progress handling is tied too closely to Bubble Tea, automated runs and tests may become brittle.

## Implementation Signals

- `NEW-CONFIG` — Introduce the initial Go module, source root, and build/test entrypoints required for local and CI execution.
- `NEW-API` — Define constructor-based command wiring and the shared runtime interfaces consumed by later epics.
- `NEW-ENTITY` — Add first-class runtime models for scan requests, run context, progress events, and pipeline stages.

## Success Criteria

### Measurable Outcomes

- **SC-001** [OBJ1]: The repository builds a focalytics binary from code under /src using a single documented Go build command on a clean checkout.
- **SC-002** [OBJ2]: Foundational packages compile while depending on declared interfaces for progress and pipeline stages rather than direct concrete implementations.
- **SC-003** [OBJ3]: A baseline `go test` command passes and fails deterministically when command validation or exit-code behavior is intentionally broken.

## Glossary

| Term | Definition |
|------|------------|
| Command constructor | A function that returns a configured Cobra command without executing package-global registration side effects. |
| Progress event | A runtime message describing execution state in a UI-agnostic form. |
| Run context | Shared execution state and services passed through the focalytics pipeline for a single invocation. |

## Compliance Check

- **Status**: PASS
- **Policy**: Aligns with local-first safety, /src source-root enforcement, modular pipeline design, and cross-platform viability from project-instructions.md.
- **Notes**: The spec preserves offline execution, avoids persistent state, and keeps later pipeline concerns separated behind explicit runtime contracts.