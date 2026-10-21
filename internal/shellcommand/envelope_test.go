package shellcommand

import "testing"

var envelopeForms = []struct {
	name string
	form EnvelopeForm
}{
	{"CommandStruct", CommandStruct},
	{"InputMap", InputMap},
	{"CwdInputMap", CwdInputMap},
}

// assertDecodeErrors grades that every form returns a structural decode error for each input.
func assertDecodeErrors(t *testing.T, inputs []string) {
	t.Helper()
	for _, f := range envelopeForms {
		for _, in := range inputs {
			if _, err := DecodeEnvelope([]byte(in), f.form); err == nil {
				t.Errorf("%s: DecodeEnvelope(%q) error = nil, want a decode error", f.name, in)
			}
		}
	}
}

func TestEnvelopeInvalidJSON(t *testing.T) {
	assertDecodeErrors(t, []string{"", "not json", `{"tool_input":`})
}

func TestEnvelopeRootShape(t *testing.T) {
	assertDecodeErrors(t, []string{`[]`, `7`, `"git status"`, `true`})
}

// inputValues is the projection of one decoded envelope that the merge cases compare.
type inputValues struct {
	command, cwd, inputCommand, inputPath string
}

func readInputValues(t *testing.T, in string, form EnvelopeForm) inputValues {
	t.Helper()
	e, err := DecodeEnvelope([]byte(in), form)
	if err != nil {
		t.Fatalf("DecodeEnvelope(%q) error = %v", in, err)
	}
	v := inputValues{command: e.Command(), cwd: e.Cwd()}
	if form == CommandStruct {
		return v
	}
	for name, dst := range map[string]*string{"command": &v.inputCommand, "file_path": &v.inputPath} {
		value, _, err := e.InputString(name)
		if err != nil {
			t.Fatalf("InputString(%q) for %q error = %v", name, in, err)
		}
		*dst = value
	}
	return v
}

// TestEnvelopeInputShape grades non-object tool_input errors and ordered repeated-object merging.
func TestEnvelopeInputShape(t *testing.T) {
	assertDecodeErrors(t, []string{`{"tool_input":[]}`, `{"tool_input":7}`, `{"tool_input":"git status"}`})
	// Every merge case keeps the same final members, so one expectation per form serves all three.
	want := map[EnvelopeForm]inputValues{
		CommandStruct: {command: "git reset --hard"},
		InputMap:      {inputCommand: "git reset --hard", inputPath: "README.md"},
		CwdInputMap:   {cwd: "/repo", inputCommand: "git reset --hard", inputPath: "README.md"},
	}
	for _, in := range []string{
		`{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"TOOL_INPUT":{"other":1},"cwd":"/repo"}`,
		`{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"tool_input":{},"cwd":"/repo"}`,
		`{"tool_input":{"command":"git status","file_path":"old.md"},"tool_input":{"command":"git reset --hard","file_path":"README.md"},"cwd":"/repo"}`,
	} {
		for _, f := range envelopeForms {
			if got := readInputValues(t, in, f.form); got != want[f.form] {
				t.Errorf("%s %q: got %+v, want %+v", f.name, in, got, want[f.form])
			}
		}
	}
}

// TestEnvelopeAbsentAndNullObjects grades that absent and null objects decode to zero values.
func TestEnvelopeAbsentAndNullObjects(t *testing.T) {
	for _, in := range []string{`{}`, `null`, `{"tool_input":null}`, `{"tool_input":{}}`} {
		for _, f := range envelopeForms {
			e, err := DecodeEnvelope([]byte(in), f.form)
			if err != nil {
				t.Fatalf("%s: DecodeEnvelope(%q) error = %v", f.name, in, err)
			}
			if e.Command() != "" || e.Cwd() != "" {
				t.Errorf("%s %q: Command %q, Cwd %q, want zero values", f.name, in, e.Command(), e.Cwd())
			}
			if value, present, err := e.InputString("command"); value != "" || present || err != nil {
				t.Errorf("%s %q: InputString = %q, %v, %v, want a missing key", f.name, in, value, present, err)
			}
		}
	}
}

// TestEnvelopeIndependentFields grades that an invalid cwd type leaves both command forms readable.
func TestEnvelopeIndependentFields(t *testing.T) {
	in := []byte(`{"cwd":7,"tool_input":{"command":"git status"}}`)
	e, err := DecodeEnvelope(in, CommandStruct)
	if err != nil || e.Command() != "git status" {
		t.Errorf("CommandStruct: Command %q, error %v, want git status and no error", e.Command(), err)
	}
	e, err = DecodeEnvelope(in, InputMap)
	if err != nil {
		t.Fatalf("InputMap: error = %v", err)
	}
	if value, present, err := e.InputString("command"); value != "git status" || !present || err != nil {
		t.Errorf("InputMap: InputString = %q, %v, %v, want git status", value, present, err)
	}
	if _, err := DecodeEnvelope(in, CwdInputMap); err == nil {
		t.Error("CwdInputMap: error = nil, want the retained cwd type error")
	}
}

// TestEnvelopeNullString grades that a null map string is present and empty, and a number is a type error.
func TestEnvelopeNullString(t *testing.T) {
	for _, form := range []EnvelopeForm{InputMap, CwdInputMap} {
		e, err := DecodeEnvelope([]byte(`{"tool_input":{"command":null,"file_path":7}}`), form)
		if err != nil {
			t.Fatalf("form %d: error = %v", form, err)
		}
		if value, present, err := e.InputString("command"); value != "" || !present || err != nil {
			t.Errorf("form %d: null InputString = %q, %v, %v, want empty, present, no error", form, value, present, err)
		}
		if _, present, err := e.InputString("file_path"); !present || err == nil {
			t.Errorf("form %d: numeric InputString present %v, error %v, want present with a type error", form, present, err)
		}
	}
}

func TestEnvelopeUnknownForm(t *testing.T) {
	if _, err := DecodeEnvelope([]byte(`{}`), CwdInputMap+1); err == nil {
		t.Error("DecodeEnvelope with an unknown form error = nil, want an error")
	}
}
