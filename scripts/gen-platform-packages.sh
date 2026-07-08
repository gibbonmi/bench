#!/usr/bin/env bash
# Generate the four @redbench/<os>-<arch> platform packages from the single matrix
# source (scripts/platforms.json) plus the wrapper version (package.json). Package
# names, os/cpu fields, and version all derive here — there is no hand-maintained
# per-package metadata to drift. The release workflow runs it for real; the gate runs
# it into a temp dir and asserts shape + idempotency.
#
# Repo-only: not in package.json files[]. The generated packages are never committed
# (committing them would be a second source for the matrix — see the gate contract).
#
#   Usage: gen-platform-packages.sh <output-dir>
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
out="${1:?usage: gen-platform-packages.sh <output-dir>}"

node - "$here" "$out" <<'NODE'
const fs = require("fs"), path = require("path");
const [here, out] = process.argv.slice(2);
const matrix = JSON.parse(fs.readFileSync(path.join(here, "scripts/platforms.json"), "utf8"));
const wrapper = JSON.parse(fs.readFileSync(path.join(here, "package.json"), "utf8"));

for (const p of matrix) {
  const dir = path.join(out, "@redbench", `${p.os}-${p.arch}`);
  fs.mkdirSync(path.join(dir, "bin"), { recursive: true });
  // Fixed key order + fixed matrix order → byte-identical output on re-run.
  const pkg = {
    name: `@redbench/${p.os}-${p.arch}`,
    version: wrapper.version,
    description: `benchkit prebuilt binary for ${p.os}-${p.arch}`,
    // A string `bin` is linked by npm under the unscoped package name (linux-x64),
    // never `bench` — so a platform package can never collide with the wrapper's
    // own `bench` bin. The wrapper's launcher execs the binary by path regardless.
    bin: "bin/bench",
    os: [p.os],
    cpu: [p.arch],
    license: wrapper.license,
  };
  fs.writeFileSync(path.join(dir, "package.json"), JSON.stringify(pkg, null, 2) + "\n");
}
NODE
