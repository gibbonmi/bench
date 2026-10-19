package ownership

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"unicode"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/tickets"
)

// reusedTestSubject is the decoded subject of one reused-test premise.
type reusedTestSubject struct {
	Path           string `json:"path"`
	Function       string `json:"function"`
	RecordSpec     string `json:"record_spec"`
	SourceCommit   string `json:"source_commit"`
	EvidenceCommit string `json:"evidence_commit"`
	Chunk          string `json:"chunk"`
	Verification   string `json:"verification"`
	Mutation       string `json:"mutation"`
}

// reusedTestExpected is the decoded expected value of one reused-test premise.
type reusedTestExpected struct {
	RedNeedle string `json:"red_needle"`
}

// reusedTest grades one reused-test premise against its committed test declaration and its
// source-bound review record. A premise turns green only when the test exists at the
// recorded source commit and current source tip, the record passes CheckTrees for the
// named chunk at that source commit, and the selected verification carries a restored,
// biting behavioral probe of the named test subject with red_needle in its native excerpt.
func (c *census) reusedTest(item PremiseItem) (string, error) {
	subject, expected, err := decodeReusedTest(item)
	if err != nil {
		return "", err
	}
	if err := c.checkTestSubject(subject); err != nil {
		return "", err
	}
	if err := reviewrecord.CheckTrees(c.root, subject.SourceCommit, subject.EvidenceCommit, subject.SourceCommit, subject.RecordSpec, subject.Chunk, false); err != nil {
		return "", fmt.Errorf("record %s chunk %s at %s failed check for source_commit %s: %w", subject.RecordSpec, subject.Chunk, subject.EvidenceCommit, subject.SourceCommit, err)
	}
	record, err := reviewrecord.ReadTree(c.root, subject.EvidenceCommit, subject.RecordSpec)
	if err != nil {
		return "", fmt.Errorf("record %s does not read at %s: %w", subject.RecordSpec, subject.EvidenceCommit, err)
	}
	sourceDigest, err := reviewrecord.SourceDigest(c.root, subject.SourceCommit, subject.RecordSpec)
	if err != nil {
		return "", fmt.Errorf("source digest of %s at %s does not read: %w", subject.RecordSpec, subject.SourceCommit, err)
	}
	chunk, ok := selectChunk(record, subject.Chunk)
	if !ok {
		return "", fmt.Errorf("record %s has no chunk %s", subject.RecordSpec, subject.Chunk)
	}
	if chunk.Tip != subject.SourceCommit || chunk.SourceDigest != sourceDigest {
		return "", fmt.Errorf("record %s chunk %s binds to %s (%s), not %s (%s)", subject.RecordSpec, subject.Chunk, chunk.Tip, chunk.SourceDigest, subject.SourceCommit, sourceDigest)
	}
	for _, confirmation := range chunk.Confirmations {
		if confirmation.ID == subject.Verification {
			return "", fmt.Errorf("%s in record %s chunk %s is a coordinator confirmation, not a verification", subject.Verification, subject.RecordSpec, subject.Chunk)
		}
	}
	verification, ok := selectVerification(chunk.Verification, subject.Verification)
	if !ok {
		return "", fmt.Errorf("record %s chunk %s has no active verification %s", subject.RecordSpec, subject.Chunk, subject.Verification)
	}
	if err := checkVerificationProbe(verification, subject, expected, sourceDigest); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s (%s) is certified by %s chunk %s verification %s at %s", subject.Path, subject.Function, subject.RecordSpec, subject.Chunk, subject.Verification, subject.EvidenceCommit), nil
}

// decodeReusedTest decodes the subject and expected objects of a reused-test premise and
// refuses any empty or whitespace-only string cell.
func decodeReusedTest(item PremiseItem) (reusedTestSubject, reusedTestExpected, error) {
	var expected reusedTestExpected
	subject, err := testSubject(item.Subject, "subject: ")
	var expectedFields map[string]json.RawMessage
	if err == nil {
		expectedFields, err = typed(item.Expected, "expected: ", &expected, "red_needle")
	}
	if err == nil {
		err = requireNonEmptyStrings(expectedFields, "expected: ", "red_needle")
	}
	return subject, expected, err
}

// testSubject decodes one reused-test subject, which a posture premise's writer probe also
// takes, and refuses any empty or whitespace-only string cell.
func testSubject(raw json.RawMessage, where string) (reusedTestSubject, error) {
	var subject reusedTestSubject
	object, err := typed(raw, where, &subject, reusedTestKeys...)
	if err == nil {
		err = requireNonEmptyStrings(object, where, reusedTestKeys...)
	}
	return subject, err
}

func requireNonEmptyStrings(object map[string]json.RawMessage, where string, keys ...string) error {
	for _, key := range keys {
		var value *string
		if err := decodeInto(object, where, key, &value); err != nil {
			return err
		}
		if value == nil || strings.TrimSpace(*value) == "" {
			return fmt.Errorf("%s%s is empty", where, key)
		}
	}
	return nil
}

// checkTestSubject verifies that source_commit and evidence_commit resolve in the
// repository and that the named _test.go file declares the top-level test function at
// source_commit without drifting at the current source tip.
func (c *census) checkTestSubject(subject reusedTestSubject) error {
	if !tickets.RepositoryPath(subject.Path) || !strings.HasSuffix(subject.Path, "_test.go") {
		return fmt.Errorf("subject: path %q is no clean repository _test.go path", subject.Path)
	}
	if !strings.HasPrefix(subject.Function, "Test") {
		return fmt.Errorf("subject: function %q is no test function", subject.Function)
	}
	if err := c.exactCommit("source_commit", subject.SourceCommit); err != nil {
		return err
	}
	if err := c.exactCommit("evidence_commit", subject.EvidenceCommit); err != nil {
		return err
	}
	sourceBytes, err := git.ReadTreeFile(c.root, subject.SourceCommit, subject.Path)
	if err != nil {
		return fmt.Errorf("%s does not read at source_commit %s: %w", subject.Path, subject.SourceCommit, err)
	}
	if !declaresTestFunction(subject.Path, sourceBytes, subject.Function) {
		return fmt.Errorf("%s declares no test function %s at source_commit %s", subject.Path, subject.Function, subject.SourceCommit)
	}
	if c.source.Tip != "" {
		tipBytes, err := git.ReadTreeFile(c.root, c.source.Tip, subject.Path)
		if err != nil {
			return fmt.Errorf("%s does not read at the source tip %s: %w", subject.Path, c.source.Tip, err)
		}
		if !bytes.Equal(sourceBytes, tipBytes) {
			return fmt.Errorf("%s at the source tip %s differs from source_commit %s", subject.Path, c.source.Tip, subject.SourceCommit)
		}
	}
	return nil
}

func (c *census) exactCommit(label, commit string) error {
	if c.root == "" {
		return fmt.Errorf("no repository root resolves %s %q", label, commit)
	}
	resolved, err := git.Output("-C", c.root, "rev-parse", "--verify", "--quiet", commit+"^{commit}")
	if err != nil || resolved != commit {
		return fmt.Errorf("%s %q does not resolve to an exact commit", label, commit)
	}
	return nil
}

func declaresTestFunction(path string, src []byte, name string) bool {
	file, err := parser.ParseFile(token.NewFileSet(), path, src, 0)
	if err != nil {
		return false
	}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Recv == nil && fn.Name != nil && fn.Name.Name == name {
			return true
		}
	}
	return false
}

func selectChunk(record reviewrecord.Record, id string) (reviewrecord.Chunk, bool) {
	for _, chunk := range record.Chunks {
		if chunk.ID == id {
			return chunk, true
		}
	}
	return reviewrecord.Chunk{}, false
}

// selectVerification answers the active verification entry matching name by requirement or
// occurrence ID. Superseded occurrences for the same requirement are ignored so stale
// probe attempts cannot satisfy a premise.
func selectVerification(items []reviewrecord.Verification, name string) (reviewrecord.Verification, bool) {
	latest := map[string]int{}
	for i, item := range items {
		latest[item.Requirement] = i
	}
	for i, item := range items {
		if (item.Requirement == name || item.ID == name) && latest[item.Requirement] == i {
			return item, true
		}
	}
	return reviewrecord.Verification{}, false
}

// checkVerificationProbe grades the selected verification and its probe excerpt against
// the expected mutation, restoration, subject identity, and behavioral failure needle.
func checkVerificationProbe(verification reviewrecord.Verification, subject reusedTestSubject, expected reusedTestExpected, sourceDigest string) error {
	if verification.SourceDigest != sourceDigest {
		return fmt.Errorf("verification %s source digest %s does not match %s", subject.Verification, verification.SourceDigest, sourceDigest)
	}
	if verification.State != "completed" || verification.Outcome != "pass" || verification.ExitCode == nil || *verification.ExitCode != 0 {
		return fmt.Errorf("verification %s is not a completed passing run", subject.Verification)
	}
	if verification.Probe == nil {
		return fmt.Errorf("verification %s carries no recorded probe", subject.Verification)
	}
	probe := *verification.Probe
	if probe.Mutation != subject.Mutation {
		return fmt.Errorf("verification %s probe mutation is %q, not %q", subject.Verification, probe.Mutation, subject.Mutation)
	}
	bit, restored := reviewrecord.ProbeBites(probe)
	if !bit {
		return fmt.Errorf("verification %s probe did not bite: outcome %q, exit code %d", subject.Verification, probe.Outcome, probe.ExitCode)
	}
	if !restored {
		return fmt.Errorf("verification %s probe restoration is %q, not \"pass\"", subject.Verification, probe.Restore)
	}
	excerpt := probe.NativeRef.Excerpt
	if excerpt == "" || reviewrecord.Digest([]byte(excerpt)) != probe.NativeRef.Digest {
		return fmt.Errorf("verification %s probe native excerpt is missing or corrupt", subject.Verification)
	}
	if !witnessesSubject(excerpt, subject.Function) {
		return fmt.Errorf("verification %s probe excerpt does not witness test subject %s", subject.Verification, subject.Function)
	}
	if !strings.Contains(excerpt, expected.RedNeedle) {
		return fmt.Errorf("verification %s probe excerpt does not hold red_needle %q", subject.Verification, expected.RedNeedle)
	}
	return nil
}

// witnessesSubject reports whether excerpt names function as the failing test subject and
// names no other Test identifier.
func witnessesSubject(excerpt, function string) bool {
	names := testIdentifiers(excerpt)
	if len(names) == 0 {
		return false
	}
	for _, name := range names {
		if name != function {
			return false
		}
	}
	return true
}

func testIdentifiers(text string) []string {
	var names []string
	for i := 0; i < len(text); i++ {
		if !strings.HasPrefix(text[i:], "Test") {
			continue
		}
		if i > 0 && isIdentByte(text[i-1]) {
			continue
		}
		end := i + len("Test")
		for end < len(text) && isIdentByte(text[end]) {
			end++
		}
		if end > i+len("Test") {
			names = append(names, text[i:end])
		}
		i = end - 1
	}
	return names
}

func isIdentByte(b byte) bool {
	return b == '_' || unicode.IsLetter(rune(b)) || unicode.IsDigit(rune(b))
}
