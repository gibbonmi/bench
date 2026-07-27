#!/usr/bin/env bash
# The build phase's probe requires this file, so a fixture grading that phase has to
# carry one. It is the kit helper reduced to the part the fixture needs: a build of
# ./cmd/bench into the path the phase names. The kit's own flag set is deliberately
# absent — this tree has no package.json and no requirements registry to read it from,
# and the fixture grades whether the phase reds on a broken package, not how it builds.
#
# The refusal line is the fixture's own, and it is what the EXPECT matches. A compiler
# diagnostic would tie the expectation to Go's wording, where a release that rephrases
# it turns the fixture into a "did not bite" report about the toolchain rather than
# about the phase. No other phase runs this script, so the line attributes exactly.
set -euo pipefail
cd "${1:?usage: go-build.sh <module-root> <output-path>}"
if ! go build -o "${2:?usage: go-build.sh <module-root> <output-path>}" ./cmd/bench; then
  echo "canary: go-build.sh refused this tree"
  exit 1
fi
