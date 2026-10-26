package gate

import (
	"context"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing/published"
)

// executeCompleteCheckpoint grades the tree that the landing publishes from HEAD of root.
// The landing publishes only committed bytes, so committedSource refuses an uncommitted
// change before the oracle runs. The landing's closure transform applies to the HEAD tree,
// and the oracle grades the result under the completion obligation for spec at HEAD. Each
// refusal here names the route of its cause.
func executeCompleteCheckpoint(ctx context.Context, root, spec string, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode) Result {
	tip, sourceTree, refusal, refused := committedSource(ctx, root, stderr, mode)
	if refused {
		return refusal
	}
	tree, err := published.Tree(root, sourceTree, spec, tip)
	if err != nil {
		return refuse(ctx, root, stderr, mode, faceComposition, "complete checkpoint unavailable: "+err.Error())
	}
	return executeTreeWithOwner(WithCompletion(ctx, spec, tip), root, tree, stdout, stderr, arm, mode, nil)
}

// faceBrokerUnavailable is the face of a complete checkpoint whose landing broker does not
// authenticate.
var faceBrokerUnavailable = checkpointFace{name: "checkpoint-broker-unavailable"}

// execBroker replaces this process with broker under argv, in the inherited working
// directory and environment. It returns only when the replacement fails. It is a seam, so a
// test observes the call with no process replacement.
var execBroker = func(broker string, argv []string) error {
	return syscall.Exec(broker, argv, os.Environ())
}

// runUnderBroker moves a complete checkpoint under the landing broker, so its green answers
// the landing. In a repository that declares Bench build inputs, it reads the primary
// checkout's broker manifest. A running executable whose digest the manifest binds keeps the
// run in this process. Any other executable execs the broker with the same arguments in the
// same working directory. A manifest that does not read, a broker whose digest differs from
// its manifest, and a failed exec each refuse, so no gate runs under an unknown executable.
// A chunk checkpoint and a repository without build inputs keep the running executable. It
// answers whether the run ended here, and its result.
func runUnderBroker(ctx context.Context, root string, checkpoint Checkpoint, stderr io.Writer, mode runMode) (Result, bool) {
	if !checkpoint.Complete {
		return Result{}, false
	}
	primary, located := git.PrimaryCheckout(root)
	if located == nil && !benchfreshness.DeclaresBuildInputs(primary) {
		return Result{}, false
	}
	unavailable := func(cause string) (Result, bool) {
		return refuse(ctx, root, stderr, mode, faceBrokerUnavailable, "complete checkpoint broker unavailable: "+cause), true
	}
	if located != nil {
		return unavailable("the primary checkout does not resolve: " + located.Error())
	}
	manifest := brokermanifest.Path(primary)
	fields, err := brokermanifest.Read(manifest)
	if err != nil {
		return unavailable("the primary broker manifest does not read: " + err.Error())
	}
	bound := fields[brokermanifest.DigestField]
	if runningExecutableDigest() == bound {
		return Result{}, false
	}
	broker := fields["path"]
	if digest, err := brokermanifest.Digest(broker); err != nil || digest != bound {
		return unavailable(fmt.Sprintf("the broker %q does not match the digest %s that the primary broker manifest %s binds", broker, bound, manifest))
	}
	if err := execBroker(broker, append([]string{broker}, os.Args[1:]...)); err != nil {
		return unavailable(fmt.Sprintf("exec the broker %q: %v", broker, err))
	}
	return Result{}, true
}
