# FT173 contextual-help inventory

Subject: `974020e4af8de5ed75098c4c5934a8907952bb2b`

This is the ticket #2 research object for `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`. It
inventories the current executable surface, the state each call already owns at
return time, and the next command that state can make concrete. It does not choose
the compatibility boundary or implement `help[]`.

## Evidence boundary

The executable census has three current owners:

- the Go registry has 48 root names (`cmd/bench/main.go:66-118`);
- the wrapper adds the public no-argument/help view, `--version`/`-v` aliases,
  and wrapper-only `repair`, but does not route the registry-only `freshness-publish` name
  (`bin/bench.sh:284-390`); and
- nested grammars add eight worktree modes, three ordinary spec operations, nine
  spec-build operations, five release operations, two release-preflight modes,
  two gate-go modes, one gate `pin` operation, and two worktree-hook modes
  (`internal/usage/worktree.go:5-26`; `internal/spec/spec.go:291-302`;
  `internal/spec/build.go:16-59`; `internal/publication/command.go:15-61`;
  `internal/gate/gate.go:231-240`).

`bench commands --brief` is not an inventory source: it prints only `version`,
`commands --brief`, and `status` from a hard-coded body
(`cmd/bench/main.go:207-225`). The tables below therefore start from the production
registry and wrapper, then expand every declared nested grammar. No command is
omitted merely because a log search did not find it.

The available usage evidence is:

- 100 Claude project transcripts under the local Bench project archive. The
  acceptance-driving example is Claude session
  `0b2828ac-933f-404e-8de4-c785735406ac`: a successful reduced gate was followed
  by `bench spec build start`, which refused with “run bench gate” even though
  only `bench gate --fresh` could establish the required evidence (JSONL lines
  257 and 266). The reviewer then required follow-up commands for every Bench CLI
  call and an unsampled command inventory (line 288).
- 626 locally available Codex rollout files, of which 91 contain a recorded
  `CommandExecution` with a cwd under `/home/devuser/workspace/bench`. These are
  project-scoped corroboration, not the promised named corpus.
- unavailable: no reviewer-supplied Codex usage-log manifest, session-id list, or
  archive was named after the Claude note saying those logs would be shared. This
  asset names that missing source instead of substituting selected Codex sessions
  or interpreting a zero local match as evidence that a command is unused.

Log marks below are `Cl` for a direct Claude Bash invocation, `Cx` for a direct
Codex command execution, and `—` when the project-scoped extraction found none.
They are lower-bound occurrence marks, not popularity scores: variable-held
executables, hook-owned calls, and commands inside a child process are not
reconstructed. Case marks are `S` success, `E` empty/no-op, `R` refusal or usage,
`St` stale/drift, and `Rc` recovery. “Yes” means contextual disclosure can remove
a discovery or wrong-remedy turn; “Direct” means the next action still needs its
own call but `help[]` can make that call copy-paste exact; “No” means the answer is
terminal or the next action belongs to a caller rather than this command. An
omitted case mark means the current owner has no distinct semantic projection for
that case; it does not mean the logs proved the case impossible.

## Home, query, capture, and diagnostics

| command | log | cases | state already known at return | exact contextual action and turn effect |
|---|---|---|---|---|
| `bench` / `bench help` | — | R | Wrapper executable and the duplicated public catalog; no repository state | **Yes:** `bench status`; the AXI home must replace the current help-only no-argument result rather than suggest a guessed leaf. |
| `bench version` | Cl | S | Version, OS, architecture, and whether the build is unstamped `dev` | **Yes only for `dev`:** `bench doctor`; otherwise **No**, the version answer is terminal. |
| `bench commands --brief` | Cl+Cx | S,R | Only its three hard-coded probe lines; it does not know the live registry | **Yes:** `bench status`; more importantly, deriving the catalog from the registry removes a separate discovery failure rather than adding hints to an incomplete list. |
| `bench status [--all]` | Cl+Cx | S,E,St,Rc | Every signal, severity, detail, and existing `action`; overflow count; gate freshness; dirty paths; worktree and drain state | **Yes:** append each invokable action already known, including `bench link`, `bench gate --fresh`, `$bench-what-next`, `bench spec retire <slug>`, and `bench worktree clean <path>`. Compound prose such as `commit on green / /bench-final-check / push` is not one executable template. |
| `bench handoff [--harness <name>] [--next <command>]` | Cl+Cx | S,E,R | Repo, branch, tip, dirty state, translated harness, and the one derived or overridden next command | **Direct:** emit the exact `## Next command` value once; a second generic hint would duplicate the owner. On refusal, `bench status` is the safe recovery view. |
| `bench dashboard [--stdout]` | — | S,E,R | Board snapshot and either the generated destination or complete stdout body | **No** on success; **Yes** on render/write refusal with `bench status`, because the command has no repair operation. |
| `bench anchors <path>` | Cl | S,E,R | Exact queried path and every anchor row | **No** on success/empty: it is a terminal lookup and cannot infer the consumer. On usage, `bench anchors <path>` is sufficient help with `<path>` left open. |
| `bench learnings` | Cl+Cx | S,E,R | Open-entry count, dates, and titles | **Yes:** when non-empty, `$bench-what-next`; when empty, **No**. |
| `bench maps [--count or --template]` | Cl+Cx | S,E,R,St | Map path, ticket, type, frontier/blocked/invalid state, and blockers; template/count mode | **Yes:** frontier → `$bench-shape-idea`; invalid → `bench maps --template` plus the named diagnostic path; empty → **No**. The current output gives the fact but no executable next step. |
| `bench guards [--brief]` | — | S,E,R,St | Guard names, deny boundaries, wiring, and scan completeness | **Yes:** stale/unwired managed guards → `bench link`; complete inventory → **No**. Hook-driven `guards --brief` is not evidence of an agent discovery turn. |
| `bench diff [--full] [--commit <sha>]` | Cl+Cx | S,E,R,St | Review base, checkout facts, file inventory, counts, whitespace state, and the supplied commit when present | **Yes:** bounded orientation → `bench diff --full`; landed subject → `bench diff --commit <sha>` with only `<sha>` open; clean/full output → **No**. Mid-read drift must suggest retrying the same exact invocation. |
| `bench coverage [--check] <spec>` | Cl+Cx | S,E,R,St | Spec identity, map state, every row, and whether validation was requested | **Yes:** unvalidated map → `bench coverage --check <spec>` carrying the resolved spec; invalid/stale → repair the named row then retry that exact command; complete checked map → **No**. |
| `bench test [--full] [package]` | Cl | S,E,R | Package selector, package outcomes, first failures/skips, truncation, and whether full detail was requested | **Yes:** truncated red → repeat the exact selector with `--full`; bounded red → retry the exact same selector after repair; green/empty package set → **No**. |
| `bench structure [--since <base>]` | Cx | S,E,R | Every finding, path, budget, and supplied comparison base | **Direct:** findings → repair the named objects then rerun the same invocation; clean → **No**. No safe automatic repair command exists. |
| `bench models` | Cl | S,E,R | Candidate model ids, harness source, and availability state | **Yes:** `bench resolve-model --harness <harness>` with only the runtime harness open; no candidates → inspect the named unavailable source rather than guess. |
| `bench outline [--full] [path]` | Cl | S,E,R | Resolved path, emitted/omitted counts, and whether the row limit applied | **Yes only when bounded:** `bench outline --full <path>` carrying a supplied path; complete or empty → **No**. |
| `bench idea "<text>"` | Cl | S,R | Captured text and destination | **Yes:** `$bench-what-next`; it is the only workflow that reviews the new capture. |
| `bench roadmap` | Cl+Cx | S,E,R,St | Roadmap body, drain counts, and its current `Next action` section | **Direct:** pending drain → `$bench-what-next`; otherwise carry the current row's literal entry phase when one exists. Do not derive a second priority order. |
| `bench roadmap --context [--full]` | Cl+Cx | S,E,R,St | Schema/trust flag, every source state, row truncation, counts, and recommended sequence | **Yes:** truncated trusted snapshot → `bench roadmap --context --full`; untrusted source → retry after the named source is repaired; trusted full snapshot → `$bench-what-next`. |
| `bench canary [root]` | Cl+Cx | S,E,R | Root, complete registered fixture/check inventory, and every fixture-specific failure | **Direct:** red → repair the named fixture owner then rerun `bench canary <root>` carrying a supplied root; green → **No**. It must not suggest a gate. |

## Adoption and installation

| command | log | cases | state already known at return | exact contextual action and turn effect |
|---|---|---|---|---|
| `bench setup [--plan or --yes]` | — | S,E,R,St | Adoption plan, confirmation mode, kit version, and conflicting footprint | **Yes:** plan → `bench setup --yes`; converged → `bench status`; conflict → the exact non-interactive setup invocation after the named conflict is resolved. |
| `bench link [copy or symlink]` | Cl | S,E,R,St | Link mode, managed manifest, kit version, and changed/stale paths | **Direct:** success/no-op → `bench status`; invalid mode → `bench link copy` or `bench link symlink`, without inventing a default. |
| `bench init` | — | S,E,R | Whether the project gate already exists and which scaffold was written | **Direct:** success/no-op → `bench setup`; refusal → fix the named existing gate conflict, then retry `bench init`. |
| `bench doctor [--fix]` | Cl | S,E,R,St,Rc | PATH-shim state (`healthy`, `stale`, `foreign`, or `missing`) and fixability | **Yes:** fixable stale/missing → `bench doctor --fix`; healthy/foreign → **No** because overwriting a foreign shim is not authorized. |
| `bench repair [--prune]` | Cl | S,E,R,St,Rc | Platform, pinned version, binary-cache state, and repair/prune outcome | **Direct:** repaired → `bench doctor`; stale cache → `bench repair --prune`; unsupported/offline refusal → no guessed install command. |
| `bench unlink [--dry-run]` | — | S,E,R,St | Manifest-owned removal set and dry-run/apply mode | **Direct:** dry run → `bench unlink`; completed/no-op → **No**; drift → rerun `bench unlink --dry-run` after resolving the named manifest conflict. |
| `bench upgrade [--check] [--force]` | — | S,E,R,St | Installed/source versions, upgrade plan, downgrade posture, and managed-file drift | **Yes:** check → `bench upgrade`; refused downgrade → `bench upgrade --force` only because the command already proved that exact condition; converged → `bench status`. |

## Gate, work execution, and worktrees

| command | log | cases | state already known at return | exact contextual action and turn effect |
|---|---|---|---|---|
| `bench gate` | Cl+Cx | S,E,R,St,Rc | Exact subject, reusable/inherited/executed phase evidence, first red phase, and verdict freshness | **Yes:** pending/absent/stale → retry `bench gate`; reusable exact green → the caller's landing command; a reusable green that cannot satisfy an explicit fresh boundary → `bench gate --fresh`. Never advertise `--fresh` merely because a gate is red. |
| `bench gate --fresh` | Cl+Cx | S,R | Exact subject and a forced complete execution outcome | **Direct:** green → the caller's landing command; red → the focused command for the named first failing phase when one exists, otherwise no guessed repair. |
| `bench gate pin` | Cl | S,R,St | HEAD and the exact `.bench` subject being pinned | **Direct:** success → `git push`; drift → rerun `bench gate pin` only after the named `.bench` change is resolved. |
| `bench worktree [--refresh] [objective]` | — | S,E,R,St,Rc | Objective, lease, selected pool slot, child result, and cleanup/recovery outcome | **Direct:** retained work → `bench worktree recovery <ref>`; capacity refusal → `bench worktree list`; clean release → **No**. |
| `bench worktree list` | Cl+Cx | S,E,R,St,Rc | Assignment id, path, branch/state, ownership, and out-of-pool rows | **Yes:** active row → `bench worktree path <id>` or `bench worktree exec <id> -- <command>`; orphan/retained row → `bench worktree clean <path>`; empty → **No**. |
| `bench worktree path <target>` | Cl+Cx | S,R,St | Target and safe resolved path, or proof it is not an active owned worktree | **Yes:** success → `bench worktree exec <target> -- <command>` carrying the target; refusal → `bench worktree list`. |
| `bench worktree exec <target> -- <command> [args...]` | Cl | S,R,St | Target, resolved path, and the child argv/exit | **No** for child success/failure: the child owns its next action. Resolution refusal → `bench worktree list`. |
| `bench worktree create [--refresh] --request <id> --label <item>` | Cl+Cx | S,E,R,St | Request, assignment id/path/branch/base, slot capacity, and refresh state | **Yes:** success → `bench worktree exec <assignment> -- <command>`; capacity refusal → `bench worktree list`; refresh refusal → repeat with the command's exact required receipt. |
| `bench worktree release --request <id> <path>` | Cl+Cx | S,E,R,St,Rc | Request, assignment, cleanup action, retained path, and recovery ref if work survives | **Yes:** retained work → the exact `bench worktree clean <path>` or `bench worktree recovery <ref>` already named by the error; completed/no-op → `bench worktree list`. |
| `bench worktree clean [flags] <path> [--apply <fingerprint>]` | Cl+Cx | S,E,R,St,Rc | Exact target, discard flags, action, tracked/ignored facts, recovery ref, and plan fingerprint | **Direct:** plan → repeat the full invocation with `--apply <fingerprint>` carrying path, flags, and fingerprint; stale fingerprint → rerun the plan without `--apply`; retain/error → do not suggest destructive flags the plan did not authorize. |
| `bench worktree recovery <ref> [--apply or --discard <fingerprint>]` | Cl | S,E,R,St,Rc | Exact ref, payload count, landed state, change summary, authorized action, and fingerprint | **Direct:** plan → exact apply or discard command carrying ref/fingerprint; stale fingerprint → plan again; completed/no-op → `bench worktree list`. |
| `bench shift [--refresh] "<objective>"` | — | S,R,St,Rc | Objective, assignment/worktree, iteration, gate outcome, commit, and recovery ref | **Yes:** interrupted/retained → `bench worktree recovery <ref>`; stale assignment → repeat the objective with `--refresh`; terminal green → **No**. |
| `bench commit -m <msg> [--spec <slug>] <path>...` | Cl+Cx | S,E,R,St | Exact named paths, dirty/foreign state, gate subject/verdict, commit id, and optional spec slug | **Yes:** green landing → `$bench-final-check`; stale/absent gate evidence → `bench gate`; foreign dirty paths → retry the same named-path command after the reported owner resolves them. It must not suggest adding unrelated paths. |

## Spec and spec-build lifecycle

Direct command extraction found all nine spec-build operations in Claude logs;
Codex logs found every operation except `abandon` and `reclaim`. The lower-bound
direct counts were especially high for `assign`, `checkpoint`, `integrate`, and
`status`, so this family is the strongest observed target for exact state-carrying
help. The current `Status.Next` owner already derives some commands, but also emits
non-invokable prose (`release assignment`, `delegate assignment`, and gate-diagnosis
labels) (`internal/specbuild/state.go:74-114`).

| command | log | cases | state already known at return | exact contextual action and turn effect |
|---|---|---|---|---|
| `bench spec implemented <slug>` | Cl | S,E,R,St | Slug, current status, and whether the implemented transition is valid | **Direct:** success/no-op → `bench spec retire <slug>` only when the command also proves retirement is reachable; otherwise `bench spec history <slug>` is the terminal inspection path. |
| `bench spec retire <slug>` | Cl+Cx | S,E,R,St,Rc | Slug, staged/implemented/merged state, owned review artifact, and removal eligibility | **Direct:** success/no-op → `bench spec history <slug>`; refusal → the exact missing precondition command when derivable, never generic deletion advice. |
| `bench spec history <slug>` | Cl+Cx | S,E,R | Slug and all matching retire/delete commits | **No** on success/empty; it is terminal provenance. Usage carries `bench spec history <slug>`. |
| `bench spec build` with no operation | Cl+Cx | R | Complete ordered operation grammar but no slug | **Yes:** `bench spec build status <slug> --full`; this is the family home template, with the genuinely unknown slug left open. |
| `bench spec build start <slug>` | Cl+Cx | S,E,R,St | Slug, staged spec, current tip, exact-green evidence class, existing/terminal run, and candidate | **Yes:** success/no-op → `bench spec build status <slug> --full`; exact-green refusal → the command that can actually satisfy the observed evidence class. The Claude FT164 trace proves that advertising plain `bench gate` when only `bench gate --fresh` works causes a wrong-remedy turn. |
| `bench spec build assign <slug> --ticket <ticket> --request <id> [--refresh <receipt>]` | Cl+Cx | S,E,R,St,Rc | Slug, ticket/fence/digest, request, capacity, candidate, assignment id/path/base, and refresh receipt state | **Direct:** success → `bench spec build checkpoint <slug> --assignment <id> --evidence <receipt>` with slug/id filled and only future evidence open; capacity/cleanup refusal → exact integrate/release action for the named assignment; stale repair → exact refresh form. |
| `bench spec build checkpoint <slug> --assignment <id> --evidence <receipt>` | Cl+Cx | S,E,R,St,Rc | Slug, assignment, fence, ticket digest, candidate/base, provisional commit, and evidence receipt result | **Direct:** success/no-op → `bench spec build integrate <slug> --assignment <id>`; stale/invalid evidence → repeat the same checkpoint command after replacing only the rejected receipt. |
| `bench spec build integrate <slug> --assignment <id>` | Cl+Cx | S,E,R,St,Rc | Slug, assignment, checkpoint/ref, candidate before/after, compare-and-swap outcome, and cleanup state | **Yes:** success → exact next assign command when tickets remain, otherwise `bench spec build review <slug> --evidence <receipt>`; moved candidate → refresh/recheckpoint the named assignment rather than generic retry. |
| `bench spec build review <slug> --evidence <receipt>` | Cl+Cx | S,E,R,St,Rc | Slug, exact candidate, all three axes/findings/dispositions, and receipt freshness | **Direct:** clean/accepted residual risk → `bench spec build promote <slug>`; accepted issue needing repair → `bench spec build assign <slug> --ticket <ticket> --request <id>` with only the new ticket/request unknown; stale candidate → rerun review, never promote. |
| `bench spec build status <slug> [--full]` | Cl+Cx | S,E,R,St,Rc | Slug, run/candidate/terminal state, assignments, checkpoints, cleanup, review, promotion disposition, and prepared operations | **Yes:** emit invokable templates derived from those facts: assign, checkpoint, integrate, review, promote, release/repair, or resume the exact prepared command. Replace orchestration labels with typed actions rather than copying prose into `help[]`. No-run status should suggest `bench spec build start <slug>`. |
| `bench spec build promote <slug>` | Cl+Cx | S,E,R,St,Rc | Slug, exact candidate/review, working tip, recomposition, gate disposition/evidence, and terminal landing state | **Yes:** terminal success → the workflow-owned retro/finalization action; moved candidate → review the newly composed candidate; candidate/inherited/infrastructure gate failure → its fenced repair or exact retry command. It must not advertise another gate because promotion owns the sole prospective gate. |
| `bench spec build abandon <slug> [--apply <fingerprint>]` | Cl | S,E,R,St,Rc | Slug, terminal run, every owned worktree/ref/checkpoint/recovery ref, plan fingerprint, and partial cleanup state | **Direct:** plan → `bench spec build abandon <slug> --apply <fingerprint>` carrying both values; stale/partial apply → plan again; completed → `bench spec build status <slug> --full`. |
| `bench spec build reclaim <slug> [--apply <fingerprint>]` | Cl | S,E,R,St,Rc | Slug, terminal run, deletable/retained ref classes, plan fingerprint, applied set, and spent partial receipt | **Direct:** plan → exact apply command; interrupted partial apply → `bench spec build reclaim <slug>` to re-plan, as the current refusal already says; completed/no-op → terminal status/history. |

## Release surfaces

The Claude archive contains direct `bench release` help/refusal probes but no
directly extracted `prepare`, `submit`, `promote`, `rollback`, or `status` operation.
The operations remain in the inventory because they are declared by the production
dispatcher, not because logs sampled them.

| command | log | cases | state already known at return | exact contextual action and turn effect |
|---|---|---|---|---|
| `bench prep-release` | Cl | S,E,R,St | Dev-green subject, artifact/preflight/canary results, package version, and ship-tier eligibility | **Direct:** green → `bench release prepare --version <version>` with the derived version filled; refusal → the exact failed prerequisite command, not an unconditional fresh gate. |
| `bench release-preflight --mode verify [--profile <profile>] [--phase <name>]` | — | S,E,R,St | Mode, profile, focused/full scope, subject, and release authorization result | **Direct:** full verify success → `bench release-preflight --mode publish --profile <profile>`; focused success → **No**, because it cannot authorize publication. |
| `bench release-preflight --mode publish [--profile <profile>]` | — | S,E,R,St | Publish authority, profile, approved index, and exact subject | **Direct:** success → `bench release prepare --version <version> --profile <profile>` carrying known values; stale/refusal → repeat the publish preflight only after its named prerequisite. |
| `bench release prepare --version <version> [...]` | — | S,E,R,St | Version, approved artifact set, release-index digest, and current `next_action=submit` | **Yes:** replace the bare word with `bench release submit --version <version> --profile <profile> --registry <base-url> --path <path>`, carrying every supplied/derived fixed value. |
| `bench release submit --version <version> [...]` | — | S,E,R,St,Rc | Version/profile/path, durable transition record, package states, registry, and next action | **Yes:** emit the exact resume/status/promote command from the durable record with version/profile/path/registry carried forward. |
| `bench release promote --version <version> [...]` | — | S,E,R,St,Rc | Version/profile, live verification, tags, record, and whether latest moved | **Direct:** success/no-op → `bench release status --version <version> --root <root>`; refusal → exact submit/preflight recovery derived from the record. |
| `bench release rollback --version <version> [...]` | — | S,E,R,St,Rc | Version/profile, tags, deprecations, record, message, and partial rollback state | **Direct:** success/partial result → `bench release status --version <version> --root <root>`; lock/stale refusal → retry the same rollback only after the named condition clears. |
| `bench release status --version <version> [...]` | — | S,E,R,St,Rc | Durable record, transitions, profile/path/result, and current `next_action` | **Yes:** expand the current word into the exact prepare/submit/promote/rollback template carrying record values; no record → `bench release prepare --version <version>`. |

## Hook, adapter, and registry-only plumbing

These commands are exhaustive because they are in the Go registry or a declared
nested grammar. `.bench/BENCH-reference.md:115-121` says sessions never type the
plumbing set. Their caller normally owns sequencing, so forcing agent-facing
`help[]` onto successful plumbing output would duplicate control flow. Refusals
should still name a public recovery command when the plumbing layer already knows
one.

| command | log | cases | state already known at return | exact contextual action and turn effect |
|---|---|---|---|---|
| `bench tree-hash` | Cx | S,R,St | Exact tree/hash input and resulting digest | **No**; terminal machine input for its caller. |
| `bench resolve-model --harness <harness>` | Cl | S,E,R,St | Harness, tier request, binding file, and resolved/unbound model | **Yes on unbound:** `bench models`; success is caller-terminal. |
| `bench worktree-pool` | — | S,E,R,St | Pool root, capacity, leases, and slot facts | **No**; the worktree caller chooses create/release. |
| `bench worktree-lease-file <path>` | — | S,R | Exact worktree and lease-file path | **No**; terminal adapter lookup. |
| `bench resume-clean` | Cl | S,E,R,St,Rc | All recoverable registrations, removed/recovered counts, and preserved refs | **Yes:** each preserved row already names `bench worktree recovery <ref>`; render it structurally instead of prose. |
| `bench session-inspect` | — | S,E,R,St,Rc | Resume result, status signals, guard summary, and deadline/partial state | **Yes through composition:** forward public actions from `status`/resume; do not invent a separate session-inspect remedy. |
| `bench gate-go gofmt [root]` | Cl | S,R | Root and exact formatting result | **No**; owning gate phase decides continuation. |
| `bench gate-go test [root]` | Cl | S,R | Root and exact ordinary test result | **No**; owning gate phase decides continuation. |
| `bench guard-git` | — | S,R | Parsed Git command and dangerous-operation class | **No**; a refusal intentionally requires changed intent, not a bypass command. |
| `bench check-agent-line --harness <harness>` | — | S,E,R,St | Harness, armed shift, requested tier/model, and binding | **Yes on mismatch/unbound:** `bench resolve-model --harness <harness>`; success is hook-terminal. |
| `bench worktree-hook create` / `bench worktree-hook remove` | — | S,E,R,St,Rc | Hook event, worktree path, assignment/lease, and cleanup result | **Yes only for preserved work:** `bench worktree recovery <ref>`; otherwise caller-terminal. |
| `bench gate-run [--fresh]` | Cl | S,E,R,St,Rc | Same gate subject/evidence as the public wrapper plus run ownership | **No separate help owner:** compose the public `bench gate` action derivation. |
| `bench gate-pin` | Cl | S,R,St | Same pin facts as public `bench gate pin` | **No separate help owner:** compose the public pin result; success action is `git push`. |
| `bench gate-phases` | Cl+Cx | S,E,R,St | Root, selected executable, phase table, dependency outcomes, and first red | **No**; `gate-run` owns the public recovery projection. |
| `bench freshness-check <root>` | Cl+Cx | S,R,St | Root, executable identity, source hash, and stale reason | **No**; the owning build/gate chooses rebuild. A direct agent call cannot safely infer the required artifact owner. |
| `bench freshness-publish <root> <output-path>` | — | S,R,St | Root, selected executable, output path, and published identity | **No**; registry-only build plumbing and not routed by the wrapper. Its presence is an inventory inconsistency to resolve, not a public AXI migration target by assumption. |
| `bench stop-verdict` | — | S,E,R,St | Armed shift state, exact gate verdict, and whether stopping would leave red work | **Yes on refusal:** `bench gate` when evidence is absent/stale, otherwise continue the armed shift; success is hook-terminal. |

## Cross-command findings

1. **The first reusable owner should be typed actions, not rendered prose.** Status,
   handoff, spec-build status, release status, worktree resume/cleanup, and
   structured errors already know identifiers needed by the next command. They
   independently flatten those facts into `action`, `next`, `next_action`, detail
   text, or stderr. A `help[]` renderer cannot recover fixed arguments reliably
   after that flattening.
2. **The Claude FT164 trace is the decisive wrong-remedy proof.** A structured
   error alone was insufficient: the command knew the evidence class but named a
   command that could not satisfy it. Contextual disclosure must be derived from
   the same typed precondition result as the refusal, not from a generic retry
   suffix.
3. **Plan/apply calls need exact carry-forward, not fewer mutations.** Worktree
   clean/recovery and spec-build abandon/reclaim must remain two calls. Their
   `help[]` value is eliminating command reconstruction: the plan owns every fixed
   flag, target, and fingerprint and should emit the one authorized apply template.
4. **Terminal queries should not invent busywork.** Anchors, completed history,
   complete outline, version, and successful plumbing lookups have no universally
   useful follow-up. “Every command is inventoried” does not mean every success
   should advertise `bench status`; empty `help[0]` is more honest than a low-value
   loop.
5. **Lifecycle and release words are not commands.** `release assignment`,
   `delegate assignment`, `retry promote`, `submit`, and `promote` need a typed
   executable-template representation before they can satisfy AXI contextual
   disclosure. Copying the existing cells into `help[]` would preserve ambiguity.
6. **Log absence cannot set scope.** The public release operations, adoption
   commands, default worktree subshell, shift, and most plumbing commands were not
   directly observed in one or both available corpora. They remain contract
   surfaces because the live dispatcher exposes them. The missing reviewer-named
   Codex corpus remains a disclosed research limitation for compatibility ticket
   #8, not a reason to sample away commands here.

## Sources

- `cmd/bench/main.go:66-225` and `cmd/bench/command_registry.go:23-73` — production
  root registry, attachment boundary, and incomplete `commands --brief` body.
- `bin/bench.sh:211-217,284-390` — wrapper-only repair/help and public routing.
- `internal/usage/worktree.go:5-26`, `internal/spec/spec.go:291-302`,
  `internal/spec/build.go:16-89`, and `internal/publication/command.go:15-113` —
  nested public operation grammars.
- `internal/specbuild/state.go:74-114`, `internal/status/status.go:133-176`,
  `internal/status/handoff.go`, `internal/handoff/facts.go:159-209`, and
  `internal/publication/command.go:293-335` — existing next-action derivations.
- Claude session `0b2828ac-933f-404e-8de4-c785735406ac`, JSONL lines 257, 266,
  and 288 — wrong exact-green remediation and reviewer acceptance boundary.
- Local Claude Bench project archive (100 JSONL transcripts) and local Codex
  session archive (626 rollouts; 91 with project-scoped command executions) —
  lower-bound direct invocation marks, inspected 2026-08-09.
