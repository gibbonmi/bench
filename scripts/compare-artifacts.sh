#!/usr/bin/env bash
# Compare two complete release-bound artifact sets and write a deterministic record.
set -euo pipefail

left="${1:?usage: compare-artifacts.sh <first-dir> <second-dir> <record> [first-root second-root [first-final-evidence second-final-evidence]]}"
right="${2:?usage: compare-artifacts.sh <first-dir> <second-dir> <record> [first-root second-root [first-final-evidence second-final-evidence]]}"
record="${3:?usage: compare-artifacts.sh <first-dir> <second-dir> <record> [first-root second-root [first-final-evidence second-final-evidence]]}"
root="$(node -e 'const fs=require("node:fs"),path=require("node:path");process.stdout.write(path.dirname(path.dirname(fs.realpathSync(process.argv[1]))))' "${BASH_SOURCE[0]}")"
# shellcheck source=scripts/lib/search.sh
. "$root/scripts/lib/search.sh"
left_root="${4:-$root}"
right_root="${5:-$root}"
left_final="${6:-}"
right_final="${7:-}"
[[ -z "$left_final" && -z "$right_final" || -n "$left_final" && -n "$right_final" ]] || {
  printf 'reproducibility comparison requires both final-evidence directories\n' >&2
  exit 1
}
left_root_physical="$(cd "$left_root" && pwd -P)"
right_root_physical="$(cd "$right_root" && pwd -P)"
[[ "$left_root_physical" != "$right_root_physical" ]] || {
  printf 'reproducibility comparison requires isolated source roots\n' >&2
  exit 1
}
[[ -d "$left" && ! -L "$left" && -d "$right" && ! -L "$right" ]] || {
  printf 'reproducibility comparison input is not a real directory\n' >&2
  exit 1
}

names_file="$(mktemp)"
evidence_names=""
final_names=""
trap 'rm -f "$names_file" "$evidence_names" "$final_names"' EXIT
version="$(node -p 'require(process.argv[1]).version' "$root/package.json")"
node "$root/scripts/release-plan.mjs" "$root" artifact-names "$version" > "$names_file"

for directory in "$left" "$right"; do
  while IFS= read -r -d '' entry; do
    name="$(basename "$entry")"
    if ! bench_search_fixed -qx -- "$name" "$names_file"; then
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
node "$root/scripts/release-plan.mjs" "$left_root" evidence-names > "$evidence_names"
while IFS= read -r name; do
  [[ -f "$left_root/$name" && ! -L "$left_root/$name" ]] || { printf 'reproducibility comparison missing first-build release evidence: %s\n' "$name" >&2; exit 1; }
  [[ -f "$right_root/$name" && ! -L "$right_root/$name" ]] || { printf 'reproducibility comparison missing second-build release evidence: %s\n' "$name" >&2; exit 1; }
  if ! cmp -s "$left_root/$name" "$right_root/$name"; then
    printf 'reproducibility release-evidence mismatch: %s\n' "$name" >&2
    exit 1
  fi
done < "$evidence_names"

if [[ -n "$left_final" ]]; then
  [[ -d "$left_final" && ! -L "$left_final" && -d "$right_final" && ! -L "$right_final" ]] || {
    printf 'reproducibility final-evidence input is not a real directory\n' >&2
    exit 1
  }
  final_names="$names_file.final"
  while IFS= read -r -d '' entry; do
    name="$(basename "$entry")"
    [[ -f "$entry" && ! -L "$entry" ]] || { printf 'reproducibility comparison contains unsafe first-build final evidence: %s\n' "$name" >&2; exit 1; }
    printf '%s\n' "$name"
  done < <(find "$left_final" -mindepth 1 -maxdepth 1 -print0) | LC_ALL=C sort > "$final_names"
  for required in release-index.json SHA256SUMS; do
    bench_search_fixed -qx -- "$required" "$final_names" || { printf 'reproducibility comparison missing first-build final evidence: %s\n' "$required" >&2; exit 1; }
  done
  while IFS= read -r name; do
    [[ -f "$left_final/$name" && ! -L "$left_final/$name" ]] || { printf 'reproducibility comparison contains unsafe first-build final evidence: %s\n' "$name" >&2; exit 1; }
    [[ -f "$right_final/$name" && ! -L "$right_final/$name" ]] || { printf 'reproducibility comparison missing second-build final evidence: %s\n' "$name" >&2; exit 1; }
    if ! cmp -s "$left_final/$name" "$right_final/$name"; then
      printf 'reproducibility final-evidence mismatch: %s\n' "$name" >&2
      exit 1
    fi
  done < "$final_names"
  while IFS= read -r -d '' entry; do
    name="$(basename "$entry")"
    [[ -f "$entry" && ! -L "$entry" ]] || { printf 'reproducibility comparison contains unsafe second-build final evidence: %s\n' "$name" >&2; exit 1; }
    bench_search_fixed -qx -- "$name" "$final_names" || { printf 'reproducibility comparison found unexpected final evidence: %s\n' "$name" >&2; exit 1; }
  done < <(find "$right_final" -mindepth 1 -maxdepth 1 -print0)
fi

node - <<'NODE' "$left" "$record" "$left_root" "$evidence_names" "${entries[@]}"
const fs = require("node:fs");
const crypto = require("node:crypto");
const path = require("node:path");
const [directory, record, root, evidenceFile, ...names] = process.argv.slice(2);
const artifacts = names.map(name => {
  const data = fs.readFileSync(`${directory}/${name}`);
  return {name, size: data.length, sha256: crypto.createHash("sha256").update(data).digest("hex"), match: true};
});
const evidenceNames = fs.readFileSync(evidenceFile, "utf8").trimEnd().split("\n");
const evidence = evidenceNames.map(name => { const data = fs.readFileSync(path.join(root, name)); return {name, size: data.length, sha256: crypto.createHash("sha256").update(data).digest("hex"), match: true}; });
const payload = JSON.stringify({schema_version: 1, status: "green", builds: 2, artifacts, evidence}) + "\n";
fs.mkdirSync(require("node:path").dirname(record), {recursive: true});
const temporary = `${record}.tmp-${process.pid}`;
fs.writeFileSync(temporary, payload, {mode: 0o644});
fs.renameSync(temporary, record);
NODE
