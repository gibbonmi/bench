package gate

// The background grade: one grade of one lane-passed tip of a kit worktree. It grades the chunk
// subject of the open build of the tip, or the landing subject of a tip with no open build, writes
// its green to the evidence store, and records its outcome in the background record of the
// worktree. It holds a claim on its grade key in place
// of the checkout lock, and it writes no gate cache.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/runbinary"
	specdoc "github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// backgroundVerb is the plumbing entry that runs one background grade.
const backgroundVerb = "gate-background"

// backgroundUsage is the grammar of the background grade entry.
const backgroundUsage = "usage: bench " + backgroundVerb + " <worktree> <tip>"

// backgroundNice is the nice value of a background grade and of each process that it starts.
const backgroundNice = 19

// backgroundResolvedEnv marks the process that the wrapper of the target worktree resolved. The
// entry consumes it, so no child of the grade inherits it.
const backgroundResolvedEnv = "BENCH_GATE_BACKGROUND"

// backgroundQueueEnv names the Bench home of the machine queue that started the grade. The entry
// consumes it, and only a grade that a queue started takes the next grade of that queue.
const backgroundQueueEnv = "BENCH_GATE_BACKGROUND_QUEUE"

// backgroundPassEnv marks a lane pass and names the Bench home of its machine queue. The entry
// consumes it, so the grade that the pass starts does not inherit it.
const backgroundPassEnv = "BENCH_GATE_BACKGROUND_PASS"

// backgroundGrade is the state of one background grade: the worktree of the grade, the grade key
// of the subject that the run accepts, the claim that it holds on that key, the git common
// directory of that claim, and whether its oracle run reached a verdict. A nil grade is a
// foreground run, so each method below keeps the foreground behavior for it.
type backgroundGrade struct {
	worktree, key, common string
	claim                 backgroundstate.Claim
	verdict               bool
}

// reenterImage replaces the image of this process with the program name, run with args and env
// in the directory dir. It is a seam, so a test can run the replacement as a child.
var reenterImage = func(dir, name string, args, env []string) error {
	if err := os.Chdir(dir); err != nil {
		return err
	}
	return subprocess.Replace(name, args, env)
}

// BackgroundCommand runs `bench gate-background <worktree> <tip>`: one background grade of tip
// in worktree, in the foreground of this process. The chunk checkpoint of worktree runs as the
// executable that the wrapper of worktree resolves, and its grade key frames the digest of that
// executable. So the entry first replaces its image with that wrapper, with no inherited wrapper
// routing and no inherited run binary, and the executable that the wrapper resolves runs the
// grade. A grade that the machine queue started then starts the oldest queued grade. A lane pass
// runs that executable too, so its grade key is the key of the grade: it applies its tip to the
// queue, starts the oldest queued grade when the slot is free, and then exits.
func BackgroundCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) != 2 || args[0] == "" || args[1] == "" {
		fmt.Fprintln(stderr, backgroundUsage)
		return 2
	}
	root, tip := args[0], args[1]
	if os.Getenv(backgroundResolvedEnv) == "" {
		resolved := append(env.WithoutWrapperRouting(os.Environ(), runbinary.Env, backgroundResolvedEnv), backgroundResolvedEnv+"=1")
		err := reenterImage(root, backgroundWrapper(root), append([]string{backgroundVerb}, args...), resolved)
		fmt.Fprintf(stderr, "bench %s: %v\n", backgroundVerb, err)
		return 1
	}
	// A landing subject runs as the broker that authenticates, so the entry execs that broker
	// with its own arguments and its markers, and the broker runs the pass or the grade.
	var moved notBrokerError
	if _, err := backgroundSubjectOf(root, tip); errors.As(err, &moved) {
		if err := execBroker(moved.broker, append([]string{backgroundVerb}, args...), os.Environ()); err != nil {
			fmt.Fprintf(stderr, "bench %s: exec the broker %q: %v\n", backgroundVerb, moved.broker, err)
		} else {
			return 0
		}
	}
	home, queued := os.LookupEnv(backgroundQueueEnv)
	passHome, pass := os.LookupEnv(backgroundPassEnv)
	for _, name := range []string{backgroundResolvedEnv, backgroundQueueEnv, backgroundPassEnv} {
		if err := os.Unsetenv(name); err != nil {
			fmt.Fprintf(stderr, "bench %s: %v\n", backgroundVerb, err)
			return 1
		}
	}
	if pass {
		if err := passLane(passHome, root, tip); err != nil {
			fmt.Fprintf(stderr, "background grade: not queued: %v\n", err)
			return 1
		}
		return 0
	}
	// The grade runs at the lowest priority that the platform allows, so a refused step keeps
	// the grade at the priority that it has.
	if err := lowerPriority(); err != nil {
		fmt.Fprintf(stderr, "bench %s: priority not lowered: %v\n", backgroundVerb, err)
	}
	// A stop can come at any step, so the signals end the grade from its start, and the grade
	// releases each resource that it holds before it exits.
	ctx, stop := notifyGateSignals(context.Background())
	code := gradeInBackground(ctx, root, tip, nil, stdout, stderr)
	stop()
	// A grade that a queue started starts the next grade of that queue, and then it exits.
	if queued {
		if err := advanceQueue(home, os.Getpid()); err != nil {
			fmt.Fprintf(stderr, "background grade: next grade not started: %v\n", err)
		}
	}
	return code
}

// gradeInBackground grades tip in the worktree root as one background grade, and arm runs once
// the grade holds its claim. The run writes its output to the running log of the worktree, so the
// log of the latest record stays whole while the grade runs. The grade then records its end
// through WriteEnd, which names a stopped grade stopped and moves the running log over the log of
// the record. The grade prints its outcome and its log, and it answers the exit of the run.
func gradeInBackground(ctx context.Context, root, tip string, arm postAcquireContextArm, stdout, stderr io.Writer) int {
	gitdir, err := benchgit.AdminDir(root)
	if err != nil {
		fmt.Fprintf(stderr, "background grade: git directory unavailable: %v\n", err)
		return 1
	}
	log, err := os.OpenFile(backgroundstate.RunningLogPath(gitdir), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		fmt.Fprintf(stderr, "background grade: log unavailable: %v\n", err)
		return 1
	}
	defer log.Close()
	grade := &backgroundGrade{worktree: root}
	commit, result := grade.run(ctx, root, tip, arm, log)
	if err := grade.release(); err != nil {
		fmt.Fprintf(stderr, "background grade: claim release failed: %v\n", err)
		return 1
	}
	// A tip that does not resolve records the tip as the caller named it.
	if commit == "" {
		commit = tip
	}
	// The owner of the grade is this process. A boot that does not read names no owner, so no
	// stop marker names the grade.
	owner, _ := backgroundstate.ClaimFor(grade.key, root)
	record, err := backgroundstate.WriteEnd(durablefile.Native(), gitdir, owner, backgroundstate.Record{Outcome: grade.outcome(result), Key: grade.key, Tip: commit, Log: backgroundstate.LogPath(gitdir)})
	if err != nil {
		fmt.Fprintf(stderr, "background grade: record persistence failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "background grade: %s (log %s)\n", record.Outcome, record.Log)
	return result.ActionExit
}

// outcome answers the outcome of the run of the grade: green for a run that ends green, red for
// a run whose oracle reached a verdict and refuses the tree, and aborted for each other end. A
// run that ends before it knows its grade key, a timeout, an interrupt, a settle refusal, and a
// subject that the grade cannot hold each end with no verdict.
func (g *backgroundGrade) outcome(result Result) backgroundstate.Outcome {
	switch {
	case g.key == "":
		return backgroundstate.Aborted
	case result.ActionExit == 0:
		return backgroundstate.Green
	case g.verdict && result.GateExit != 0:
		return backgroundstate.Red
	}
	return backgroundstate.Aborted
}

// settleVerdict runs settleRun over a finished oracle run. When the settle step accepts the run,
// the background grade that ctx carries marks that its oracle reached a verdict. A settle refusal
// gives no verdict, so it marks nothing. A foreground run marks nothing.
func settleVerdict(ctx context.Context, evaluation executionEvaluation, runtimeRoot, storageRoot string, declaredPaths []string, plan subject, run processGroupResult, release func() error, stderr io.Writer) (subject, Result, bool) {
	after, refusal, refused := settleRun(evaluation, runtimeRoot, storageRoot, declaredPaths, plan, run, release, stderr)
	if g := backgroundGradeOf(ctx); g != nil && !refused {
		g.verdict = true
	}
	return after, refusal, refused
}

// run grades tip in root as the subject of tip: the chunk subject of its open build, or with no
// open build the landing subject, which runs in the primary checkout as a landing runs. The chunk
// id comes from the latest chunk entry of the review record at tip, and with no entry the run
// takes the whole schedule. It answers the commit of tip and the result of the run.
func (g *backgroundGrade) run(ctx context.Context, root, tip string, arm postAcquireContextArm, out io.Writer) (string, Result) {
	subject, err := backgroundSubjectOf(root, tip)
	if err != nil {
		return subject.source.tip, operational(root, 0, out, "background grade: "+err.Error())
	}
	if subject.landing != nil {
		ctx = WithCheckpoint(ctx, Checkpoint{Background: g})
		return subject.source.tip, executeTreeWithOwner(ctx, subject.landing.primary, subject.tree, out, out, arm, reuseFreshGreen, nil)
	}
	ctx = WithCheckpoint(ctx, Checkpoint{Spec: subject.spec, Chunk: latestChunk(root, subject.source.tree, subject.spec), Background: g})
	return subject.source.tip, executeStrippedSubject(ctx, root, subject.tree, subject.source, out, out, arm, reuseFreshGreen)
}

// gradeSubject is the subject of a background grade of a tip: the source of the tip, the spec of
// its open build, the tree that the grade runs, and the landing of a tip with no open build.
type gradeSubject struct {
	source     chunkSource
	spec, tree string
	landing    *landingSubject
}

// backgroundSubjectOf answers the subject of a background grade of tip in root. A tip with an
// open build of a staged spec grades its tree without the review record of that spec. A tip with
// no open build in a linked worktree grades the tree that its landing publishes. The lane pass
// keys that tree and the grade runs it, so one derivation serves both. A source that it cannot
// read keeps each part that it resolved.
func backgroundSubjectOf(root, tip string) (gradeSubject, error) {
	commit, err := benchgit.ResolveCommit(root, tip)
	if err != nil {
		return gradeSubject{}, fmt.Errorf("tip unavailable: %w", err)
	}
	tree, err := benchgit.Output("-C", root, "rev-parse", "--verify", commit+"^{tree}")
	if err != nil {
		return gradeSubject{source: chunkSource{tip: commit}}, fmt.Errorf("tip unavailable: %w", err)
	}
	subject := gradeSubject{source: chunkSource{tip: commit, tree: tree}}
	spec, open := openBuild(root, tree)
	if !open {
		subject.landing, subject.tree, err = landingSubjectOf(root, subject.source)
		return subject, err
	}
	subject.spec = spec
	if subject.tree, err = reviewrecord.SourceDigest(root, tree, spec); err != nil {
		return subject, fmt.Errorf("subject unavailable: %w", err)
	}
	return subject, nil
}

// openBuild answers the spec of the open build of tree: the one spec with a Status: staged line
// whose review record tree holds as a regular file. The review record owner names both paths, so
// a spec is each file of tree that has a record path. A tree with no such spec, or with more than
// one, holds no open build. A wrong answer costs only a lookup miss, because a grade records its
// green under the key that it ran.
func openBuild(root, tree string) (string, bool) {
	entries, err := benchgit.ListTree(root, "-r", tree)
	if err != nil {
		return "", false
	}
	regular := map[string]bool{}
	for _, entry := range entries {
		regular[entry.Path] = entry.IsRegularFile()
	}
	var open []string
	for _, entry := range entries {
		if record, err := reviewrecord.RecordPath(entry.Path); err != nil || !regular[record] {
			continue
		}
		if data, err := benchgit.ReadTreeFile(root, tree, entry.Path); err == nil {
			if _, err := specdoc.Implemented(data); err == nil {
				open = append(open, entry.Path)
			}
		}
	}
	if len(open) != 1 {
		return "", false
	}
	return open[0], true
}

// latestChunk answers the id of the latest chunk entry of the review record of spec in tree, or
// the empty id when no entry resolves.
func latestChunk(root, tree, spec string) string {
	record, err := reviewrecord.ReadTree(root, tree, spec)
	if err != nil || len(record.Chunks) == 0 {
		return ""
	}
	return record.Chunks[len(record.Chunks)-1].ID
}

// executeChunkSubject grades the chunk subject of spec at source for the chunk checkpoint: the tree
// of source without the review record of spec, which executeStrippedSubject runs.
func executeChunkSubject(ctx context.Context, root, spec string, source chunkSource, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode) Result {
	stripped, err := reviewrecord.SourceDigest(root, source.tree, spec)
	if err != nil {
		return sourceUnavailable(ctx, root, stderr, mode, err)
	}
	return executeStrippedSubject(ctx, root, stripped, source, stdout, stderr, arm, mode)
}

// executeStrippedSubject grades stripped, the tree of source without its review record. That tree
// runs through the prospective owner in a private checkout, with root as the storage and identity
// root, so a record-only commit keeps the subject of an earlier chunk green. That checkout is also
// the baseline root. The chunk checkpoint and the background grade both run their subject here.
func executeStrippedSubject(ctx context.Context, root, stripped string, source chunkSource, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode) Result {
	ctx = withCheckoutBaseline(context.WithValue(ctx, deltaSourceKey{}, source))
	if checkpoint, _ := ctx.Value(checkpointKey{}).(Checkpoint); !checkpoint.gradesSubjectAlone() {
		ctx = context.WithValue(ctx, chunkSourceKey{}, source)
	}
	return executeTreeWithOwner(ctx, root, stripped, stdout, stderr, arm, mode, nil)
}

// backgroundGradeOf answers the background grade that ctx carries, or nil for a foreground run.
func backgroundGradeOf(ctx context.Context) *backgroundGrade {
	checkpoint, _ := ctx.Value(checkpointKey{}).(Checkpoint)
	return checkpoint.Background
}

// hold writes the claim of the grade on the grade key of plan through replacer, in the git common
// directory of root, for the worktree of the grade, or for root when the grade names none. The
// grade keeps that key before the claim, so a claim that fails still records its outcome under the
// key. A foreground run holds none.
func (g *backgroundGrade) hold(replacer durablefile.Replacer, root string, plan subject) error {
	if g == nil {
		return nil
	}
	g.key = evidenceName(plan)
	common, err := benchgit.CommonDir(root)
	if err != nil {
		return err
	}
	worktree := g.worktree
	if worktree == "" {
		worktree = root
	}
	claim, err := backgroundstate.ClaimFor(g.key, worktree)
	if err != nil {
		return err
	}
	if err := backgroundstate.WriteClaim(replacer, common, claim); err != nil {
		return err
	}
	g.common, g.claim = common, claim
	return nil
}

// release removes the claim that the grade holds, if it holds one.
func (g *backgroundGrade) release() error {
	if g.claim.Key == "" {
		return nil
	}
	return backgroundstate.ReleaseClaim(g.common, g.claim)
}

// oracleEnv answers the oracle environment of a run over base. A background grade adds the width
// input.
func (g *backgroundGrade) oracleEnv(base []string) []string {
	if g == nil {
		return base
	}
	return mergeEnv(base, []string{widthEnv + "=1"})
}

// persist installs rec as the gate cache of gitdir, as persistRecovering does. A background grade
// writes no gate cache.
func (g *backgroundGrade) persist(replacer durablefile.Replacer, gitdir string, rec, pending verdictRecord) error {
	if g != nil {
		return nil
	}
	return persistRecovering(replacer, gitdir, rec, pending)
}
