#!/usr/bin/env bash
# The build phase's probe requires this file, so a fixture grading that phase has to
# carry one. It is the kit helper reduced to the part the fixture needs: a build of
# ./cmd/bench into the path the phase names. The kit's own flag set is deliberately
# absent — this tree has no package.json and no requirements registry to read it from,
# and the fixture grades whether the phase reds on a broken package, not how it builds.
set -euo pipefail
cd "${1:?usage: go-build.sh <module-root> <output-path>}"
go build -o "${2:?usage: go-build.sh <module-root> <output-path>}" ./cmd/bench
