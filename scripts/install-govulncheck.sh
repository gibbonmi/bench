#!/usr/bin/env bash
# This script is the repository's single govulncheck setup and version pin. Release
# workflows use it.
set -euo pipefail

go install golang.org/x/vuln/cmd/govulncheck@v1.6.0
