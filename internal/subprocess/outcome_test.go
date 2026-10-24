package subprocess_test

import (
	"strings"
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

// TestUncertifiedNamesEveryUnresolvedPart is PL116 over the protection and mixed cases. A
// protected disposition with no error still names its users, and one diagnostic names the
// cleanup obligation, the failed stream, and the retained protection together.
func TestUncertifiedNamesEveryUnresolvedPart(t *testing.T) {
	none := subprocess.Stream{Status: subprocess.StreamNone}
	o := subprocess.Outcome{
		Process:    subprocess.Process{Started: true, Waited: true},
		Streams:    subprocess.StreamSet{Stdin: none, Stdout: subprocess.Stream{Status: subprocess.StreamFailed}, Stderr: none},
		Cleanup:    subprocess.Cleanup{Status: subprocess.CleanupIncomplete, Obligations: []subprocess.Obligation{subprocess.ObligationWaiter}},
		Protection: protection.Disposition{State: protection.Protected, Users: []string{"user-a"}},
	}
	protected := o
	protected.Streams.Stdout.Status, protected.Cleanup = subprocess.StreamComplete, subprocess.Cleanup{}
	for name, tc := range map[string]struct {
		outcome subprocess.Outcome
		want    []string
	}{
		"protection": {protected, []string{"process protection protected", "user-a"}},
		"mixed":      {o, []string{"process cleanup incomplete", string(subprocess.ObligationWaiter), "process stream failed", "process protection protected", "user-a"}},
	} {
		err := tc.outcome.Uncertified()
		if err == nil || strings.Contains(err.Error(), "<nil>") {
			t.Fatalf("%s: diagnostic = %v, want every unresolved part and no empty cause", name, err)
		}
		for _, want := range tc.want {
			if strings.Count(err.Error(), want) != 1 {
				t.Errorf("%s: diagnostic %q names %q other than once", name, err, want)
			}
		}
	}
}
