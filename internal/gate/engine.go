package gate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	"github.com/gibbonmi/bench/internal/gate/prospectiveartifact"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// EvidenceInspection reports whether retained gate evidence can authorize one tree.
type EvidenceInspection struct {
	Tree          string
	Oracle        string
	RecordedAt    time.Time
	ReusableGreen bool
	Reason        string
}

// InspectTree reports retained exact evidence for an unpublished Git tree.
func InspectTree(root, tree string) EvidenceInspection {
	return inspectProspective(root, tree, time.Now())
}

// ExecuteTree runs or reuses the gate for an unpublished Git tree.
func ExecuteTree(ctx context.Context, root, tree string, stdout, stderr io.Writer) Result {
	ctx, finishSpan := beginGateSpan(ctx, root, "prospective")
	ctx, finishLog := beginGateRunLog(ctx, root, stderr, "prospective")
	result := executeTreeWithOwner(ctx, root, tree, stdout, stderr, nil)
	finishLog(result)
	finishSpan(result)
	return result
}

func executeTreeWithOwner(ctx context.Context, root, tree string, stdout, stderr io.Writer, owner runBinaryOwner) Result {
	artifacts, err := openProspectiveArtifacts(root, tree)
	if err != nil {
		fmt.Fprintln(stderr, "prospective gate subject unavailable")
		return Result{ActionExit: 1}
	}
	defer artifacts.Close()
	checkout := artifacts.Checkout()
	if owner == nil {
		owner = prospectiveRunBinaryOwnerAt(checkout, artifacts.Root())
	}
	evaluation := newProspectiveTreeEvaluation(checkout, root, tree)
	return executeSubjectWithRunBinary(ctx, checkout, root, stdout, stderr, nil, reuseFreshGreen, evaluation, owner, root)
}

// ValidateProjectGreen reports whether branch's tip and marker have retained exact green evidence.
func ValidateProjectGreen(root, branch string) EvidenceInspection {
	current, err := benchgit.Output("-C", root, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil || current != branch {
		return EvidenceInspection{Reason: "working branch changed"}
	}
	tip, err := benchgit.Output("-C", root, "rev-parse", "HEAD")
	if err != nil {
		return EvidenceInspection{Reason: "working tip unavailable"}
	}
	marker, present, err := greenmarker.Read(root, branch)
	if err != nil || !present || marker != tip {
		return EvidenceInspection{Reason: "project-green marker changed"}
	}
	plan, err := buildSubject(root)
	if err != nil {
		return EvidenceInspection{Reason: "subject unavailable"}
	}
	return inspectEvidence(root, plan, time.Now())
}

func inspectProspective(root, tree string, now time.Time) EvidenceInspection {
	artifacts, err := openProspectiveArtifacts(root, tree)
	if err != nil {
		return EvidenceInspection{Reason: "subject unavailable"}
	}
	defer artifacts.Close()
	plan, err := buildProspectiveSubjectFor(artifacts.Checkout(), root)
	if err != nil || plan.Tree != tree {
		return EvidenceInspection{Reason: "subject unavailable"}
	}
	return inspectEvidence(root, plan, now)
}

func inspectEvidence(root string, plan subject, now time.Time) EvidenceInspection {
	return inspectEvidenceWindowed(root, plan, now, true)
}

func inspectEvidenceWindowed(root string, plan subject, now time.Time, expires bool) EvidenceInspection {
	inspection := EvidenceInspection{Tree: plan.Tree, Oracle: plan.Oracle, Reason: plan.Reason}
	if !plan.Closed {
		return inspection
	}
	gitdir, err := benchgit.CommonDir(root)
	if err != nil {
		inspection.Reason = "evidence unavailable"
		return inspection
	}
	loaded := loadVerdict(evidencePath(gitdir, plan), now)
	if loaded.state != Ready || loaded.record.Status != "green" {
		if loaded.state == Ready {
			inspection.Reason = "recorded " + loaded.record.Status
		} else if loaded.state != Absent {
			inspection.Reason = "evidence unavailable"
		} else {
			inspection.Reason = "evidence absent"
		}
		return inspection
	}
	if loaded.record.Tree != plan.Tree || loaded.record.Oracle != plan.Oracle {
		inspection.Reason = "evidence changed"
		return inspection
	}
	recorded, _ := time.Parse(time.RFC3339, loaded.record.RecordedAt)
	inspection.RecordedAt = recorded
	if expires && now.Sub(recorded) >= freshness {
		inspection.Reason = "verdict expired"
		return inspection
	}
	inspection.ReusableGreen, inspection.Reason = true, ""
	return inspection
}

func evidencePath(gitdir string, plan subject) string {
	return filepath.Join(gitdir, "bench-gate-evidence", evidenceName(plan))
}

func evidenceName(plan subject) string {
	h := sha256.New()
	frame(h, plan.Tree)
	frame(h, plan.Oracle)
	return hex.EncodeToString(h.Sum(nil))
}

// openProspectiveArtifacts opens one artifact bundle for root and materializes tree in
// its private checkout. Every prospective producer -- the full gate, evidence
// inspection, and the fast lane -- enters its checkout through here, so the tree
// validation, the owner record, and the teardown scope have one source. The caller
// closes the returned owner.
func openProspectiveArtifacts(root, tree string) (*prospectiveartifact.Owner, error) {
	if !treeHashRE.MatchString(tree) {
		return nil, errors.New("invalid tree")
	}
	artifacts, err := prospectiveartifact.Open(root)
	if err != nil {
		return nil, err
	}
	if err := artifacts.Materialize(tree); err != nil {
		_ = artifacts.Close()
		return nil, err
	}
	return artifacts, nil
}

func durableReplaceAt(dir, name string, rec verdictRecord) error {
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return durableReplaceRecordAt(dir, name, data)
}

// durableReplaceRecordAt installs data as the store entry named name. It writes a private
// temporary file beside the entry, syncs it, and renames it over whatever was there. It
// then syncs the directory so the rename itself survives a crash. Every record class in
// the store publishes through here, so a reader never observes a half-written or
// world-readable record, whatever class wrote it.
func durableReplaceRecordAt(dir, name string, data []byte) error {
	data = append(data, '\n')
	tmp, err := os.CreateTemp(dir, ".bench-gate-evidence-")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = tmp.Close(); _ = os.Remove(tmpName) }()
	if err := tmp.Chmod(0o600); err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, filepath.Join(dir, name)); err != nil {
		return err
	}
	file, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Sync()
}

// reusableEvidence projects the verdict state from the plan the caller already accepted.
// Both call sites hold the current subject, so rebuilding one here would add a second,
// independent capture beside the generation that authorized the plan.
func reusableEvidence(root string, plan subject, now time.Time) Inspection {
	projection := inspectSubjectAt(root, plan, now)
	if projection.State == Pending || projection.State == Invalid || projection.State == Unavailable {
		return Inspection{}
	}
	evidence := inspectEvidence(root, plan, now)
	if !evidence.ReusableGreen {
		return Inspection{}
	}
	return Inspection{State: Ready, Status: "green", CachedTree: plan.Tree, CurrentTree: plan.Tree, RecordedAt: evidence.RecordedAt, ReusableGreen: true}
}

func runResolved(ctx context.Context, root string, res Resolution, env []string, stdout, stderr io.Writer, processGroup bool) processGroupResult {
	cmd := res.command(root)
	if cmd == nil {
		return processGroupResult{Code: 3}
	}
	cmd.Dir, cmd.Stdout, cmd.Stderr, cmd.Env = root, stdout, stderr, withGateSpanEnv(ctx, append([]string(nil), env...))
	if processGroup {
		return runProcessGroupCommand(ctx, cmd)
	}
	if err := cmd.Run(); err != nil {
		if cmd.ProcessState != nil {
			if code := cmd.ProcessState.ExitCode(); code > 0 {
				return processGroupResult{Code: code}
			}
		}
		return processGroupResult{Code: 1, StartErr: err}
	}
	return processGroupResult{}
}

type controlSafeWriter struct{ io.Writer }

func (w controlSafeWriter) Write(p []byte) (int, error) {
	safe := make([]byte, 0, len(p))
	for _, b := range p {
		if (b >= 0x20 && b != 0x7f) || b == '\n' || b == '\r' || b == '\t' {
			safe = append(safe, b)
		}
	}
	if _, err := w.Writer.Write(safe); err != nil {
		return 0, err
	}
	return len(p), nil
}
