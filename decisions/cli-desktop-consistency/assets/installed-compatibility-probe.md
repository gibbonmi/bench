# Installed CLI and desktop compatibility probe

The installed interfaces do not currently provide consistent command execution.
The CLI passed two normal-permission commands; the desktop failed before either command started.
This negative result supports recovery design, but it does not qualify the complete Bench workflow.

Date: 2026-10-02
Consumed by: decision tickets #10, #11, and #12
Drift: Repeat affected probes after a runtime restart, update, configuration change, or workspace change.
Retire when: An implementation qualification report supersedes this incident evidence.

## Question and evidence boundary

Can each installed interface start a shell and invoke Bench with its normal sandbox permissions?
The author observed the desktop calls directly.
The reviewer supplied the CLI transcript from an existing WSL Codex CLI chat.
The author did not independently observe that CLI chat.

The calls used the same repository and absolute Bench command.
They did not use identical runtime builds or configuration homes.
Those differences describe the comparison; they do not establish the cause.

## Tested results

Each command ran as a separate tool call with normal sandbox permissions.
The desktop calls used `exec_command` with the repository as `workdir`.
The calls omitted `sandbox_permissions`, which selects the normal sandbox.
The CLI transcript explicitly states that neither call used escalation.

Repository: `/home/mgibs/workspace/bench`

| Invocation | WSL CLI | Desktop with WSL agent |
|---|---|---|
| `pwd` | Exit 0; `/home/mgibs/workspace/bench` | Process creation failed; no command exit code |
| `/home/mgibs/workspace/bench/bin/bench.sh version` | Exit 0; `bench 0.2.0 (linux/amd64)` | Process creation failed; no command exit code |

The desktop returned this error for both calls:

```text
Failed to create unified exec process: No such file or directory (os error 2)
```

The reviewer supplied this CLI evidence:

```text
pwd: /home/mgibs/workspace/bench — exit 0.
/home/mgibs/workspace/bench/bin/bench.sh version:
bench 0.2.0 (linux/amd64) — exit 0.
No startup errors. Both ran separately with normal sandbox permissions; no escalation.
```

The desktop Node tool separately returned this startup error:

```text
js: codex/sandbox-state-meta: sandboxCwd is not a local file URI: file:///home/mgibs/workspace/bench
```

Source: Actual desktop tool responses and the reviewer's CLI transcript in the authoring chat on the recorded date.

## Observed identity and configuration

Read-only diagnostics inspected process executables, selected environment values, configuration files, and held lock paths.
These diagnostics required approved unsandboxed execution because the desktop's normal shell could not start.
They do not count as a successful normal-permission probe.

| Field | WSL CLI | Desktop with WSL agent |
|---|---|---|
| Launcher or bundled version | `codex-cli 0.160.0` | `codex-cli 0.159.0-alpha.12.1` |
| Observed server executable release | `0.159.0-x86_64-unknown-linux-musl` | Bundled binary below |
| Configuration home | Linux default; no process `CODEX_HOME` override | `/mnt/c/Users/gibbo/.codex` |
| User configuration | `/home/mgibs/.codex/config.toml` | `/mnt/c/Users/gibbo/.codex/config.toml` |
| Hook feature setting | Explicitly true | Not explicitly set |
| Temporary launcher lock paths | Two observed paths existed | Observed held lock path was absent |

The desktop bundled executable was:

```text
/mnt/c/Users/gibbo/.codex/bin/wsl/3ac368078cf7546b/codex
```

The desktop server PID was `3083902`.
The CLI server PIDs were `3410389` and `3410483`.
The desktop held this absent lock path:

```text
/mnt/c/Users/gibbo/.codex/tmp/arg0/codex-arg0q0yovE/.lock
```

Both user configurations contained repository hook trust records, with different recorded state.
A missing explicit feature setting does not prove that hooks are disabled.
The repository hook declarations also do not prove live enforcement.

Source: Selected values from `/proc/<pid>/exe`, `/proc/<pid>/environ`, `/proc/<pid>/fd`, and the named configuration files on the recorded date.

## Diagnostic controls

Approved unsandboxed execution could run commands and read the repository.
A fresh sandbox process from the desktop binary also ran each command successfully.
The exact fresh-process invocations were:

```sh
/mnt/c/Users/gibbo/.codex/bin/wsl/3ac368078cf7546b/codex sandbox -P :read-only -C /home/mgibs/workspace/bench -- /bin/pwd
/mnt/c/Users/gibbo/.codex/bin/wsl/3ac368078cf7546b/codex sandbox -P :read-only -C /home/mgibs/workspace/bench -- /home/mgibs/workspace/bench/bin/bench.sh version
```

Both controls returned exit 0 with the expected directory or Bench version.
These controls used separate processes, not the failed chat tool path.
Neither control establishes that the desktop chat recovered.

An earlier attempt to recreate the absent launcher directory returned `File exists`.
No successful private-runtime repair resulted from that attempt.
The author did not restart the active app during this probe.

Source: Diagnostic tool responses in the same authoring chat on the recorded date.

## Workflow coverage

| Capability | Evidence state |
|---|---|
| Shell and absolute Bench invocation | CLI passed the two commands; desktop failed before command execution |
| Bench resolution through PATH | Not qualified across the actual interfaces |
| File reads and writes with normal permissions | Not qualified; desktop shell access failed |
| Rules and skill behavior | This session received instructions; equivalent behavior remains unqualified |
| Startup, stop, and command hooks | Declarations and trust records inspected; live enforcement remains unqualified |
| Required permission behavior | Normal command results observed; the broader permission contract remains unqualified |
| Worktree creation, reviews, and concurrent writers | Not qualified across both interfaces |
| Recovery and retest through the failed interface | No successful recovery observed |
| Desktop Node tool | Startup failed; required operations and equivalent routes remain to be qualified |
| Other desktop-specific tools | Not qualified; availability alone does not establish a Bench requirement |

## Inferences and unknowns

The failure occurs before Bench starts, so this observation does not establish a Bench command defect.
The successful controls narrow the incident to the active desktop execution path or its associated state.
They do not identify the exact missing component.
The absent launcher path is a diagnostic lead, not a proven cause.

The deletion actor, timing, and relationship to the startup failure remain unknown.
A runtime version difference alone does not explain the incident.
The Node URI failure may have a separate cause.
No current evidence proves full workflow compatibility or automatic recovery.

## Validation plan

1. Repeat both commands through the actual desktop chat after an authorized recovery.
2. Record the active runtime identity and effective configuration for each interface.
3. Exercise each unqualified workflow row through both interfaces before claiming complete compatibility.
4. Verify hooks through harmless rejection cases and permitted controls.
5. Exercise concurrent writers in isolated Bench worktrees.
6. Confirm recovery through the failed interface with its normal permissions.

The probe used inline commands and retained no disposable implementation code.
