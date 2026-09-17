package preflight

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
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
	if line := oversizedOperand(args); line != "" {
		return boundResponse(line+"\n", 2)
	}
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return boundResponse(line+"\n", code)
	}
	mode, slug := parsed.Positionals[0], parsed.Positionals[1]
	op, line := selectOperation(mode, parsed.Flags)
	if line != "" {
		return boundResponse(line+"\n", 2)
	}
	out, code := dispatch(version, op, slug, parsed.Flags, args)
	if op.bounded {
		return boundResponse(out, code)
	}
	return out, code
}

func dispatch(version string, op operation, slug string, flags map[string]string, args []string) (string, int) {
	base, sourceTip, ticket := flags[flagBase], flags[flagTip], flags[flagTicket]
	_, full := flags[flagFull]
	quota := uint64(chargeevidence.DefaultQuota)
	if text, ok := flags[flagQuota]; ok {
		value, valid := chargeevidence.ParseDecimal(text)
		if !valid || value == 0 {
			return toon.Usage(grammar.Cmd, flagQuota+" needs a positive decimal byte count within the unsigned 64-bit range") + "\n", 2
		}
		quota = value
	}
	if op.kind == opReadEvidence {
		if line := evidenceOperandRefusal(slug, flags[flagCursor]); line != "" {
			return line + "\n", 2
		}
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	if err := unrepresentableCell("--source-tip", sourceTip); err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	switch op.kind {
	case opLegacyCharge:
		return chargeCommand(root, op.mode, slug, base, sourceTip, ticket, full, version, args)
	case opProposal:
		return proposeWritesCommand(root, op.mode, slug, base, sourceTip, ticket, args)
	case opPrepareEvidence:
		return prepareEvidenceCommand(root, slug, base, sourceTip, ticket, quota, args)
	case opReadEvidence:
		return readEvidenceCommand(root, slug, flags[flagCursor])
	}
	return verdictCommand(root, op.mode, slug, base, sourceTip, args)
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

// maxOperandBytes bounds every preflight operand before the parser can echo it. Every valid
// selector, pin, identifier, cursor, and quota is far shorter.
const maxOperandBytes = 1024

// oversizedOperand refuses an operand too long to echo, naming it only by length and digest.
func oversizedOperand(args []string) string {
	for _, arg := range args {
		if len(arg) > maxOperandBytes {
			return toon.Usage(grammar.Cmd, boundedOperand("oversized operand", arg))
		}
	}
	return ""
}

// boundedOperand identifies a hostile operand by type, length, and digest, never by value.
func boundedOperand(kind, value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%s bytes=%d sha256=%s", kind, len(value), hex.EncodeToString(sum[:]))
}

// boundResponse is the shared final guard for bounded forms and every usage line: a response
// above chargeevidence.ResponseLimit encoded bytes becomes one bounded operational refusal.
func boundResponse(out string, code int) (string, int) {
	if len(out) <= chargeevidence.ResponseLimit {
		return out, code
	}
	return toon.Errorf("response bound exceeded", fmt.Sprintf("the response held %d bytes above the %d-byte limit; report this defect", len(out), chargeevidence.ResponseLimit)) + "\n", 1
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
