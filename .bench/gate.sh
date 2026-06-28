#!/usr/bin/env bash
# .bench/gate.sh — the oracle for benchkit. Exits 0 when the kit is shippable.
#
# benchkit has no test/type/lint stack; it is shell + markdown + JSON consumed by
# harnesses. "Shippable" therefore means kit conformance: scripts parse, JSON is
# valid, skills carry frontmatter, the published file list resolves, and the
# AGENTS.md index stays in sync with the skills/commands on disk.
#
# This file is NOT in package.json files[], so it never ships to consumers.
set -uo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null)" || { echo "gate: not in a git repo" >&2; exit 3; }
cd "$root"

fail=0
err() { echo "gate: $*" >&2; fail=1; }

# 1. Shell scripts parse.
for f in bin/bench.sh .claude/hooks/*.sh; do
  bash -n "$f" || err "bash syntax error in $f"
done

# 1b. Scripts the harness/CLI exec by path are executable in git. A fresh clone or
#     npm install gets the git index mode; if it is 100644 the hooks fail at runtime
#     with "permission denied" (exit 126) on every Stop and every Bash tool call.
for f in bin/bench.sh .claude/hooks/*.sh; do
  mode="$(git ls-files -s "$f" | awk '{print $1}')"
  [ -z "$mode" ] && continue   # untracked — not part of what ships
  [ "$mode" = "100755" ] || err "$f is not executable in git (mode $mode); the harness runs it as a command path"
done

# 1c. The CLI must name the gate/done files that actually exist (.bench/gate.sh,
#     .bench/done.sh). An extensionless reference routes `bench gate` to auto-detect
#     instead of the repo's oracle, so it is a hard error.
bad_refs="$(grep -oE '\.bench/(gate|done)(\.sh)?' bin/bench.sh | grep -vE '\.sh$' | sort -u || true)"
[ -z "$bad_refs" ] || err "bin/bench.sh has extensionless gate/done refs ($(echo "$bad_refs" | tr '\n' ' ')); the contract is .sh"

# 2. JSON is valid.
for f in package.json .claude/settings.json; do
  node -e 'JSON.parse(require("fs").readFileSync(process.argv[1],"utf8"))' "$f" \
    || err "invalid JSON in $f"
done

# 3. Every skill carries YAML frontmatter (first line is the --- fence).
for f in .claude/skills/*/SKILL.md; do
  [ "$(head -1 "$f")" = "---" ] || err "$f missing frontmatter"
done

# 4. Every path in package.json files[] exists (npm-pack integrity).
node -e '
  const fs = require("fs"), p = require("./package.json");
  let bad = 0;
  for (const f of p.files) if (!fs.existsSync(f)) { console.error("gate: package.json files[] missing " + f); bad = 1; }
  process.exit(bad);
' || fail=1

# 5. Kit conformance — AGENTS.md index stays in sync with disk, both directions.
#    a) every skill dir is referenced in AGENTS.md
for d in .claude/skills/*/; do
  name="$(basename "$d")"
  grep -q "$name" AGENTS.md || err "skill '$name' on disk but not referenced in AGENTS.md"
done
#    b) every skill the index names exists on disk
for name in $(grep -oE 'skills/[a-z0-9-]+/SKILL\.md' AGENTS.md | sed -E 's#skills/([a-z0-9-]+)/SKILL\.md#\1#' | sort -u); do
  [ -d ".claude/skills/$name" ] || err "AGENTS.md indexes skill '$name' with no .claude/skills/$name on disk"
done
#    c) every command file is referenced as /name in AGENTS.md
for f in .claude/commands/*.md; do
  name="$(basename "$f" .md)"
  grep -q "/$name" AGENTS.md || err "command '/$name' on disk but not referenced in AGENTS.md"
done

# 6. shellcheck — stronger shell lint, best-effort (runs only when installed).
if command -v shellcheck >/dev/null 2>&1; then
  shellcheck -S warning bin/bench.sh .claude/hooks/*.sh || err "shellcheck reported issues"
fi

if [ "$fail" -eq 0 ]; then echo "gate: green"; else echo "gate: red" >&2; fi
exit "$fail"
