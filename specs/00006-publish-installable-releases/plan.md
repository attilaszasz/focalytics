# Implementation Plan: Publish Installable Releases

**Branch**: `00006-publish-installable-releases` | **Date**: 2026-04-05 | **Spec**: `specs/00006-publish-installable-releases/spec.md`

## Summary

**Goal**: Publish one immutable cross-platform focalytics artifact set with checksums and package-manager update inputs derived from those exact release assets.  
**Approach**: Centralize artifact naming and manifest generation in a Go release package under `/src`, use a thin shell packager to archive binaries locally and in CI, and wire GitHub Actions to validate, publish, verify, and expose package-manager metadata as release assets.  
**Key Constraint**: Preserve one canonical artifact contract across GitHub Releases, Homebrew inputs, and WinGet inputs without channel-specific rebuild drift.

## Technical Context

**Language/Version**: Go 1.24, POSIX shell for CI orchestration  
**Primary Dependencies**: Go stdlib, GitHub Actions, gh CLI on hosted runners  
**Storage**: Release assets and checksum manifest in GitHub Releases only  
**Testing**: `go test -race -count=1 ./...`, `go test -tags=integration -count=1 ./...`, `golangci-lint run ./...`, `govulncheck ./...`, targeted release contract tests  
**Target Platform**: GitHub-hosted CI producing macOS, Linux, and Windows archives  
**Project Type**: single  
**Project Mode**: greenfield extension  
**Performance Goals**: Validation should stay within normal CI limits; metadata generation should remain negligible compared with build time  
**Constraints**: Source logic under `/src`, immutable release assets, checksum traceability, cross-platform parity, no package-manager rebuilds  
**Scale/Scope**: One release workflow, six target archives, two downstream metadata outputs

## Instructions Check

| Principle | Status | Plan Response |
|-----------|--------|---------------|
| Narrow Product Scope | PASS | Limits work to delivery automation and release metadata, not runtime product expansion. |
| Local-First Safety | PASS | Adds build-time release automation only; runtime remains offline and non-mutating. |
| Modular Pipeline Design | PASS | Isolates release contract and metadata generation in `/src/internal/release` and `/src/tools/releasegen`. |
| Honest Data Reporting | PASS | Publishes explicit checksums and fails on contract drift instead of hiding inconsistencies. |
| Cross-Platform Release Quality | PASS | Defines one artifact matrix spanning macOS, Windows, and Linux with shared validation gates. |

## Architecture

```mermaid
flowchart TD
    Tag[Git tag] --> Verify[Validate source module]
    Verify --> Build[Build and archive target matrix]
    Build --> Checksums[Generate checksum manifest]
    Checksums --> Release[Upload canonical release assets]
    Release --> Meta[Generate Homebrew and WinGet inputs]
    Meta --> PublishMeta[Upload package-manager inputs]

    subgraph Repo[/Users/attila/Repos/focalytics/src]
        Contract[internal/release]
        Tool[tools/releasegen]
    end

    Build --> Contract
    Checksums --> Contract
    Meta --> Tool
    Tool --> Contract
```

## Architecture Decisions

| ID | Decision | Options Considered | Chosen | Rationale |
|----|----------|--------------------|--------|-----------|
| AD-001 | Where should the release contract live? | YAML-only workflow logic / shell-only scripts / Go package under `/src` | Go package under `/src` | Keeps artifact naming, checksum verification, and metadata generation in one typed source of truth. |
| AD-002 | How should archives be packaged? | Full third-party releaser / custom shell around `go build` / custom Go archiver | Thin shell packager around `go build` plus Go contract lookups | Reuses the native toolchain while keeping naming rules centralized in Go. |
| AD-003 | What should downstream channels consume? | Channel-specific rebuilds / handwritten notes / machine-readable release-asset metadata | Machine-readable metadata derived from GitHub Release URLs and checksums | Satisfies immutable artifact and checksum traceability requirements. |
| AD-004 | How should drift be detected? | Manual review / workflow naming conventions only / explicit checksum contract verification | Explicit checksum contract verification | Fails deterministically before downstream metadata is published. |

## Data Model Summary

| Entity | Key Fields | Relationships | Notes |
|--------|------------|---------------|-------|
| Target | goos, goarch, archive_ext, winget_architecture | creates ReleaseArtifact | Defines one supported output in the canonical matrix. |
| ReleaseArtifact | archive_name, binary_name, url, sha256 | belongs to Target; recorded in ChecksumManifest | Immutable deliverable for one OS/arch pair. |
| ChecksumManifest | asset_name -> sha256 | covers ReleaseArtifact | Canonical release proof used by metadata generation. |
| HomebrewFormulaInput | formula_name, version, artifacts | derived from ReleaseArtifact | Maintainer input, not the tap formula itself. |
| WinGetManifestInput | package_identifier, package_version, installers | derived from ReleaseArtifact | Maintainer input, not a submission workflow. |

## API Surface Summary

| Method | Path | Purpose | Auth | Req/Res Types |
|--------|------|---------|------|---------------|
| Go | `DefaultTargets` | Return the canonical release matrix | none | `[]Target` |
| Go | `ReleaseAssetName` | Produce the versioned archive name for one target | none | `name, version, target -> string` |
| Go | `ParseChecksums` | Read a checksum manifest into typed data | none | `io.Reader -> map[string]string` |
| Go | `VerifyChecksumContract` | Fail on missing or unexpected assets | none | `name, version, targets, checksums -> error` |
| CLI | `go run ./tools/releasegen metadata` | Generate Homebrew and WinGet input JSON | none | `flags -> files` |
| CLI | `go run ./tools/releasegen verify` | Validate checksum coverage against the canonical contract | none | `flags -> exit code` |

## Testing Strategy

| Tier | Tool | Scope | Mock Boundary | Install |
|------|------|-------|---------------|---------|
| Unit | `go test -race -count=1 ./...` | Release contract naming, checksum parsing, metadata generation, existing CLI packages | none | configured |
| Integration | `go test -tags=integration -count=1 ./...` | Existing tagged CLI integration tests remain part of release gating | Temp filesystem only | configured |
| Static Analysis | `golangci-lint run ./...` | Release package, tool, and existing source packages | none | configured |
| Security | `govulncheck ./...` | Reachable vulnerability scan for the Go module before release | none | `go install golang.org/x/vuln/cmd/govulncheck@latest` |
| Coverage | `go test -coverprofile=coverage.out -coverpkg=./... ./...` | Total module coverage including release automation code | none | configured |

## Error Handling Strategy

- Contract verification errors must report missing and unexpected artifact names together so release drift is immediately actionable.
- Metadata generation must fail on unknown repository, missing checksum entries, or unsupported target mappings instead of emitting partial files.
- The release workflow must stop before upload if validation or packaging fails, and before metadata upload if checksum verification fails.

## Integration Points

| Spec Reference | System/Service | Technical Approach | Contract |
|----------------|----------------|--------------------|----------|
| IP-001 | E001 source module | Reuse `/src` build and validation commands without redefining module layout | `src/go.mod`, E001 testing strategy |
| IP-002 | GitHub Actions | Tag-triggered workflow validates, builds, releases, and uploads metadata | `.github/workflows/release.yml` |
| IP-003 | Homebrew maintainer flow | JSON input references release assets and checksums for formula updates | `homebrew-formula.json` release asset |
| IP-004 | WinGet maintainer flow | JSON input references release assets and checksums for manifest updates | `winget-manifests.json` release asset |

## Risk Mitigation

| Risk (from spec) | Likelihood | Impact | Mitigation | Owner |
|-------------------|------------|--------|------------|-------|
| Workflow/tool drift | Medium | High | Centralize asset naming and checksum verification in the Go release package used by CI and local tooling. | release package |
| Checksum coverage gaps | Medium | High | Verify checksum manifests against the exact expected asset matrix before metadata generation and upload. | release workflow |
| Platform matrix creep | Low | Medium | Keep the supported target list in one function and cover it with tests. | release package |

## Requirement Coverage Map

| Req ID | Component(s) | File Path(s) | Notes |
|--------|--------------|--------------|-------|
| TR-001 | target contract, packaging script, release workflow | `src/internal/release/contract.go`, `scripts/release/build-archive.sh`, `.github/workflows/release.yml` | Canonical asset naming and checksum manifest contract. |
| TR-002 | metadata generator, JSON outputs | `src/internal/release/metadata.go`, `src/tools/releasegen/main.go`, `.github/workflows/release.yml` | Package-manager inputs derive from GitHub Release asset URLs and checksums. |
| TR-003 | release validation workflow | `.github/workflows/release.yml` | Build, test, lint, and vulnerability scan run before publishing. |
| TR-004 | checksum parser, verifier, tests | `src/internal/release/checksums.go`, `src/internal/release/release_test.go`, `.github/workflows/release.yml` | Drift causes explicit failure before metadata upload. |

## Project Structure

### Source Code

```text
.github/
  workflows/
    release.yml
scripts/
  release/
    build-archive.sh
src/
  go.mod
  tools/
    releasegen/
      main.go
  internal/
    release/
      checksums.go
      contract.go
      metadata.go
      release_test.go
```

## Implementation Hints

- **[HINT-001]** Keep version normalization in the Go release package so tags with and without a leading `v` behave consistently.
- **[HINT-002]** Let the shell script ask the Go tool for archive names instead of duplicating naming logic in shell.
- **[HINT-003]** Use GitHub Release download URLs as metadata outputs so downstream channels reference immutable assets.
- **[HINT-004]** Verify checksum coverage before metadata generation and again before metadata upload to keep failure causes obvious.