package evidencecmd

import (
	"slices"
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

// formName is the name one registered form carries: its mode and the selectors that choose
// it. The registry gives no two forms the same pair, so the name is the form's identity.
func (op Operation) formName() string {
	return strings.Join(append([]string{op.Mode}, op.selectors...), " ")
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
			Name:     op.formName(),
			Mode:     op.Mode,
			Operand:  modeOperands[op.Mode],
			Flags:    append(append([]string{}, op.selectors...), op.required...),
			Optional: append([]string{}, op.optional...),
		})
	}
	return forms
}

// TestBoundedFormsProjectEveryBoundedOperation pins the projection to the registry it
// projects. The case list of the response bound test derives from BoundedForms, so a form
// the projection drops takes its behavior case with it and no other check moves. The names
// are compared in registry order against the registry's own bounded rows, so a drop, an
// addition, and a reorder all red here, and no count stands in for the names.
func TestBoundedFormsProjectEveryBoundedOperation(t *testing.T) {
	var want []string
	for _, op := range operations {
		if op.Bounded {
			want = append(want, op.formName())
		}
	}
	var got []string
	for _, form := range BoundedForms() {
		got = append(got, form.Name)
	}
	if !slices.Equal(got, want) {
		t.Errorf("BoundedForms names = %v, want the registry's bounded rows %v", got, want)
	}
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
		// A mode is registered with the operand it takes, and an operand-less mode is
		// registered with an empty one. An unregistered mode would take an operand no
		// usage line advertises.
		if _, ok := modeOperands[op.Mode]; !ok {
			t.Errorf("mode %q is not registered in the operand table", op.Mode)
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
