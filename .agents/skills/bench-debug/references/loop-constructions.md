# Loop constructions

Ten ways to build the Phase 1 feedback loop, tried in roughly this order. Stop
at the first one that yields a red-capable command; drop down the list only
when the one above genuinely can't reach the bug.

1. **Failing test** at whatever seam reaches the bug — unit, integration, e2e.
2. **Curl / HTTP script** against a running dev server.
3. **CLI invocation** with a fixture input, diffing stdout against a
   known-good snapshot.
4. **Headless browser script** (Playwright / Puppeteer) — drives the UI,
   asserts on DOM, console, or network.
5. **Trace replay.** Save a real request, payload, or event log to disk;
   replay it through the code path in isolation.
6. **Throwaway harness.** Spin up a minimal subset of the system — one
   service, mocked deps — that exercises the bug path with a single call.
7. **Property / fuzz loop.** For "sometimes wrong output": run a thousand
   random inputs and look for the failure mode.
8. **Bisection harness.** If the bug appeared between two known states
   (commit, dataset, version), automate "boot at state X, check, repeat" so
   `git bisect run` can consume it.
9. **Differential loop.** Run the same input through old version and new
   (or two configs) and diff the outputs.
10. **Structured human-in-the-loop.** Last resort, when a human must click or
    observe. Keep the loop structured by scripting the human's side: a small
    script prints the exact steps to perform, waits for them to paste or
    confirm what they saw, records that output to a file, and repeats — so
    every round is captured and feeds back into diagnosis instead of
    evaporating in conversation. Write the script for the project at hand; the
    kit ships no template.
