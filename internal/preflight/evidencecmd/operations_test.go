package evidencecmd

import "testing"

// TestOperationFlagsRegistered proves that every flag an operation names comes from the flag
// table with the kind its position needs, and that every mode names its operand. A selector
// the table does not list as a switch, or a valued flag listed as one, fails here.
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
			if placeholder, ok := placeholders[name]; !ok || placeholder != "" {
				t.Errorf("%s: selector %s is not a registered switch", op.usageLine(), name)
			}
		}
		for _, name := range append(append([]string{}, op.required...), op.optional...) {
			if placeholders[name] == "" {
				t.Errorf("%s: flag %s is not a registered valued flag", op.usageLine(), name)
			}
		}
	}
}
