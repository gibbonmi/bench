package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/signal"
	"syscall"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/specbuild"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/worktree"
)

// specBuildCommand exposes the spec-build lifecycle without duplicating its durable state.
func specBuildCommand(args []string) (string, int) {
	invocation, out, code := spec.ParseBuild(args)
	if out != "" {
		return out, code
	}
	root, err := git.Root()
	if err != nil {
		return buildError(errors.New("Git repository unavailable"), "run this command inside the working checkout")
	}
	service := specbuild.New(root, productionGateOwner{}, productionWorktreeOwner{})
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()
	return executeBuild(ctx, service, invocation.Operation, invocation.Slug, invocation.Flags)
}

type buildService interface {
	Start(context.Context, string) (specbuild.Status, error)
	Assign(context.Context, string, string, string) (specbuild.Assignment, specbuild.Status, error)
	Checkpoint(context.Context, string, string, string) (specbuild.Status, error)
	Integrate(context.Context, string, string) (specbuild.Status, error)
	Review(context.Context, string, string) (specbuild.Status, error)
	Status(string) (specbuild.Status, error)
	FullStatus(string) (specbuild.FullStatus, error)
	Promote(context.Context, string) (specbuild.Status, error)
	Abandon(context.Context, string) (specbuild.AbandonmentPlan, error)
	ApplyAbandon(context.Context, string, string) (specbuild.Status, error)
	Reclaim(string) (specbuild.ReclamationPlan, error)
	ApplyReclaim(string, string) (specbuild.ReclamationPlan, error)
}

func dispatchSpec(args []string, stdout io.Writer) int {
	var out string
	var code int
	if len(args) > 0 && args[0] == "build" {
		out, code = specBuildCommand(args[1:])
	} else {
		out, code = spec.Command(args)
	}
	fmt.Fprint(stdout, out)
	return code
}

func executeBuild(ctx context.Context, service buildService, operation, slug string, flags map[string]string) (string, int) {
	var status specbuild.Status
	var err error
	switch operation {
	case "start":
		status, err = service.Start(ctx, slug)
	case "assign":
		var assignment specbuild.Assignment
		assignment, status, err = service.Assign(ctx, slug, flags["--ticket"], flags["--request"])
		if err == nil {
			return renderAssignment(status, assignment)
		}
	case "checkpoint":
		status, err = service.Checkpoint(ctx, slug, flags["--assignment"], flags["--evidence"])
	case "integrate":
		status, err = service.Integrate(ctx, slug, flags["--assignment"])
	case "review":
		status, err = service.Review(ctx, slug, flags["--evidence"])
	case "status":
		if _, full := flags["--full"]; full {
			projection, fullErr := service.FullStatus(slug)
			if fullErr != nil {
				return buildError(fullErr, "inspect the named spec and retained build state")
			}
			return renderFullStatus(projection)
		}
		status, err = service.Status(slug)
	case "promote":
		status, err = service.Promote(ctx, slug)
	case "abandon":
		if fingerprint, apply := flags["--apply"]; apply {
			status, err = service.ApplyAbandon(ctx, slug, fingerprint)
		} else {
			plan, planErr := service.Abandon(ctx, slug)
			if planErr != nil {
				return buildError(planErr, "resolve the reported ownership or state drift, then plan again")
			}
			return renderAbandonment(plan)
		}
	case "reclaim":
		var plan specbuild.ReclamationPlan
		if fingerprint, apply := flags["--apply"]; apply {
			plan, err = service.ApplyReclaim(slug, fingerprint)
		} else {
			plan, err = service.Reclaim(slug)
		}
		if err != nil {
			return reclamationFailure(plan, err)
		}
		return renderReclamation(plan)
	}
	if err != nil {
		return buildError(err, "resolve the reported lifecycle precondition, then retry the same command")
	}
	return renderStatus(status)
}

func buildError(err error, hint string) (string, int) {
	return toon.Errorf(sanitize.Controls(err.Error()), sanitize.Controls(hint)) + "\n", 1
}

func renderStatus(status specbuild.Status) (string, int) {
	out, err := toon.Table("spec_build", []string{"slug", "state", "subject", "next"}, [][]string{{status.Slug, status.State, status.Subject, status.Next}})
	if err != nil {
		return buildError(err, "remove control text from Git-sourced lifecycle identities")
	}
	return out, 0
}

func renderAssignment(status specbuild.Status, assignment specbuild.Assignment) (string, int) {
	out, code := renderStatus(status)
	if code != 0 {
		return out, code
	}
	assignmentOut, err := toon.Table("assignment", []string{"id", "path", "base"}, [][]string{{assignment.ID, assignment.Path, assignment.Base}})
	if err != nil {
		return buildError(err, "remove control text from assignment identities")
	}
	return out + assignmentOut, 0
}

func renderFullStatus(full specbuild.FullStatus) (string, int) {
	out, code := renderStatus(full.Status)
	if code != 0 {
		return out, code
	}
	rows := make([][]string, 0, len(full.Assignments))
	for _, assignment := range full.Assignments {
		rows = append(rows, []string{assignment.ID, assignment.Ticket, assignment.TicketDigest, assignment.Base, assignment.Checkpoint, assignment.CheckpointRef, assignment.Integrated, assignment.ReceiptDigest, assignment.Cleanup})
	}
	assignments, err := toon.Table("assignments", []string{"id", "ticket", "ticket_digest", "base", "checkpoint", "checkpoint_ref", "integrated", "receipt_digest", "cleanup"}, rows)
	if err != nil {
		return buildError(err, "remove control text from retained assignment evidence")
	}
	out += assignments
	if full.Review == nil {
		review, _ := toon.Table("review", []string{"candidate", "axis", "finding", "disposition", "receipt_digest"}, nil)
		return out + review, 0
	}
	var reviewRows [][]string
	for _, axis := range full.Review.Axes {
		if len(axis.Findings) == 0 {
			reviewRows = append(reviewRows, []string{full.Review.Candidate, axis.Axis, "", "", full.Review.Digest})
		}
		for _, finding := range axis.Findings {
			reviewRows = append(reviewRows, []string{full.Review.Candidate, axis.Axis, finding.ID, finding.Disposition, full.Review.Digest})
		}
	}
	review, err := toon.Table("review", []string{"candidate", "axis", "finding", "disposition", "receipt_digest"}, reviewRows)
	if err != nil {
		return buildError(err, "remove control text from retained review evidence")
	}
	return out + review, 0
}

func renderAbandonment(plan specbuild.AbandonmentPlan) (string, int) {
	rows := make([][]string, 0, len(plan.Worktrees))
	for _, item := range plan.Worktrees {
		rows = append(rows, []string{item.ID, item.Path, item.Request, item.OwnerFingerprint})
	}
	out, err := toon.Table("abandon", []string{"fingerprint", "worktrees", "provisional_refs", "checkpoints", "recovery_refs"}, [][]string{{plan.Fingerprint, fmt.Sprint(len(plan.Worktrees)), fmt.Sprint(len(plan.ProvisionalRefs)), fmt.Sprint(len(plan.UnintegratedCheckpoints)), fmt.Sprint(len(plan.RecoveryRefs))}})
	if err != nil {
		return buildError(err, "remove control text from abandonment evidence")
	}
	worktrees, err := toon.Table("abandon_worktrees", []string{"id", "path", "request", "owner_fingerprint"}, rows)
	if err != nil {
		return buildError(err, "remove control text from abandonment evidence")
	}
	return out + worktrees, 0
}

// reclamationFailure keeps the receipt an interrupted apply spent. Git deletes one ref at a
// time with no transaction across the set, so a refusal partway through has already taken the
// refs the receipt names, and nothing can re-read them afterwards — the receipt is their only
// account. It prints ahead of the refusal, which still carries the exit status and the
// convergence hint, so the spent set never reads as a completed operation. A failure that took
// no ref carries no receipt and stays a refusal alone: an empty table would report a deletion
// set for an operation that deleted nothing.
func reclamationFailure(receipt specbuild.ReclamationPlan, err error) (string, int) {
	refusal, code := buildError(err, "plan again with bench spec build reclaim <slug>, then apply its exact fingerprint")
	if len(receipt.Refs) == 0 {
		return refusal, code
	}
	spent, spentCode := renderReclamation(receipt)
	if spentCode != 0 {
		return spent, spentCode
	}
	return spent + refusal, code
}

// renderReclamation reports every class the plan carries, not only the deletable one, so a
// maintainer sees what was left behind and why alongside what the apply would take.
func renderReclamation(plan specbuild.ReclamationPlan) (string, int) {
	headers := []string{"slug", "fingerprint", "applied"}
	summary := []string{plan.Slug, plan.Fingerprint, fmt.Sprint(plan.Applied)}
	for _, class := range plan.Classes() {
		headers = append(headers, class.Name)
		summary = append(summary, fmt.Sprint(class.Count))
	}
	out, err := toon.Table("reclaim", headers, [][]string{summary})
	if err != nil {
		return buildError(err, "remove control text from reclamation evidence")
	}
	rows := make([][]string, 0, len(plan.Refs))
	for _, ref := range plan.Refs {
		rows = append(rows, []string{ref.Name, ref.Object, ref.Assignment, ref.Disposition})
	}
	refs, err := toon.Table("reclaim_refs", []string{"name", "object", "assignment", "disposition"}, rows)
	if err != nil {
		return buildError(err, "remove control text from reclamation evidence")
	}
	return out + refs, 0
}

type productionWorktreeOwner struct{}

func (productionWorktreeOwner) Create(_ context.Context, root, request, label, start string) (specbuild.OwnedWorktree, error) {
	created, err := worktree.Create(root, request, label, nil, start)
	return specbuild.OwnedWorktree{ID: created.Assignment.ID, Path: created.Path, Branch: created.Assignment.Branch}, err
}

func (productionWorktreeOwner) Release(_ context.Context, root, request, path string, evidence specbuild.ReleaseEvidence) error {
	return worktree.ReleaseProvisional(root, request, path, worktree.ProvisionalEvidence{
		Base: evidence.Base, CheckpointRef: evidence.CheckpointRef, Checkpoint: evidence.Checkpoint, IntegratedRef: evidence.IntegratedRef, Integrated: evidence.Integrated,
	})
}

func (productionWorktreeOwner) PlanAbandon(_ context.Context, root, request, path string) (string, error) {
	return worktree.PlanAbandon(root, request, path)
}

func (productionWorktreeOwner) ApplyAbandon(_ context.Context, root, request, path, fingerprint string) error {
	_, err := worktree.ApplyAbandon(root, request, path, fingerprint)
	return err
}

type productionGateOwner = specbuild.AuthorizationGate
