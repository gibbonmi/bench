# Refuse a State that breaks the grammar

Blocked by: keep-the-next-command-and-refuse-a-stale-state.md, read-a-legacy-document-as-main.md
Writes: internal/handoffdoc, internal/handoff/handoff.go, internal/handoff/render.go, internal/handoff/render_test.go, internal/handoff/state_scan.go, internal/handoff/state_scan_test.go
Covers: HS27, HS28, HS29, HS31, HS32

## What to build

Repair for review findings C1, C2, and C4. Verify the premise first.
`splitSections` carries the fence state across sections. `bench handoff`
accepts any State bytes. `scanState` splits on newlines with no fence
tracking.

Then make `Parse` refuse a fence still open at the end of the file. The
refusal names the file and the line that opened the fence. Make
`bench handoff` refuse the owned section's State before the write in two
cases: an unterminated fence, or a line that opens a level-two heading
outside a fence. Print the line in each refusal.

Make the State scan skip fenced lines. When both `cat-file` probes fail
and git names the token ambiguous, refuse with an ambiguity reason.

## Acceptance

- [ ] A State that opens a fence and never closes it exits 1, prints the line, and leaves the file unchanged.
- [ ] `Parse` refuses a document with an open fence at end of file, naming the file and the opening line.
- [ ] A State line `## Open questions` outside a fence exits 1 and prints the line.
- [ ] A real off-ancestry commit inside a fenced block in State exits 0.
- [ ] An ambiguous 7-hex abbreviation in State exits 1 with an ambiguity reason.
- [ ] Self-probe: drop the fence skip in the scan, and report the fenced-commit test red.
