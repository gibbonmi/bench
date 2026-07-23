#!/usr/bin/env bash
# .bench/skills-index.sh — the skills index in .bench/BENCH-reference.md is
# generated, not hand-maintained. Source of truth is each skill's frontmatter
# (`index:` trigger phrase, optional `index-note:` suffix); the generated block
# lives between the bench:skills-index markers, alphabetical by skill directory.
#
#   --check   (default) verify the committed block equals the generated one;
#             attributed errors on stderr, nonzero exit on drift. Run by the gate.
#   --write   regenerate the block in place (temp file + mv).
#
# Command-adapter skills (a same-named file in .agents/commands/) stay unindexed.
# Kit-only: lives beside gate.sh, never ships in package.json files[].
set -uo pipefail

mode="${1:---check}"
bench_md=".bench/BENCH-reference.md"
start_marker='<!-- bench:skills-index:start -->'
end_marker='<!-- bench:skills-index:end -->'
regen_hint="(regenerate: .bench/skills-index.sh --write)"
fail=0
err() { echo "$*" >&2; fail=1; }

frontmatter_field() { # file, key — first value of key inside the --- fence
  awk -v key="$2" -F': *' '
    /^---$/ { fence++; next }
    fence == 1 && $1 == key { sub(/^[^:]*: */, ""); print; exit }
  ' "$1"
}

# Skill sources the consumer-payload allowlist withholds. The allowlist is the one
# source of who receives an asset, so the index marks the same rows rather than
# carrying a second list of kit-only skills. Its presence is not re-asserted here —
# the Go core embeds the same file, so a kit without it does not build — and a tree
# that carries no allowlist simply has no kit-only rows to mark. Read at top level:
# expected_lines runs in a command substitution, where an err() would be lost.
#
# The two fields are matched independently, so a row that orders "audience" before
# "source" still resolves. A shell reader cannot parse JSON in general: a file that
# declares kit-only rows this reader finds no source for is a parse failure, not an
# empty result, and is reported below rather than silently generating an index with
# every kit-only marker missing.
payload=".bench/consumer-payload.json"
kit_only_sources=""
payload_unparsed=""
if [ -f "$payload" ]; then
  kit_only_sources="$(awk '
    /"audience"[[:space:]]*:[[:space:]]*"kit-only"/ && match($0, /"source"[[:space:]]*:[[:space:]]*"[^"]*"/) {
      field = substr($0, RSTART, RLENGTH)
      sub(/^"source"[[:space:]]*:[[:space:]]*"/, "", field)
      sub(/"$/, "", field)
      print field
    }' "$payload")"
  if grep -q '"kit-only"' "$payload" && [ -z "$kit_only_sources" ]; then
    payload_unparsed=1
  fi
fi

# Generated lines, one per indexed skill, alphabetical by directory (glob order).
# Validation stays out of this function: it runs in a command substitution, where
# an err() would set fail=1 in a subshell and be lost.
expected_lines() {
  local d name trigger note marker source
  for d in .agents/skills/*/; do
    [ -f "$d/SKILL.md" ] || continue
    name="$(basename "$d")"
    [ -f ".agents/commands/$name.md" ] && continue
    trigger="$(frontmatter_field "$d/SKILL.md" index)"
    [ -z "$trigger" ] && continue
    note="$(frontmatter_field "$d/SKILL.md" index-note)"
    marker=""
    while IFS= read -r source; do
      [ "$source" = ".agents/skills/$name" ] && marker=" (kit-only)"
    done <<EOF
$kit_only_sources
EOF
    if [ -n "$note" ]; then
      printf -- '- %s → `.agents/skills/%s/SKILL.md`%s + %s\n' "$trigger" "$name" "$marker" "$note"
    else
      printf -- '- %s → `.agents/skills/%s/SKILL.md`%s\n' "$trigger" "$name" "$marker"
    fi
  done
}

# Every indexed skill must declare its trigger — checked in parent scope so the
# error survives even when the committed block happens to match.
for d in .agents/skills/*/; do
  [ -f "$d/SKILL.md" ] || continue
  name="$(basename "$d")"
  [ -f ".agents/commands/$name.md" ] && continue
  [ -n "$(frontmatter_field "$d/SKILL.md" index)" ] \
    || err "skill '$name' missing index: frontmatter (the skills index is generated)"
done

[ -f "$bench_md" ] || { err "$bench_md missing (skills index unverifiable)"; exit 1; }
if ! grep -qF "$start_marker" "$bench_md" || ! grep -qF "$end_marker" "$bench_md"; then
  err "$bench_md skills-index markers missing (bench:skills-index)"
  exit 1
fi

# Fail before generating, not after: a --write that ran on an unparsed allowlist would
# commit a block with every kit-only marker missing, and the drift it caused would be
# attributed to the skills rather than to the reader that dropped them.
if [ -n "$payload_unparsed" ]; then
  err "$payload declares kit-only rows this reader could not resolve to a source (each row's \"source\" and \"audience\" must share a line)"
  exit 1
fi

expected="$(expected_lines)"
actual="$(awk -v s="$start_marker" -v e="$end_marker" '
  $0 == e { on = 0 } on { print } $0 == s { on = 1 }
' "$bench_md")"

if [ "$mode" = "--write" ]; then
  tmp="$(mktemp)"
  awk -v s="$start_marker" -v e="$end_marker" -v block="$expected" '
    $0 == s { print; print block; skip = 1; next }
    $0 == e { skip = 0 }
    !skip { print }
  ' "$bench_md" > "$tmp"
  # mktemp creates 0600; the reference file is a shipped 0644 asset whose mode the
  # release-evidence registry checks, so restore it before the file is replaced.
  chmod 644 "$tmp"
  mv "$tmp" "$bench_md"
  exit "$fail"
fi

[ "$expected" = "$actual" ] && exit "$fail"

# Attribute the drift per skill before the catch-all.
attributed=0
while IFS= read -r line; do
  [ -n "$line" ] || continue
  name="$(printf '%s' "$line" | sed -n 's#.*\.agents/skills/\([a-z0-9-]*\)/SKILL\.md.*#\1#p')"
  if ! grep -qF -- "$line" <<EOF
$actual
EOF
  then
    if printf '%s\n' "$actual" | grep -qF ".agents/skills/$name/SKILL.md"; then
      err "skills index entry for '$name' drifted from its frontmatter $regen_hint"
    else
      err "skills index missing entry for skill '$name' $regen_hint"
    fi
    attributed=1
  fi
done <<EOF
$expected
EOF
while IFS= read -r line; do
  [ -n "$line" ] || continue
  name="$(printf '%s' "$line" | sed -n 's#.*\.agents/skills/\([a-z0-9-]*\)/SKILL\.md.*#\1#p')"
  [ -n "$name" ] || continue
  if ! printf '%s\n' "$expected" | grep -qF ".agents/skills/$name/SKILL.md"; then
    err "skills index entry '$name' has no indexed .agents/skills/$name on disk $regen_hint"
    attributed=1
  fi
done <<EOF
$actual
EOF
[ "$attributed" = "1" ] || err "skills index block drifted from generated form $regen_hint"
exit 1
