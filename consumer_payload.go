// Package bench embeds the consumer-payload allowlist. This allowlist is the single
// canonical source for what a linked repo, the npm wrapper tarball, and the
// release-evidence pipeline may ship. The tracked bytes live at
// .bench/consumer-payload.json, the one copy, so shell and Node readers reach it
// without a Go build. go:embed cannot match a pattern against a dot-prefixed
// directory's contents. But it can embed a file named explicitly.
//
// This root-level package, the only one at the module root, names the path
// directly. It gives internal/adopt the allowlist compiled into the binary,
// instead of resolved against whatever kit directory is on disk.
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

// PayloadRow is one allowlist entry. It names a source path relative to the kit
// root, its shipped file mode, whether it names a directory to walk, and its
// audience.
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

// PayloadRows parses the embedded allowlist. buildLinkPlan is the Go caller. The Node
// release-evidence builder and the package-shipped-surface conformance suite read the
// same tracked bytes directly, as a JSON file rather than a Go call. This direct read
// avoids hand-listing the payload a second time.
func PayloadRows() ([]PayloadRow, error) {
	return PayloadRowsFrom(consumerPayloadJSON)
}

// PayloadRowsFrom is the one parser for allowlist bytes, from any source. Decoding
// and row validation happen as a single step, so no reader can decode the JSON and
// then act on rows the allowlist forbids. Empty bytes count as a present, unusable
// allowlist, and this parser refuses them. Only a caller that never reaches this
// parser, for an absent file, may treat the allowlist as optional.
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

// PayloadRowsAt reads the allowlist from the filesystem, for every consumer that
// reads the tracked file instead of the embedded copy. The path is attacker-shaped
// input to checks that gate what ships, so PayloadRowsAt classifies it before
// opening it. A link is refused rather than followed. A FIFO cannot block a check
// inside open(2). The classified bytes then go through the same parser the embedded
// copy uses. absent reports separately, because only one caller may continue on it.
//
// The skills index is that one caller. It withholds nothing when there is no
// allowlist.
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

// validatePayloadRows fails closed on row shapes that downstream readers cannot
// resolve safely. Every destination joins Source onto a root it owns, for example a
// linked repo, a staged tarball, or an evidence bundle. So an absolute path, a
// backslash separator, or a ".." segment would write outside the tree the caller
// consented to. validatePayloadRows also rejects a source named twice. Two rows
// naming the same source can disagree on mode or audience. Read order, not the
// allowlist, would then decide which one wins.
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

// PayloadKitOnlyPrefixes returns each kit-only row's source path. A caller walking a
// consumer-audience directory uses these paths to exclude a matching file or tree. A
// file excludes by exact match; a tree excludes by prefix match.
func PayloadKitOnlyPrefixes(rows []PayloadRow) []string {
	var out []string
	for _, r := range rows {
		if r.Audience == PayloadAudienceKitOnly {
			out = append(out, r.Source)
		}
	}
	return out
}

// PayloadExcluded reports whether sourcePath, a kit-relative forward-slash path,
// falls under one of the kit-only prefixes. It matches an exact withheld file or a
// path inside a withheld tree.
func PayloadExcluded(sourcePath string, kitOnlyPrefixes []string) bool {
	for _, ex := range kitOnlyPrefixes {
		if sourcePath == ex || len(sourcePath) > len(ex) && sourcePath[:len(ex)+1] == ex+"/" {
			return true
		}
	}
	return false
}
