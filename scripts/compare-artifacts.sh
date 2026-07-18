#!/usr/bin/env bash
# Compare two complete release-bound artifact sets and write a deterministic record.
set -euo pipefail

left="${1:?usage: compare-artifacts.sh <first-dir> <second-dir> <record> }"
right="${2:?usage: compare-artifacts.sh <first-dir> <second-dir> <record> }"
record="${3:?usage: compare-artifacts.sh <first-dir> <second-dir> <record> }"
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
[[ -d "$left" && ! -L "$left" && -d "$right" && ! -L "$right" ]] || {
  printf 'reproducibility comparison input is not a real directory\n' >&2
  exit 1
}

names_file="$(mktemp)"
trap 'rm -f "$names_file"' EXIT
node -e '
  const rows = require(process.argv[1]);
  const version = require(process.argv[2]).version;
  if (!Array.isArray(rows) || rows.length !== 4) throw new Error("canonical platform matrix must contain exactly four rows");
  const names = [`redbench-${version}.tgz`, ...rows.flatMap(row => [
    `redbench-${row.os}-${row.arch}-${version}.tgz`,
    `redbench-${version}-${row.os}-${row.arch}.tar.gz`,
  ])];
  names.sort((a, b) => Buffer.compare(Buffer.from(a), Buffer.from(b)));
  process.stdout.write(names.join("\n") + "\n");
' "$root/scripts/platforms.json" "$root/package.json" > "$names_file"

for directory in "$left" "$right"; do
  while IFS= read -r -d '' entry; do
    name="$(basename "$entry")"
    if ! node -e 'const fs=require("fs"); const rows=require(process.argv[1]); const v=require(process.argv[2]).version; const names=new Set([`redbench-${v}.tgz`,...rows.flatMap(p=>[`redbench-${p.os}-${p.arch}-${v}.tgz`,`redbench-${v}-${p.os}-${p.arch}.tar.gz`])]); process.exit(names.has(process.argv[3])?0:1)' "$root/scripts/platforms.json" "$root/package.json" "$name"; then
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

node - <<'NODE' "$left" "$record" "${entries[@]}"
const fs = require("node:fs");
const crypto = require("node:crypto");
const [directory, record, ...names] = process.argv.slice(2);
const artifacts = names.map(name => {
  const data = fs.readFileSync(`${directory}/${name}`);
  return {name, size: data.length, sha256: crypto.createHash("sha256").update(data).digest("hex"), match: true};
});
const payload = JSON.stringify({schema_version: 1, status: "green", builds: 2, artifacts}) + "\n";
fs.mkdirSync(require("node:path").dirname(record), {recursive: true});
const temporary = `${record}.tmp-${process.pid}`;
fs.writeFileSync(temporary, payload, {mode: 0o644});
fs.renameSync(temporary, record);
NODE
