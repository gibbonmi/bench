#!/usr/bin/env bash
set -euo pipefail

root="${1:?usage: gate-prospective.sh <root>}"
exec "$root/.bench/gate.sh"
