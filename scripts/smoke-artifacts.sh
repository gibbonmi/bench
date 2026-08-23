#!/usr/bin/env bash
# Execute the host-native tarball through the same installed wrapper a user
# receives.
set -euo pipefail
artifacts="${1:?usage: smoke-artifacts.sh <artifact-dir> [release-evidence-dir]}"
evidence_dir="${2:-}"
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
version="$(node -p 'require(process.argv[1]).version' "$root/package.json")"
npm_install_flags=()
while IFS= read -r arg; do npm_install_flags+=("$arg"); done < <(node -e 'for (const arg of require(process.argv[1]).toolchains.find(tool => tool.name === "npm").operations.install) console.log(arg)' "$root/internal/releaseevidence/requirements.json")
case "$(uname -s)" in Darwin) host_os=darwin ;; Linux) host_os=linux ;; *) host_os=unsupported ;; esac
case "$(uname -m)" in arm64|aarch64) host_arch=arm64 ;; x86_64|amd64) host_arch=x64 ;; *) host_arch=unsupported ;; esac
target_row="$(node "$root/scripts/release-plan.mjs" "$root" target "$host_os" "$host_arch")" || true
target="$(cut -f1-2 <<< "$target_row" | tr '\t' '-')"
[[ -n "$target" ]] || { printf 'bench artifacts: unsupported smoke host %s/%s\n' "$(uname -s)" "$(uname -m)" >&2; exit 2; }
wrapper="$artifacts/redbench-$version.tgz"
native="$artifacts/redbench-$target-$version.tgz"
[[ -f "$wrapper" && -f "$native" ]] || { printf 'bench artifacts: host tarballs missing for %s\n' "$target" >&2; exit 1; }
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/app" "$tmp/home/cache/bin/$version/$target"
mkdir -p "$tmp/native"
tar -xzf "$native" -C "$tmp/native" package/bin/bench
native_binary="$tmp/native/package/bin/bench"
[[ -x "$native_binary" && -s "$native_binary" ]] || { printf 'bench artifacts: %s binary mode or size invalid\n' "$target" >&2; exit 1; }
format="$(file "$native_binary")"
case "$target" in
  darwin-arm64) [[ "$format" == *"Mach-O 64-bit"* && "$format" == *"arm64"* ]] ;;
  darwin-x64) [[ "$format" == *"Mach-O 64-bit"* && "$format" == *"x86_64"* ]] ;;
  linux-arm64) [[ "$format" == *"ELF 64-bit"* && "$format" == *"ARM aarch64"* && "$format" == *"statically linked"* ]] ;;
  linux-x64) [[ "$format" == *"ELF 64-bit"* && "$format" == *"x86-64"* && "$format" == *"statically linked"* ]] ;;
esac || { printf 'bench artifacts: %s format mismatch: %s\n' "$target" "$format" >&2; exit 1; }
printf '#!/bin/sh\nprintf "poisoned cache\\n"\nexit 99\n' > "$tmp/home/cache/bin/$version/$target/bench"
chmod 0755 "$tmp/home/cache/bin/$version/$target/bench"
printf '{"private":true}\n' > "$tmp/app/package.json"
HOME="$tmp/home" BENCH_HOME="$tmp/home/.bench" BENCH_NO_REPAIR=1 npm_config_cache="$tmp/npm-cache" npm_config_offline=true npm_config_registry=http://127.0.0.1:9 npm install --offline --prefix "$tmp/app" "${npm_install_flags[@]}" "$wrapper" "$native" >/dev/null
installed="$tmp/app/node_modules/redbench/bin/bench.sh"
out="$(HOME="$tmp/home" BENCH_HOME="$tmp/home/.bench" BENCH_NO_REPAIR=1 bash "$installed" version)"
case "$out" in
	"bench $version "*) ;;
  *) printf 'bench artifacts: native version mismatch: %s\n' "$out" >&2; exit 1 ;;
esac
HOME="$tmp/home" BENCH_HOME="$tmp/home/.bench" BENCH_NO_REPAIR=1 bash "$installed" help >/dev/null
HOME="$tmp/home" BENCH_HOME="$tmp/home/.bench" npm_config_registry=http://127.0.0.1:9 npm_config_offline=true npm uninstall --offline --prefix "$tmp/app" redbench "@redbench/$target" >/dev/null
[[ ! -e "$tmp/app/node_modules/redbench" && ! -e "$tmp/home/.bench/cache" ]] || { printf 'bench artifacts: local npm uninstall left residue\n' >&2; exit 1; }
if [[ -n "$evidence_dir" ]]; then
  bash "$root/scripts/smoke-offline.sh" "$artifacts" "$evidence_dir"
fi
printf 'bench artifacts: %s selected %s\n' "$target" "$out"
