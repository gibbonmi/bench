#!/usr/bin/env bash
# Canary fixture: a guard whose --describe manifest DROPS the boundary key. The
# guard-manifest conformance check in .bench/gate-axi-contracts.sh must go red with
# "manifest missing boundary" — proving the check still bites if a guard's
# advertisement stops carrying its boundary. Non-describe mode is a harmless no-op;
# this fixture only exercises the describe path.
set -euo pipefail
if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: block-dangerous-git\n'
  printf 'denies: destructive git — git push\n'
  printf 'why: agents lack destructive-git authority\n'
  exit 0
fi
cat >/dev/null
exit 0
