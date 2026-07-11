#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  roadmap) printf 'context[1]{schema,full}:\n  1,false\n' ;;
  *) echo ok ;;
esac
