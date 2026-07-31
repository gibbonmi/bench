# Sweep kit guidance for unbound model tokens

Blocked by: Resolve the line through the harness matrix

## What to build

One new conformance owner sweeps kit guidance prose so a hard-coded binding
cannot rot back in where nothing grades it. It discovers entries under
`.agents/commands` and `.agents/skills` at run time, sorts them, traverses every
regular file, and rejects a non-regular discovered entry before any read. It
rejects a literal matching the tier-model token grammar that names no cell in
the parsed matrix, and any retired `BENCH_TIER_*` or `BENCH_ALIAS_*` schema key,
naming the file, the token, and the binding source. The same change rewrites
`craft-line`, `/bench-setup-repo`, and the profile so they declare and recommend
in the current harness's family — the rule that produced this bite cannot
regenerate it.

Covers stories 7 and 9.

## Acceptance

- [ ] Every regular file discovered under the two directories is scanned,
      including one newly added after the check was written.
- [ ] One unbound model literal inserted into an otherwise clean copy emits its
      file-and-token diagnostic, with no other diagnostic present.
- [ ] Each of the six retired schema keys in guidance prose emits its own
      diagnostic.
- [ ] Every bound matrix token and ordinary non-model prose token stays accepted.
- [ ] Each of four non-regular entry kinds — FIFO, character device, socket, and
      symlink — is rejected before it is read, with a FIFO having no writer
      returning the discovery diagnostic without blocking.
- [ ] `craft-line`, `/bench-setup-repo`, and the profile carry no retired schema
      key and name the matrix.
