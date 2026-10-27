package gate

// The background width: one Bench-owned input that narrows each Go test phase of the gate.

import "os"

// widthEnv selects the background width for the one gate-phases process that a background
// grade launches. Any non-empty value selects it. It is not a member of the grade key, so a
// green at the background width answers a run of the same subject at another width.
const widthEnv = "BENCH_GATE_WIDTH"

// backgroundWidth is the -p and -parallel value of each Go test phase under the width input.
const backgroundWidth = "2"

// widthFlags answers the width flags of a Go test phase: -p and -parallel at the background
// width when this process carries the width input, and no flag otherwise, so a run without
// the input keeps its argv.
func widthFlags() []string {
	if os.Getenv(widthEnv) == "" {
		return nil
	}
	return []string{"-p=" + backgroundWidth, "-parallel=" + backgroundWidth}
}

// gateTestArgv is the argv of a Go test phase of the gate over args: BaseTestArgv with the
// width flags before args.
func gateTestArgv(args ...string) []string {
	return BaseTestArgv("", append(widthFlags(), args...)...)
}
