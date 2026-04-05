# Research: Bootstrap CLI Core
> Feature E001 | Date: 2026-04-05 | Purpose: Lock the bootstrap architecture, testing baseline, and QC toolchain for the first focalytics implementation wave.

## CLI Package Layout
- **Decision**: Keep a thin executable under /src with constructor-based Cobra commands and internal runtime packages for command flow, pipeline contracts, and progress orchestration.
- **Rationale**: This preserves the modular monolith boundary from the SAD without letting Cobra wiring become the business layer.
- **Rejected**: Package-global `init` registration and a fat `main` package because they hide dependencies and increase coupling.
- **Pitfalls**: Do not export pipeline internals prematurely or let command files own scan orchestration.
- **Sources**: https://go.dev/doc/modules/layout, https://cobra.dev/docs/how-to-guides/working-with-commands/

## Progress Integration Boundary
- **Decision**: Model progress as typed runtime events with a sink interface; keep Bubble Tea optional and adapter-scoped.
- **Rationale**: The runtime must work in CI, tests, and non-interactive shells without terminal coupling.
- **Rejected**: Embedding orchestration state in Bubble Tea update loops because it makes the UI the source of truth.
- **Pitfalls**: Avoid TTY-only assumptions and keep event payloads UI-agnostic.
- **Sources**: https://pkg.go.dev/github.com/charmbracelet/bubbletea, https://pkg.go.dev/github.com/charmbracelet/bubbles/progress

## Foundational Testing Baseline
- **Decision**: Use stdlib-first package tests with `go test` for unit and integration-tagged runs, plus small CLI test helpers built around temp directories and injected streams.
- **Rationale**: The bootstrap epic needs stable regression coverage without introducing test framework overhead.
- **Rejected**: Third-party test runners because native Go tools already cover the required baseline.
- **Pitfalls**: Avoid global mutable state and full end-to-end UI dependencies in foundational tests.
- **Sources**: https://pkg.go.dev/testing, https://go.dev/doc/tutorial/add-a-test

## QC Toolchain
- **Decision**: Standardize on `go test` for unit and integration tiers, `golangci-lint` plus `gofmt` for lint/static analysis, `govulncheck` for required vulnerability scanning, and native Go coverage tooling with `-coverpkg=./...`.
- **Rationale**: This satisfies the project quality policy with low-friction, widely adopted Go-native tooling.
- **Rejected**: Custom wrappers and security-only linting because they either add overhead or miss reachable dependency vulnerabilities.
- **Pitfalls**: Do not rely on cached CI tests, plain `-cover` without cross-package coverage, or `go test`'s built-in vet pass as the only static analysis stage.
- **Sources**: https://pkg.go.dev/cmd/go#hdr-Test_packages, https://go.dev/doc/security/vuln/

## Summary
| Topic | Decision | Rationale |
|-------|----------|-----------|
| CLI Package Layout | Thin /src entrypoint plus internal runtime packages | Keeps command wiring separate from domain orchestration |
| Progress Integration Boundary | Typed progress events with optional Bubble Tea adapter | Preserves non-interactive execution and testability |
| Foundational Testing Baseline | Native `go test` with temp-dir helpers and tags | Provides enough structure without extra framework cost |
| QC Toolchain | `gofmt`, `golangci-lint`, `govulncheck`, native coverage | Matches project instructions and Go ecosystem norms |

## Sources Index
| URL | Topic | Fetched |
|-----|-------|---------|
| https://go.dev/doc/modules/layout | CLI Package Layout | 2026-04-05 |
| https://cobra.dev/docs/how-to-guides/working-with-commands/ | CLI Package Layout | 2026-04-05 |
| https://pkg.go.dev/github.com/charmbracelet/bubbletea | Progress Integration Boundary | 2026-04-05 |
| https://pkg.go.dev/github.com/charmbracelet/bubbles/progress | Progress Integration Boundary | 2026-04-05 |
| https://pkg.go.dev/testing | Foundational Testing Baseline | 2026-04-05 |
| https://go.dev/doc/tutorial/add-a-test | Foundational Testing Baseline | 2026-04-05 |
| https://pkg.go.dev/cmd/go#hdr-Test_packages | QC Toolchain | 2026-04-05 |
| https://go.dev/doc/security/vuln/ | QC Toolchain | 2026-04-05 |

