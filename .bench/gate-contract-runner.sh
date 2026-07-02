# Shared fixture harness for the gate's contract fragments — the one source for
# the provision/report/cleanup block a fixture contract needs (one source per
# fact: N pasted mktemp/subshell/cleanup harnesses drift; this one cannot).
# Sourced by gate.sh after err() is defined and before any fragment. Fragments
# call:
#
#   contract "<label>" [--space-path] [--no-repo] <<'BODY'
#     ...assertions...
#   BODY
#
# The BODY runs under `set -u` in a subshell whose cwd is a fresh throwaway
# fixture directory, git-inited unless --no-repo. $root (the kit repo) and $tmp
# (the fixture dir) are in scope. A body fails by echoing a diagnostic and
# exiting nonzero; any nonzero body reports "<label> failed" through err(). The
# fixture is removed on success and failure alike — no leaked scratch dirs.
#
#   --space-path   provision the fixture under a parent directory containing a
#                  space, for path-quoting contracts
#   --no-repo      skip git init, for contracts asserting outside-a-repo
#                  behavior (the body may init its own repos)

contract() {
  local _label="$1"; shift
  local _space=0 _repo=1
  while [ $# -gt 0 ]; do
    case "$1" in
      --space-path) _space=1 ;;
      --no-repo)    _repo=0 ;;
      *) err "contract '$_label': unknown option $1"; return 1 ;;
    esac
    shift
  done
  local _body _parent tmp
  _body="$(cat)"
  _parent="$(mktemp -d)"
  if [ "$_space" = 1 ]; then
    tmp="$_parent/space dir"
    mkdir -p "$tmp"
  else
    tmp="$_parent"
  fi
  (
    set -u
    cd "$tmp" || exit 1
    if [ "$_repo" = 1 ]; then git init -q || exit 1; fi
    eval "$_body"
  ) || err "$_label failed"
  rm -rf "$_parent"
}
