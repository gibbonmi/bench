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

## 2026-07-15 — authored a decision map directly from same-session reviewer decisions  [open]

- **What happened:** FT95 was reviewer-prioritized ("fix asap") in the session
  that discovered it; the three load-bearing forks were put to the reviewer and
  closed in-conversation. `/bench-write-spec` requires a `decisions/<topic>.md`
  map, but `/bench-shape-idea` targets ideas needing multi-session decision
  work. I wrote `decisions/worktree-lifecycle.md` directly as the record of the
  closed forks and compiled the spec from it, instead of stopping to run a
  shape phase over decisions already made.
- **Right behavior:** unsure — the map gate's letter is satisfied (named map,
  placeholder-free Handoff), but no canonical phase produced the map.
- **Proposed rule change:** name the fast path in `/bench-write-spec`: forks
  closed by the reviewer in-session may be recorded straight into a map,
  flagged in-spec for veto, when nothing is left open.
