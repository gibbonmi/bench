// Package bench embeds the consumer-payload allowlist, the single canonical source of
// what a linked repo, the npm wrapper tarball, and the release-evidence pipeline are
// permitted to ship. The tracked bytes live at .bench/consumer-payload.json — the one
// copy — so shell and Node readers reach it without a Go build. go:embed cannot match a
// pattern against a dot-prefixed directory's contents, but it can embed a file named
// explicitly, so this root-level package (the module root carries no other .go files)
// names the path directly and gives internal/adopt the allowlist compiled into the
// binary rather than resolved against whatever kit directory happens to be on disk.
package bench

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

//go:embed .bench/consumer-payload.json
var consumerPayloadJSON []byte

// PayloadRow is one allowlist entry: a source path relative to the kit root, its
// shipped file mode, whether it names a directory to walk, and who receives it.
type PayloadRow struct {
	Source   string `json:"source"`
	Mode     string `json:"mode"`
	Tree     bool   `json:"tree"`
	Audience string `json:"audience"`
}

const (
	// PayloadAudienceConsumer marks a row every linked repo and packaged artifact receives.
	PayloadAudienceConsumer = "consumer"
	// PayloadAudienceKitOnly marks a row withheld from every destination: linked repos,
	// the npm wrapper tarball, and the release-evidence pipeline.
	PayloadAudienceKitOnly = "kit-only"
)

// PayloadRows parses the embedded allowlist. buildLinkPlan is the Go caller; the Node
// release-evidence builder and the package-shipped-surface conformance suite read the
// same tracked bytes directly (a JSON file, not a Go call) rather than hand-listing the
// payload a second time.
func PayloadRows() ([]PayloadRow, error) {
	return PayloadRowsFrom(consumerPayloadJSON)
}

// PayloadRowsFrom is the one parser for allowlist bytes, wherever they came from:
// decoding and row validation are a single step so that no reader can decode the JSON
// and then act on rows the allowlist forbids. Empty bytes are a present, unusable
// allowlist and are refused here; only a caller that never reaches this parser — an
// absent file — may treat the allowlist as optional.
func PayloadRowsFrom(data []byte) ([]PayloadRow, error) {
	var rows []PayloadRow
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("consumer payload allowlist is invalid: %w", err)
	}
	if err := validatePayloadRows(rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// PayloadRowsAt reads the allowlist from the filesystem for every consumer that reads
// the tracked file rather than the embedded copy. The path is attacker-shaped input to
// checks that gate what ships, so it is classified before it is opened — a link is
// refused rather than followed and a FIFO cannot block a check in open(2) — and the
// classified bytes then go through the same parser the embedded copy uses. absent is
// reported separately because only one caller (the skills index, which withholds
// nothing when there is no allowlist) is allowed to continue on it.
func PayloadRowsAt(path string) (rows []PayloadRow, absent bool, err error) {
	classified := bounds.ClassifyNoFollow(path)
	switch classified.State {
	case bounds.StateAbsent:
		return nil, true, nil
	case bounds.StateParsed:
	default:
		reason := classified.Reason
		if reason == "" {
			reason = string(classified.State)
		}
		return nil, false, fmt.Errorf("%s unreadable: %s", path, reason)
	}
	rows, err = PayloadRowsFrom(classified.Data)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", path, err)
	}
	return rows, false, nil
}

// validatePayloadRows fails closed on the row shapes the readers downstream cannot
// resolve safely. Every destination joins Source onto a root it owns — a linked repo, a
// staged tarball, an evidence bundle — so an absolute path, a backslash separator, or a
// ".." segment would write outside the tree the caller consented to. A source named
// twice is rejected on the same footing: the two rows can disagree on mode or audience,
// and which one wins would then depend on read order rather than on the allowlist.
func validatePayloadRows(rows []PayloadRow) error {
	seen := make(map[string]bool, len(rows))
	for _, r := range rows {
		if r.Source == "" || r.Audience != PayloadAudienceConsumer && r.Audience != PayloadAudienceKitOnly {
			return fmt.Errorf("consumer payload allowlist row is invalid: %+v", r)
		}
		if !payloadSourceSafe(r.Source) {
			return fmt.Errorf("consumer payload allowlist row names an unsafe source path: %q", r.Source)
		}
		if seen[r.Source] {
			return fmt.Errorf("consumer payload allowlist names the source %q twice", r.Source)
		}
		seen[r.Source] = true
	}
	return nil
}

// payloadSourceSafe reports whether source is a kit-relative forward-slash path that
// stays inside the kit root.
func payloadSourceSafe(source string) bool {
	if strings.HasPrefix(source, "/") || strings.ContainsRune(source, '\\') {
		return false
	}
	for _, segment := range strings.Split(source, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return false
		}
	}
	return true
}

// PayloadConsumerRows returns only the rows every destination is allowed to ship.
func PayloadConsumerRows(rows []PayloadRow) []PayloadRow {
	var out []PayloadRow
	for _, r := range rows {
		if r.Audience == PayloadAudienceConsumer {
			out = append(out, r)
		}
	}
	return out
}

// PayloadKitOnlyPrefixes returns each kit-only row's source path, for callers that need
// to exclude a matching file (exact match) or tree (prefix match) while walking a
// consumer-audience directory that contains both.
func PayloadKitOnlyPrefixes(rows []PayloadRow) []string {
	var out []string
	for _, r := range rows {
		if r.Audience == PayloadAudienceKitOnly {
			out = append(out, r.Source)
		}
	}
	return out
}

// PayloadExcluded reports whether sourcePath (a kit-relative path, forward-slash
// separated) falls under one of the kit-only prefixes: an exact match (a withheld
// file) or a path inside a withheld tree.
func PayloadExcluded(sourcePath string, kitOnlyPrefixes []string) bool {
	for _, ex := range kitOnlyPrefixes {
		if sourcePath == ex || len(sourcePath) > len(ex) && sourcePath[:len(ex)+1] == ex+"/" {
			return true
		}
	}
	return false
}
