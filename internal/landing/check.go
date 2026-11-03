package landing

import "fmt"

// CheckRefusal is the refusal of a candidate commit that the request's check did not pass.
// Land refuses it before the ref update, so no ref moved and the candidate object stays
// unreferenced. A caller routes the refusal by the type and never by the sentence.
type CheckRefusal struct {
	Candidate string
	Err       error
}

func (r CheckRefusal) Error() string {
	return fmt.Sprintf("the pre-publication check refused candidate %s: %v", r.Candidate, r.Err)
}
func (r CheckRefusal) Unwrap() error { return r.Err }

// checkCandidate runs the request's check on the candidate commit, after the commit object
// exists and before any ref names it. A request with no check passes.
func checkCandidate(r Request, candidate string) error {
	if r.Check == nil {
		return nil
	}
	if err := r.Check(candidate); err != nil {
		return CheckRefusal{Candidate: candidate, Err: err}
	}
	return nil
}
