# Runtime Filter Contracts

## Scan Request Contract

| Field | Purpose | Default | Notes |
|-------|---------|---------|-------|
| `ArchiveRoot` | Select scan target | required | Existing contract remains unchanged |
| `Interactive` | Select TTY-aware progress behavior | inferred | Existing contract remains unchanged |
| `IgnorePhonePhotos` | Enable phone-photo exclusion for affected analytics | `false` | New opt-in filter surface |

## Device Classification Contract

| Surface | Input | Output | Contract |
|---------|-------|--------|----------|
| `metadata.Service.Recover` | recovered camera identity metadata | `DeviceClassification` on each fact | Only trusted make/model evidence may classify `phone`; ambiguity yields `unknown` |
| `metadata.lookupCameraProfile` | normalized camera model | fallback profile or phone marker | Existing phone detection logic remains reusable but must not become the only user-facing filter evidence |

## Filtered Aggregate Contract

| Surface | Input | Output | Contract |
|---------|-------|--------|----------|
| `aggregate.Service.Aggregate` | metadata facts + `AnalysisFilter` | full-archive totals plus filtered gear/technical summaries | Timeline and total-photo views remain whole-archive in this increment |
| `aggregate.Result` | filtered and unfiltered counters | `FilteredScopeSummary` | Filtered-photo counts must match render and terminal disclosure |

## Output Disclosure Contract

| Event | Output Channel | Contract | Notes |
|------|----------------|----------|-------|
| Successful unfiltered run | stdout | report path only | Existing durable output |
| Successful filtered run | stdout + stderr/UI | stdout still reports path only; stderr/UI reports active filter and filtered count | Preserves scriptability |
| Filtered report render | HTML overview + affected sections | always-visible scope note plus section notes | Essential scope context must not hide behind optional disclosure |