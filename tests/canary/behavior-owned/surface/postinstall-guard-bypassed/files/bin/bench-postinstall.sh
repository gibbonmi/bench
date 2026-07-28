#!/usr/bin/env bash
# Canary fixture: postinstall runs the fix on a global install WITHOUT the .git guard
# (stories 10/11) — so it mutates in a dev checkout / npm link. The postinstall
# contract must go red on the .git-present permutation.
set -u
advice() { printf "bench: run 'bench doctor --fix' to install the PATH shim\n"; }
case "${npm_config_global:-}" in true|1|yes) ;; *) advice; exit 0 ;; esac
here="$(cd "$(dirname "$0")/.." && pwd)"
# No `.git` guard — the regression.
bash "$here/bin/bench.sh" doctor --fix || advice
exit 0
