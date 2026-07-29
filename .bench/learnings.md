# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-what-next` verdicts every open entry in its reviewed batch
diff: work-shaped and rule-shaped entries become roadmap items (rule-shaped ones
built later under the synthesis discipline), the rest are dismissed with one line
of why. A resolved entry leaves this file, and its verdict is recorded in the
drain commit. The journal holds open entries only; history lives in git. Never
rewrite a kit rule yourself — that is the whole point of capturing here instead.

Format per entry. Heading: `## YYYY-MM-DD — short title  [open]`

- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

- 2026-07-29  Treated a session-opening venue directive ("I want this session to
  build it with a codex reviewer") as a batch approval and entered
  /bench-implement-spec right after emitting the spec approval table, without
  waiting for sign-off. Right behavior: a directive that names the venue and
  reviewer setup for a build is a plan for the session, not approval of a spec
  that does not yet exist — spec sign-off stays a hard stop unless the reviewer
  has approved a batch plan in terms that cover unseen specs ("roll the
  roadmap"). Proposed rule: none needed; BENCH.md already states it — the miss
  was classification, not a rule gap.
- 2026-07-29  Landing a ticket's `bench commit` while a concurrent write-delegate ran heavy focused contract tests produced a transient full-gate red (phase output lost to a `tail` filter; serial rerun green). Right behavior: treat the whole-tree gate as a serialized resource against delegates' *test* phases too, not just other gates — either stagger the landing until the delegate's test phase is done, or capture the full gate log before retrying so a real red is attributable. Proposed rule: none yet; second occurrence should promote it to craft-delegate's serialization clause.
- 2026-07-29  Second unattributed transient full-gate red at a ticket landing (no concurrent delegate this time; serial rerun green, phase name lost to a tail filter both times). Right behavior: pipe every `bench commit`/`bench gate` landing run to a log file so a red is attributable before retrying. Proposed rule: /bench-implement-spec's landing step should say "capture the full gate log"; two data points now.
- 2026-07-29  The two unattributed transient gate reds above are now attributed: `TestExecuteDeadlineRecordsDistinctTimeout` (internal/gate/resource_bounds_test.go) races a 50 ms stubbed gateTimeout against a 500 ms parent context; under the gate's own concurrent contract phase (surface/artifact ~119 s) on WSL2 the parent context can expire first (GateExit 130, not 124). Solo: 10/10 green. Pre-existing, load-coupled, not introduced by ft91-gate-fastpath (first occurrence predates the short-circuit landing). Fix requires editing a gate test's timing margin — reviewer decision; flagged in the build's exit report. ~3 reds in ~12 gate runs today.
