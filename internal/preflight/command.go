package preflight

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/chargeevidence"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/preflight/staleness"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

const (
	modeBuild  = evidencecmd.ModeBuild
	modeReview = evidencecmd.ModeReview
)

// Command is the legacy adapter for `bench preflight review <slug>` and `bench
// preflight build <slug>`. It is the CLI-contract seam. Grammar and usage errors ride
// usage.Parse and the operation registry (exit 2). A not-in-repo cwd or a bootstrap
// failure is one toon.Errorf line (exit 1). Otherwise the verdict renders as TOON and the
// exit code follows Verdict.Red (0 green, 1 red).
func Command(args []string) (string, int) {
	return CommandWithVersion("")(args)
}

// CommandWithVersion returns the preflight command bound to one executable version.
// Review charges pass that version to the consumer citation owner.
func CommandWithVersion(version string) func([]string) (string, int) {
	return func(args []string) (string, int) { return command(version, args) }
}

func command(version string, args []string) (string, int) {
	if line := evidencecmd.OversizedOperand(args); line != "" {
		return evidencecmd.Bound(line+"\n", 2)
	}
	parsed, line, code := usage.Parse(evidencecmd.Grammar, args)
	if line != "" {
		return evidencecmd.Bound(line+"\n", code)
	}
	mode := parsed.Positionals[0]
	op, line := evidencecmd.Select(mode, parsed.Flags)
	if line != "" {
		return evidencecmd.Bound(line+"\n", 2)
	}
	// The registry owns how many operands each mode takes, so a mode that takes none
	// refuses a second operand and a mode that takes one refuses its absence.
	slug, line := evidencecmd.SelectOperand(op, parsed.Positionals)
	if line != "" {
		return evidencecmd.Bound(line+"\n", 2)
	}
	out, code := dispatch(version, op, slug, parsed.Flags, args)
	if op.Bounded {
		return evidencecmd.Bound(out, code)
	}
	return out, code
}

func dispatch(version string, op evidencecmd.Operation, slug string, flags map[string]string, args []string) (string, int) {
	base, sourceTip, ticket := flags[evidencecmd.FlagBase], flags[evidencecmd.FlagTip], flags[evidencecmd.FlagTicket]
	quota, line := evidencecmd.SelectQuota(op, slug, flags)
	if line != "" {
		return line + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	if err := unrepresentableCell("--source-tip", sourceTip); err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	switch op.Kind {
	case evidencecmd.KindPrepareEvidence:
		return prepareEvidenceCommand(root, op.Kind, slug, base, sourceTip, ticket, quota, args)
	case evidencecmd.KindPrepareReviewEvidence:
		return prepareReviewEvidenceCommand(root, op.Kind, version, slug, base, sourceTip, quota, args)
	case evidencecmd.KindReadEvidence:
		return evidencecmd.Read(root, slug, flags)
	case evidencecmd.KindVerifyEvidence:
		return evidencecmd.Verify(root, slug)
	case evidencecmd.KindCurrentEvidence:
		return currentEvidenceCommand(root, op.Kind, slug, args)
	case evidencecmd.KindExportEvidence:
		return evidencecmd.Export(root, slug, flags[evidencecmd.FlagTo], chargeevidence.WriteAll)
	case evidencecmd.KindCleanPlan:
		return evidencecmd.CleanPlan(root, flags)
	case evidencecmd.KindCleanApply:
		return evidencecmd.CleanApply(root, flags)
	}
	// The audit form gathers its staleness facts in the same movement-checked attempt.
	gatherFacts := GatherPinned
	if op.Kind == evidencecmd.KindAuditSpec {
		gatherFacts = gatherAudit
	}
	return verdictCommand(root, op.Kind, op.Mode, args, func() (Facts, *BootstrapFailure) {
		return gatherFacts(root, op.Mode, slug, base, sourceTip)
	})
}

// prepareEvidenceCommand runs the movement-checked build preparation and hands each
// attempt's validated pack, or its refusal, to the evidence publisher. kind is the
// operation that the dispatcher selected, and the facts carry it into Decide.
func prepareEvidenceCommand(root string, kind evidencecmd.Kind, slug, base, sourceTip, name string, quota uint64, args []string) (string, int) {
	return evidencecmd.Prepare(root, quota, func(stage func(*chargeevidence.Pack, string, string) string) (string, int) {
		return preparedAttempts(root, modeBuild, slug, base, sourceTip, "charge", args, func(facts Facts) (string, int) {
			facts.Operation = kind
			pack, refusal := buildChargePack(root, facts, boundedVerdict(Decide(facts)), name, buildSourcePolicy())
			if refusal == "" {
				if err := (commitrepo.Store{Root: root}).Ready(facts.SpecPath); err != nil {
					refusal = chargeRefusal("commitment", err.Error(), "run bench commitment start for this deliverable")
				}
			}
			if refusal = stage(pack, facts.AssignmentTarget, refusal); refusal != "" {
				return refusal, 1
			}
			return "", 0
		})
	})
}

// prepareReviewEvidenceCommand runs the movement-checked review preparation on the same
// publisher the build uses. Each attempt runs the collectors once and stages its validated
// pack; a refused attempt publishes no handle, and only the unmoved final attempt publishes.
// kind is the operation that the dispatcher selected, and the facts carry it into Decide.
func prepareReviewEvidenceCommand(root string, kind evidencecmd.Kind, version, slug, base, sourceTip string, quota uint64, args []string) (string, int) {
	return evidencecmd.Prepare(root, quota, func(stage func(*chargeevidence.Pack, string, string) string) (string, int) {
		return preparedAttempts(root, modeReview, slug, base, sourceTip, "charge", args, func(facts Facts) (string, int) {
			facts.Operation = kind
			pack, refusal := reviewChargePack(root, version, facts, boundedVerdict(Decide(facts)))
			if refusal = stage(pack, facts.AssignmentTarget, refusal); refusal != "" {
				return refusal, 1
			}
			return "", 0
		})
	})
}

// boundedVerdict identifies each check detail too long for a bounded response by type,
// length, and digest. The check names, verdicts, and remedies stay unchanged.
func boundedVerdict(verdict Verdict) Verdict {
	checks := make([]CheckResult, len(verdict.Checks))
	for i, check := range verdict.Checks {
		check.Detail = evidencecmd.BoundedDiagnostic("detail", check.Detail)
		checks[i] = check
	}
	verdict.Checks = checks
	return verdict
}

// verdictCommand renders the ordinary, plan-only, or audit verdict of the facts that gathered
// returns. kind is the operation that the dispatcher selected, and the facts carry it into
// Decide. The audit form prints its staleness tables after the checks. A supported audit red
// exits 1 like a readiness red, and a diagnostic-only audit exits 0 with manual review still
// required. The audit decides only after the commitment check, so its omission needs that row
// green too. args is the invocation that a snapshot-drift refusal retries. Nil args name no
// public invocation, so the refusal prints no retry action.
func verdictCommand(root string, kind evidencecmd.Kind, mode string, args []string, gathered func() (Facts, *BootstrapFailure)) (string, int) {
	facts, bootErr := gathered()
	if bootErr != nil {
		if bootErr.Kind == "snapshot drift" && args != nil {
			return snapshotDriftRefusal(args, bootErr.Hint), 1
		}
		return toon.Errorf(bootErr.Kind, bootErr.Hint) + "\n", 1
	}
	if err := unrepresentableChangedPath(facts.ChangedPaths); err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	facts.Operation = kind

	verdict := Decide(facts)
	if mode == modeBuild && kind != evidencecmd.KindPlanOnly && !verdict.Red {
		if err := (commitrepo.Store{Root: root}).Ready(facts.SpecPath); err != nil {
			verdict.Checks = append(verdict.Checks, CheckResult{Check: "commitment", Verdict: verdictRed, Detail: err.Error(), Next: "bench commitment start"})
			verdict.Red = true
		}
	}
	checks, err := renderChecks(verdict)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}

	var b strings.Builder
	fmt.Fprintf(&b, "phase: %s\n", mode)
	fmt.Fprintf(&b, "spec: %s\n", facts.SpecPath)
	if facts.SourceBase != "" {
		source, err := toon.Table("source", []string{"base", "tip"}, [][]string{{facts.SourceBase, facts.SourceTip}})
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		b.WriteString(source)
	}
	b.WriteString(checks)

	exit := 0
	if verdict.Red {
		exit = 1
	}
	if kind == evidencecmd.KindAuditSpec {
		audit := staleness.Decide(facts.Audit, verdict.requiredGreen())
		tables, err := audit.Render()
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		b.WriteString(tables)
		if audit.Red() {
			exit = 1
		}
	}
	return b.String(), exit
}

// renderChecks prints one count line for the verdict and a table of the red rows
// only, so a green verdict costs one line and a red row keeps its detail and next.
// The charge forms do not use this render: a red check refuses the charge, and a
// green charge prints no check row.
func renderChecks(verdict Verdict) (string, error) {
	counts := map[string]int{}
	for _, c := range verdict.Checks {
		counts[c.Verdict]++
	}
	summary := fmt.Sprintf("checks{green=%d,not_applicable=%d,red=%d}\n",
		counts[verdictGreen], counts[verdictNA], counts[verdictRed])
	red := redRows(verdict)
	if len(red) == 0 {
		return summary, nil
	}
	tbl, err := toon.Table("checks", checkFields, red)
	if err != nil {
		return "", err
	}
	return summary + tbl, nil
}

// checkFields are the columns of every rendered red check row, and redRows is the row of
// each red check of verdict, in check order.
var checkFields = []string{"check", "verdict", "detail", "next"}

func redRows(verdict Verdict) [][]string {
	var red [][]string
	for _, c := range verdict.Checks {
		if c.Verdict == verdictRed {
			red = append(red, []string{c.Check, c.Verdict, c.Detail, c.Next})
		}
	}
	return red
}

func snapshotDriftRefusal(args []string, hint string) string {
	invocation := make([]axi.InvocationArgument, 0, len(args)+1)
	invocation = append(invocation, axi.KnownArgument("preflight"))
	for _, arg := range args {
		invocation = append(invocation, axi.KnownArgument(arg))
	}
	help, err := axi.RenderHelp([]axi.Action{axi.RetryInvocation(invocation...)})
	if err != nil {
		return toon.RenderError(err) + "\n"
	}
	return toon.Errorf("snapshot drift", hint) + "\n" + help
}

// unrepresentableChangedPath refuses a changed path spec-TOON cannot render as a
// cell, before the verdict table is ever built. PF7's contract is unconditional.
// A path carrying a control byte exits 1 the same way in every case: a
// green row (never rendered) or a red row's detail cell. The refusal
// never depends on which row a later check sorts it into.
func unrepresentableChangedPath(paths []string) error {
	for _, p := range paths {
		if err := unrepresentableCell("changed path", p); err != nil {
			return err
		}
	}
	return nil
}

// unrepresentableCell is that refusal for one value. The changed-path
// sweep and --source-tip both share it: a pin reaches a detail cell and
// the snapshot-drift retry action. So a control byte in it gets refused,
// not rendered. The %q quoting keeps the offending byte out of the
// message it explains.
func unrepresentableCell(what, value string) error {
	if toon.Representable(value) {
		return nil
	}
	return fmt.Errorf("%s %q contains a byte spec-TOON cannot represent", what, value)
}
