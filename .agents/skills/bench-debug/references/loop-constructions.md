# Loop constructions

This list gives ten ways to build the Phase 1 feedback loop, in roughly this
order. Stop at the first construction that yields a red-capable command. Move
down the list only when the construction above cannot reach the bug.

1. **Failing test** at whatever seam reaches the bug — unit, integration, e2e.
2. **Curl / HTTP script** against a running dev server.
3. **CLI invocation** with a fixture input, diffing stdout against a
   known-good snapshot.
4. **Headless browser script** (Playwright / Puppeteer) — drives the UI,
   asserts on DOM, console, or network.
5. **Trace replay.** Save a real request, payload, or event log to disk.
   Replay it through the code path in isolation.
6. **Throwaway harness.** Spin up a minimal subset of the system: one
   service, with mocked dependencies. This subset exercises the bug path
   with a single call.
7. **Property / fuzz loop.** Use this construction for a "sometimes wrong
   output" symptom. Run a thousand random inputs. Look for the failure mode.
8. **Bisection harness.** If the bug appeared between two known states
   (commit, dataset, version), automate "boot at state X, check, repeat".
   Then `git bisect run` can consume it.
9. **Differential loop.** Run the same input through the old version and the
   new version (or two configs). Diff the outputs.
10. **Structured human-in-the-loop.** Last resort: use this construction only
    when a human must click or observe. Keep the loop structured by
    scripting the human's side of the loop. The script prints the exact
    steps to perform, then waits for the human to paste or confirm what they
    saw. It records that output to a file, repeats, and feeds every round
    back into diagnosis so the conversation does not lose it. Write the
    script for the project at hand. The kit ships no template.
