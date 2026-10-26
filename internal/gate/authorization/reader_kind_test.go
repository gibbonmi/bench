// The reader kinds: attribution of a gate run that ran the reader set.
package authorization

import (
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

// readerTree is the composed tree of the reader rows. A reader set writes no verdict, so the
// inspection after it never reads reusable green.
const readerTree = "0123456789abcdef0123456789abcdef01234567"

// TestAuthorizeAttributesAReaderGreen is NG55. A green reader set left no reusable green
// evidence, so an authorizer that reads only the verdict cache attributes infrastructure.
func TestAuthorizeAttributesAReaderGreen(t *testing.T) {
	got := attribute(gate.Result{Schedule: gate.ScheduleReader}, gate.EvidenceInspection{}, false, readerTree)
	if got.Kind != ReaderGreen || got.Evidence != "" || got.Reason != "" {
		t.Fatalf("result = %+v, want kind %q with no evidence token and no reason", got, ReaderGreen)
	}
}

// TestAuthorizeAttributesAReaderRedToTheCandidate is NG56. A red reader set writes no
// verdict, so the chain proof stands in for the green baseline, with or without inherited
// green proof. A timeout or a cancel graded nothing, so it attributes infrastructure.
func TestAuthorizeAttributesAReaderRedToTheCandidate(t *testing.T) {
	for _, inheritedGreen := range []bool{false, true} {
		for exit, want := range map[int]Kind{1: Candidate, 124: Infrastructure, 130: Infrastructure} {
			run := gate.Result{GateExit: exit, ActionExit: exit, Schedule: gate.ScheduleReader}
			if got := attribute(run, gate.EvidenceInspection{}, inheritedGreen, readerTree); got.Kind != want {
				t.Fatalf("exit %d, inherited green %t: result = %+v, want kind %q", exit, inheritedGreen, got, want)
			}
		}
	}
}
