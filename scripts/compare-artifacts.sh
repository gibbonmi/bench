#!/usr/bin/env bash
# Compare two complete release-bound artifact sets and write a deterministic record.
set -euo pipefail

left="${1:?usage: compare-artifacts.sh <first-dir> <second-dir> <record> [first-root second-root]}"
right="${2:?usage: compare-artifacts.sh <first-dir> <second-dir> <record> [first-root second-root]}"
record="${3:?usage: compare-artifacts.sh <first-dir> <second-dir> <record> [first-root second-root]}"
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
left_root="${4:-$root}"
right_root="${5:-$root}"
[[ -d "$left" && ! -L "$left" && -d "$right" && ! -L "$right" ]] || {
  printf 'reproducibility comparison input is not a real directory\n' >&2
  exit 1
}

names_file="$(mktemp)"
evidence_names=""
trap 'rm -f "$names_file" "$evidence_names"' EXIT
version="$(node -p 'require(process.argv[1]).version' "$root/package.json")"
node "$root/scripts/release-plan.mjs" "$root" artifact-names "$version" > "$names_file"

for directory in "$left" "$right"; do
  while IFS= read -r -d '' entry; do
    name="$(basename "$entry")"
    if ! rg -Fqx -- "$name" "$names_file"; then
      printf 'reproducibility comparison found unexpected artifact: %s\n' "$name" >&2
      exit 1
    fi
  done < <(find "$directory" -mindepth 1 -maxdepth 1 -type f -print0)
done

entries=()
while IFS= read -r name; do
  [[ -f "$left/$name" && ! -L "$left/$name" ]] || { printf 'reproducibility comparison missing first-build artifact: %s\n' "$name" >&2; exit 1; }
  [[ -f "$right/$name" && ! -L "$right/$name" ]] || { printf 'reproducibility comparison missing second-build artifact: %s\n' "$name" >&2; exit 1; }
  if ! cmp -s "$left/$name" "$right/$name"; then
    printf 'reproducibility mismatch: %s\n' "$name" >&2
    exit 1
  fi
  entries+=("$name")
done < "$names_file"

evidence_names="$names_file.evidence"
node - <<'NODE' "$left_root" > "$evidence_names"
const fs = require("node:fs");
const path = require("node:path");
const root = process.argv[2];
const requirements = JSON.parse(fs.readFileSync(path.join(root, "internal/releaseevidence/requirements.json")));
const names = ["internal/releaseevidence/requirements.json", "scripts/release-plan.json", ...requirements.records.filter(record => record.package_mode).map(record => record.path)];
if (new Set(names).size !== names.length) throw new Error("release-bound evidence inventory has duplicate paths");
process.stdout.write(names.sort((left, right) => Buffer.compare(Buffer.from(left), Buffer.from(right))).join("\n") + "\n");
NODE
while IFS= read -r name; do
  [[ -f "$left_root/$name" && ! -L "$left_root/$name" ]] || { printf 'reproducibility comparison missing first-build release evidence: %s\n' "$name" >&2; exit 1; }
  [[ -f "$right_root/$name" && ! -L "$right_root/$name" ]] || { printf 'reproducibility comparison missing second-build release evidence: %s\n' "$name" >&2; exit 1; }
  if ! cmp -s "$left_root/$name" "$right_root/$name"; then
    printf 'reproducibility release-evidence mismatch: %s\n' "$name" >&2
    exit 1
  fi
done < "$evidence_names"

node - <<'NODE' "$left" "$record" "$left_root" "${entries[@]}"
const fs = require("node:fs");
const crypto = require("node:crypto");
const path = require("node:path");
const [directory, record, root, ...names] = process.argv.slice(2);
const artifacts = names.map(name => {
  const data = fs.readFileSync(`${directory}/${name}`);
  return {name, size: data.length, sha256: crypto.createHash("sha256").update(data).digest("hex"), match: true};
});
const requirements = JSON.parse(fs.readFileSync(path.join(root, "internal/releaseevidence/requirements.json")));
const evidenceNames = ["internal/releaseevidence/requirements.json", "scripts/release-plan.json", ...requirements.records.filter(record => record.package_mode).map(record => record.path)].sort((a, b) => Buffer.compare(Buffer.from(a), Buffer.from(b)));
const evidence = evidenceNames.map(name => { const data = fs.readFileSync(path.join(root, name)); return {name, size: data.length, sha256: crypto.createHash("sha256").update(data).digest("hex"), match: true}; });
const payload = JSON.stringify({schema_version: 1, status: "green", builds: 2, artifacts, evidence}) + "\n";
fs.mkdirSync(require("node:path").dirname(record), {recursive: true});
const temporary = `${record}.tmp-${process.pid}`;
fs.writeFileSync(temporary, payload, {mode: 0o644});
fs.renameSync(temporary, record);
NODE
