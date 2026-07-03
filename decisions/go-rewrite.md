# Go rewrite — move the kit's executable logic from shell to Go

Graduated from the roadmap 2026-07-03. The parked line weighed two paths —
swap grep→rg in the scripts, or rewrite the shell in Go. The bootstrap grill
rejected the rg swap (#1) and made this map the Go-rewrite decision. Evidence
gathered at bootstrap: 30 shell files / ~5.7k LOC plus a 418-line Python hook
analyzer; a full gate run costs 82s wall / 49s CPU, paid per shift iteration.

What Go would buy (assessed at bootstrap, packaging excluded): a unit-test
seam for the parsers/emitters the canary layer currently compensates for; the
shell-footgun classes of the hostile-input checklist dissolve structurally;
one-source-per-fact becomes enforceable across the CLI/hook boundary (the
gate_tree_hash and hooks_dir mirrors are live defects); one language absorbs
git-guard.py; gate loop ~2–4× faster (fork/exec churn vanishes, git subprocess
time remains); cheap runway for the five parked parser candidates. The main
non-packaging counterweight: agents currently read and hot-patch the kit as
plain text — Go inserts a compile step and makes the core opaque.

## #1: Is the grep→rg swap worth doing?

Type: Grill

### Question
The parked idea's first path: replace grep with ripgrep in the kit scripts
for speed and consistency.

### Answer
**Rejected.** Every one of the ~100+ grep call sites is a single-file
assertion (`grep -qF needle file`) or a pipe filter — zero recursive tree
searches, the only place ripgrep is faster, so the speed win is nil. And rg
would be the kit's first hard non-POSIX dependency, shipped to every linked
repo (the gate deliberately keeps even shellcheck best-effort to avoid this).
Do not reopen for consistency's sake; reopen only if recursive search
appears in a hot path.

## #2: What is candidate for rewrite, and what stays plain text regardless?

Type: Grill

### Question
Scope boundary of the language question itself.

### Answer
Candidate: executable logic only — `bin/*.sh`, the gate fragments, `.bench`
lib and hooks, and `git-guard.py` (~6.1k LOC across two languages). All
markdown content (skills, commands, BENCH.md, docs, profiles) and the JSON
adapters stay plain text under every outcome — the content surface is the
product's portable half and is never in question.

## #3: Does the kit accept a compiled core?

Type: Grill

### Question
The product-identity call, and the one that can kill the whole idea: today
"plain files are the product" — agents read the scripts as source, patch them,
and the edit takes effect immediately; the gate's parse layer (`bash -n`) and
canaries assume text. A Go core makes the executable half opaque to agents,
inserts compile between edit and effect, and changes what "reading the kit"
means for every future session. Is that identity change acceptable in
principle? If no, the map closes as rejected and #4–#6 die.

### Answer
— (open — grill pending; reviewer's call, hard stop)

## #4: How would consumers get the binary, and what does the toolchain cost?

Type: Research

### Question
The packaging reality deferred at bootstrap. npm currently ships shell that
`bench link` copies into consumer repos and hooks exec by path. For a Go
core: prebuilt per-platform binaries vs `go install` vs source build;
whether the Go toolchain becomes a consumer or CI dependency; what the
pre-push hook execs on a machine with no toolchain; what `bench link` copies
and how binary/asset version skew is detected; what the gate's parse layer
becomes (go build/vet/test in place of `bash -n`). Output: a short summary
asset with a recommended distribution shape and its hard dependencies.

### Answer
— (open — research session)

## #5: Go or no-go?

Blocked by: #3, #4
Type: Grill

### Question
The decision this map exists for: do the bootstrap-assessed benefits
(testability, footgun elimination, unification, gate speed) justify the
identity change (#3) at the distribution cost (#4)? A ~6k-LOC rewrite with
the black-box gate as the only regression net is the migration risk to
price in.

### Answer
— (open — blocked)

## #6: Rewrite scope, migration order, and testing shape

Blocked by: #5
Type: Grill

### Question
Only if #5 is go: big-bang vs strangler (seam-by-seam behind the existing
CLI dispatch); which seam first (the AXI query surface and parsers are the
highest-value, lowest-risk candidates; the gate itself migrates last, under
the old gate's watch); what the unit layer covers vs what stays black-box
gate contract; what the canary layer becomes.

### Answer
— (open — blocked)
