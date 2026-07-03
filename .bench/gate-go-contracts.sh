# Go layer for the benchkit gate: the compiled core's build/vet/test/cross-compile
# authority, the version-routing runtime contracts, the platform-package generator
# contract, and the release-workflow structure asserts. Sourced by .bench/gate.sh.
#
# $root is the repo under grade — the real kit normally, a canary fixture during the
# canary sweep. $gate_dir is the real kit's .bench/ (fragments are sourced from there
# even when grading a fixture), so realkit is the real kit root: the home of the
# shared build helper the fixtures deliberately do not carry.
realkit="$(cd "$gate_dir/.." && pwd)"
gobuild="$realkit/scripts/go-build.sh"

# ---- compiled core: gofmt / build / vet / test / cross-compile --------------
# Graded only where a module exists. A minimal fixture with no go.mod has no core to
# grade and skips the whole layer; the two Go canary fixtures carry a go.mod and so
# trip these checks, proving they bite.
if [ -f "$root/go.mod" ]; then
  if ! command -v go >/dev/null 2>&1; then
    # Hard dependency, unlike shellcheck's best-effort posture: once go.mod exists the
    # core is load-bearing, and a gate that silently skips it is an always-pass.
    # (Consumers never run this gate — a missing toolchain only blocks kit dev.)
    err "go.mod present but no Go toolchain on PATH — the compiled core is load-bearing; install Go"
  else
    # gofmt drift. The build canary is a type error (valid syntax), so it stays
    # gofmt-clean and this check attributes to the build, not to formatting.
    unformatted="$(cd "$root" && gofmt -l . 2>/dev/null)"
    [ -z "$unformatted" ] || err "gofmt: unformatted Go files: $(echo "$unformatted" | tr '\n' ' ')"

    # Rebuild the stamped dev binary into dist/ BEFORE the routing contracts exec it —
    # this is what stops a torn or stale dist/ from ever surviving into an assertion.
    if [ -f "$gobuild" ]; then
      ( bash "$gobuild" "$root" "$root/dist/bench" ) >/dev/null 2>&1 || err "go build failed"
    else
      err "go build helper missing ($gobuild)"
    fi

    ( cd "$root" && go vet ./... ) >/dev/null 2>&1 || err "go vet failed"
    ( cd "$root" && go test ./... ) >/dev/null 2>&1 || err "go test failed"

    # Cross-compile every matrix target so a broken GOOS/GOARCH is caught here, not at
    # tag time. Matrix is the single source (scripts/platforms.json); a fixture without
    # it skips cross-compile (it targets build/test, not portability).
    if [ -f "$root/scripts/platforms.json" ]; then
      xtmp="$(mktemp -d)"
      while read -r goos goarch; do
        [ -n "$goos" ] || continue
        ( GOOS="$goos" GOARCH="$goarch" bash "$gobuild" "$root" "$xtmp/bench-$goos-$goarch" ) >/dev/null 2>&1 \
          || err "cross-compile failed: $goos/$goarch"
      done < <(node -e 'const fs=require("fs");for(const p of JSON.parse(fs.readFileSync(process.argv[1],"utf8")))console.log(p.goos,p.goarch);' "$root/scripts/platforms.json")
      rm -rf "$xtmp"
    fi
  fi
fi

# ---- version-routing seam: fabricated-layout sandbox ------------------------
# Runs only against a full kit checkout with a freshly built dev binary (the normal
# gate). Canary fixtures have no bin/bench.sh or dist/bench and skip this whole block.
if [ -f "$root/bin/bench.sh" ] && [ -x "$root/dist/bench" ]; then
  # Executable one-source check: `bench version` == benchkit <pkg.json version>
  # (<host GOOS/GOARCH>). Version drift or an unstamped build breaks the equality.
  want="benchkit $(node -e 'const fs=require("fs");process.stdout.write(JSON.parse(fs.readFileSync(process.argv[1],"utf8")).version)' "$root/package.json") ($(cd "$root" && go env GOOS)/$(cd "$root" && go env GOARCH))"
  got="$(bash "$root/bin/bench.sh" version)"; rc=$?
  [ "$rc" -eq 0 ] || err "bench version exited $rc (want 0)"
  [ "$got" = "$want" ] || err "bench version output '$got' != expected '$want'"
  { [ "$(printf '%s' "$got" | wc -l)" -eq 0 ] && [ -n "$got" ]; } || err "bench version was not exactly one line"

  # Callable from a cwd outside any git repo (humans and hooks run it anywhere).
  ( cd "$(mktemp -d)" && bash "$root/bin/bench.sh" version ) >/dev/null 2>&1 || err "bench version failed outside a git repo"

  vt="$(mktemp -d)"
  (
    set -u
    kit="$vt/a b/kit"                 # spaced path — exercises router quoting
    mkdir -p "$kit"
    cp -r "$root/bin" "$kit/bin"
    # Host package name in npm spelling — node's platform/arch match the router's
    # uname mapping, so this names the exact path the router resolves to (no second
    # copy of the mapping).
    hp="$(node -e 'process.stdout.write("@benchkit/"+process.platform+"-"+process.arch)')"
    stub() { mkdir -p "$(dirname "$1")"; printf '#!/bin/sh\necho %s\n' "$2" > "$1"; chmod +x "$1"; }
    run() { bash "$kit/bin/bench.sh" version; }

    # Precedence: dev build > bundled platform pkg > hoisted global sibling.
    stub "$kit/dist/bench" devbuild
    stub "$kit/node_modules/$hp/bin/bench" bundled
    stub "$kit/../$hp/bin/bench" hoisted
    [ "$(run)" = devbuild ] || { echo "dev build not preferred over bundled"; exit 1; }
    rm "$kit/dist/bench"
    [ "$(run)" = bundled ] || { echo "bundled not preferred over hoisted"; exit 1; }
    rm "$kit/node_modules/$hp/bin/bench"
    [ "$(run)" = hoisted ] || { echo "hoisted sibling not resolved"; exit 1; }
    rm "$kit/../$hp/bin/bench"

    # No binary anywhere → named-package remedy, exit 127.
    out="$(run 2>&1)"; rc=$?
    [ "$rc" -eq 127 ] || { echo "missing-binary exit $rc (want 127)"; exit 1; }
    printf '%s' "$out" | grep -qF "$hp" || { echo "missing-binary error did not name the package"; exit 1; }
    printf '%s' "$out" | grep -qF 'npm install' || { echo "missing-binary error lacked npm install remedy"; exit 1; }

    # Present-but-empty (but executable) binary → -s guard treats it as missing.
    printf '' > "$kit/dist/bench"; chmod +x "$kit/dist/bench"
    rc=0; run >/dev/null 2>&1 || rc=$?
    [ "$rc" -eq 127 ] || { echo "empty binary was exec'd (exit $rc, want 127)"; exit 1; }
    # Present-but-non-executable binary → -x guard treats it as missing.
    printf '#!/bin/sh\necho nope\n' > "$kit/dist/bench"; chmod -x "$kit/dist/bench"
    rc=0; run >/dev/null 2>&1 || rc=$?
    [ "$rc" -eq 127 ] || { echo "non-executable binary was exec'd (exit $rc, want 127)"; exit 1; }
    rm -f "$kit/dist/bench"

    # Symlink invocation resolves the kit from the link's target, not its location.
    stub "$kit/dist/bench" devbuild
    ln -s "$kit/bin/bench.sh" "$vt/benchlink"
    [ "$(bash "$vt/benchlink" version)" = devbuild ] || { echo "symlink invocation did not resolve the kit root"; exit 1; }
    rm "$kit/dist/bench"

    # Off-matrix platform → exit 2, and never names a package (must not alias story 4).
    mkdir -p "$vt/fakebin"
    printf '#!/bin/sh\ncase "$1" in -s) echo Plan9;; -m) echo sparc64;; *) exec /usr/bin/uname "$@";; esac\n' > "$vt/fakebin/uname"
    chmod +x "$vt/fakebin/uname"
    out="$(PATH="$vt/fakebin:$PATH" run 2>&1)"; rc=$?
    [ "$rc" -eq 2 ] || { echo "unsupported platform exit $rc (want 2)"; exit 1; }
    printf '%s' "$out" | grep -qi 'unsupported platform' || { echo "unsupported-platform message missing"; exit 1; }
    if printf '%s' "$out" | grep -qF '@benchkit/'; then echo "unsupported platform named a package"; exit 1; fi
  ) || err "version-routing seam contract failed"
  rm -rf "$vt"
fi

# ---- packaging generator seam: dry-run into a temp dir ----------------------
if [ -f "$root/scripts/gen-platform-packages.sh" ]; then
  gtmp="$(mktemp -d)"
  bash "$root/scripts/gen-platform-packages.sh" "$gtmp/a" >/dev/null 2>&1 || err "platform-package generator failed"
  bash "$root/scripts/gen-platform-packages.sh" "$gtmp/b" >/dev/null 2>&1 || err "platform-package generator (2nd run) failed"
  diff -r "$gtmp/a" "$gtmp/b" >/dev/null 2>&1 || err "platform-package generator is not idempotent"
  ROOT="$root" GEN="$gtmp/a" node <<'NODE' || err "platform-package generator output contract failed"
const fs = require("fs"), path = require("path");
const root = process.env.ROOT, gen = process.env.GEN;
const matrix = JSON.parse(fs.readFileSync(path.join(root, "scripts/platforms.json"), "utf8"));
const wrapper = JSON.parse(fs.readFileSync(path.join(root, "package.json"), "utf8"));
let bad = 0; const e = m => { console.error("gate: " + m); bad = 1; };
for (const p of matrix) {
  const name = `@benchkit/${p.os}-${p.arch}`;
  const pj = path.join(gen, "@benchkit", `${p.os}-${p.arch}`, "package.json");
  if (!fs.existsSync(pj)) { e(`generator did not emit ${name}`); continue; }
  const g = JSON.parse(fs.readFileSync(pj, "utf8"));
  if (g.name !== name) e(`${name}: name is ${g.name}`);
  if (g.version !== wrapper.version) e(`${name}: version ${g.version} != wrapper ${wrapper.version}`);
  if (JSON.stringify(g.os) !== JSON.stringify([p.os])) e(`${name}: os ${JSON.stringify(g.os)}`);
  if (JSON.stringify(g.cpu) !== JSON.stringify([p.arch])) e(`${name}: cpu ${JSON.stringify(g.cpu)}`);
  if (!g.bin) e(`${name}: missing bin field`);
}
const want = Object.fromEntries(matrix.map(p => [`@benchkit/${p.os}-${p.arch}`, wrapper.version]));
const got = wrapper.optionalDependencies || {};
if (JSON.stringify(Object.keys(got).sort()) !== JSON.stringify(Object.keys(want).sort()))
  e(`wrapper optionalDependencies ${JSON.stringify(Object.keys(got))} != matrix`);
for (const k of Object.keys(want))
  if (got[k] !== want[k]) e(`optionalDependency ${k} is ${got[k]}, want ${want[k]}`);
process.exit(bad);
NODE
  rm -rf "$gtmp"
fi

# ---- release workflow: structure asserts (the gate cannot execute Actions) --
# Gated on the matrix source so the Go canary fixtures (go.mod, no matrix) skip it and
# stay attributable to their build/test failure.
if [ -f "$root/scripts/platforms.json" ]; then
  wf="$root/.github/workflows/release.yml"
  if [ ! -f "$wf" ]; then
    err "release workflow missing (.github/workflows/release.yml)"
  else
    grep -qE '^\s*tags:' "$wf" || err "release workflow does not trigger on tags"
    grep -qF 'scripts/platforms.json' "$wf" || err "release workflow does not derive targets from the matrix (scripts/platforms.json)"
    grep -qF 'scripts/gen-platform-packages.sh' "$wf" || err "release workflow does not run the platform-package generator"
    grep -qF 'npm publish' "$wf" || err "release workflow does not publish to npm"
    grep -qF 'provenance' "$wf" || err "release workflow does not publish with provenance"
  fi
fi
