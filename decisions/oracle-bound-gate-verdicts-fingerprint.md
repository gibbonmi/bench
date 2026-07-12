# Oracle fingerprint inventory

The resolver has six executable kinds after `None`: an executable
`.bench/gate.sh`, the exact `BENCH_GATE` shell string, then fixed pnpm, npm,
Python, and Cargo commands selected in that order. Resolution and command
construction have one owner in `internal/gate`; every command runs from the
repository root
([resolver](../internal/gate/gate.go#L75),
[commands](../internal/gate/gate.go#L103),
[cwd](../internal/gate/gate.go#L146)).

## Common subject

Every reusable kind needs one canonical subject digest over:

- verdict schema and fingerprint-policy version;
- resolved kind, exact executable/argv or shell string, and repository-root cwd;
- the working-tree hash;
- the resolved launcher closure that the policy recognizes: path, symlink
  disposition, executable mode, and content identity;
- one canonical digest over declared environment, tool, ignored-file, and other
  local inputs; and
- the freshness-policy identity once ticket #4 settles it.

The working-tree hash already stages tracked plus untracked, unignored files in
a throwaway Git index, so repository content, Git executable mode, and symlink
entries join that hash; ignored files and the content of an external symlink
target do not
([tree-hash scope](../internal/git/tree.go#L11),
[tree construction](../internal/git/tree.go#L24)). An unavailable tree hash is
not a subject and therefore cannot be reused.

The current runner passes the full inherited environment after removing only
`BENCH_KIT` and `BENCH_WRAPPER`
([environment construction](../internal/gate/gate.go#L125),
[run attachment](../internal/gate/gate.go#L146)). Consequently, an explicit
input declaration does not close the current oracle by itself: undeclared
environment still reaches the gate and can change its result. Ticket #3 must
either narrow execution to a declared passlist or make any undeclared inherited
state disable reuse.

## Kind-specific closure

| Kind | Resolver-owned material | Additional material needed for reuse |
|---|---|---|
| `.bench/gate.sh` | `GateSh`, absolute script path, tree hash | Final executable/interpreter closure and every ignored or external input the script reaches. A managed template may supply this declaration; arbitrary project scripts cannot be inferred. |
| `BENCH_GATE` | `BenchGate`, exact shell string, `bash -c`, tree hash | Resolved Bash identity plus declared executables, files, and environment. The command is opaque; no declaration means rerun. |
| pnpm | `Pnpm`, fixed typecheck/test/lint string, `pnpm-lock.yaml` selection, tree hash | Resolved Bash and pnpm identities plus declared external/ignored inputs. Repository scripts/config are already in the tree when unignored. |
| npm | `Npm`, fixed typecheck/test/lint string, `package.json` selection, tree hash | Resolved Bash and npm identities plus declared external/ignored inputs. |
| Python | `Pyproject`, fixed mypy/pytest/ruff string, `pyproject.toml` selection, tree hash | Resolved Bash and the three tool identities plus declared external/ignored inputs. |
| Cargo | `Cargo`, fixed test/clippy string, `Cargo.toml` selection, tree hash | Resolved Bash, cargo, rustc/toolchain and clippy identities plus declared external/ignored inputs. |

Selection files and precedence are fixed by the resolver
([selection table](../internal/gate/gate.go#L87)); their bytes are already in
the tree when unignored. A symlink qualifies today because production probing
uses `os.Stat`, which follows it; a future subject must record the link and the
recognized final target rather than silently inheriting that behavior
([production probe](../internal/gate/gate.go#L60)).

The kit's own gate demonstrates why ignored runtime inputs need declarations:
the tracked script executes the wrapper, whose binary resolver may choose the
ignored `dist/bench` before package or cache candidates
([kit gate](../.bench/gate.sh#L11),
[binary selection](../bin/bench.sh#L119)). The tree hash alone cannot identify
that binary.

## Runnable mutation probes

These focused probes use the real tree-hash and resolver surfaces and were
green on 2026-07-12:

```sh
go test -count=1 ./internal/git -run '^TestTreeHash'
go test -count=1 ./internal/gate -run '^(TestResolvePrecedence|TestGateEnvStripsWrapperRoutingInternals)$'
go test -count=1 ./internal/contract/runtime -run '^TestRuntimeGateContracts$'
```

The tree probe mutates tracked content and adds untracked content and observes a
different subject
([mutation cases](../internal/git/git_test.go#L48)). The resolver probe changes
gate-script presence, `BENCH_GATE`, and `package.json`, then observes the real
CLI select gate script, environment command, auto-detect, and finally no gate
([runtime mutation](../internal/contract/runtime/runtime_gate_test.go#L68)). The
verdict probe changes `.bench/gate.sh` from green to red and observes the cached
tree/status change
([verdict mutation](../internal/contract/runtime/runtime_gate_test.go#L141)).

## Result

The smallest honest manifest is extensible data, not an extensible public Go
interface: a fixed canonical envelope plus sorted declared input entries. It
stores one aggregate subject digest and non-sensitive diagnostics, not raw
environment values or command-derived secrets. A kind is reusable only when
the resolver can produce the complete manifest required by ticket #1; inability
to read, resolve, or digest any required entry makes the verdict non-reusable.
