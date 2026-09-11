package assessment

import (
	"fmt"
	"time"
)

func Validate(r Run) error {
	if r.Version != 1 {
		return fmt.Errorf("unsupported assessment version")
	}
	if !safeID.MatchString(r.RunID) || r.RepoKey == "" || r.Source == "" || r.TaskID == "" || r.Condition == "" || !state(r.State) {
		return fmt.Errorf("invalid run identity or state")
	}
	if err := interval(r.StartedAt, r.EndedAt); err != nil {
		return err
	}
	ids := map[string]bool{}
	events := map[string]string{}
	for _, a := range r.Attempts {
		if !safeID.MatchString(a.AttemptID) || ids[a.AttemptID] {
			return fmt.Errorf("invalid or duplicate attempt ID")
		}
		ids[a.AttemptID] = true
		if a.ChunkID == "" || a.SessionID == "" || a.Model == "" || a.Effort == "" || !state(a.State) {
			return fmt.Errorf("incomplete attempt identity")
		}
		switch a.Role {
		case "implementation", "repair", "verification", "review", "diagnostic":
		default:
			return fmt.Errorf("unknown performer role")
		}
		if err := interval(a.StartedAt, a.EndedAt); err != nil {
			return err
		}
		for _, e := range a.Usage {
			if e.EventID == "" || e.SessionID != a.SessionID || e.Epoch < 0 || e.Sequence < 0 || e.Counter == "" || !validReference(e.Reference) {
				return fmt.Errorf("missing native usage provenance")
			}
			key := e.SessionID + "/" + e.EventID
			if owner, ok := events[key]; ok && owner != a.AttemptID {
				return fmt.Errorf("event belongs to multiple attempts")
			}
			events[key] = a.AttemptID
		}
		if _, err := Estimate(a); err != nil {
			return err
		}
		for _, ref := range a.Evidence {
			if !validReference(ref) {
				return fmt.Errorf("missing evidence provenance")
			}
		}
	}
	for _, ref := range r.Evidence {
		if !validReference(ref) {
			return fmt.Errorf("missing evidence provenance")
		}
	}
	for _, v := range r.Quality {
		if v != nil && (!finite(*v) || *v < 0) {
			return fmt.Errorf("invalid quality measure")
		}
	}
	return nil
}
func state(s string) bool {
	switch s {
	case "running", "succeeded", "failed", "cancelled", "incomplete":
		return true
	}
	return false
}
func interval(start, end *time.Time) error {
	if start != nil && start.IsZero() || end != nil && end.IsZero() {
		return fmt.Errorf("zero timestamp")
	}
	if end != nil && (start == nil || end.Before(*start)) {
		return fmt.Errorf("invalid time interval")
	}
	return nil
}
