package durablefiletest

import (
	"fmt"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// Fault names one replacement by its destination and by its ordinal among the
// replacements of that destination, counting from 1.
type Fault struct {
	Destination string
	Ordinal     int
}

// String names f as its destination and ordinal, such as "policy#2".
func (f Fault) String() string { return fmt.Sprintf("%s#%d", f.Destination, f.Ordinal) }

// Router routes each replacement of several destinations through its own Recorder, which
// fails the stages that Fail names for that replacement. So one fault and an independent
// later fault can fail one caller operation. The replacement is known from its rename, so
// a fault fails the rename and the stages after it.
type Router struct {
	// Destination names the destination of the renamed path.
	Destination func(path string) string
	Fail        map[Fault]map[string]error
	// Trace holds each replacement in the order of its rename.
	Trace []Fault
	paths map[string]string
}

// Count returns how many replacements of destination reached their rename.
func (r *Router) Count(destination string) int {
	count := 0
	for _, fault := range r.Trace {
		if fault.Destination == destination {
			count++
		}
	}
	return count
}

// Replacer returns the native replacer with each rename and the stages after it routed
// through the Recorder of its replacement.
func (r *Router) Replacer() durablefile.Replacer {
	native := durablefile.Native()
	current := native
	return durablefile.Replacer{
		CreateTemp: native.CreateTemp,
		Rename: func(oldpath, newpath string) error {
			destination := r.Destination(newpath)
			fault := Fault{Destination: destination, Ordinal: r.Count(destination) + 1}
			r.Trace = append(r.Trace, fault)
			if r.paths == nil {
				r.paths = map[string]string{}
			}
			r.paths[destination] = newpath
			current = (&Recorder{Fail: r.Fail[fault]}).Replacer()
			return current.Rename(oldpath, newpath)
		},
		OpenDirectory: func(name string) (durablefile.Directory, error) { return current.OpenDirectory(name) },
		Remove:        native.Remove,
	}
}

// Failure is the message of the injected cause at stage of the last replacement of
// destination.
func (r *Router) Failure(destination, stage string, published bool, cause error) string {
	return (&durablefile.Error{Path: r.paths[destination], Stage: stage, Published: published, Err: cause}).Error()
}
