# Package and installable-surface contracts for the benchkit gate.

# Every path in package.json files[] must exist.
node -e '
  const fs = require("fs"), p = require("./package.json");
  let bad = 0;
  for (const f of p.files) if (!fs.existsSync(f)) { console.error("gate: package.json files[] missing " + f); bad = 1; }
  process.exit(bad);
' || fail=1

pack_json="$(npm_config_cache="${TMPDIR:-/tmp}/bench-npm-cache" npm pack --dry-run --json 2>/dev/null)" || {
  err "npm pack --dry-run failed"
  pack_json="[]"
}
printf '%s' "$pack_json" | node -e '
  const fs = require("fs");
  const packs = JSON.parse(fs.readFileSync(0, "utf8"));
  const files = new Set((packs[0]?.files ?? []).map(f => f.path));
  let bad = 0;
  for (const required of [
    "bin/bench.sh",
    "bin/bench-postinstall.sh",
    ".agents/commands/bench-implement-spec.md",
    ".agents/skills/bench-craft-seams/SKILL.md",
    ".agents/skills/bench-implement-spec/SKILL.md",
    ".agents/skills/bench-implement-spec/agents/openai.yaml",
    ".bench/BENCH.md",
    ".bench/BENCH-reference.md",
    ".bench/adapters/claude",
    ".bench/adapters/codex",
    ".bench/adapters/opencode",
    ".bench/hooks/stop.sh",
    ".bench/lib/resolve-bench.sh",
    ".claude/README.md",
    ".codex/hooks.json",
  ]) {
    if (!files.has(required)) {
      console.error("gate: npm package missing " + required);
      bad = 1;
    }
  }
  for (const forbidden of [".claude/settings.local.json"]) {
    if (files.has(forbidden)) {
      console.error("gate: npm package includes local-only file " + forbidden);
      bad = 1;
    }
  }
  process.exit(bad);
' || fail=1

PACK_JSON="$pack_json" node <<'NODE' || fail=1
const fs = require("fs");

const packs = JSON.parse(process.env.PACK_JSON || "[]");
const files = new Set((packs[0]?.files ?? []).map(f => f.path));
let bad = 0;
const err = msg => { console.error("gate: " + msg); bad = 1; };

const repoOnlyPaths = ["projects/", "specs/", "decisions/", "tests/"];
const packageClaim = /\b(ship|ships|shipped|shipping|package|packaged|tarball|installable|included|includes)\b/i;
const packageClaimHeading = /^#{1,6}\s+.*\b(ship|ships|shipped|shipping|package|packaged|tarball|installable|included|includes|surfaces?)\b/i;
const explicitRepoOnly = /\b(repo-only|development context|local development|not shipped|not in the npm package|not part of the npm package)\b/i;

for (const file of [...files].filter(f => f.endsWith(".md")).sort()) {
  if (!fs.existsSync(file)) continue;
  const lines = fs.readFileSync(file, "utf8").split(/\n/);
  let inPackageClaimSection = false;
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (/^#{1,6}\s+/.test(line)) inPackageClaimSection = packageClaimHeading.test(line);
    if (!(inPackageClaimSection || packageClaim.test(line)) || explicitRepoOnly.test(line)) continue;
    for (const repoOnlyPath of repoOnlyPaths) {
      if (line.includes(repoOnlyPath)) {
        err(`${file}:${i + 1} claims repo-only path '${repoOnlyPath}' is shipped/package content; label it repo-only development context`);
      }
    }
  }
}

process.exit(bad);
NODE
