#!/usr/bin/env bash
# Build the complete npm artifact set from explicit regular-file inputs. Nothing is
# written to the requested output until every wrapper and platform tarball is ready.
set -euo pipefail

source_root="${1:?usage: build-artifacts.sh <source-root> <output-dir>}"
output="${2:?usage: build-artifacts.sh <source-root> <output-dir>}"
source_root="$(cd "$source_root" && pwd)"
parent="$(dirname "$output")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.bench-artifacts.XXXXXX")"
trap 'rm -rf "$stage"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP
wrapper="$stage/wrapper"
packages="$stage/packages"
artifacts="$stage/artifacts"
mkdir -p "$wrapper" "$packages" "$artifacts"

node - "$source_root" "$wrapper" "$packages" <<'NODE'
const fs = require("fs"), path = require("path");
const [root, wrapperDir, packagesDir] = process.argv.slice(2);
const matrix = JSON.parse(fs.readFileSync(path.join(root, "scripts/platforms.json"), "utf8"));
const assets = JSON.parse(fs.readFileSync(path.join(root, "scripts/wrapper-assets.json"), "utf8"));
const sourcePackage = JSON.parse(fs.readFileSync(path.join(root, "package.json"), "utf8"));

function copyRegular(src, dst, mode) {
  const stat = fs.lstatSync(src);
  if (!stat.isFile()) throw new Error(`wrapper allowlist input is not a regular file: ${src}`);
  fs.mkdirSync(path.dirname(dst), {recursive: true});
  fs.copyFileSync(src, dst);
  fs.chmodSync(dst, Number.parseInt(mode, 8));
}

for (const asset of assets) {
  const src = path.join(root, asset.source);
  if (!asset.tree) {
    copyRegular(src, path.join(wrapperDir, asset.source), asset.mode);
    continue;
  }
  const stat = fs.lstatSync(src);
  if (!stat.isDirectory()) throw new Error(`wrapper allowlist tree is not a directory: ${src}`);
  const walk = (dir) => {
    for (const entry of fs.readdirSync(dir, {withFileTypes: true})) {
      const child = path.join(dir, entry.name);
      if (entry.isDirectory()) walk(child);
      else if (entry.isFile()) copyRegular(child, path.join(wrapperDir, path.relative(root, child)), asset.mode);
      else throw new Error(`wrapper allowlist input is not a regular file: ${child}`);
    }
  };
  walk(src);
}

const optionalDependencies = Object.fromEntries(matrix.map(p => [`@redbench/${p.os}-${p.arch}`, sourcePackage.version]));
const wrapperPackage = {...sourcePackage, optionalDependencies};
wrapperPackage.scripts = {...(sourcePackage.scripts || {})};
delete wrapperPackage.scripts.prepare;
wrapperPackage.files = assets.map(a => a.source);
fs.writeFileSync(path.join(wrapperDir, "package.json"), JSON.stringify(wrapperPackage, null, 2) + "\n", {mode: 0o644});

for (const p of matrix) {
  const dir = path.join(packagesDir, `${p.os}-${p.arch}`);
  fs.mkdirSync(path.join(dir, "bin"), {recursive: true});
  const pkg = {
    name: `@redbench/${p.os}-${p.arch}`,
    version: sourcePackage.version,
    description: `benchkit prebuilt binary for ${p.os}-${p.arch}`,
    bin: "bin/bench",
    os: [p.os],
    cpu: [p.arch],
    license: sourcePackage.license,
  };
  fs.writeFileSync(path.join(dir, "package.json"), JSON.stringify(pkg, null, 2) + "\n", {mode: 0o644});
}
NODE

while IFS=$'\t' read -r os arch goos goarch; do
  binary="$packages/$os-$arch/bin/bench"
  if [[ "$goos" == linux ]]; then
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" bash "$source_root/scripts/go-build.sh" "$source_root" "$binary"
  else
    GOOS="$goos" GOARCH="$goarch" bash "$source_root/scripts/go-build.sh" "$source_root" "$binary"
  fi
  chmod 0755 "$binary"
  npm pack "$packages/$os-$arch" --pack-destination "$artifacts" --ignore-scripts --silent >/dev/null
done < <(node -e 'for (const p of require(process.argv[1])) console.log([p.os,p.arch,p.goos,p.goarch].join("\t"))' "$source_root/scripts/platforms.json")

npm pack "$wrapper" --pack-destination "$artifacts" --ignore-scripts --silent >/dev/null
expected="$(node -e 'process.stdout.write(String(require(process.argv[1]).length + 1))' "$source_root/scripts/platforms.json")"
actual="$(find "$artifacts" -maxdepth 1 -type f -name '*.tgz' -print | wc -l | tr -d ' ')"
[[ "$actual" == "$expected" ]] || { printf 'bench artifacts: emitted %s tarballs, expected %s\n' "$actual" "$expected" >&2; exit 1; }

rm -rf "$output"
mv "$artifacts" "$output"
trap - EXIT INT TERM HUP
rm -rf "$stage"
