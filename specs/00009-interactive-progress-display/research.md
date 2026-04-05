# Research: Interactive Progress Display
> Feature | 2026-04-05 | Purpose: inform TTY-aware progress UX design for the existing Go CLI runtime

## TTY Mode Detection
- **Decision**: Use explicit terminal detection and enable redraw-based UI only for interactive terminals.
- **Rationale**: This preserves scriptable runs and prevents control sequences from leaking into redirected output.
- **Rejected**: Always-on TUI because it breaks piping and CI logs.
- **Pitfalls**: Detecting the wrong stream or treating any stderr writer as interactive will corrupt non-TTY output.
- **Sources**: https://pkg.go.dev/golang.org/x/term, https://pkg.go.dev/os

## Bubble Tea Runtime Integration
- **Decision**: Keep the pipeline event sink UI-agnostic and add Bubble Tea as a consumer for TTY sessions.
- **Rationale**: The existing sink pattern already matches Bubble Tea's message-driven model and avoids moving pipeline logic into the UI layer.
- **Rejected**: Rewriting stage execution around direct terminal writes.
- **Pitfalls**: Coupling stage logic to terminal rendering would make non-interactive runs and tests harder to preserve.
- **Sources**: https://github.com/charmbracelet/bubbletea, https://pkg.go.dev/github.com/charmbracelet/bubbletea

## Progress Signal Granularity
- **Decision**: Publish stage lifecycle and aggregate counters rather than rendering one update per discovered file.
- **Rationale**: Coarse-grained signals keep large runs readable and reduce redraw noise while still showing forward movement.
- **Rejected**: Per-file status lines as the primary user-facing feedback path.
- **Pitfalls**: Over-publishing events can recreate the current noise problem inside the TUI and in tests.
- **Sources**: https://github.com/charmbracelet/bubbles, https://github.com/charmbracelet/bubbletea

## Summary
| Topic | Decision | Rationale |
|-------|----------|-----------|
| TTY Mode Detection | Enable interactive UI only for terminals | Keeps redirected output script-friendly |
| Bubble Tea Runtime Integration | Add Bubble Tea above the sink interface | Preserves modular runtime boundaries |
| Progress Signal Granularity | Use stage and metric events, not per-file output | Improves readability on large scans |

## Sources Index
| URL | Topic | Fetched |
|-----|-------|---------|
| https://pkg.go.dev/golang.org/x/term | TTY Mode Detection | 2026-04-05 |
| https://pkg.go.dev/os | TTY Mode Detection | 2026-04-05 |
| https://github.com/charmbracelet/bubbletea | Bubble Tea Runtime Integration | 2026-04-05 |
| https://pkg.go.dev/github.com/charmbracelet/bubbletea | Bubble Tea Runtime Integration | 2026-04-05 |
| https://github.com/charmbracelet/bubbles | Progress Signal Granularity | 2026-04-05 |
