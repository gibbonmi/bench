During a ticket run focused checks at its declared seams. The path-scoped `bench commit`
is the only per-ticket full-project-gate boundary and commits atomically only on green.
If it goes red, repair from that output and retry; the normal green path is one full
gate. `/bench-final-check` still runs the final full gate over the composed feature.
