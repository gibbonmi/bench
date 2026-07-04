# Go behavior contracts for the benchkit gate: version-routing runtime contracts
# and the platform-package generator contract. Root-grading compiled-core and
# release-structure checks live in the Go conformance suite.
#
# $root is the repo under grade — the real kit normally, a canary fixture during the
# canary sweep. $gate_dir is the real kit's .bench/ (fragments are sourced from there
# even when grading a fixture), so realkit is the real kit root: the home of the
# shared build helper the fixtures deliberately do not carry.
realkit="$(cd "$gate_dir/.." && pwd)"
gobuild="$realkit/scripts/go-build.sh"

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

  # Linked-worktree resolution: dist/ and node_modules/ are untracked, so a git
  # worktree (the .claude/worktrees delegate case) carries the tracked wrapper but
  # no binary. The resolver re-anchors the kit dir under the main tree and retries
  # there; a worktree-local dev build still takes precedence. Without this, every
  # hook in a harness worktree hit the 127 rim and delegates lost read-only git.
  wt="$(mktemp -d)"
  (
    set -u
    gci() { git -c user.email=bench@local -c user.name=bench "$@"; }
    stub() { mkdir -p "$(dirname "$1")"; printf '#!/bin/sh\necho %s\n' "$2" > "$1"; chmod +x "$1"; }
    main="$wt/main"; mkdir -p "$main"; cd "$main"
    git init -q
    cp -r "$root/bin" bin
    gci add -A; gci commit -q -m init
    stub "$main/dist/bench" mainbuild
    gci worktree add -q --detach "$wt/linked" HEAD
    [ "$(bash "$wt/linked/bin/bench.sh" version)" = mainbuild ] \
      || { echo "worktree wrapper did not resolve the main tree's binary"; exit 1; }
    stub "$wt/linked/dist/bench" localbuild
    [ "$(bash "$wt/linked/bin/bench.sh" version)" = localbuild ] \
      || { echo "worktree-local build not preferred over the main tree's"; exit 1; }
  ) || err "linked-worktree binary-resolution contract failed"
  rm -rf "$wt"
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
