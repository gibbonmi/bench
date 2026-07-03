#!/usr/bin/env bash
# Canary fixture: postinstall propagates a doctor --fix write failure as a nonzero
# exit (story 13) — breaking the never-fail-the-install invariant. The postinstall
# contract must go red on the read-only-target permutation.
set -u
advice() { printf "bench: run 'bench doctor --fix' to install the PATH shim\n"; }
case "${npm_config_global:-}" in true|1|yes) ;; *) advice; exit 0 ;; esac
here="$(cd "$(dirname "$0")/.." && pwd)"
[ -e "$here/.git" ] && { advice; exit 0; }
bash "$here/bin/bench.sh" doctor --fix || exit 1   # BUG: fails the global install
exit 0
