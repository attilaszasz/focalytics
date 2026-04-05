---
feature_branch: "00006-publish-installable-releases"
created: "2026-04-05"
input: "E006 Publish Installable Releases"
spec_type: "technical"
spec_maturity: "draft"
epic_id: "E006"
epic_sources: "{PRD:CAP-007}{SAD:ADR-005}"
---

# Feature Specification: Publish Installable Releases

**Feature Branch**: `00006-publish-installable-releases`  
**Created**: 2026-04-05  
**Status**: Draft  
**Spec Type**: technical  
**Spec Maturity**: draft  
**Epic ID**: E006  
**Epic Sources**: {PRD:CAP-007}{SAD:ADR-005}  
**Product Document**: specs/prd.md

## Problem Statement

focalytics has a buildable CLI, but it does not yet provide a trustworthy release path for end users. Without one canonical artifact set, published checksums, and package-manager update inputs derived from those exact release assets, cross-platform installation remains manual and channel drift becomes likely.

## Scope

### Included

- Define the canonical release target matrix, archive naming contract, and checksum manifest naming for GitHub Releases.
- Add CI automation that runs validation, builds cross-platform archives, publishes release assets, and uploads published checksum data.
- Generate Homebrew and WinGet update inputs from the published release asset URLs and checksum manifest.
- Add drift detection so release automation fails when the canonical artifact contract is incomplete or inconsistent.

### Excluded

- Implementing archive discovery, metadata parsing, aggregation, or report rendering behavior.
- Owning downstream Homebrew tap commits or WinGet submission automation in this repository.
- Adding runtime network behavior, cloud services, or mutable release channels.

### Edge Cases & Boundaries

- Re-running a release workflow for an existing tag must reuse the same canonical asset names instead of creating alternate channel-specific outputs.
- Windows archives must preserve an `.exe` binary name while still following the same versioned archive contract as macOS and Linux.
- Package-manager inputs must fail generation if a required release asset checksum is missing or an unexpected asset appears in the manifest.

## Technical Objectives

### Objective 1 - Publish Canonical Release Assets (Priority: P1)

Add release automation that validates the codebase, builds one cross-platform artifact matrix from the tagged source revision, and publishes those archives plus a checksum manifest to GitHub Releases.

**Why this priority**: The immutable release artifact set is the delivery root for every installation channel.

**Rationale**: If the release workflow and naming contract are not fixed first, downstream installation metadata cannot stay auditable.

**Deliverables**:
- Canonical target matrix and artifact naming contract
- GitHub Actions workflow for validation, build, checksum, and release upload
- Reusable local packaging script aligned with the same contract

**Validation Criteria**:
1. **Given** a version tag, **When** the release workflow runs, **Then** it uploads macOS, Linux, and Windows release archives plus a checksum manifest built from that tagged revision.
2. **Given** a rerun or later inspection, **When** release asset names are compared against the contract, **Then** they match the expected version, OS, and architecture layout exactly.

### Objective 2 - Derive Package-Manager Inputs From Release Assets (Priority: P1)

Generate Homebrew formula input data and WinGet manifest input data from the published GitHub Release URLs and recorded checksums instead of from separate per-channel builds.

**Why this priority**: Package-manager trust depends on downstream channels pointing to the same immutable assets users can inspect directly.

**Rationale**: One artifact root reduces supply-chain ambiguity and removes duplicated packaging logic.

**Deliverables**:
- Release metadata generator under `/src`
- Homebrew input JSON derived from release asset URLs and checksums
- WinGet input JSON derived from release asset URLs and checksums

**Validation Criteria**:
1. **Given** a checksum manifest for a tagged release, **When** package-manager inputs are generated, **Then** every referenced URL points to a GitHub Release asset and carries the matching checksum.
2. **Given** the published release assets, **When** a maintainer updates Homebrew or WinGet, **Then** no extra rebuild step is required to obtain artifact URLs or hashes.

### Objective 3 - Enforce Release Drift Detection (Priority: P1)

Fail release automation immediately when validation commands, expected archives, checksum coverage, or package-manager metadata drift from the canonical release contract.

**Why this priority**: Silent drift would undermine checksum traceability and create inconsistent install channels.

**Rationale**: Automated contract verification is cheaper than manual release audits and scales with future targets.

**Deliverables**:
- Checksum parser and contract verifier
- Tests covering asset naming and metadata generation
- Workflow failure path tied to the verifier before package metadata upload

**Validation Criteria**:
1. **Given** a missing or extra archive checksum, **When** the verifier runs, **Then** the workflow exits non-zero before metadata generation succeeds.
2. **Given** the expected artifact matrix, **When** release metadata is generated, **Then** the same verifier confirms the manifest is complete before any package-manager input is published.

### Technical Constraints

- Release automation logic must remain modular and live under `/src`, with shell automation limited to orchestration around the Go-based contract.
- The workflow must preserve immutable GitHub Release assets as the distribution root for all downstream channels.
- Validation must reuse the baseline Go build, test, lint, and vulnerability commands from the existing source module.
- The implementation must keep macOS, Windows, and Linux outputs in one canonical contract without channel-specific rebuild variants.

## Integration Points

- **IP-001**: E001 provides the binary name, source root, and reusable test commands that the release workflow must invoke without redefining them.
- **IP-002**: GitHub Actions acts as the CI release orchestrator and needs permissions to publish release assets and metadata.
- **IP-003**: Homebrew maintainers consume generated formula input data that references the canonical GitHub Release assets.
- **IP-004**: WinGet maintainers consume generated manifest input data that references the canonical GitHub Release assets.

## Requirements

### Technical Requirements

- **TR-001**: System MUST define one canonical cross-platform release contract for focalytics archives, including versioned asset names and a checksum manifest.
- **TR-002**: System MUST derive Homebrew formula data and WinGet manifest data from GitHub Release asset URLs and their published checksums rather than from channel-specific rebuilds.
- **TR-003**: System MUST run build, test, lint, and vulnerability validation before publishing release assets.
- **TR-004**: System MUST fail clearly when expected release assets, checksums, or package-manager inputs drift from the canonical release contract.

### Key Entities

- **Release Artifact**: A versioned archive containing one focalytics binary for one operating-system and architecture pair.
- **Checksum Manifest**: The published SHA-256 file that records every canonical release artifact hash.
- **Homebrew Formula Data**: Machine-readable input that maps canonical release asset URLs and checksums to Homebrew-relevant targets.
- **WinGet Manifest Data**: Machine-readable input that maps canonical release asset URLs and checksums to WinGet installer entries.

## Assumptions & Risks

### Assumptions

- GitHub-hosted runners can build the focalytics Go module for the target matrix using the existing Go toolchain.
- Package-manager maintainers are willing to consume generated input files rather than bespoke manual release notes.
- GitHub Releases remain the canonical immutable distribution root for the first installable release flow.

### Risks

- **Workflow/tool drift** *(likelihood: medium, impact: high)*: Validation commands or packaging steps could diverge from the source module contract and break releases unexpectedly.
- **Checksum coverage gaps** *(likelihood: medium, impact: high)*: Missing hashes would make downstream installation metadata unverifiable.
- **Platform matrix creep** *(likelihood: low, impact: medium)*: Adding targets ad hoc could break naming consistency if the contract is not centralized.

## Implementation Signals

- `NEW-CONFIG` — Add GitHub Actions release automation and a reusable packaging script.
- `NEW-MODULE` — Introduce a dedicated release contract package under `/src/internal/release`.
- `NEW-TOOL` — Add a Go-based metadata generator for package-manager inputs and contract verification.

## Success Criteria

### Measurable Outcomes

- **SC-001** [OBJ1]: A tagged CI run builds and publishes the canonical macOS, Linux, and Windows release archives plus a checksum manifest from one source revision.
- **SC-002** [OBJ2]: Generated Homebrew and WinGet input files reference only GitHub Release asset URLs and matching checksums from the canonical manifest.
- **SC-003** [OBJ3]: Contract verification fails deterministically when release assets, checksum contents, or generated package-manager metadata drift from the expected matrix.

## Glossary

| Term | Definition |
|------|------------|
| Canonical artifact set | The complete versioned archive matrix and checksum manifest published for one focalytics release. |
| Drift detection | Automated verification that the published artifacts and derived metadata still match the declared release contract. |
| Package-manager input | Machine-readable data maintainers use to update a downstream install channel without rebuilding binaries. |

## Compliance Check

- **Status**: PASS
- **Policy**: Aligns with local-first safety, modular `/src` boundaries, and cross-platform release quality requirements from project-instructions.md.
- **Notes**: The scope is delivery automation only; it does not alter runtime architecture or introduce hosted runtime dependencies.