# Local assessment records

Use `bench assessment record --input <file>` to import one versioned JSON record.
Use `bench assessment list` to find a run. Use `bench assessment show <run-id>` to inspect all fields and derived totals.
The commands never launch a model. The user controls paid trials and changes to model defaults.

## Record fields

The `Run`, `Attempt`, `Event`, `Usage`, `Cost`, and `Rates` types define the JSON schema.
Unknown JSON fields cause a refusal. A missing numeric value is unknown. A numeric zero means measured zero.

The store uses the repository pool key and keeps records outside the worktree pool.
Records have no automatic expiry. To remove a record, name its exact stored path before explicit file removal.

Each run names its source, condition, task, held-out status, state, timestamps, attempts, evidence, and quality values.
Each attempt names its chunk, role, session, model, effort, state, timestamps, usage events, costs, and evidence.
The roles are `implementation`, `repair`, `verification`, `review`, and `diagnostic`.
States are `running`, `succeeded`, `failed`, `cancelled`, and `incomplete`.

Each usage event names its native event ID, session ID, epoch, sequence, mode, counter, and reference.
A reference names the producer and a native evidence location. A reference never causes a URL fetch or command execution.

The `usage` object inside each event holds `input_uncached`, `input_cached`, `output`, and optional `input_total`.
Declare `total_semantics` as `inclusive` or `exclusive` when a total exists.
Inclusive totals contain cached input. Exclusive totals contain only uncached input.

Declare `mode` as `delta` for event quantities or `cumulative` for snapshots from an epoch's zero baseline.
Use increasing sequence numbers within each cumulative epoch. A counter reset requires a new epoch.
Unknown or conflicting semantics cannot produce a complete total. Repeated event identities count once.

A record update preserves the existing array order. It can append entries and fill missing evidence, but cannot erase known values.

## Cost and time

The `cost.estimated` object holds three token rates, currency, unit scale, source, date, and pricing conditions.
The `cost.other` array holds applicable tool and other estimated charges. A missing applicable amount makes the estimate partial.

The `cost.actual` array holds authoritative billing charges. A token estimate never fills this array.
Each charge names its kind, amount, currency, and reference. Reports keep currencies separate.
Known components remain visible when other components are unknown.

Wall time uses run endpoints when available. Otherwise, it uses the union of complete attempt intervals.
The separate effort duration sums those attempt intervals. Concurrent review intervals do not multiply wall time.
Missing spans make observed interval coverage incomplete.

## Synthetic input

This example records unknown usage and unknown charges. It launches no model.
Replace the repository key with the target repository's pool key before import.

```json
{"version":1,"run_id":"example","repo_key":"example-123","source":"synthetic","condition":"bench","task_id":"fixture","holdout":false,"state":"running","attempts":[],"evidence":[],"quality":{}}
```

The package tests provide fixed native event fragments and unequal synthetic rates.
The prices in those tests are arithmetic fixtures. They are not current provider prices or actual account charges.
