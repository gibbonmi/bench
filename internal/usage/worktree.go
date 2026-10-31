package usage

import (
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
)

const (
	WorktreeList        = "bench worktree list"
	WorktreeListPaths   = WorktreeList + " --view paths --target <target> [--target <target>]..."
	WorktreeListRelated = WorktreeList + " --spec <slug-or-path> [--offset <byte-offset> --source <identity>]"
	WorktreePath        = "bench worktree path <target>"
	WorktreeExec        = "bench worktree exec <target> [--env KEY=VALUE]... -- <command> [args...]"
	WorktreeShow        = "bench worktree show <target> <rev>:<path>"
	WorktreeRead        = "bench worktree read [--in <label|primary>] <path> [--section <literal-H2-title> | --lines <first>:<last>] [--offset <byte-offset> --source <identity>]"
	WorktreeBuild       = "bench worktree build <target>"
	WorktreeCreate      = "bench worktree create [--refresh] --request <opaque-id> --label <work-item> [--from <target>]"
	WorktreeRelease     = "bench worktree release --request <opaque-id> <path>"
	WorktreeClean       = "bench worktree clean [--discard-ignored] [--discard-branch] [--full] (<path> | --landed | --target <target>...) [--apply <fingerprint>] | bench worktree clean --discard-branch --unclaimed [--apply <fingerprint> | --apply-current]"
	WorktreeReclaim     = "bench worktree reclaim [--apply <fingerprint>]"
	WorktreeReauthorize = "bench worktree reauthorize --assignment <assignment-id> --request <opaque-id> --base <commit> --source-tip <commit> <path>"
	WorktreeMerge       = "bench worktree merge --from <commit|target> <target>"
	WorktreeReset       = "bench worktree reset (--to <commit> | --restore <ref>) <target> [--apply <fingerprint>]"
	WorktreeLand        = "bench worktree land --request <opaque-id> --base <commit> --source-tip <commit> [--spec <slug>] -m <message> <path>"
	WorktreeLandResume  = "bench worktree land --resume <published-commit> --request <opaque-id> --base <commit> --source-tip <commit> [--spec <slug>] <path>"
	WorktreeLandStack   = "bench worktree land --stack <file> [--source-tip <request>=<commit>]..."
	// WorktreeLandStackResume is the stack form of WorktreeLandResume.
	WorktreeLandStackResume = "bench worktree land --resume <published-commit> --stack <file> [--source-tip <request>=<commit>]..."
)

// WorktreeExecGate is the one gate form a worktree reader gets. It is an exec composition
// rather than a worktree route, so it trails the grammar list instead of joining it, and
// both this usage trailer and the `bench help` row read it from here.
const WorktreeExecGate = "bench worktree exec <target> -- bench gate"

var worktreeCommands = []string{
	WorktreeList,
	WorktreeListPaths,
	WorktreeListRelated,
	WorktreePath,
	WorktreeExec,
	WorktreeShow,
	WorktreeRead,
	WorktreeBuild,
	WorktreeCreate,
	WorktreeRelease,
	WorktreeClean,
	WorktreeReclaim,
	WorktreeReauthorize,
	WorktreeMerge,
	WorktreeReset,
	WorktreeLand,
	WorktreeLandResume,
	WorktreeLandStack,
	WorktreeLandStackResume,
}

func WorktreeUsage() string {
	return "usage: bench worktree [--refresh] [objective...]\n       " + strings.Join(worktreeCommands, "\n       ") + "\n       " + WorktreeExecGate + "\n"
}

// WorktreeBuildInstall is the install form of the build verb. The build verb's help prints
// it after WorktreeBuild, and WorktreeBuild stays the one form that each rebuild remedy fills.
const WorktreeBuildInstall = "bench worktree build --install <target>"

// WorktreeBuildFor is WorktreeBuild with target in its `<target>` slot, so each remedy that
// names the build of one worktree fills the grammar one way. The caller quotes target when
// its sink needs a shell word.
func WorktreeBuildFor(target string) string { return withTarget(WorktreeBuild, target) }

// WorktreeBuildInstallFor is WorktreeBuildInstall with target in its `<target>` slot, under
// the quoting rule of WorktreeBuildFor.
func WorktreeBuildInstallFor(target string) string {
	return withTarget(WorktreeBuildInstall, target)
}

// WorktreeBuildRestore is the restore form of the build verb, whose operand is the spec slug
// of a stopped build. The build verb's help prints it after WorktreeBuildInstall.
const WorktreeBuildRestore = "bench worktree build --restore <slug>"

// WorktreeBuildRestoreFor is WorktreeBuildRestore with slug in its `<slug>` slot, under the
// quoting rule of WorktreeBuildFor.
func WorktreeBuildRestoreFor(slug string) string {
	return strings.Replace(WorktreeBuildRestore, "<slug>", slug, 1)
}

// withTarget fills the `<target>` slot of one build grammar.
func withTarget(grammar, target string) string {
	return strings.Replace(grammar, "<target>", target, 1)
}

// PrimaryCheckoutRefusal is the one refusal a Bench write verb prints from the primary
// checkout. Within Bench, main receives writes only through landings, so every verb
// that writes the tree redirects to a worktree with this line.
func PrimaryCheckoutRefusal() string {
	return toon.Errorf("primary checkout is read-only for Bench phases", "run "+WorktreeCreate)
}
