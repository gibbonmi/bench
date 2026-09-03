# shellcheck shell=sh
# Canary fixture: the shell word test has lost its xargs routine-prefix arm, so
# "xargs -- bench help" reads as a non-Bench call. The shared table beside it still
# declares that row true. The guard-classifier-table conformance check must go red
# naming that one row.
# This is the shared shell classification library for every hook shim
# (.bench/hooks/*.sh) and shift adapter (.bench/adapters/*). It carries the one
# source of the wrapper search order — repo `.bench/bin/bench.sh`, then kit
# `bin/bench.sh`, then `bench` on PATH — and the shell Bench-call word test. It
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
# xargs. It resolves no path and reads no wrapper string, so a `bench` reached
# through a symlink or through `bash -c` is outside it. The two derivations are
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
      xargs-dropped)
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

# bench_is_assignment reports whether a word is a portable NAME=VALUE assignment.
bench_is_assignment() {
  case $1 in *=*) ;; *) return 1 ;; esac
  _bench_name=${1%%=*}
  case $_bench_name in
    ''|[!A-Za-z_]*|*[!A-Za-z0-9_]*) return 1 ;;
  esac
  return 0
}
