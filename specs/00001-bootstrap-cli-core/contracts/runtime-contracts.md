# Runtime Contracts

## CLI Surface

| Surface | Input | Output | Notes |
|---------|-------|--------|-------|
| `focalytics <archive-root>` | filesystem path argument | process exit code, stderr diagnostics, optional stdout report path | Primary user-facing contract for one run. |

## Command Constructors

| Symbol | Input | Output | Contract |
|--------|-------|--------|----------|
| `NewRootCommand` | app dependencies, IO streams | configured root command | Must not execute package-global side effects. |
| `NewRunCommand` | runtime runner, exit policy, IO streams | configured run subcommand | Must validate arguments before invoking downstream stages. |

## Runtime Interfaces

| Interface | Method | Input | Output | Contract |
|-----------|--------|-------|--------|----------|
| `Runner` | `Run` | `context.Context`, `ScanRequest` | `RunResult`, `error` | Coordinates stages and maps failures to exit policy. |
| `ProgressSink` | `Publish` | `ProgressEvent` | `error` or none | UI-agnostic observer that must tolerate non-interactive execution. |
| `Stage` | `Run` | `context.Context`, `RunContext` | `StageResult`, `error` | One pipeline boundary for discovery, metadata, aggregation, or rendering. |

## Exit Contract

| Outcome | Exit Code Class | Notes |
|---------|-----------------|-------|
| Success | `0` | Reserved for fully successful command completion. |
| Invalid input | non-zero stable validation code | Used for missing path, unreadable root, or bad arguments. |
| Fatal runtime failure | non-zero stable runtime code | Used when execution cannot finish trustworthily. |
