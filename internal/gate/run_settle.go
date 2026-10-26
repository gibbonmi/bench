package gate

import (
	"fmt"
	"io"
	"strings"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// settleRun checks a finished oracle run before any verdict: the subject after the run is
// the subject that it graded, a failed run left a safe protection state, and the run's
// resources release. It answers the subject after the run, or the refusal and true.
func settleRun(evaluation executionEvaluation, runtimeRoot, storageRoot string, declaredPaths []string, plan subject, run processGroupResult, release func() error, stderr io.Writer) (subject, Result, bool) {
	rc := run.Code
	after, err := evaluation.capturePost()
	if err != nil || !sameSubject(plan, after) {
		fmt.Fprintln(stderr, subjectChangeDiagnostic(runtimeRoot, declaredPaths, plan, after, err))
		inspection := inspectAt(storageRoot, time.Now().UTC())
		if err == nil {
			inspection = inspectSubjectAt(storageRoot, after, time.Now().UTC())
		}
		return after, Result{GateExit: rc, ActionExit: 1, Inspection: inspection}, true
	}
	if rc != 0 && run.Outcome.Process.Started && run.Outcome.Protection.State != protection.Safe {
		return after, unprotectedFailure(storageRoot, plan, rc, stderr, run), true
	}
	// A nested user that outlives a zero exit keeps a resource; its release failure is no verdict.
	if err := release(); err != nil {
		return after, retainedFailure(storageRoot, plan, rc, stderr, "gate resource release failed: "+err.Error()), true
	}
	return after, Result{}, false
}

func sameSubject(a, b subject) bool {
	return a.Tree == b.Tree && a.Oracle == b.Oracle && a.Resolution == b.Resolution && a.Closed == b.Closed && a.Reason == b.Reason
}

func subjectChangeDiagnostic(root string, declaredPaths []string, before, after subject, err error) string {
	if err == nil && before.Tree != after.Tree {
		if changed, ok := benchgit.ChangedPathsBetweenTrees(root, before.Tree, after.Tree); ok {
			for _, path := range declaredPaths {
				for _, moved := range changed {
					if moved == path || strings.HasPrefix(moved, path+"/") {
						return "gate subject changed during execution: declared input " + path
					}
				}
			}
		}
		return "gate subject changed during execution: tree " + after.Tree
	}
	return "gate subject changed during execution"
}
