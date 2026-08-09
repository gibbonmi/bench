package freshness

import "fmt"

// SelectionInput is the immutable digest pair a freshness adapter has loaded.
type SelectionInput struct {
	StoredSource, CurrentSource         string
	StoredExecutable, CurrentExecutable string
}

// Selection is the branch-native freshness decision.
type Selection struct {
	Accepted bool
	Reason   string
}

// Select accepts only complete, well-formed, matching source and executable pairs.
func Select(input SelectionInput) Selection {
	for _, digest := range []struct{ name, value string }{
		{name: "stored source", value: input.StoredSource},
		{name: "current source", value: input.CurrentSource},
		{name: "stored executable", value: input.StoredExecutable},
		{name: "current executable", value: input.CurrentExecutable},
	} {
		if !isDigest(digest.value) {
			return Selection{Reason: fmt.Sprintf("%s digest is malformed", digest.name)}
		}
	}
	if input.StoredExecutable != input.CurrentExecutable {
		return Selection{Reason: "seal executable digest does not match binary contents"}
	}
	if input.StoredSource != input.CurrentSource {
		return Selection{Reason: "seal source digest does not match current build inputs"}
	}
	return Selection{Accepted: true}
}
