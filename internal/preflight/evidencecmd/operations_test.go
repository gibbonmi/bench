package evidencecmd

import (
	"strings"
	"testing"
)

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
