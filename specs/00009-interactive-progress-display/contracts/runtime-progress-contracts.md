# Runtime Progress Contracts

## Output Mode Contract

| Mode | Trigger | User-Facing Behavior | Notes |
|------|---------|----------------------|-------|
| Interactive | Terminal-attached run | Bubble Tea progress UI on stderr plus final report path on stdout | Default for local terminal sessions |
| Non-interactive | Redirected or piped run | No redraw UI, no routine progress chatter, warnings/errors only, final report path on stdout | Preserves scriptability |

## Progress Sink Contract

| Surface | Input | Output | Contract |
|---------|-------|--------|----------|
| `progress.Sink.Publish` | `progress.Event` | `error` or none | Must remain UI-agnostic and safe for both TTY and non-TTY execution |
| `progress.TUISink.Publish` | `progress.Event` | `error` or none | Forwards progress events into the Bubble Tea program during interactive runs |

## Progress Event Contract

| Field / Kind | Purpose | Producer |
|--------------|---------|----------|
| `Stage` | Identifies the active pipeline stage | pipeline runner |
| `Kind=status` | Reports aggregate progress within a stage | discovery, metadata, runner |
| `Kind=warning` | Reports recoverable issues that must stay visible | discovery, metadata, runner |
| `Kind=metric` | Reports processed-versus-total counters for long-running work | metadata |
| `Kind=stage_start` / `stage_end` | Marks stage transitions for the UI | pipeline runner |

## Stdout Contract

| Event | Stdout Output | Notes |
|-------|---------------|-------|
| Successful run | report path only | Durable success output for users and scripts |
| Warning during run | none | Warnings stay on stderr/UI channel |
| Fatal failure | none | Failure remains signaled by stderr and exit code |
