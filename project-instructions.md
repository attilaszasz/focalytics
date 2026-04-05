<!-- template-version: 2 -->
# focalytics Project Instructions

## Core Principles

### I. Narrow Product Scope

focalytics MUST remain a single-purpose archive-insight CLI; cloud sync, collaboration, editing workflows, metadata write-back, and hosted dashboards stay out of scope unless the Product Document is amended first. — This protects the product promise of fast, zero-config archive analysis instead of drifting into a general photo-management suite.

### II. Local-First Safety

Runtime MUST operate offline, MUST never mutate source archives, and MUST treat files, sidecars, directory names, and metadata as untrusted input. — The product's trust model depends on privacy-first local execution and defensive parsing of user-owned archives.

### III. Modular Pipeline Design

Implementation MUST keep command handling, discovery, metadata recovery, aggregation, rendering, and release automation in separate modules under /src with explicit package boundaries. — A local modular monolith only stays maintainable if the core pipeline does not collapse into a single orchestration package.

### IV. Honest Data Reporting

The product MUST surface fallbacks, exclusions, and parse failures explicitly in logs and generated reports, and SHOULD degrade per file or metric instead of failing the whole run whenever a trustworthy partial result is possible. — Users need to understand what the report represents and where archive data was incomplete or inferred.

### V. Cross-Platform Release Quality

Changes MUST preserve macOS, Windows, and Linux viability and MUST keep GitHub Release artifacts, checksums, and package-manager metadata consistent across channels. — Cross-platform distribution is part of the product value, not an afterthought.

### VI. Agent Output Style

All agent output MUST be concise and outcome-oriented. This principle supersedes any verbose defaults.

- **Progress reports**: Facts and outcomes only — no narration, no restating the task.
- **Artifacts**: Emit required sections only — no preamble paragraphs, no summary epilogues.
- **Reasoning**: Omit unless the user asks "why" or the decision is non-obvious.
- **Errors / blockers**: State the problem, the attempted fix, and the result — nothing else.
- **Phase-boundary reports**: ≤ 5 bullet points.
- **Preserve without compressing**: Artifact template structure and required sections; explicit decision / registration / validation guidance in shared skills; delegation constraints and sub-agent role definitions; existing size limits (spec ≤ 10 KB, research ≤ 4 KB, stories ≤ 200 words).

## Technology Stack

- **Language/Runtime**: Go 1.24
- **Frameworks**: Cobra, Bubble Tea, Bubbles, go-exif, Go html/template, Go embed
- **Storage**: none
- **Infrastructure**: local only; GitHub Actions for release automation

## Testing & Quality Policy

<!-- QC extracts enforcement rules from this section. Use the keywords below so automated checks activate correctly. -->
<!-- Keywords recognised by QC: lint, static analysis, code quality, coverage, security, vulnerability, OWASP, WCAG, accessibility, benchmark, performance -->

- **Coverage Target**: 80%
- **Required QC Categories**: linting, security scanning, coverage
- **Test Strategy**: Unit, integration, golden-file report tests, and malformed-file fixture tests after implementation
- **Linting / Formatting**: golangci-lint + gofmt

## Source Code Layout

- **Policy**: ENFORCE_SRC_ROOT
- **Convention**: Source code under /src; bootstrap artifacts under /specs; repository config at root

## Development Workflow

- **Branching**: Feature branches from main with squash merge
- **Commit Convention**: Conventional Commits
- **CI Requirements**: All tests pass, lint clean, security scanning clean, and release build checks stay green before merge

## Governance

- Project instructions supersede all other documentation and practices.
- Amendments require a version bump with ISO-dated changelog entry.
- All implementations MUST pass the Instructions Check gate during planning.
- Complexity beyond these principles MUST be justified and documented.
- Registered Product Document, Technical Context Document, Deployment & Operations Document, and Project Plan paths in .github/sddp-config.md MUST be preserved unless intentionally replaced.
- The canonical implementation path is the earliest unchecked P1 epic in specs/project-plan.md unless the user explicitly reprioritizes.

**Version**: 1.0.0 | **Last Amended**: 2026-04-05
