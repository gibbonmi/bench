# Local assessment records

Use `bench assessment record --input <file>` to import one versioned JSON record.
Use `bench assessment list` to find a run. Use `bench assessment show <run-id>` to inspect all fields and derived totals.
The commands never launch a model. The user controls paid trials and changes to model defaults.

## Record fields

The `Run`, `Attempt`, `Event`, `Usage`, `Cost`, and `Rates` types define the JSON schema.
Unknown JSON fields cause a refusal. A missing numeric value is unknown. A numeric zero means measured zero.

The store uses the repository pool key and keeps records outside the worktree pool.
Records have no automatic expiry. To remove a record, name its exact stored path before explicit file removal.

Quality entries use `Measure`: a nullable value and its producer/native reference.
Known timestamps require `time_reference`; applicable charges require references even when amounts are unknown.

Each run names its source, condition, task, held-out status, state, timestamps, attempts, evidence, and quality values.
Each attempt names its chunk, role, session, model, effort, state, timestamps, usage events, costs, and evidence.
[Roles and States](vocabulary.go) define the accepted versioned vocabulary.

Each usage event names its native event ID, session ID, epoch, sequence, mode, counter, and reference.
A reference names the producer and a native evidence location. A reference never causes a URL fetch or command execution.

The `usage` object inside each event holds `input_uncached`, `input_cached`, `output`, and optional `input_total`.
Declare `total_semantics` as `inclusive` or `exclusive` when a total exists.
Inclusive totals contain cached input. Exclusive totals contain only uncached input.

[DeltaMode and CumulativeMode](vocabulary.go) define the counter mode values and their meanings.
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

## Selected native evidence

`record --input` accepts optional `bench_inputs` and `harness_inputs` on the normalized Run.
[The collection types](collection.go) define selectors and their explicit attempt, chunk, and role mapping.
The expected assignment must match the selected trace's assignment attribute and each census event's assignment.
Historical worktree spans can use their assignment subject. Published Git subjects never identify assignments.

Selection uses identities. It does not discover sessions from timestamps.

`trace_ids` contains selectors whose `id` is a trace ID. Finished spans become observed attempt intervals and provenance references.
The reader streams the existing OTEL file, retaining bounded selected evidence. Start records with a completed counterpart do not remain unfinished.
Nested and disjoint spans use the same interval-union calculation as other attempt timing.

`census_event_ids` contains selectors whose `id` combines the assignment ID, a colon, and the one-based record position.
The existing census owner decodes each selected row. Each observed raw command has a referenced measure and its native verb head.
Before release removes the assignment's census file, collect census evidence.

[The harness mapper](harness.go) owns the supported native format and counter mapping.
It accepts an explicitly supplied, bounded Codex token-count JSON fragment; it does not read a whole session automatically.
Inclusive native input subtracts cached input. Reasoning output is already part of output and is not added again.
Unsupported counter semantics refuse the import. Missing or malformed native fragments remain diagnostics with unknown measurements.

The `measures` map retains other explicitly supplied numeric observations with a reference per measure.
When the selected producer does not expose these measures, they remain unknown:

- Tools
- Read paths
- Turns
- Iterations

Use native evidence references for read paths. Do not infer a read from arbitrary shell command text.

The `intervals` array retains observed start/end pairs and their native references. It does not replace explicit attempt endpoints.
The `diagnostics` array preserves incomplete-input evidence. Diagnostic references remain part of the record when later evidence arrives.
Unknown native coverage never becomes a zero-valued measurement.
