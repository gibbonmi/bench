# shellcheck shell=sh
# This is the shared shell classification library for every hook shim
# (.bench/hooks/*.sh) and shift adapter (.bench/adapters/*). It carries the one
# source of the wrapper search order — repo `.bench/bin/bench.sh`, then kit
# `bin/bench.sh`, then `bench` on PATH — the shell Bench-call word test, the
# PreToolUse envelope read, and the shell spelling of the rebuild invocation. It
# carries no fail posture: what a shim does with either answer is the shim's own.
#
# Each shim owns its own fail posture when the search comes up empty. The git and
# agent-line guards refuse or warn-and-allow. `stop` and `session-start` warn and
# fail open. The adapters fail closed. No shim sources this blindly: each
# source-guards it under that posture, because a missing lib that let the shim
# error before its rims run would turn a fail-closed guard fail-open.
#
# This file is POSIX sh; the adapters, which must stay POSIX-clean, source it with
# `.`. It is meant to be sourced, never run directly. Its variables are
# underscore-prefixed, so sourcing it cannot clobber an adapter's locals. It defines
# bench_resolve_wrapper: this prints the wrapper path and returns 0, or returns 1
# when none of the three candidates resolves to an executable.
bench_resolve_wrapper() {
  _bench_root=$(git rev-parse --show-toplevel 2>/dev/null) || _bench_root=
  if [ -n "$_bench_root" ]; then
    for _bench_candidate in "$_bench_root/.bench/bin/bench.sh" "$_bench_root/bin/bench.sh"; do
      if [ -x "$_bench_candidate" ]; then
        printf '%s\n' "$_bench_candidate"
        return 0
      fi
    done
  fi
  command -v bench 2>/dev/null || return 1
}

# bench_invokes_bench reports 0 when the Bash command text in $1 runs Bench, and 1
# when it does not. It is the shell derivation of internal/benchguard.InvokesBench,
# which a shim reaches for when the core binary cannot answer at all.
#
# Its reach is the head word of every control-operator-delimited segment, after the
# shell assignments and the routine prefixes env, command, nohup, timeout, and
# xargs. Each word is unquoted first, so a quoted or backslash-escaped head reads as
# the name the shell would run. It resolves no path and reads no wrapper string, so a
# `bench` reached through a symlink or through `bash -c` is outside it. So is a head
# inside a subshell or a command substitution, which is grammar rather than quoting
# and needs the lexer this half does not have. The two derivations are
# pinned row by row by internal/conformance/guard_classifier_table_test.go, whose
# table therefore holds resolver-independent, single-level rows only.
bench_invokes_bench() {
  _bench_stream=$(printf '%s' "$1" | tr ';|&' '\n')
  _bench_ifs=$IFS
  _bench_verdict=1
  IFS='
'
  set -f
  for _bench_segment in $_bench_stream; do
    if bench_segment_runs_bench "$_bench_segment"; then
      _bench_verdict=0
      break
    fi
  done
  set +f
  IFS=$_bench_ifs
  return $_bench_verdict
}

# bench_segment_runs_bench reports whether one simple command's head is Bench.
bench_segment_runs_bench() {
  # The caller splits segments on a newline IFS; words split on the default one.
  _bench_outer_ifs=$IFS
  unset IFS
  set -f
  # The split is the point: the segment is command text, not one word.
  # shellcheck disable=SC2086
  set -- $1
  set +f
  IFS=$_bench_outer_ifs
  # Word splitting keeps the quotes it split on, so every word is unquoted before the
  # reads below. Rotating the list replaces each word in place.
  _bench_word_budget=$#
  while [ "$_bench_word_budget" -gt 0 ]; do
    set -- "$@" "$(bench_unquote_word "$1")"
    shift
    _bench_word_budget=$(( _bench_word_budget - 1 ))
  done
  while [ "$#" -gt 0 ]; do
    if bench_is_assignment "$1"; then
      shift
      continue
    fi
    case ${1##*/} in
      env)
        shift
        while [ "$#" -gt 0 ]; do
          if [ "$1" = "--" ]; then
            shift
            break
          fi
          if bench_is_assignment "$1"; then
            shift
            continue
          fi
          case $1 in -?*) ;; *) break ;; esac
          case $1 in
            -u|--unset|-C|--chdir) shift; if [ "$#" -gt 0 ]; then shift; fi ;;
            *) shift ;;
          esac
        done
        ;;
      command)
        shift
        while [ "$#" -gt 0 ]; do
          if [ "$1" = "--" ]; then
            shift
            break
          fi
          case $1 in -?*) ;; *) break ;; esac
          # `command -v` and `command -V` query rather than execute.
          case ${1#-} in *[vV]*) return 1 ;; esac
          shift
        done
        ;;
      nohup)
        shift
        while [ "$#" -gt 0 ]; do
          if [ "$1" = "--" ]; then
            shift
            break
          fi
          case $1 in -?*) shift ;; *) break ;; esac
        done
        ;;
      timeout)
        shift
        while [ "$#" -gt 0 ]; do
          if [ "$1" = "--" ]; then
            shift
            break
          fi
          case $1 in -?*) ;; *) break ;; esac
          case $1 in
            -s|--signal|-k|--kill-after) shift; if [ "$#" -gt 0 ]; then shift; fi ;;
            *) shift ;;
          esac
        done
        # The duration operand sits between the options and the command.
        if [ "$#" -gt 0 ]; then shift; fi
        ;;
      xargs)
        shift
        while [ "$#" -gt 0 ]; do
          if [ "$1" = "--" ]; then
            shift
            break
          fi
          case $1 in -?*) ;; *) break ;; esac
          case $1 in
            -E|-I|-L|-P|-d|-n|-s|--eof|--replace|--max-lines|--max-procs|--delimiter|--max-args|--max-chars)
              shift; if [ "$#" -gt 0 ]; then shift; fi ;;
            *) shift ;;
          esac
        done
        ;;
      *)
        break
        ;;
    esac
  done
  [ "$#" -gt 0 ] || return 1
  case ${1##*/} in bench|bench.sh) return 0 ;; esac
  return 1
}

# bench_unquote_word prints $1 with its shell quoting removed, so the word test reads
# the name the shell would run rather than its spelling. A single-quoted run stands as
# written, a double-quoted run drops its delimiters and honours the four escapes that
# are special inside it, and a backslash elsewhere escapes the character after it.
#
# This is the shell derivation of the quote folding internal/shellcommand's tokenizer
# performs. It reads one word, never a command line, so it decides no grammar: a
# parenthesis or a dollar sign stays an ordinary character.
#
# A word whose quoting does not close stands unchanged. The caller split the command
# text on whitespace before this ran, so a quoted space arrives as two half-quoted
# words. Folding those halves would read `bench" "gate` as a Bench call, which the
# tokenizer reads as the single word `bench gate`.
bench_unquote_word() {
  _bench_quote_in=$1
  _bench_quote_word=''
  _bench_quote_state=bare
  while [ -n "$_bench_quote_in" ]; do
    _bench_quote_next=${_bench_quote_in#?}
    _bench_quote_char=${_bench_quote_in%"$_bench_quote_next"}
    _bench_quote_in=$_bench_quote_next
    case $_bench_quote_state in
      single)
        if [ "$_bench_quote_char" = "'" ]; then
          _bench_quote_state=bare
          continue
        fi
        ;;
      double)
        if [ "$_bench_quote_char" = '"' ]; then
          _bench_quote_state=bare
          continue
        fi
        if [ "$_bench_quote_char" = '\' ] && [ -n "$_bench_quote_in" ]; then
          _bench_quote_next=${_bench_quote_in#?}
          case ${_bench_quote_in%"$_bench_quote_next"} in
            '"'|'\'|'$'|'`')
              _bench_quote_char=${_bench_quote_in%"$_bench_quote_next"}
              _bench_quote_in=$_bench_quote_next
              ;;
          esac
        fi
        ;;
      *)
        case $_bench_quote_char in
          "'")
            _bench_quote_state=single
            continue
            ;;
          '"')
            _bench_quote_state=double
            continue
            ;;
          '\')
            [ -n "$_bench_quote_in" ] || continue
            _bench_quote_next=${_bench_quote_in#?}
            _bench_quote_char=${_bench_quote_in%"$_bench_quote_next"}
            _bench_quote_in=$_bench_quote_next
            ;;
        esac
        ;;
    esac
    _bench_quote_word=$_bench_quote_word$_bench_quote_char
  done
  if [ "$_bench_quote_state" != bare ]; then
    printf '%s' "$1"
    return 0
  fi
  printf '%s' "$_bench_quote_word"
}

# bench_is_assignment reports whether a word is a portable NAME=VALUE assignment.
bench_is_assignment() {
  case $1 in *=*) ;; *) return 1 ;; esac
  _bench_name=${1%%=*}
  case $_bench_name in
    ''|[!A-Za-z_]*|*[!A-Za-z0-9_]*) return 1 ;;
  esac
  return 0
}

# bench_envelope_command prints the string value of the PreToolUse envelope in $1, at
# its `tool_input.command` field. It returns 1 when the envelope carries no readable
# value:
#   - the object is absent
#   - the field is absent
#   - the value is not a string
#   - the string is unterminated
#   - an escape cannot be decoded
# It carries no fail posture. The git guard refuses on a 1, and the follow-on shim
# stays open on one, so neither reads an unreadable envelope as an empty command.
#
# The read is scoped under `"tool_input"`, so a sibling field such as `cwd` that
# follows that object never joins the command.
bench_envelope_command() {
  _bench_env_rest=${1#*'"tool_input"'}
  [ "$_bench_env_rest" = "$1" ] && return 1
  _bench_env_before=$_bench_env_rest
  _bench_env_rest=${_bench_env_rest#*'"command"'}
  [ "$_bench_env_rest" = "$_bench_env_before" ] && return 1
  bench_env_skip_space
  case $_bench_env_rest in :*) ;; *) return 1 ;; esac
  _bench_env_rest=${_bench_env_rest#:}
  bench_env_skip_space
  case $_bench_env_rest in '"'*) ;; *) return 1 ;; esac
  _bench_env_rest=${_bench_env_rest#\"}
  _bench_env_value=''
  while [ -n "$_bench_env_rest" ]; do
    _bench_env_tail=${_bench_env_rest#?}
    _bench_env_char=${_bench_env_rest%"$_bench_env_tail"}
    if [ "$_bench_env_char" = '"' ]; then
      printf '%s' "$_bench_env_value"
      return 0
    fi
    if [ "$_bench_env_char" != '\' ]; then
      _bench_env_value=$_bench_env_value$_bench_env_char
      _bench_env_rest=$_bench_env_tail
      continue
    fi
    _bench_env_rest=$_bench_env_tail
    [ -n "$_bench_env_rest" ] || return 1
    _bench_env_tail=${_bench_env_rest#?}
    _bench_env_escape=${_bench_env_rest%"$_bench_env_tail"}
    _bench_env_rest=$_bench_env_tail
    case $_bench_env_escape in
      n) _bench_env_octal=012 ;;
      t) _bench_env_octal=011 ;;
      r) _bench_env_octal=015 ;;
      b|f) _bench_env_octal=040 ;;
      '"') _bench_env_octal=042 ;;
      '\') _bench_env_octal=134 ;;
      /) _bench_env_octal=057 ;;
      u) bench_env_unicode_octal || return 1 ;;
      *) return 1 ;;
    esac
    # The sentinel survives the command substitution's trailing-newline strip, so a
    # decoded newline reaches the value rather than vanishing from it.
    _bench_env_value=$(printf '%s%bX' "$_bench_env_value" "\\0$_bench_env_octal")
    _bench_env_value=${_bench_env_value%X}
  done
  return 1
}

# bench_env_skip_space advances the reader past JSON whitespace.
bench_env_skip_space() {
  while :; do
    case $_bench_env_rest in
      [[:space:]]*) _bench_env_rest=${_bench_env_rest#?} ;;
      *) return 0 ;;
    esac
  done
}

# bench_env_unicode_octal consumes a four-hex-digit Unicode escape and sets
# _bench_env_octal to the byte it decodes to.
#
# An escape in the ASCII range decodes to its own byte. Go's encoding/json escapes &,
# <, and > by default, so a control operator or a command name reaches this decoder
# escaped as often as literal. A placeholder there would hide the operator the word
# test looks for.
#
# Above ASCII the placeholder `_` stands, because no such rune is a shell operator. A
# NUL escape refuses, because the shell cannot carry NUL, so decoding it would
# silently truncate the command the caller is deciding about.
bench_env_unicode_octal() {
  case $_bench_env_rest in
    [0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f]*) ;;
    *) return 1 ;;
  esac
  _bench_env_after=${_bench_env_rest#????}
  _bench_env_hex=${_bench_env_rest%"$_bench_env_after"}
  _bench_env_rest=$_bench_env_after
  case $_bench_env_hex in
    00[0-7][0-9A-Fa-f])
      _bench_env_code=$(( 0x$_bench_env_hex ))
      [ "$_bench_env_code" -eq 0 ] && return 1
      _bench_env_octal=$(printf '%03o' "$_bench_env_code")
      ;;
    *) _bench_env_octal=137 ;;
  esac
  return 0
}

# bench_shell_quote prints $1 as one single-quoted shell word. bench_rebuild_action
# composes the one rebuild command for the kit root $1. Both are the shell spelling of
# internal/freshness's shellQuote and RebuildAction, whose Go half is the source of the
# invocation. This copy exists only because the invocation is needed at the one moment
# the Go binary cannot run, so nothing can ask it for the string. internal/systemtest
# asserts the printed line against freshness.RebuildAction, so a drift between the two
# reds.
bench_shell_quote() {
  _bench_quote_rest=$1
  _bench_quote_out=''
  while :; do
    case $_bench_quote_rest in
      *\'*)
        _bench_quote_out=$_bench_quote_out${_bench_quote_rest%%\'*}"'\\''"
        _bench_quote_rest=${_bench_quote_rest#*\'}
        ;;
      *)
        _bench_quote_out=$_bench_quote_out$_bench_quote_rest
        break
        ;;
    esac
  done
  printf "'%s'" "$_bench_quote_out"
}

bench_rebuild_action() {
  printf 'cd %s && bash scripts/go-build.sh %s %s' "$(bench_shell_quote "$1")" "$(bench_shell_quote "$1")" "$(bench_shell_quote "$1/dist/bench")"
}
