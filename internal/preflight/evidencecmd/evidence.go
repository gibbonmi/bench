package evidencecmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// evidenceRecovery is the next action for each store refusal class that has one. A class
// without a row keeps the generic retry action.
var evidenceRecovery = map[string]string{
	chargeevidence.RefuseAbsent:   "run the build preparation that prints this identifier",
	chargeevidence.RefuseUnsafe:   "repair the evidence store by hand; Bench removed nothing",
	chargeevidence.RefuseExisting: "inspect the corrupt published artifact; Bench replaced nothing",
	chargeevidence.RefuseReplaced: "rerun the exact read after the store stops changing",
	chargeevidence.RefuseCursor:   "rerun the read with the next command a previous page printed",
}

// evidenceStore opens the repository-common store for root. It creates nothing.
func evidenceStore(root string) (*chargeevidence.Store, string) {
	common, err := git.CommonDir(root)
	if err != nil {
		return nil, toon.Errorf("evidence store unavailable", err.Error()) + "\n"
	}
	return chargeevidence.OpenStore(common, chargeevidence.StoreOptions{Pause: chargeevidence.PauseFromEnvironment()}), ""
}

// Admit validates op's evidence operands before any repository access: the quota operand,
// then a read's identifier and cursor. It returns the selected quota, or the usage line.
func Admit(op Operation, identity string, flags map[string]string) (uint64, string) {
	quota := uint64(chargeevidence.DefaultQuota)
	if text, ok := flags[flagQuota]; ok {
		value, valid := chargeevidence.ParseDecimal(text)
		if !valid || value == 0 {
			return 0, toon.Usage(Grammar.Cmd, flagQuota+" needs a positive decimal byte count within the unsigned 64-bit range")
		}
		quota = value
	}
	if op.Kind == KindReadEvidence {
		return quota, operandRefusal(identity, flags[flagCursor])
	}
	return quota, ""
}

// Prepare publishes one immutable build evidence artifact. run executes the
// movement-checked preparation attempts. Each attempt that reaches its build passes stage
// its validated pack and assignment, or its refusal. stage discards the previous attempt's
// candidate and stages this one's verified temporary pack; only the unmoved final attempt
// publishes, and every other staged candidate is discarded.
func Prepare(root string, quota uint64, run func(stage func(pack *chargeevidence.Pack, assignment, refusal string) string) (string, int)) (string, int) {
	store, refusal := evidenceStore(root)
	if refusal != "" {
		return refusal, 1
	}
	var (
		staged     *chargeevidence.Staged
		pack       *chargeevidence.Pack
		assignment string
		attempt    int
	)
	out, code := run(func(built *chargeevidence.Pack, target, refusal string) string {
		attempt++
		staged.Discard()
		staged, pack = nil, nil
		if refusal != "" {
			return refusal
		}
		candidate, err := store.Stage(built, quota, attempt)
		if err != nil {
			return storeRefusal(err)
		}
		staged, pack, assignment = candidate, built, target
		return ""
	})
	if code != 0 {
		staged.Discard()
		return out, code
	}
	// The response renders before publication, so a response that cannot render
	// publishes nothing.
	m := pack.Manifest()
	text, err := chargeevidence.Prepared{
		Evidence: pack.Identity(), Mode: m.Selection.Mode, Base: m.Selection.Base, SourceTip: m.Selection.SourceTip,
		Assignment: assignment, Sources: len(m.Sources), Pages: len(m.Pages), ManifestBytes: len(pack.ManifestBytes()),
		Next: evidenceInvocation(pack.Identity(), ""),
	}.Encode()
	if err != nil {
		staged.Discard()
		return toon.RenderError(err) + "\n", 1
	}
	if _, err := staged.Publish(attempt); err != nil {
		return storeRefusal(err), 1
	}
	return text, 0
}

// operandRefusal validates the read operands before any path use. A refusal names a
// hostile operand only by length and digest.
func operandRefusal(identity, cursor string) string {
	if !chargeevidence.ValidIdentity(identity) {
		return toon.Usage(Grammar.Cmd, boundedOperand("invalid evidence identifier", identity))
	}
	if cursor == "" {
		return ""
	}
	if _, err := chargeevidence.ParseCursor(cursor, identity); err != nil {
		return toon.Usage(Grammar.Cmd, boundedOperand("invalid cursor", cursor)+" "+refusalClass(err))
	}
	return ""
}

// Read prints one bounded fragment of the default evidence stream at the position the
// cursor flag names. It keeps no reading state: the cursor alone names the position.
func Read(root, identity string, flags map[string]string) (string, int) {
	store, refusal := evidenceStore(root)
	if refusal != "" {
		return refusal, 1
	}
	artifact, err := store.Open(identity)
	if err != nil {
		return storeRefusal(err), 1
	}
	defer artifact.Close()
	cursor := artifact.First()
	if cursorText := flags[flagCursor]; cursorText != "" {
		if cursor, err = chargeevidence.ParseCursor(cursorText, identity); err != nil {
			return storeRefusal(err), 1
		}
	}
	fragment, err := artifact.Read(cursor)
	if err != nil {
		return storeRefusal(err), 1
	}
	next := ""
	if fragment.Next != nil {
		next = evidenceInvocation(identity, fragment.Next.String())
	}
	text, err := fragment.Encode(identity, next)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return text, 0
}

// evidenceInvocation is the exact read command for one artifact position. An empty cursor
// is the manifest-first command.
func evidenceInvocation(identity, cursor string) string {
	args := []string{"bench", "preflight", modeEvidence, identity}
	if cursor != "" {
		args = append(args, flagCursor, cursor)
	}
	for i := range args {
		args[i] = axi.ShellQuote(args[i])
	}
	return strings.Join(args, " ")
}

// storeRefusal renders one bounded store refusal and its recovery action.
func storeRefusal(err error) string {
	var capacity *chargeevidence.CapacityError
	if errors.As(err, &capacity) {
		kind := fmt.Sprintf("evidence %s: the store holds %d bytes, the candidate needs %d bytes, and the quota is %d bytes",
			chargeevidence.RefuseCapacity, capacity.Observed, capacity.Candidate, capacity.Quota)
		if capacity.Required == 0 {
			return toon.Errorf(kind, "no representable quota admits the candidate; remove evidence explicitly") + "\n"
		}
		return toon.Errorf(kind, fmt.Sprintf("retry with %s %d", flagQuota, capacity.Required)) + "\n"
	}
	next, ok := evidenceRecovery[refusalClass(err)]
	if !ok {
		next = "repair the reported evidence condition and rerun the exact command"
	}
	return toon.Errorf("evidence "+err.Error(), next) + "\n"
}

func refusalClass(err error) string {
	var refusal *chargeevidence.Refusal
	if errors.As(err, &refusal) {
		return refusal.Class
	}
	return ""
}
