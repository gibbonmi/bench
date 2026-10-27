// Stack file reader for the stack landing: the strict file, its members, and the source-tip overrides.
package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/usage"
)

// The member bound of one stack. The ceiling keeps the worst bisect at five whole gates.
const (
	stackMemberFloor   = 2
	stackMemberCeiling = 16
)

// StackMember is one member of a stack file. The fields other than EffectiveTip hold the
// file values as written; Spec is empty when the file omits it. EffectiveTip is the
// override tip when an override names Request, and SourceTip otherwise.
type StackMember struct {
	Path, Request, Base, SourceTip, Spec, Message string
	EffectiveTip                                  string
}

// StackOverride is one parsed `<request>=<commit>` override value.
type StackOverride struct{ Request, Commit string }

// StackInput is a read stack: its members in file order and its overrides in argument
// order.
type StackInput struct {
	Members   []StackMember
	Overrides []StackOverride
}

// StackInputError is the one refusal of ReadStack. Member is the 1-based index of the
// refused member, or zero when the file or an override is at fault.
type StackInputError struct {
	Member int
	Detail string
}

func (e *StackInputError) Error() string {
	if e.Member == 0 {
		return e.Detail
	}
	return fmt.Sprintf("stack member %d: %s", e.Member, e.Detail)
}

// stackDocument is the stack file schema. Each member field is a pointer so that an
// absent field stays absent in the member's argument vector, where the landing grammar
// refuses it as it refuses an absent flag.
type stackDocument struct {
	Members []stackDocumentMember `json:"members"`
}

type stackDocumentMember struct {
	Path      *string `json:"path"`
	Request   *string `json:"request"`
	Base      *string `json:"base"`
	SourceTip *string `json:"source_tip"`
	Spec      *string `json:"spec"`
	Message   *string `json:"message"`
}

// landArgs is the argument vector of the one landing that the member stands for.
func (m stackDocumentMember) landArgs() []string {
	var args []string
	for _, flag := range []struct {
		name  string
		value *string
	}{
		{"--request", m.Request}, {"--base", m.Base}, {"--source-tip", m.SourceTip}, {"--spec", m.Spec}, {"-m", m.Message},
	} {
		if flag.value != nil {
			args = append(args, flag.name, *flag.value)
		}
	}
	if m.Path != nil {
		args = append(args, *m.Path)
	}
	return args
}

// ReadStack reads the stack file at path, which resolves against workingDir when it is
// relative, and applies the override values to its members. It reads no repository
// state. Every refusal is a *StackInputError.
func ReadStack(path, workingDir string, overrides []string) (StackInput, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(workingDir, path)
	}
	// The type check comes before any open: an open of a FIFO with no writer blocks.
	if !isRegularFile(path) {
		return StackInput{}, &StackInputError{Detail: fmt.Sprintf("stack file %q is absent or not a regular file", path)}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return StackInput{}, &StackInputError{Detail: fmt.Sprintf("read stack file: %v", err)}
	}
	var doc stackDocument
	if err := jsonfile.DecodeExactDocument(data, &doc); err != nil {
		return StackInput{}, &StackInputError{Detail: fmt.Sprintf("stack file %q: %v", path, err)}
	}
	if n := len(doc.Members); n < stackMemberFloor || n > stackMemberCeiling {
		return StackInput{}, &StackInputError{Detail: fmt.Sprintf("stack file %q holds %d members, want %d to %d", path, n, stackMemberFloor, stackMemberCeiling)}
	}
	members := make([]StackMember, len(doc.Members))
	for i, member := range doc.Members {
		parsed, line, _ := usage.Parse(landGrammar, member.landArgs())
		if line != "" {
			// A refusal for an empty or unknown value names that value but not the
			// grammar, so the grammar line rides beside it.
			if line != landGrammar.Help {
				line += "; " + landGrammar.Help
			}
			return StackInput{}, &StackInputError{Member: i + 1, Detail: line}
		}
		members[i] = StackMember{
			Path: parsed.Positionals[0], Request: parsed.Flags["--request"], Base: parsed.Flags["--base"],
			SourceTip: parsed.Flags["--source-tip"], Spec: parsed.Flags["--spec"], Message: parsed.Flags["-m"],
			EffectiveTip: parsed.Flags["--source-tip"],
		}
	}
	parsedOverrides, err := applyStackOverrides(members, overrides)
	if err != nil {
		return StackInput{}, err
	}
	return StackInput{Members: members, Overrides: parsedOverrides}, nil
}

// applyStackOverrides parses each `<request>=<commit>` value and sets the effective tip
// of each member whose request it names.
func applyStackOverrides(members []StackMember, values []string) ([]StackOverride, error) {
	known := make(map[string]bool, len(members))
	for _, member := range members {
		known[member.Request] = true
	}
	tips := make(map[string]string, len(values))
	parsed := make([]StackOverride, 0, len(values))
	for _, value := range values {
		request, commit, ok := strings.Cut(value, "=")
		if !ok || request == "" || commit == "" {
			return nil, &StackInputError{Detail: fmt.Sprintf("override %q is not <request>=<commit>", value)}
		}
		if !known[request] {
			return nil, &StackInputError{Detail: fmt.Sprintf("override %q names no stack member", value)}
		}
		if _, repeated := tips[request]; repeated {
			return nil, &StackInputError{Detail: fmt.Sprintf("override %q is a second override for request %q", value, request)}
		}
		tips[request] = commit
		parsed = append(parsed, StackOverride{Request: request, Commit: commit})
	}
	for i := range members {
		if tip, ok := tips[members[i].Request]; ok {
			members[i].EffectiveTip = tip
		}
	}
	return parsed, nil
}
