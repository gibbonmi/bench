package conformance

import (
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"reflect"
	"runtime"
	"strings"
)

// checkBinding is the executable half of a registry row. The executable map repeats only the facts an
// independent mutation oracle needs: the name-to-function binding, tier, and subject. Registry
// order, meta membership, and inputs stay single-sourced in registry.Checks. Keeping them
// independent makes the named CM5/CM6/CM7 mutations red when a function swaps, or the advertised
// tier or subject drifts while the binding is unchanged.
type checkBinding struct {
	implementation any
	tier           registry.Tier
	subject        registry.Subject
}

func (b checkBinding) identity() string {
	fn := runtime.FuncForPC(reflect.ValueOf(b.implementation).Pointer())
	if fn == nil {
		return ""
	}
	name := fn.Name()
	return name[strings.LastIndex(name, ".")+1:]
}

func (b checkBinding) runsAt(tier registry.Tier) bool {
	return b.tier == registry.Dev || tier == registry.Ship
}

func (b checkBinding) run(root, kitRoot string, tier registry.Tier) []string {
	subject := root
	if b.subject == registry.SubjectKitRoot {
		subject = kitRoot
	}
	switch run := b.implementation.(type) {
	case func(string) []string:
		return run(subject)
	case func(string, string) []string:
		return run(root, kitRoot)
	case func(string, registry.Tier) []string:
		return run(subject, tier)
	default:
		return []string{"conformance check carries an unsupported executable binding"}
	}
}
