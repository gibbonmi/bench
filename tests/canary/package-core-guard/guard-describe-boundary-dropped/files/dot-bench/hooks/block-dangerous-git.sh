#!/usr/bin/env bash
# name: block-dangerous-git
# denies: destructive git operations
# why: agents lack destructive-git authority
# Canary fixture: a static guard header that DROPS the boundary key. The
# guard-manifest conformance check must go red with "manifest missing boundary".
set -euo pipefail
exit 0
