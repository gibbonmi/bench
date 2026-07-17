#!/usr/bin/env bash
# The single repository-owned govulncheck setup and version pin used by release workflows.
set -euo pipefail

go install golang.org/x/vuln/cmd/govulncheck@v1.6.0
