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


## 2026-07-27 — a drain promoted a capture's diagnosis without checking it  [open]

- **What happened:** The drain I ran turned an `IDEAS.md` capture into roadmap row
  FT146 and carried its diagnosis forward verbatim, including the half that said
  the artifact contract tests resolve their output directory to the graded root.
  The very next phase read the tree and found that claim false — both call sites
  in `internal/contract/surface/artifact/artifact_offline_test.go` already use
  `t.TempDir()`. Half a HIGH row was fiction, and it was fiction I wrote, an hour
  earlier, from a capture I had every opportunity to check. It cost the following
  phase a false start into `/bench-write-spec` before the mismatch surfaced.
- **Right behavior:** A capture records a *symptom* reliably and a *diagnosis*
  only as the capturing session's best guess under time pressure. The drain is
  where those two get separated: promote the symptom as the row's evidence, and
  either verify the diagnosis against the tree before writing it as fact or mark
  it as the capture's unverified reading. I verified the destructive script half
  (I read the `mv`/`rm -rf` lines) and did not apply the same standard to the
  test half of the same capture.
- **Proposed rule change:** none new — FT126 already carries this exact gap, and
  its body explicitly parks it: "Not carried here: verifying what a row *claims*
  about the code … That is per-row semantic checking with no mechanical source."
  What this entry adds is a demonstrated cost, which that row lacked: the gap bit
  inside one phase transition rather than eventually. Worth attaching to FT126 as
  evidence, and worth weighing against the cheap partial version the row does not
  currently consider — not a general prober, just a `/bench-what-next` rule that a
  drained diagnosis is either checked or written as a claim rather than a fact.
