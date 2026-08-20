#!/usr/bin/env bash
# name: block-dangerous-git
# boundary: PreToolUse:Bash
# denies: destructive git operations
# why: agents lack destructive-git authority; merge and history rewrites belong to the reviewer, and discarding work detaches the gate verdict from the tree
# PreToolUse guard: the agent has no destructive git authority. This makes
# invariant #4 ("you assist, you don't decide where a decision is mine")
# enforceable for the operations that can silently destroy a shift's work or
# bypass the merge — not just aspirational.
#
# Threat model: this is an honest-mistake layer, not an evasion-resistant
# boundary. It stops a well-meaning agent from reflexively running destructive
# git; it does not try to survive deliberate evasion. Wrapper scanning goes
# exactly one level deep by design (see internal/gitguard). The backstops for a
# misaligned agent are the git pre-push hook and bench's pooled-worktree
# isolation, not this script.
#
# Note the boundary: this intercepts the AGENT's Bash tool calls. Bench's own
# controlled rollback inside `bench shift` runs in-process (not through the
# agent's shell), so the harness can still reset/clean a failed iteration while
# the agent itself cannot. That asymmetry is the point.
#
# This is a thin shim over the Go core: it resolves the bench wrapper, pipes the
# PreToolUse envelope to `bench guard-git`, and passes the verdict through. Every
# classification the guard makes with a reachable core (tokenize, scan, verdict, the
# BLOCKED message) lives in internal/gitguard. The shim owns exactly two fail-closed
# rims — core unresolvable/missing, and core errored — and the first carries its own
# coarser restatement of the question, because it runs precisely when the core cannot
# answer it. Wire under
# hooks.PreToolUse with matcher "Bash". Exit 2 blocks and returns the message to
# the agent.
set -uo pipefail

# resolve_wrapper echoes the bench.sh wrapper path (repo-local first, then a
# global `bench`), or fails when none is reachable. The ~8-line search is inlined
# rather than shared with .bench/hooks/stop.sh: sourcing a shared lib would give
# this hook a new fail-OPEN mode (missing lib → the shim errors before its rims
# run, and a non-2 PreToolUse exit is a non-blocking error that silently grants).
# The conformance check in internal/conformance (checkGuardResolverOrderDrift)
# reds if this inline's search order ever drifts from .bench/lib/resolve-bench.sh.
resolve_wrapper() {
  local root candidate
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  if [[ -n "$root" ]]; then
    for candidate in "$root/.bench/bin/bench.sh" "$root/bin/bench.sh"; do
      [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
    done
  fi
  command -v bench 2>/dev/null || return 1
}

input="$(cat)"

# fail_closed_no_core is the "cannot classify" rim. It runs only when the core is
# unreachable, so it cannot delegate to internal/gitguard and restates the question in
# shell — deliberately coarser: it refuses *every* git invocation rather than only the
# destructive ones, a strictly wider verdict than the core's, which is why it can stay
# this small. It decides from the envelope's `tool_input.command` field, so `git` inside
# a path, an argument, or another envelope field no longer refuses an ordinary read
# during the one session that has to recover the core. An envelope with no readable
# command field is refused, so the fail-closed posture survives the narrowing. Threat
# model is unchanged: honest mistakes, with the pre-push hook and pooled-worktree
# isolation as the backstops.
fail_closed_no_core() {
  local command_text
  if ! command_text="$(envelope_command)"; then
    echo "BLOCKED: guard degraded (bench core missing) — no readable command in this tool call, so refusing it. Restore the bench core (bench link) or hand back." >&2
    exit 2
  fi
  if invokes_git "$command_text" 1; then
    echo "BLOCKED: guard degraded (bench core missing) — can't classify commands, refusing anything that invokes git. Restore the bench core (bench link) or hand back." >&2
    exit 2
  fi
  exit 0
}

# envelope_command prints the string value of the PreToolUse envelope's
# `tool_input.command` and fails when the envelope carries no readable one — object
# absent, field absent, value not a string, string unterminated, or an escape it cannot
# decode. Failing is a refusal upstream, so an envelope this cannot read is never an
# empty command.
envelope_command() {
  local rest before value ch esc i n hex dec oct chr
  rest=${input#*'"tool_input"'}
  [[ "$rest" == "$input" ]] && return 1
  before=$rest
  rest=${rest#*'"command"'}
  [[ "$rest" == "$before" ]] && return 1
  while [[ "$rest" == [[:space:]]* ]]; do rest=${rest#?}; done
  [[ "$rest" == :* ]] || return 1
  rest=${rest#:}
  while [[ "$rest" == [[:space:]]* ]]; do rest=${rest#?}; done
  [[ "$rest" == '"'* ]] || return 1
  rest=${rest#\"}
  value=''
  n=${#rest}
  i=0
  while (( i < n )); do
    ch=${rest:i:1}
    if [[ "$ch" == '"' ]]; then
      printf '%s' "$value"
      return 0
    fi
    if [[ "$ch" == '\' ]]; then
      esc=${rest:i+1:1}
      case "$esc" in
        n) value=$value$'\n' ;;
        t) value=$value$'\t' ;;
        r) value=$value$'\r' ;;
        b | f) value="$value " ;;
        # A \uXXXX escape in the ASCII range decodes to its own byte: Go's
        # encoding/json escapes & < > by default, so a control operator or a
        # command name reaches this decoder escaped as often as literal, and a
        # placeholder there hides the operator the scan below looks for. Above
        # ASCII the placeholder stands, since no such rune is a shell operator.
        # \u0000 refuses: bash cannot carry NUL, so decoding it would silently
        # truncate the command this guard is deciding about.
        u)
          hex=${rest:i+2:4}
          [[ "$hex" == [0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f] ]] || return 1
          if [[ "$hex" == 00[0-7][0-9A-Fa-f] ]]; then
            printf -v dec '%d' "0x$hex"
            (( dec == 0 )) && return 1
            printf -v oct '%03o' "$dec"
            printf -v chr "\\$oct"
            value=$value$chr
          else
            value="${value}_"
          fi
          i=$(( i + 4 ))
          ;;
        '"' | '\' | /) value="$value$esc" ;;
        *) return 1 ;;
      esac
      i=$(( i + 2 ))
      continue
    fi
    value="$value$ch"
    i=$(( i + 1 ))
  done
  return 1
}

# invokes_git succeeds when the command text runs `git` in command position. $2 enables
# the wrapper recursion, which goes exactly one level deep — the depth internal/gitguard
# documents. The scan walks command by command, skipping the honest-mistake prefixes an
# agent reflexively types, so `git` as an argument or a path component is not a match.
invokes_git() {
  local allow_wrapper=$2 flag
  local -a tok
  lex_command "$1"
  tok=()
  (( ${#LEXED[@]} )) && tok=("${LEXED[@]}")
  local n=${#tok[@]} i=0 expect=1 word base end j k
  while (( i < n )); do
    word=${tok[$i]}
    if is_control_op "$word"; then
      expect=1
      i=$(( i + 1 ))
      continue
    fi
    if (( expect )) && is_keyword "$word"; then
      i=$(( i + 1 ))
      continue
    fi
    if (( expect )); then
      end=$i
      while (( end < n )) && ! is_control_op "${tok[$end]}"; do end=$(( end + 1 )); done
      j=$(resolve_prefixes "$i" "$end" "${tok[@]}")
      if (( j < end )); then
        base=${tok[$j]##*/}
        if [[ "$base" == git ]]; then
          return 0
        fi
        if (( allow_wrapper )) && [[ "$base" == sh || "$base" == bash || "$base" == zsh ]]; then
          k=$(( j + 1 ))
          while (( k < end )); do
            flag=${tok[$k]}
            if [[ "$flag" == -* && "${flag#-}" == *c* && "${flag#-}" != *[!A-Za-z]* ]]; then
              (( k + 1 < end )) && invokes_git "${tok[$(( k + 1 ))]}" 0 && return 0
              break
            fi
            k=$(( k + 1 ))
          done
        fi
      fi
      expect=0
      i=$end
      continue
    fi
    i=$(( i + 1 ))
  done
  return 1
}

# is_control_op reports whether a token ends the current command, so the next word is a
# fresh command position.
is_control_op() {
  case "$1" in
    ';' | '&&' | '||' | '|' | '&' | '(' | ')') return 0 ;;
    *) return 1 ;;
  esac
}

# is_keyword reports whether a token is a shell keyword skipped in command position, so
# the verb after it is the one that gets read (`if git …`).
is_keyword() {
  case "$1" in
    if | then | elif | else | do | while | until | '!' | '{') return 0 ;;
    *) return 1 ;;
  esac
}

# resolve_prefixes echoes the index of the real verb in tokens $3.. between $1 and $2,
# stepping past leading environment assignments and the command wrappers an agent types
# in front of the verb it means (env/command/nohup/timeout/xargs).
resolve_prefixes() {
  local i=$1 end=$2
  shift 2
  local -a tok=("$@")
  while (( i < end )); do
    if [[ "${tok[$i]}" == [A-Za-z_]*=* && "${tok[$i]%%=*}" != *[!A-Za-z0-9_]* ]]; then
      i=$(( i + 1 ))
      continue
    fi
    case "${tok[$i]##*/}" in
      env)
        i=$(( i + 1 ))
        while (( i < end )) && [[ "${tok[$i]}" == [A-Za-z_]*=* && "${tok[$i]%%=*}" != *[!A-Za-z0-9_]* ]]; do i=$(( i + 1 )); done
        ;;
      command | nohup) i=$(( i + 1 )) ;;
      timeout)
        i=$(( i + 1 ))
        while (( i < end )) && [[ "${tok[$i]}" == -* ]]; do i=$(( i + 1 )); done
        (( i < end )) && i=$(( i + 1 ))
        ;;
      xargs)
        i=$(( i + 1 ))
        while (( i < end )) && [[ "${tok[$i]}" == -* ]]; do i=$(( i + 1 )); done
        ;;
      *) break ;;
    esac
  done
  printf '%s' "$i"
}

# lex_command splits a command line into LEXED: words with quotes and escapes folded in
# so a wrapper's `-c` string survives as one token, and runs of the shell operator
# characters as their own tokens. A bare newline lexes as an operator, so a multi-line
# block scans as separate commands. Unbalanced quoting falls back to a plain split that
# still honors newline boundaries.
lex_command() {
  local text=$1 cur='' active=0 i=0 start n ch
  LEXED=()
  n=${#text}
  while (( i < n )); do
    ch=${text:i:1}
    case "$ch" in
      ' ' | $'\t' | $'\r')
        (( active )) && { LEXED+=("$cur"); cur=''; active=0; }
        i=$(( i + 1 ))
        ;;
      '(' | ')' | ';' | '<' | '>' | '|' | '&' | $'\n')
        (( active )) && { LEXED+=("$cur"); cur=''; active=0; }
        start=$i
        while (( i < n )); do
          case "${text:i:1}" in
            '(' | ')' | ';' | '<' | '>' | '|' | '&' | $'\n') i=$(( i + 1 )) ;;
            *) break ;;
          esac
        done
        LEXED+=("$(collapse_operator "${text:start:i-start}")")
        ;;
      "'")
        active=1
        i=$(( i + 1 ))
        while (( i < n )) && [[ "${text:i:1}" != "'" ]]; do
          cur="$cur${text:i:1}"
          i=$(( i + 1 ))
        done
        (( i >= n )) && { lex_fallback "$text"; return 0; }
        i=$(( i + 1 ))
        ;;
      '"')
        active=1
        i=$(( i + 1 ))
        while (( i < n )) && [[ "${text:i:1}" != '"' ]]; do
          if [[ "${text:i:1}" == '\' ]] && [[ "${text:i+1:1}" == '"' || "${text:i+1:1}" == '\' ]]; then
            cur="$cur${text:i+1:1}"
            i=$(( i + 2 ))
            continue
          fi
          cur="$cur${text:i:1}"
          i=$(( i + 1 ))
        done
        (( i >= n )) && { lex_fallback "$text"; return 0; }
        i=$(( i + 1 ))
        ;;
      '\')
        (( i + 1 >= n )) && { lex_fallback "$text"; return 0; }
        active=1
        cur="$cur${text:i+1:1}"
        i=$(( i + 2 ))
        ;;
      *)
        active=1
        cur="$cur$ch"
        i=$(( i + 1 ))
        ;;
    esac
  done
  (( active )) && LEXED+=("$cur")
  return 0
}

# collapse_operator drops the newlines from an operator-only token so `&&`+newline reads
# as the control operator it is, and a bare newline reads as a command boundary.
collapse_operator() {
  local op=${1//$'\n'/}
  [[ -z "$op" ]] && op=';'
  printf '%s' "$op"
}

# lex_fallback is the unbalanced-quoting path: split each line on whitespace and keep
# every newline as a command boundary, so a multi-line block still scans command by
# command. It drops quote folding, exactly as internal/gitguard's tokenizer does.
lex_fallback() {
  local line first=1
  local -a words
  LEXED=()
  while IFS= read -r line || [[ -n "$line" ]]; do
    (( first )) || LEXED+=(';')
    first=0
    words=()
    read -r -a words <<< "$line"
    (( ${#words[@]} )) && LEXED+=("${words[@]}")
  done <<< "$1"
}

# Rim 1: core unresolvable. No wrapper on disk or PATH → cannot classify.
cmd="$(resolve_wrapper)" || fail_closed_no_core

# Hand the envelope to the core. The core writes its own BLOCKED message to stderr
# and exits 2 on a block, 0 on allow; the wrapper exits 127 when no binary is
# installed for this platform.
rc=0
printf '%s' "$input" | "$cmd" guard-git || rc=$?
case "$rc" in
  0 | 2) exit "$rc" ;;                 # allow / block — the core owns the verdict + message
  127) fail_closed_no_core ;;             # rim 1: binary missing for this platform
  *) echo "BLOCKED: guard analyzer error — failing closed; rephrase the command." >&2; exit 2 ;;  # rim 2: core errored
esac
