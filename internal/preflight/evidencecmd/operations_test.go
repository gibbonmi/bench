package evidencecmd

import (
	"strings"
	"testing"
)

// BoundedForm is one registered bounded form's invocation shape: the name its mode and
// selectors give it, the mode it names, the operand placeholder that mode declares, the
// flags an invocation must state, and the flags it also accepts.
type BoundedForm struct {
	Name, Mode, Operand string
	Flags               []string
	Optional            []string
}

// BoundedForms projects every registered form that declares the shared response bound, in
// registry order. The command tests derive their bounded case list from this projection,
// so a bounded form registered later arrives with its own behavior case.
func BoundedForms() []BoundedForm {
	var forms []BoundedForm
	for _, op := range operations {
		if !op.Bounded {
			continue
		}
		forms = append(forms, BoundedForm{
			Name:     strings.Join(append([]string{op.Mode}, op.selectors...), " "),
			Mode:     op.Mode,
			Operand:  modeOperands[op.Mode],
			Flags:    append(append([]string{}, op.selectors...), op.required...),
			Optional: append([]string{}, op.optional...),
		})
	}
	return forms
}

// TestOperationFlagsRegistered proves that every flag an operation names comes from the flag
// table with the kind its position needs, and that every mode names its operand. An
// unregistered flag fails here, a required or optional flag that takes no value fails here,
// and a valued selector that hides its operand from the usage line fails here.
func TestOperationFlagsRegistered(t *testing.T) {
	placeholders := map[string]string{}
	for _, flag := range flagTable {
		placeholders[flag.name] = flag.placeholder
	}
	for _, op := range operations {
		if modeOperands[op.Mode] == "" {
			t.Errorf("mode %q has no registered operand", op.Mode)
		}
		for _, name := range op.selectors {
			placeholder, ok := placeholders[name]
			if !ok {
				t.Errorf("%s: selector %s is not a registered flag", op.usageLine(), name)
				continue
			}
			if placeholder != "" && !strings.Contains(op.usageLine(), name+" "+placeholder) {
				t.Errorf("%s: selector %s does not advertise its operand", op.usageLine(), name)
			}
		}
		for _, name := range append(append([]string{}, op.required...), op.optional...) {
			if placeholders[name] == "" {
				t.Errorf("%s: flag %s is not a registered valued flag", op.usageLine(), name)
			}
		}
	}
}
