# Doc-conformance contracts for the benchkit gate: living docs name real commands,
# cold-pickup docs list real CLI commands, command-first anchors, acceptance and
# edge-coverage anchors, and the README layout. Sourced by gate.sh so it shares
# $root, $gate_dir, err(), and $fail with the other fragments.

#    g) living docs must name commands that exist now. Historical specs/maps may
#       mention old command names only when explicitly marked on a line by itself, but
#       the cold pickup surface, live maps, and command/skill bodies must not point
#       agents at dead slash commands or Codex $bench-* adapters.
node <<'NODE' || fail=1
const fs = require("fs");
const path = require("path");

let bad = 0;
const err = msg => { console.error("gate: " + msg); bad = 1; };

const commandsDir = ".agents/commands";
const validSlash = new Set();
if (fs.existsSync(commandsDir)) {
  for (const f of fs.readdirSync(commandsDir)) {
    if (f.endsWith(".md")) validSlash.add("/" + path.basename(f, ".md"));
  }
}
for (const external of ["/model"]) validSlash.add(external);

const validCodex = new Set(
  [...validSlash]
    .filter(token => token.startsWith("/bench-"))
    .map(token => "$" + token.slice(1))
);

const files = [];
const addFile = file => {
  if (fs.existsSync(file)) files.push(file);
};
const walk = dir => {
  if (!fs.existsSync(dir)) return;
  for (const ent of fs.readdirSync(dir, { withFileTypes: true })) {
    const file = path.join(dir, ent.name);
    if (ent.isDirectory()) {
      walk(file);
    } else if (
      ent.name === "SKILL.md" ||
      /\.(md|ya?ml|json|sh)$/.test(ent.name)
    ) {
      files.push(file);
    }
  }
};

for (const f of [
  "README.md",
  "AGENTS.md",
  ".bench/BENCH.md",
  ".bench/learnings.md",
  "CONTEXT.md",
  "HANDOFF.md",
  "CHANGELOG.md",
]) addFile(f);
walk("specs");
walk("decisions");
walk(".agents");

const knownStale = new Set([
  "/resynthesize",
  "/spec",
  "/grill",
  "/start-ideation",
  "/setup",
  "/build",
  "/prep-shift",
  "/fix-bug",
  "/verify-gate",
  "/map",
  "/diagnose",
  "/review",
  "/verify",
  "/shift",
]);
const historicalMarker = /^<!-- command-currency: historical -->$/m;
const slashRef = /(^|[\s([`"'])\/([A-Za-z][A-Za-z0-9_-]*[A-Za-z0-9])/g;
const codexRef = /(^|[\s([`"'])\$([A-Za-z][A-Za-z0-9_-]*[A-Za-z0-9])/g;

for (const file of [...new Set(files)].sort()) {
  let text = fs.readFileSync(file, "utf8");
  if (historicalMarker.test(text)) continue;
  if (file === ".bench/learnings.md") {
    text = text.split("<!-- entries below -->")[0];
  }
  if (file === "CHANGELOG.md") {
    text = text.split(/\n## /)[0];
  }
  const lines = text.split(/\n/);
  for (let i = 0; i < lines.length; i++) {
    let m;
    slashRef.lastIndex = 0;
    while ((m = slashRef.exec(lines[i])) !== null) {
      const token = "/" + m[2];
      if (!validSlash.has(token) && (token.startsWith("/bench-") || knownStale.has(token))) {
        err(`stale command reference ${token} in ${file}:${i + 1}`);
      }
    }
    codexRef.lastIndex = 0;
    while ((m = codexRef.exec(lines[i])) !== null) {
      const token = "$" + m[2];
      if (token.startsWith("$bench-") && !validCodex.has(token)) {
        err(`stale Codex adapter reference ${token} in ${file}:${i + 1}`);
      }
    }
  }
}
process.exit(bad);
NODE

#    h) shipped cold-pickup docs that list CLI commands must list the real subcommands
#       from bin/bench.sh. HANDOFF.md ships in the npm package, and .bench/BENCH.md is
#       the operating guide installed into consumer repos.
node <<'NODE' || fail=1
const fs = require("fs");

let bad = 0;
const err = msg => { console.error("gate: " + msg); bad = 1; };
const bench = fs.readFileSync("bin/bench.sh", "utf8");
const commands = [...bench.matchAll(/^  ([a-z][a-z-]*)\)\s/gm)].map(m => m[1]).sort();
for (const file of ["HANDOFF.md", ".bench/BENCH.md"]) {
  if (!fs.existsSync(file)) continue;
  const text = fs.readFileSync(file, "utf8");
  for (const cmd of commands) {
    if (!text.includes(`bench ${cmd}`)) err(`${file} does not list CLI command 'bench ${cmd}'`);
  }
}
process.exit(bad);
NODE

#    h2) the project profile's AXI seam list must name the wave-2 parsers — the
#        seam list is where a cold session learns the query surface exists.
grep -qF 'bench diff' projects/benchkit.md 2>/dev/null \
  || err "projects/benchkit.md does not name bench diff on the AXI seam"
grep -qF 'bench coverage' projects/benchkit.md 2>/dev/null \
  || err "projects/benchkit.md does not name bench coverage on the AXI seam"

#    i) command-first usability anchors. The reviewer-facing README must start from
#       the command path, and every command phase must orient the reviewer at entry
#       and hand them off at exit. Keep this structural; prose quality is review.
node <<'NODE' || fail=1
const fs = require("fs");
const path = require("path");

let bad = 0;
const err = msg => { console.error("gate: " + msg); bad = 1; };

if (!fs.existsSync("README.md")) {
  err("README.md missing");
} else {
  const readme = fs.readFileSync("README.md", "utf8");
  const firstH2 = (readme.match(/^## .+$/m) || [""])[0];
  if (firstH2 !== "## Reviewer quick start") {
    err(`README first H2 is '${firstH2 || "(none)"}'; expected '## Reviewer quick start'`);
  }
}

const commandsDir = ".agents/commands";
if (fs.existsSync(commandsDir)) {
  for (const f of fs.readdirSync(commandsDir).filter(f => f.endsWith(".md")).sort()) {
    const file = path.join(commandsDir, f);
    const text = fs.readFileSync(file, "utf8");
    if (!/^## Entry orientation$/m.test(text)) err(`${file} missing Entry orientation`);
    if (!/^## Exit handoff$/m.test(text)) err(`${file} missing Exit handoff`);
  }
}

process.exit(bad);
NODE

#    j) acceptance coverage maps are now part of the feature-build workflow. These
#       anchors are intentionally structural: they prove the command/skill surfaces
#       still carry the contract, while semantic completeness stays a review/dogfood
#       responsibility.
require_anchor() {
  file="$1"
  needle="$2"
  if [ ! -f "$file" ]; then
    err "acceptance coverage anchor file missing: $file"
  elif ! grep -qF "$needle" "$file"; then
    err "$file missing acceptance coverage anchor: $needle"
  fi
}
require_anchor ".agents/commands/bench-write-spec.md" "acceptance coverage map"
require_anchor ".agents/commands/bench-write-spec.md" "why it catches the failure"
require_anchor ".agents/commands/bench-write-spec.md" "red signal"
require_anchor ".agents/skills/bench-craft-tdd/SKILL.md" "acceptance row"
require_anchor ".agents/skills/bench-craft-tdd/SKILL.md" "not TDD-able"
require_anchor ".agents/skills/bench-craft-tdd/SKILL.md" "call count"
require_anchor ".agents/commands/bench-implement-spec.md" "coverage table"
require_anchor ".agents/commands/bench-implement-spec.md" "already covered"
require_anchor ".agents/commands/bench-implement-spec.md" "turning red-to-green"
require_anchor ".agents/commands/bench-review-implementation.md" "acceptance coverage map"
require_anchor ".agents/commands/bench-review-implementation.md" "mapped behavior"
require_anchor ".agents/commands/bench-review-implementation.md" "bench diff"
#    k) the final-check phase must name the actual gate resolution chain
#       (.bench/gate.sh -> $BENCH_GATE -> auto-detect); the profile documents the
#       gate but never selects it, and the doc must not hide that seam.
require_anchor ".agents/commands/bench-final-check.md" ".bench/gate.sh"
require_anchor ".agents/commands/bench-final-check.md" "BENCH_GATE"
#    l) edge-case coverage contracts. The spec phase must generate edge cases (edge
#       inventory + won't-handle lines, fed by the profile's hostile-input checklist),
#       craft-tdd must treat stories as the breadth floor, craft-seams must bind seam
#       height to failure-mode observability, and review must hunt coverage gaps via
#       the Coverage axis. Structural anchors only; semantics stay with review/dogfood.
require_anchor ".agents/commands/bench-write-spec.md" "seam diagram"
require_anchor ".agents/commands/bench-write-spec.md" "tests attach here"
require_anchor ".agents/commands/bench-write-spec.md" "edge inventory"
require_anchor ".agents/commands/bench-write-spec.md" "Won't handle"
require_anchor ".agents/commands/bench-write-spec.md" "hostile-input checklist"
require_anchor ".agents/skills/bench-craft-tdd/SKILL.md" "floor, not the ceiling"
require_anchor ".agents/skills/bench-craft-seams/SKILL.md" "failure modes"
require_anchor ".agents/skills/bench-craft-seams/SKILL.md" "structure.budgets"
require_anchor ".agents/commands/bench-review-implementation.md" "## Coverage"
require_anchor ".agents/commands/bench-review-implementation.md" "Coverage axis"
require_anchor ".agents/commands/bench-setup-repo.md" "hostile-input checklist"
require_anchor "projects/benchkit.md" "hostile-input checklist"
#    l2) oracle-authoring anchors. The two phases where gate authoring actually
#        happens must route through craft-gate — scaffolding a first gate, and
#        stopping on a check that itself looks wrong.
require_anchor ".agents/commands/bench-setup-repo.md" "craft-gate"
require_anchor ".agents/commands/bench-final-check.md" "craft-gate"
#    l3) review-judgment anchors. The review phase routes its axis charges to
#        craft-review (one source for the judgment), and the skill must keep the
#        adversarial coverage rule that makes the Coverage axis bite.
require_anchor ".agents/commands/bench-review-implementation.md" "craft-review"
require_anchor ".agents/skills/bench-craft-review/SKILL.md" "an edge nobody decided"
#    l4) delegation anchors. The one phase that orchestrates delegates must route
#        spawning through craft-delegate, and the skill must keep the done-claim
#        rule that makes delegation safe.
require_anchor ".agents/commands/bench-review-implementation.md" "craft-delegate"
require_anchor ".agents/skills/bench-craft-delegate/SKILL.md" "a claim, not a result"
#    l5) workflow-exit anchors. The paths that used to dead-end must keep their
#        routes: a capped/unmet build, a superseded spec, and the debug phase's
#        repro-commit-before-shift ordering.
require_anchor ".agents/commands/bench-implement-spec.md" "When the build stops short"
require_anchor ".agents/commands/bench-write-spec.md" "Superseded by"
require_anchor ".agents/commands/bench-debug.md" "before launching the shift"

grep -qF 'session-start.sh' README.md || err "README layout omits .bench/hooks/session-start.sh"
grep -qF 'bench.sh' README.md || err "README layout omits the real bin/bench.sh filename"
grep -qF 'benchkit.md' README.md || err "README layout omits projects/benchkit.md"
! grep -qF '│   └── bench                 #' README.md || err "README layout still names bin/bench instead of bin/bench.sh"

#    m) the skills index is generated: --write produces a block that --check
#       accepts, and a second --write changes nothing. Behavioral contract for
#       .bench/skills-index.sh (the check side is canary-covered).
si_tmp="$(mktemp -d)"
(
  set -u; cd "$si_tmp"
  mkdir -p .bench .agents/skills/zeta-skill
  printf -- '---\nname: zeta-skill\ndescription: d\nindex: doing zeta things\n---\n' > .agents/skills/zeta-skill/SKILL.md
  printf '# Guide\n\n<!-- bench:skills-index:start -->\n<!-- bench:skills-index:end -->\n' > .bench/BENCH.md
  if bash "$root/.bench/skills-index.sh" --check >/dev/null 2>&1; then
    echo "check passed on an empty index block"; exit 1
  fi
  bash "$root/.bench/skills-index.sh" --write || { echo "--write failed"; exit 1; }
  grep -qF -- '- doing zeta things → `.agents/skills/zeta-skill/SKILL.md`' .bench/BENCH.md \
    || { echo "--write did not generate the entry from frontmatter"; exit 1; }
  bash "$root/.bench/skills-index.sh" --check >/dev/null 2>&1 || { echo "check red right after --write"; exit 1; }
  before="$(cat .bench/BENCH.md)"
  bash "$root/.bench/skills-index.sh" --write || { echo "second --write failed"; exit 1; }
  [ "$before" = "$(cat .bench/BENCH.md)" ] || { echo "--write is not idempotent"; exit 1; }
) || err "skills-index generate/verify contract failed"
rm -rf "$si_tmp"

#    n) acceptance coverage maps parse. Any specs/*.md carrying the map heading
#       must have the canonical header, five non-empty cells per row, and story
#       references that resolve (an integer within the spec's numbered stories,
#       or an edge… reference). Pre-convention specs opt out with
#       <!-- coverage-map: historical -->. Cell semantics stay with review.
#       The validation lives in `bench coverage --check` — one parser for the
#       convention (spec second-wave-parsers story 7). The CLI is resolved from
#       the gate script's own tree, never the working tree's copy: canary inner
#       runs execute the real gate against minimal fixture trees with no CLI.
if [ -d specs ]; then
  for cov_f in specs/*.md; do
    [ -f "$cov_f" ] || continue
    cov_out="$(bash "$gate_dir/../bin/bench.sh" coverage --check "$cov_f")" && cov_rc=0 || cov_rc=$?
    if [ "$cov_rc" -ne 0 ]; then
      if [ -n "$cov_out" ]; then
        while IFS= read -r cov_line; do
          err "${cov_line#error: }"
        done <<<"$cov_out"
      else
        err "$cov_f coverage --check failed (exit $cov_rc) with no message"
      fi
    fi
  done
fi
