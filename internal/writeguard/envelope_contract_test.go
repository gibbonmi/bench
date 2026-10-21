package writeguard

import (
	"encoding/json"
	"errors"
	"testing"
)

// TestWriteEnvelopeCwdContract grades the outer cwd decode before path handling and map-null clearing.
// A wrong cwd type stays a decode error even for an absolute path or a later valid cwd.
func TestWriteEnvelopeCwdContract(t *testing.T) {
	for _, tc := range []struct{ name, in string }{
		{"numeric cwd with an absolute path", `{"cwd":7,"tool_input":{"file_path":"/repo/README.md"}}`},
		{"array cwd with an absolute path", `{"cwd":["/repo"],"tool_input":{"file_path":"/repo/README.md"}}`},
		{"object cwd with an absolute path", `{"cwd":{},"tool_input":{"file_path":"/repo/README.md"}}`},
		{"boolean cwd with an absolute path", `{"cwd":true,"tool_input":{"file_path":"/repo/README.md"}}`},
		{"earlier cwd type error", `{"cwd":7,"cwd":"/repo","tool_input":{"file_path":"README.md"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path, err := PathFromEnvelope([]byte(tc.in))
			var typeErr *json.UnmarshalTypeError
			if !errors.As(err, &typeErr) || typeErr.Field != "cwd" {
				t.Errorf("PathFromEnvelope(%s) = (%q, %v), want a cwd decode error", tc.in, path, err)
			}
		})
	}
	t.Run("duplicate input null", func(t *testing.T) {
		in := `{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"tool_input":null,"cwd":"/repo"}`
		// The missing-field text is written here independently, so a respelled production error turns this row red.
		if path, err := PathFromEnvelope([]byte(in)); err == nil || err.Error() != "file_path field missing" {
			t.Errorf("PathFromEnvelope(%s) = (%q, %v), want the file_path field missing error", in, path, err)
		}
	})
}
