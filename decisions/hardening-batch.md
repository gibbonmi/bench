# Hardening Batch (FT34)

## #1: What is in the batch, and what stays out?

Blocked by: —
Type: Grill

### Question
The assessment's cheap-hardening list is mechanical but spans four areas.
Confirm the list and the boundaries so the batch doesn't grow.

### Answer
Four slices, all small: (a) SafeToken corpus gains the newline-class rows
(trailing `\n`, embedded `\n`, CR) proving the `$`-anchor property the
grammar most needs; (b) `stophook.Run` gets a seam test covering the gate
exec path — green verdict, red verdict, and the rc==3 no-gate branch that
blocks the stop — via a stub gate script in a temp repo; (c) README gains the
prerequisites it omits (Go-or-node for the from-clone path, Windows
unsupported / WSL2 note, the nvm/asdf global-install caveat) and the
`.bench/` layout block gains the shipped pieces it omits; (d) the three
longest skill descriptions (synthesis, adr, seams) drop their redundant third
trigger clause — description-only trims, no body changes. Out of scope:
anything touching enforcement semantics (that's [FT27]), status output
(that's [FT30]), and any further README restructuring.

## Handoff

1. **Module boundaries.** (a) `internal/modelid` corpus; (b) `internal/
   stophook` tests; (c) README; (d) three SKILL.md description lines plus
   their index rows in BENCH-reference.
2. **Contracts.** No behavior changes anywhere — new tests pin existing
   behavior; docs state existing facts; descriptions shorten without losing
   any trigger situation they uniquely own.
3. **Deep vs thin.** All thin; the only judgment is (d), which is
   leverage-tier prose.
4. **Black-box assertables.** Corpus rows red if the grammar regresses on
   newlines; stophook test red if rc mapping or the rc==3 block changes;
   skills-index `--check` stays green after description edits.
5. **Gate attachment.** (a) and (b) are gate-observed unit/seam tests; (c)
   and (d) ride the docs conformance scans.
6. **Hostile-input owners.** The corpus rows are the hostile input (newline
   classes); stophook's stub gate owns the missing/non-executable gate case.
7. **Uncertainty flags.** n/a — mechanical.
8. **Rejected alternatives.** Folding these into their subject-area specs
   (each is too small to carry a spec alone; the batch is the right vehicle).
9. **Domain watch-outs.** Skill descriptions are always-loaded context — (d)
   is a token-cost cut, so verify the trimmed descriptions still trigger by
   containing the concrete quoted phrases sessions actually hit.

Dependency order: n/a — single spec; slice (d) authored at top tier, (a)–(c)
delegable.
