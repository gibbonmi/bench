package gate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	"github.com/gibbonmi/bench/internal/gate/prospectiveartifact"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// EvidenceInspection reports whether retained gate evidence can authorize one tree. Grade is
// the grade key that names the evidence record.
type EvidenceInspection struct {
	Tree          string
	Grade         string
	RecordedAt    time.Time
	ReusableGreen bool
	Reason        string
}

// InspectTree reports retained exact evidence for tree. It does not execute the gate
// or change evidence.
func InspectTree(root, tree string) EvidenceInspection {
	return InspectTreeContext(context.Background(), root, tree)
}

// InspectTreeContext checks retained evidence for the same obligation used by execution.
func InspectTreeContext(ctx context.Context, root, tree string) EvidenceInspection {
	return inspectProspective(ctx, root, tree, time.Now(), true)
}

// ExecuteTree executes or reuses the gate for tree and returns its evidence inspection.
// It keeps the prospective checkout private and records evidence under root. If ctx
// ends, the incomplete run does not become verdict evidence.
func ExecuteTree(ctx context.Context, root, tree string, stdout, stderr io.Writer) Result {
	ctx, finishSpan := beginGateSpan(ctx, root, "prospective")
	ctx, finishLog := beginGateRunLog(ctx, root, stderr, "prospective")
	result := executeTreeWithOwner(ctx, root, tree, stdout, stderr, nil, reuseFreshGreen, nil)
	finishLog(result)
	finishSpan(result)
	return result
}

func executeTreeWithOwner(ctx context.Context, root, tree string, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode, owner runBinaryOwner) Result {
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
	baseline := baselineRoot(ctx, checkout, root)
	evaluation := checkpointEvaluation(ctx, newProspectiveTreeEvaluation(checkout, root, baseline, tree))
	return executeSubjectOwning(ctx, checkout, root, stdout, stderr, arm, mode, evaluation, owner, baseline, durablefile.Native(), artifacts)
}

// ValidateProjectGreen reports whether branch's current tip and marker have retained
// exact green evidence. It does not execute the gate or change evidence.
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

// inspectProspective reports retained exact evidence for tree as the prospective owner keys
// it. expires selects the freshness window, and the delta base reads with no window.
func inspectProspective(ctx context.Context, root, tree string, now time.Time, expires bool) EvidenceInspection {
	artifacts, err := openProspectiveArtifacts(root, tree)
	if err != nil {
		return EvidenceInspection{Reason: "subject unavailable"}
	}
	defer artifacts.Close()
	evaluation := checkpointEvaluation(ctx, newProspectiveTreeEvaluation(artifacts.Checkout(), root, root, tree))
	plan, err := evaluation.acceptPre()
	if err != nil || plan.Tree != tree {
		return EvidenceInspection{Reason: "subject unavailable"}
	}
	return inspectEvidenceWindowed(root, plan, now, expires)
}

func inspectEvidence(root string, plan subject, now time.Time) EvidenceInspection {
	return inspectEvidenceWindowed(root, plan, now, true)
}

func inspectEvidenceWindowed(root string, plan subject, now time.Time, expires bool) EvidenceInspection {
	inspection := EvidenceInspection{Tree: plan.Tree, Grade: evidenceName(plan), Reason: plan.Reason}
	if !plan.Closed {
		return inspection
	}
	gitdir, err := benchgit.CommonDir(root)
	if err != nil {
		inspection.Reason = "evidence unavailable"
		return inspection
	}
	loaded := loadEvidence(evidencePath(gitdir, plan), now)
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
	// A record answers only under its own name, so a copied or renamed record is changed evidence.
	if loaded.record.Tree != plan.Tree || loaded.record.Grade != inspection.Grade {
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

// evidenceName answers the grade key: the tree framed with the grade identity. The key holds
// no checkout path and no run purpose, so a green of one checkout answers a run of the same
// key from any checkout, and each verb runs its own purpose checks before the lookup.
func evidenceName(plan subject) string {
	h := sha256.New()
	frame(h, plan.Tree)
	frame(h, plan.Identity)
	return hex.EncodeToString(h.Sum(nil))
}

// evidenceFields is the exact field set of an evidence record. The record names its grade
// key in place of the per-checkout oracle digest.
var evidenceFields = []string{"grade", "recorded_at", "schema", "state", "status", "tree"}

// loadEvidence reads one evidence record under the file discipline of every store record. The
// store holds whole greens only, so the record is a ready verdict.
func loadEvidence(path string, now time.Time) loadedVerdict {
	read := readStoreRecord(path)
	loaded := loadedVerdict{state: read.state, reason: read.reason, bytes: read.bytes}
	if read.data == nil {
		return loaded
	}
	record := &loaded.record
	if strictJSON(read.data, record) != nil || requireObjectFields(read.data, evidenceFields) != nil ||
		validateStoreRecordCommon(*record) != nil || validateReadyRecord(*record, now) != nil {
		loaded.state, loaded.reason = Invalid, "invalid cache record"
		return loaded
	}
	loaded.state = Ready
	return loaded
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

func durableReplaceAt(replacer durablefile.Replacer, dir, name string, rec verdictRecord) error {
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return durableReplaceRecordAt(replacer, dir, name, data)
}

// durableReplaceRecordAt installs data with a final newline as the private store entry
// named name through replacer. Every record class in the store publishes through here,
// so a reader never observes a half-written or world-readable record, whatever class
// wrote it.
func durableReplaceRecordAt(replacer durablefile.Replacer, dir, name string, data []byte) error {
	return replacer.Replace(filepath.Join(dir, name), append(data, '\n'), 0o600)
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

func runResolved(ctx context.Context, root string, res Resolution, env []string, stdout, stderr io.Writer, local ...protection.Descriptor) processGroupResult {
	cmd := res.command(root)
	if cmd == nil {
		return processGroupResult{Code: 3}
	}
	cmd.Dir, cmd.Stdout, cmd.Stderr, cmd.Env = root, stdout, stderr, withGateSpanEnv(ctx, append([]string(nil), env...))
	return runProcessGroupCommand(ctx, cmd, local...)
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
