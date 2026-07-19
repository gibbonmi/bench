#!/usr/bin/env bash
# Validate the complete native proof set derived from the canonical matrix.
set -euo pipefail

proofs="${1:?usage: aggregate-native-proofs.sh <proof-dir>}"
root="$(node -e 'const fs=require("node:fs"),path=require("node:path");process.stdout.write(path.dirname(path.dirname(fs.realpathSync(process.argv[1]))))' "${BASH_SOURCE[0]}")"
[[ -d "$proofs" && ! -L "$proofs" ]] || { printf 'native proof directory is unsafe\n' >&2; exit 1; }
matrix="$(mktemp "${TMPDIR:-/tmp}/bench-native-matrix.XXXXXX")"
trap 'rm -f "$matrix"' EXIT
node "$root/scripts/release-plan.mjs" "$root" targets > "$matrix"
expected="$(while IFS=$'\t' read -r os arch _goos _goarch _runner; do printf '%s-%s.json\n' "$os" "$arch"; done < "$matrix" | LC_ALL=C sort)"
count=0
while IFS=$'\t' read -r os arch _goos _goarch runner; do
  file="$proofs/$os-$arch.json"
  [[ -f "$file" && ! -L "$file" && -s "$file" ]] || { printf 'native proof set is missing %s/%s\n' "$os" "$arch" >&2; exit 1; }
  # shellcheck disable=SC2016 # Node template literals are intentionally literal here.
  node -e 'const fs=require("fs"), [file,os,arch,runner]=process.argv.slice(1), proof=JSON.parse(fs.readFileSync(file)); const want=`${os}-${arch}`; if (proof.schema_version!==1 || proof.target!==want || proof.runner!==runner || proof.status!=="green" || !proof.rebuilt_sha256 || !proof.binary_sha256 || !proof.package_sha256 || !proof.archive_sha256 || proof.operations_status!=="green" || proof.strip_status!=="green" || proof.tools_status!=="green" || (os==="linux" ? proof.musl_status!=="green" : proof.musl_status!=="not_applicable")) process.exit(1)' "$file" "$os" "$arch" "$runner" || { printf 'native proof is incomplete or red for %s/%s\n' "$os" "$arch" >&2; exit 1; }
  count=$((count + 1))
done < "$matrix"
planned_count="$(wc -l < "$matrix" | tr -d '[:space:]')"
[[ "$count" == "$planned_count" ]] || { printf 'native proof set has %s targets, want %s from release plan\n' "$count" "$planned_count" >&2; exit 1; }
actual="$(find "$proofs" -mindepth 1 -maxdepth 1 -print0 | while IFS= read -r -d '' entry; do
  [[ -f "$entry" && ! -L "$entry" ]] || { printf 'native proof directory contains an unsafe entry\n' >&2; exit 1; }
  printf '%s\n' "${entry##*/}"
done | LC_ALL=C sort)"
[[ "$actual" == "$expected" ]] || { printf 'native proof set does not contain exactly the canonical proof files\n' >&2; exit 1; }
printf 'native proof set: %s canonical targets verified\n' "$planned_count"
