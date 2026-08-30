# 5. Write a record line at the span start and at the span end

Blocked by: add-the-encoder-and-the-attribute-set.md, write-the-keyed-record-file.md, move-the-bench-home-read.md
Line: opus / medium
Rows: OT5, OT6, OT7
Writes: internal/otelrecord/processor.go (new), internal/otelrecord/processor_test.go (new), internal/otelrecord/provider.go (new)

## What to build

One span processor writes at `OnStart` and again at `OnEnd`. The start line
exists before the span ends, so an interrupted run keeps its started seam. No
published exporter writes at the start, so this processor is novel work on a
sanctioned SDK hook.

The start line carries the attribute `bench.record=start`, and it carries no
end time. A consumer filters unfinished spans by that marker. The end line
carries the complete span with `endTimeUnixNano`, so a reader derives the
elapsed time from the record alone. No consumer merges the pair. This
disposition stays open to the reviewer's veto.

The package also exports the provider constructor. A verb boundary resolves the
Bench home once through `internal/benchhome`, builds the provider there, and
threads the tracer through `context.Context`. The gate threads its run log the
same way today.

## Acceptance

- [ ] OT5: a record line exists after the span starts and before the span ends.
- [ ] OT6: the start line carries the attribute `bench.record=start` and no end time.
- [ ] OT7: the end line carries `endTimeUnixNano` and the complete span.
- [ ] the provider constructor reads the home through `internal/benchhome` and returns a tracer that a caller reads back from a context.
