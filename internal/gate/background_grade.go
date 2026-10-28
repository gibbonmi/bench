package gate

// The background grade: one grade of one lane-passed tip of a kit worktree. It grades the chunk
// subject of the open build of the tip, writes its green to the evidence store, and records its
// outcome in the background record of the worktree. It holds a claim on its grade key in place
// of the checkout lock, and it writes no gate cache.

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

// backgroundGrade is the state of one background grade: the grade key of the subject that the
// run accepts, the claim that it holds on that key, the git common directory of that claim, and
// whether its oracle reached its own exit. A nil grade is a foreground run, so each method below
// keeps the foreground behavior for it.
type backgroundGrade struct {
	key, common string
	claim       backgroundstate.Claim
	verdict     bool
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
// grade.
func BackgroundCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) != 2 || args[0] == "" || args[1] == "" {
		fmt.Fprintln(stderr, backgroundUsage)
		return 2
	}
	root, tip := args[0], args[1]
	if os.Getenv(backgroundResolvedEnv) == "" {
		resolved := append(env.WithoutWrapperRouting(os.Environ(), runbinary.Env, backgroundResolvedEnv), backgroundResolvedEnv+"=1")
		err := reenterImage(root, filepath.Join(root, "bin", "bench.sh"), append([]string{backgroundVerb}, args...), resolved)
		fmt.Fprintf(stderr, "bench %s: %v\n", backgroundVerb, err)
		return 1
	}
	if err := os.Unsetenv(backgroundResolvedEnv); err != nil {
		fmt.Fprintf(stderr, "bench %s: %v\n", backgroundVerb, err)
		return 1
	}
	// The grade runs at the lowest priority that the platform allows, so a refused step keeps
	// the grade at the priority that it has.
	if err := lowerPriority(); err != nil {
		fmt.Fprintf(stderr, "bench %s: priority not lowered: %v\n", backgroundVerb, err)
	}
	return gradeInBackground(context.Background(), root, tip, notifyGateSignals, stdout, stderr)
}

// gradeInBackground grades tip in the worktree root as one background grade, and arm runs once
// the grade holds its claim. The run writes its output to the background log of the worktree. Once
// the run knows its grade key, it records its outcome in the background record. The grade then
// prints its outcome and its log, and it answers the exit of the run.
func gradeInBackground(ctx context.Context, root, tip string, arm postAcquireContextArm, stdout, stderr io.Writer) int {
	gitdir, err := benchgit.AdminDir(root)
	if err != nil {
		fmt.Fprintf(stderr, "background grade: git directory unavailable: %v\n", err)
		return 1
	}
	logPath := backgroundstate.LogPath(gitdir)
	log, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		fmt.Fprintf(stderr, "background grade: log unavailable: %v\n", err)
		return 1
	}
	defer log.Close()
	grade := &backgroundGrade{}
	commit, result := grade.run(ctx, root, tip, arm, log)
	if err := grade.release(); err != nil {
		fmt.Fprintf(stderr, "background grade: claim release failed: %v\n", err)
		return 1
	}
	if grade.key == "" {
		fmt.Fprintf(stdout, "background grade: no verdict (log %s)\n", logPath)
		return result.ActionExit
	}
	outcome := grade.outcome(result)
	record := backgroundstate.Record{Outcome: outcome, Key: grade.key, Tip: commit, Log: logPath}
	if err := backgroundstate.WriteRecord(durablefile.Native(), gitdir, record); err != nil {
		fmt.Fprintf(stderr, "background grade: record persistence failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "background grade: %s (log %s)\n", outcome, logPath)
	return result.ActionExit
}

// outcome answers the outcome of the run of the grade: green for a run that ends green, red for
// a run whose oracle reached its own exit and refuses the tree, and aborted for each other end.
// A timeout, an interrupt, and a subject that the grade cannot hold each end with no verdict.
func (g *backgroundGrade) outcome(result Result) backgroundstate.Outcome {
	switch {
	case result.ActionExit == 0:
		return backgroundstate.Green
	case g.verdict && result.GateExit != 0:
		return backgroundstate.Red
	}
	return backgroundstate.Aborted
}

// settle marks that the oracle of the grade reached its own exit. A foreground run marks nothing.
func (g *backgroundGrade) settle() {
	if g != nil {
		g.verdict = true
	}
}

// run grades tip in root as the chunk subject of the open build of tip. The chunk id comes from
// the latest chunk entry of the review record at tip, and with no entry the run takes the whole
// schedule. It answers the commit of tip and the result of the run.
func (g *backgroundGrade) run(ctx context.Context, root, tip string, arm postAcquireContextArm, out io.Writer) (string, Result) {
	commit, err := benchgit.ResolveCommit(root, tip)
	if err != nil {
		return "", operational(root, 0, out, "background grade: tip unavailable: "+err.Error())
	}
	tree, err := benchgit.Output("-C", root, "rev-parse", "--verify", commit+"^{tree}")
	if err != nil {
		return commit, operational(root, 0, out, "background grade: tip unavailable: "+err.Error())
	}
	spec, open := openBuild(root, tree)
	if !open {
		return commit, operational(root, 0, out, "background grade: the tip holds no open build of a staged spec")
	}
	ctx = WithCheckpoint(ctx, Checkpoint{Spec: spec, Chunk: latestChunk(root, tree, spec), Background: g})
	return commit, executeChunkSubject(ctx, root, spec, chunkSource{tip: commit, tree: tree}, out, out, arm, reuseFreshGreen)
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

// executeChunkSubject grades the chunk subject of spec at source: the tree of source without the
// review record of spec. That tree runs through the prospective owner in a private checkout, with
// root as the storage and identity root, so a record-only commit keeps the subject of an earlier
// chunk green. That checkout is also the baseline root. The chunk checkpoint and the background
// grade both run their subject here.
func executeChunkSubject(ctx context.Context, root, spec string, source chunkSource, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode) Result {
	stripped, err := reviewrecord.SourceDigest(root, source.tree, spec)
	if err != nil {
		return sourceUnavailable(ctx, root, stderr, mode, err)
	}
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
// directory of root. The grade keeps that key before the claim, so a claim that fails still
// records its outcome under the key. A foreground run holds none.
func (g *backgroundGrade) hold(replacer durablefile.Replacer, root string, plan subject) error {
	if g == nil {
		return nil
	}
	g.key = evidenceName(plan)
	common, err := benchgit.CommonDir(root)
	if err != nil {
		return err
	}
	claim, err := backgroundstate.ClaimFor(g.key, root)
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
