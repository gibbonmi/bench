package assessment

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/toon"
	"reflect"
	"slices"
	"time"
)

func Validate(r Run) error {
	if !safeText(reflect.ValueOf(r)) {
		return fmt.Errorf("control characters in assessment record")
	}
	if r.Version != 1 {
		return fmt.Errorf("unsupported assessment version")
	}
	if !safeID.MatchString(r.RunID) || r.RepoKey == "" || r.Source == "" || r.TaskID == "" || r.Condition == "" || !state(r.State) {
		return fmt.Errorf("invalid run identity or state")
	}
	if err := interval(r.StartedAt, r.EndedAt, r.TimeReference); err != nil {
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
		if !slices.Contains(Roles(), a.Role) {
			return fmt.Errorf("unknown performer role")
		}
		if err := interval(a.StartedAt, a.EndedAt, a.TimeReference); err != nil {
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
		if !validReference(v.Reference) || (v.Value != nil && (!finite(*v.Value) || *v.Value < 0)) {
			return fmt.Errorf("invalid quality measure")
		}
	}
	return nil
}
func state(s string) bool { return slices.Contains(States(), s) }
func interval(start, end *time.Time, ref *Reference) error {
	if (start != nil || end != nil) && (ref == nil || !validReference(*ref)) {
		return fmt.Errorf("missing time provenance")
	}
	if ref != nil && !validReference(*ref) {
		return fmt.Errorf("invalid time provenance")
	}
	if start != nil && start.IsZero() || end != nil && end.IsZero() {
		return fmt.Errorf("zero timestamp")
	}
	if end != nil && (start == nil || end.Before(*start)) {
		return fmt.Errorf("invalid time interval")
	}
	return nil
}

func safeText(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return toon.Representable(v.String())
	case reflect.Pointer, reflect.Interface:
		if !v.IsNil() {
			return safeText(v.Elem())
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !safeText(v.Field(i)) {
				return false
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !safeText(v.Index(i)) {
				return false
			}
		}
	case reflect.Map:
		it := v.MapRange()
		for it.Next() {
			if !safeText(it.Key()) || !safeText(it.Value()) {
				return false
			}
		}
	}
	return true
}
