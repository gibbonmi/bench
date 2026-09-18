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
	chargeevidence.RefuseAbsent:    "run the build preparation that prints this identifier",
	chargeevidence.RefuseUnsafe:    "repair the evidence store by hand; Bench removed nothing",
	chargeevidence.RefuseExisting:  "inspect the corrupt published artifact; Bench replaced nothing",
	chargeevidence.RefuseReplaced:  "rerun the exact read after the store stops changing",
	chargeevidence.RefuseCursor:    "rerun the read with the next command a previous page printed",
	chargeevidence.RefuseSource:    "read the manifest and name one source identifier it declares",
	chargeevidence.RefuseBusy:      "rerun the exact command after the active reader or writer finishes",
	chargeevidence.RefuseStalePlan: freshPlan,
}

// evidenceStore opens the repository-common store for root. It creates nothing.
func evidenceStore(root string) (*chargeevidence.Store, string) {
	common, err := git.CommonDir(root)
	if err != nil {
		return nil, toon.Errorf("evidence store unavailable", err.Error()) + "\n"
	}
	return chargeevidence.OpenStore(common, chargeevidence.StoreOptions{Pause: chargeevidence.PauseFromEnvironment()}), ""
}

// SelectQuota validates op's evidence operands before any repository access: the quota
// operand, then a read's identifier and cursor. It returns the selected quota, or the usage
// line.
func SelectQuota(op Operation, identity string, flags map[string]string) (uint64, string) {
	quota := uint64(chargeevidence.DefaultQuota)
	if text, ok := flags[flagQuota]; ok {
		value, valid := chargeevidence.ParseDecimal(text)
		if !valid || value == 0 {
			return 0, toon.Usage(Grammar.Cmd, flagQuota+" needs a positive decimal byte count within the unsigned 64-bit range")
		}
		quota = value
	}
	// The operation registry owns which forms read evidence, so a newly registered evidence
	// form validates its operands without a second enumeration here.
	if op.Mode == modeEvidence {
		return quota, operandRefusal(identity, flags[flagSource], flags[flagCursor])
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
		Next: evidenceInvocation(pack.Identity(), "", ""),
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
// hostile operand only by length and digest. A source read pairs its cursor with the
// selected source, so a cursor from another stream refuses before the store opens.
func operandRefusal(identity, source, cursor string) string {
	if !chargeevidence.ValidIdentity(identity) {
		return toon.Usage(Grammar.Cmd, boundedOperand("invalid evidence identifier", identity))
	}
	if source != "" && !chargeevidence.ValidSourceID(source) {
		return toon.Usage(Grammar.Cmd, boundedOperand("invalid source identifier", source))
	}
	if cursor == "" {
		return ""
	}
	position, err := chargeevidence.ParseCursor(cursor, identity)
	if err != nil {
		return toon.Usage(Grammar.Cmd, boundedOperand("invalid cursor", cursor)+" "+refusalClass(err))
	}
	if source != "" && chargeevidence.SourceID(position.Ordinal) != source {
		return toon.Usage(Grammar.Cmd, boundedOperand("invalid cursor", cursor)+" names another source than "+source)
	}
	return ""
}

// OpenEvidence opens one published artifact for a read that needs its manifest. The caller
// closes the artifact; a refusal carries its own exit code.
func OpenEvidence(root, identity string) (*chargeevidence.Artifact, string, int) {
	store, refusal := evidenceStore(root)
	if refusal != "" {
		return nil, refusal, 1
	}
	artifact, err := store.Open(identity)
	if err != nil {
		return nil, storeRefusal(err), 1
	}
	return artifact, "", 0
}

// Read prints one bounded fragment at the position the cursor flag names. Without a source
// flag it reads the default stream; with one it reads only that declared source and ends
// after it. It keeps no reading state: the flags alone name the position.
func Read(root, identity string, flags map[string]string) (string, int) {
	artifact, refusal, code := OpenEvidence(root, identity)
	if refusal != "" {
		return refusal, code
	}
	defer artifact.Close()
	source := flags[flagSource]
	cursor, refusal := readCursor(artifact, identity, source, flags[flagCursor])
	if refusal != "" {
		return refusal, 1
	}
	read := artifact.Read
	if source != "" {
		read = artifact.ReadWithin
	}
	fragment, err := read(cursor)
	if err != nil {
		return storeRefusal(err), 1
	}
	next := ""
	if fragment.Next != nil {
		next = evidenceInvocation(identity, source, fragment.Next.String())
	}
	text, err := fragment.Encode(identity, next)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return text, 0
}

// readCursor resolves the requested position: the named cursor, the first page of the
// selected source, or the first fragment of the default stream.
func readCursor(artifact *chargeevidence.Artifact, identity, source, cursorText string) (chargeevidence.Cursor, string) {
	if cursorText != "" {
		cursor, err := chargeevidence.ParseCursor(cursorText, identity)
		if err != nil {
			return chargeevidence.Cursor{}, storeRefusal(err)
		}
		return cursor, ""
	}
	if source == "" {
		return artifact.First(), ""
	}
	ordinal, err := artifact.SourceOrdinal(source)
	if err != nil {
		return chargeevidence.Cursor{}, storeRefusal(err)
	}
	return chargeevidence.Cursor{Identity: identity, Ordinal: ordinal}, ""
}

// Verify reads every stored page and source digest of one artifact and prints the complete
// verification result. It certifies stored bytes, never consumer delivery.
func Verify(root, identity string) (string, int) {
	artifact, refusal, code := OpenEvidence(root, identity)
	if refusal != "" {
		return refusal, code
	}
	defer artifact.Close()
	verified, err := artifact.Verify()
	if err != nil {
		return storeRefusal(err), 1
	}
	text, err := verified.Encode()
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return text, 0
}

// evidenceInvocation is the exact read command for one artifact position. An empty cursor
// is the manifest-first command, and a source read repeats its source flag.
func evidenceInvocation(identity, source, cursor string) string {
	args := []string{"bench", "preflight", modeEvidence, identity}
	if source != "" {
		args = append(args, flagSource, source)
	}
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
		// Cleanup is the first recovery action now that its producer exists. The larger
		// quota stays the documented alternative, and it is absent only when no
		// representable quota admits the candidate.
		if capacity.Required == 0 {
			return toon.Errorf(kind, "run "+CleanCommand+"; no representable quota admits the candidate") + "\n"
		}
		return toon.Errorf(kind, fmt.Sprintf("run %s, or retry with %s %d", CleanCommand, flagQuota, capacity.Required)) + "\n"
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
