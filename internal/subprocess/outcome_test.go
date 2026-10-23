package subprocess_test

import (
	"testing"

	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestUncertifiedRefusesEachUncertainPart is story 6 and the rule that a raw child status
// cannot certify required cleanup or publication. Each outcome carries a zero exit and
// exactly one uncertain part, and only the all-certain outcome reports success.
func TestUncertifiedRefusesEachUncertainPart(t *testing.T) {
	certain := subprocess.Outcome{
		Process:    subprocess.Process{Started: true, Waited: true},
		Streams:    subprocess.StreamSet{Stdin: subprocess.Stream{Status: subprocess.StreamNone}, Stdout: subprocess.Stream{Status: subprocess.StreamComplete}, Stderr: subprocess.Stream{Status: subprocess.StreamNone}},
		Cleanup:    subprocess.Cleanup{Status: subprocess.CleanupComplete},
		Protection: protection.Disposition{State: protection.Safe},
	}
	if err := certain.Uncertified(); err != nil {
		t.Fatalf("a certain zero exit = %v, want success", err)
	}
	cases := map[string]func(*subprocess.Outcome){
		"obligation": func(o *subprocess.Outcome) {
			o.Cleanup.Obligations = []subprocess.Obligation{subprocess.ObligationGroup}
		},
		"failed stdout": func(o *subprocess.Outcome) { o.Streams.Stdout.Status = subprocess.StreamFailed },
		"closed stderr": func(o *subprocess.Outcome) { o.Streams.Stderr.Status = subprocess.StreamClosed },
		"pending stdin": func(o *subprocess.Outcome) { o.Streams.Stdin.Status = subprocess.StreamPending },
		"protected":     func(o *subprocess.Outcome) { o.Protection.State = protection.Protected },
		"uncertain":     func(o *subprocess.Outcome) { o.Protection.State = protection.Uncertain },
	}
	for name, mutate := range cases {
		o := certain
		mutate(&o)
		if err := o.Uncertified(); err == nil {
			t.Errorf("%s: a zero exit certified %+v", name, o)
		}
	}
}
