#!/usr/bin/env bash
# shellcheck shell=bash
# Canary runner — the gate proving its own checks still bite. Sourced by
# .bench/gate.sh (benchkit) and by the gate that `bench init` scaffolds (consumers),
# so the logic has one home instead of a copy per gate. Sourced, not run.
#
# Contract: the sourcing gate must define `root`, `err`, and `fail` in scope before
# sourcing, and must NOT use `set -e` (the biting loop reads `$?` after a subshell
# that is expected to fail). The gate under test is assumed to live at
# "$root/.bench/gate.sh".
#
# For each known-broken fixture in tests/canary/, run THIS gate against it in a
# throwaway repo and assert it goes red WITH the fixture's targeted error substring. A
# fixture that stops biting means a check rotted into an always-pass. Attribution is by
# substring, not isolation: a minimal fixture over-fails on unrelated checks, and that
# is fine — we only assert the targeted message is present. BENCH_CANARY_INNER marks
# the inner run so this check skips itself and never recurses.
#
# `root` is provided by the sourcing gate, not assigned here.
# shellcheck disable=SC2154
if [ "${BENCH_CANARY_INNER:-0}" != "1" ]; then
  # A gate with no canary harness cannot prove its own checks bite — the exact hole a
  # lazy agent opens by deleting tests/canary/ to escape a failing check. Absent OR
  # empty is red: removing or emptying the harness is a red gate and a visible diff,
  # never a silent pass. (The glob stays literal when nothing matches, so `[ -d ]`
  # covers both the absent-dir and empty-dir cases.)
  have_fixture=0
  for fx in tests/canary/*/; do [ -d "$fx" ] && { have_fixture=1; break; }; done
  if [ "$have_fixture" -eq 0 ]; then
    err "canary harness absent — tests/canary/ has no fixtures; the gate cannot prove its own checks bite"
  else
    # Attribution baseline: an EXPECT that also matches a completely empty fixture
    # proves nothing about its planted regression — the canary is vacuous and the
    # check it guards can rot into an always-pass unnoticed.
    d0="$(mktemp -d)"
    ( cd "$d0" && git init -q && BENCH_CANARY_INNER=1 bash "$root/.bench/gate.sh" ) >"$d0/out" 2>&1 || true
    for fx in tests/canary/*/; do
      name="$(basename "$fx")"
      [ -f "$fx/EXPECT" ] || continue
      if grep -qF "$(cat "$fx/EXPECT")" "$d0/out"; then
        err "canary '$name' EXPECT is vacuous (also matches an empty fixture)"
      fi
    done
    rm -rf "$d0"
    for fx in tests/canary/*/; do
      name="$(basename "$fx")"
      [ -f "$fx/EXPECT" ] || { err "canary fixture '$name' has no EXPECT file"; continue; }
      [ -d "$fx/files" ]  || { err "canary fixture '$name' has no files/ tree"; continue; }
      exp="$(cat "$fx/EXPECT")"
      d="$(mktemp -d)"
      cp -r "$fx/files/." "$d/"
      # Fixtures store dot-dirs as dot-<name> so the harness doesn't load a fixture's
      # .claude/skills as real skills; restore them to .<name> for the inner gate.
      for dd in "$d"/dot-*; do [ -e "$dd" ] && mv "$dd" "$d/.${dd##*/dot-}"; done
      ( cd "$d" && git init -q && BENCH_CANARY_INNER=1 bash "$root/.bench/gate.sh" ) >"$d/out" 2>&1
      rc=$?
      if [ "$rc" -eq 0 ] || ! grep -qF "$exp" "$d/out"; then
        err "canary '$name' did not bite (want red + \"$exp\"; got exit $rc)"
      fi
      rm -rf "$d"
    done
  fi
fi
