package bench

import "testing"

// TestValidatePayloadRows pins the fail-closed shapes: the allowlist is the one source
// of what every destination may write, so a row that could escape the destination root
// or that names a source twice is rejected before any reader acts on it.
func TestValidatePayloadRows(t *testing.T) {
	consumer := func(source string) PayloadRow {
		return PayloadRow{Source: source, Mode: "0644", Audience: PayloadAudienceConsumer}
	}
	cases := []struct {
		name    string
		rows    []PayloadRow
		wantErr bool
	}{
		{name: "nested and dotted sources are ordinary", rows: []PayloadRow{consumer(".bench/gate.sh"), consumer("bin/bench.sh")}},
		{name: "traversal escapes the destination root", rows: []PayloadRow{consumer(".bench/../../x")}, wantErr: true},
		{name: "leading traversal escapes too", rows: []PayloadRow{consumer("../x")}, wantErr: true},
		{name: "absolute path ignores the destination root", rows: []PayloadRow{consumer("/etc/passwd")}, wantErr: true},
		{name: "backslash separator is not a kit-relative path", rows: []PayloadRow{consumer(`bin\bench.sh`)}, wantErr: true},
		{name: "current-directory segment is not a real source", rows: []PayloadRow{consumer("bin/./bench.sh")}, wantErr: true},
		{name: "empty segment is not a real source", rows: []PayloadRow{consumer("bin//bench.sh")}, wantErr: true},
		{name: "duplicate source has no defined winner", rows: []PayloadRow{consumer("bin/bench.sh"), consumer("bin/bench.sh")}, wantErr: true},
		{
			name:    "duplicate source disagreeing on audience has no defined winner",
			rows:    []PayloadRow{consumer("bin/bench.sh"), {Source: "bin/bench.sh", Mode: "0644", Audience: PayloadAudienceKitOnly}},
			wantErr: true,
		},
		{name: "unknown audience is invalid", rows: []PayloadRow{{Source: "bin/bench.sh", Mode: "0644", Audience: "everyone"}}, wantErr: true},
		{name: "empty source is invalid", rows: []PayloadRow{consumer("")}, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePayloadRows(tc.rows)
			if tc.wantErr && err == nil {
				t.Fatalf("validatePayloadRows(%+v) accepted a row it must refuse", tc.rows)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("validatePayloadRows(%+v) refused a valid allowlist: %v", tc.rows, err)
			}
		})
	}
}

// TestPayloadRowsAcceptsTheShippedAllowlist keeps the validator honest against the real
// tracked file: a rule tightened past what the kit itself ships would otherwise turn
// every link and pack red only once a consumer ran it.
func TestPayloadRowsAcceptsTheShippedAllowlist(t *testing.T) {
	rows, err := PayloadRows()
	if err != nil {
		t.Fatalf("the shipped consumer-payload allowlist does not satisfy its own validator: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("the shipped consumer-payload allowlist parsed to no rows")
	}
}

// TestPayloadRowsFromRejectsInvalidBytes pins the one parser every consumer reaches:
// syntax and row semantics are refused together, so a reader that decoded the JSON
// itself and skipped the row predicates cannot admit a row the allowlist forbids.
// Empty bytes are a present, unusable allowlist, not an absent one.
func TestPayloadRowsFromRejectsInvalidBytes(t *testing.T) {
	cases := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{name: "the shipped shape parses", data: `[{"source":"bin/bench.sh","mode":"0755","audience":"consumer"}]`},
		{name: "empty bytes are not an allowlist", data: "", wantErr: true},
		{name: "truncated JSON", data: `[{"source":"bin/bench.sh"`, wantErr: true},
		{name: "not an array", data: `{"source":"bin/bench.sh"}`, wantErr: true},
		{name: "unknown audience", data: `[{"source":"bin/bench.sh","audience":"everyone"}]`, wantErr: true},
		{name: "empty source", data: `[{"source":"","audience":"consumer"}]`, wantErr: true},
		{name: "unsafe source", data: `[{"source":"../x","audience":"consumer"}]`, wantErr: true},
		{name: "duplicate source", data: `[{"source":"a","audience":"consumer"},{"source":"a","audience":"kit-only"}]`, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := PayloadRowsFrom([]byte(tc.data))
			if tc.wantErr && err == nil {
				t.Fatalf("PayloadRowsFrom(%q) accepted bytes it must refuse", tc.data)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("PayloadRowsFrom(%q) refused a valid allowlist: %v", tc.data, err)
			}
		})
	}
}
