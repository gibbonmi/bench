package preflight

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

const modeBuild = evidencecmd.ModeBuild

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
	mode, slug := parsed.Positionals[0], parsed.Positionals[1]
	op, line := evidencecmd.Select(mode, parsed.Flags)
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
	_, full := flags[evidencecmd.FlagFull]
	quota, line := evidencecmd.Admit(op, slug, flags)
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
	case evidencecmd.KindLegacyCharge:
		return chargeCommand(root, op.Mode, slug, base, sourceTip, ticket, full, version, args)
	case evidencecmd.KindProposal:
		return proposeWritesCommand(root, op.Mode, slug, base, sourceTip, ticket, args)
	case evidencecmd.KindPrepareEvidence:
		return prepareEvidenceCommand(root, slug, base, sourceTip, ticket, quota, args)
	case evidencecmd.KindReadEvidence:
		return evidencecmd.Read(root, slug, flags)
	}
	return verdictCommand(root, op.Mode, slug, base, sourceTip, args)
}

// prepareEvidenceCommand runs the movement-checked build preparation and hands each
// attempt's validated pack, or its refusal, to the evidence publisher.
func prepareEvidenceCommand(root, slug, base, sourceTip, name string, quota uint64, args []string) (string, int) {
	return evidencecmd.Prepare(root, quota, func(stage func(*chargeevidence.Pack, string, string) string) (string, int) {
		return preparedAttempts(root, modeBuild, slug, base, sourceTip, "charge", args, func(facts Facts) (string, int) {
			pack, refusal := buildChargePack(root, facts, boundedVerdict(Decide(facts)), name, buildSourcePolicy())
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

func verdictCommand(root, mode, slug, base, sourceTip string, args []string) (string, int) {
	facts, bootErr := GatherPinned(root, mode, slug, base, sourceTip)
	if bootErr != nil {
		if bootErr.Kind == "snapshot drift" {
			return snapshotDriftRefusal(args, bootErr.Hint), 1
		}
		return toon.Errorf(bootErr.Kind, bootErr.Hint) + "\n", 1
	}
	if err := unrepresentableChangedPath(facts.ChangedPaths); err != nil {
		return toon.RenderError(err) + "\n", 1
	}

	verdict := Decide(facts)
	rows := make([][]string, len(verdict.Checks))
	for i, c := range verdict.Checks {
		rows[i] = []string{c.Check, c.Verdict, c.Detail, c.Next}
	}
	tbl, err := toon.Table("checks", []string{"check", "verdict", "detail", "next"}, rows)
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
	b.WriteString(tbl)

	exit := 0
	if verdict.Red {
		exit = 1
	}
	return b.String(), exit
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
