# Auto-generated skills index

## Problem
The `.bench/BENCH.md` skills index is hand-maintained prose that must stay in
sync with the skills on disk. Gate checks 5a/5b catch presence drift (a skill
missing from the index, an index line with no skill) but cannot catch wording
drift, and every skill add/rename is a two-place edit. The index's trigger
phrases live nowhere the skill itself owns.

## Solution
Each indexed skill declares its own trigger phrase in frontmatter; a kit script
generates the index block between markers in `.bench/BENCH.md` from that
frontmatter. The gate verifies the committed block equals the generated one —
presence, wording, and order — and names the drifted skill. Checks 5a/5b retire
in favor of this single generate-and-compare check. Consumers are untouched:
they receive the generated `BENCH.md` as a static file.

## User stories
1. As the kit maintainer, I want each craft skill to carry an `index:` trigger
   phrase (plus optional `index-note:` appended after the path) in its
   frontmatter, so the index text has a single owner — the skill itself.
   Line: claude-fable-5 (inline, session model) / medium. Gate-observable
   conformance edits at a known seam; run inline to avoid delegation overhead.
2. As the kit maintainer, I want `.bench/skills-index.sh --write` to rewrite the
   marked index block in `.bench/BENCH.md` from frontmatter, deterministically
   (alphabetical by skill directory), so adding or renaming a skill is a
   one-place edit followed by a regenerate.
   Line: claude-fable-5 (inline, session model) / medium. Oracle logic — a
   wrong generator silently corrupts the shipped operating guide.
3. As the kit maintainer, I want the gate to run the same script in check mode
   and go red when the committed block differs from the generated one, naming
   the specific skill (missing entry, dangling entry, or wording drift), so
   drift is attributable at a glance.
   Line: claude-fable-5 (inline, session model) / medium. Gate/conformance
   logic; correctness of the oracle matters more than speed.
4. As the kit maintainer, I want a skill without an `index:` field (or with an
   empty one) to fail the gate, so an unindexable skill cannot land silently.
   Line: claude-fable-5 (inline, session model) / medium. Same check surface as
   story 3.
5. As the kit maintainer, I want the canary fixtures that guarded 5a/5b updated
   to the new error strings (and one added for wording drift), so the new check
   is proven to bite, not assumed to.
   Line: claude-fable-5 (inline, session model) / medium. The canary is the
   gate guarding the gate; its EXPECT strings are the red signals below.

## Implementation decisions
- **Frontmatter contract:** indexed skills (every `.agents/skills/*/` dir with
  no same-named file in `.agents/commands/`) must carry `index: <trigger
  phrase>`; `index-note: <suffix>` is optional and renders after the path as
  ` + <suffix>` (used by design-system's "your project's design source").
  Command-adapter skills stay unindexed, as today.
- **Generated block:** delimited by `<!-- bench:skills-index:start -->` /
  `<!-- bench:skills-index:end -->` markers inside the existing "Skills index"
  section. Line form is unchanged: `- <trigger> → \`.agents/skills/<dir>/SKILL.md\``.
  Order is alphabetical by directory — the index is a lookup table, not a
  narrative, and a deterministic order is what makes the diff meaningful.
- **One script, two modes:** `.bench/skills-index.sh`, default `--check`
  (gate mode: compare, report per-skill errors, exit nonzero), `--write`
  (regenerate the block in place via temp file + mv). The gate calls it from
  `$gate_dir` so canary inner runs check fixtures with the real script, like
  the existing gate fragments.
- **Gate surgery:** checks 5a/5b are replaced by the script call; 5c–5f are
  untouched. The script joins the `bash -n` parse list.
- **Not shipped:** the script lives beside `gate.sh`, outside `package.json`
  `files[]`, like every other kit-only gate asset.

## Testing decisions
- Good test = plant a repo state, run the real gate (or script), assert the
  verdict and the attributed message — never read the diff and judge it.
- Seam: the gate/canary seam (`tests/canary/*` fixtures + `EXPECT` strings),
  the same seam 5a/5b are tested at today. One new runtime-style behavior
  (--write idempotency) attaches to the script directly.
- Gate command: `.bench/gate.sh` (self-proving on this repo: the committed
  index must equal the generated one for the gate to be green).

### Seam diagram

    trigger: gate run (kit repo or canary inner run) / maintainer --write
        │
        ▼
    .agents/skills/*/SKILL.md frontmatter (index:, index-note:)
        │
        ▼
    ──▶ [ .bench/skills-index.sh ] ──▶ --check: per-skill errors, exit code
    .bench/BENCH.md marked block ──▶ [ ] ──▶ --write: rewritten block
                  ◀ tests attach here: canary fixtures plant drifted repos and
                    assert the gate goes red with the attributed message

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | every indexed skill carries `index:` | gate on this repo | gate red "missing index:" until all nine craft skills gain the field | field absent → attributed error |
| 2 | --write regenerates the block; second --write is a no-op | script run in a throwaway copy (runtime contract) | contract fails before the script exists | non-idempotent or non-deterministic generation fails the double-run diff |
| 3 | skill on disk, no index entry → red naming the skill | canary `unindexed-skill` | EXPECT updated to new message → "canary did not bite" until implemented | proves the presence check still bites |
| 3 | index entry, no skill on disk → red naming the entry | canary `dangling-index` | EXPECT updated to new message → "canary did not bite" until implemented | proves the dangling check still bites |
| 3 | trigger wording differs from frontmatter → red | new canary `stale-index-wording` | new fixture fails "canary did not bite" until implemented | wording drift is the capability 5a/5b never had |
| 4 | missing/empty `index:` field → red naming the skill | new canary `missing-index-field` | new fixture fails until implemented | an unindexable skill cannot land silently |
| edge of 3 | markers missing from BENCH.md → distinct red (not a silent pass) | canary (folded into `stale-index-wording` fixture family or asserted in check logic test) | script exits nonzero with "markers" message | a consumer-mangled BENCH.md must not read as green |

### Edge inventory
- error path → coverage rows above (attributed messages).
- empty/absent input → story 4 row (absent field) and markers-missing row (absent block).
- re-run idempotency → story 2 row (double --write).
- **Won't handle:** skill directory names with spaces/globs — the frontmatter
  name contract already constrains skills to kebab-case directories.
- **Won't handle:** interrupted --write — temp-file + mv is atomic; the kit repo
  is single-writer.
- **Won't handle:** symlinked invocation/cwd below root — the gate already cds
  to the repo root before sourcing fragments; the script inherits that.
- **Won't handle:** project-added skills in linked repos — the gate (and this
  script) never ship to consumers; their BENCH.md is static content.

## Out of scope
- Auto-generating the AGENTS.md/BENCH.md *command* references (5c–5d) the same
  way — a separate capability over a different surface (commands vs skills),
  ~8 edits, ~4 gate runs if wanted later.
