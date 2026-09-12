package assessment

import (
	"fmt"
	"math"
	"reflect"
	"sort"
)

// Usage keeps absent counters distinct from measured zero.
type Usage struct {
	Unknown        []string `json:"unknown,omitempty"`
	InputUncached  *int64   `json:"input_uncached,omitempty"`
	InputCached    *int64   `json:"input_cached,omitempty"`
	Output         *int64   `json:"output,omitempty"`
	InputTotal     *int64   `json:"input_total,omitempty"`
	TotalSemantics string   `json:"total_semantics,omitempty"`
}

func NormalizeUsage(u Usage) (Usage, error) {
	for _, n := range []*int64{u.InputUncached, u.InputCached, u.Output, u.InputTotal} {
		if n != nil && *n < 0 {
			return u, fmt.Errorf("negative usage")
		}
	}
	if u.InputTotal != nil {
		var n int64
		switch u.TotalSemantics {
		case "inclusive":
			if u.InputCached == nil {
				return u, nil
			}
			n = *u.InputTotal - *u.InputCached
		case "exclusive":
			n = *u.InputTotal
		default:
			return u, fmt.Errorf("unknown input total semantics")
		}
		if n < 0 || (u.InputUncached != nil && *u.InputUncached != n) {
			return u, fmt.Errorf("conflicting input total")
		}
		u.InputUncached = &n
	}
	return u, nil
}

func UsageTotal(events []Event) (Usage, error) {
	deltas, err := eventDeltas(events)
	if err != nil {
		return Usage{}, err
	}
	return sumEvents(deltas, nil)
}

func eventDeltas(events []Event) ([]Event, error) {
	seen := map[string]Event{}
	groups := map[string][]Event{}
	epochs := map[string]int{}
	regressed := map[string]bool{}
	for _, e := range events {
		key := e.SessionID + "/" + e.EventID
		if prior, ok := seen[key]; ok {
			if !reflect.DeepEqual(prior, e) {
				return nil, fmt.Errorf("conflicting native event")
			}
			continue
		}
		seen[key] = e
		if e.Mode != DeltaMode && e.Mode != CumulativeMode {
			return nil, fmt.Errorf("unrecognized counter semantics")
		}
		stream := e.SessionID + "/" + e.Counter
		if e.Mode == CumulativeMode {
			if last, ok := epochs[stream]; ok && e.Epoch < last {
				regressed[stream] = true
			}
			epochs[stream] = e.Epoch
		}
		group := fmt.Sprintf("%s/%s/%d", e.SessionID, e.Counter, e.Epoch)
		groups[group] = append(groups[group], e)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var deltas []Event
	for _, key := range keys {
		entries := groups[key]
		sort.Slice(entries, func(i, j int) bool { return entries[i].Sequence < entries[j].Sequence })
		var previous Usage
		var group []Event
		ambiguous := entries[0].Mode == CumulativeMode && regressed[entries[0].SessionID+"/"+entries[0].Counter]
		for i, e := range entries {
			u, err := NormalizeUsage(e.Usage)
			if err != nil {
				return nil, err
			}
			if e.Mode != entries[0].Mode || (e.Mode == CumulativeMode && i > 0 && e.Sequence == entries[i-1].Sequence) {
				ambiguous = true
				break
			}
			if e.Mode == CumulativeMode {
				current := u
				if current.InputUncached == nil {
					current.InputUncached = previous.InputUncached
				}
				if current.InputCached == nil {
					current.InputCached = previous.InputCached
				}
				if current.Output == nil {
					current.Output = previous.Output
				}
				dst := []**int64{&u.InputUncached, &u.InputCached, &u.Output}
				prev := []*int64{previous.InputUncached, previous.InputCached, previous.Output}
				for j, n := range []*int64{u.InputUncached, u.InputCached, u.Output} {
					if n == nil || prev[j] == nil {
						continue
					}
					if *n < *prev[j] {
						ambiguous = true
						break
					}
					v := *n - *prev[j]
					*dst[j] = &v
				}
				previous = current
			}
			u.InputTotal = nil
			u.TotalSemantics = ""
			e.Usage = u
			e.Mode = DeltaMode
			group = append(group, e)
		}
		if ambiguous {
			for _, e := range entries {
				e.Usage = Usage{Unknown: []string{key + ": ambiguous sequence or epoch"}}
				e.Mode = DeltaMode
				deltas = append(deltas, e)
			}
		} else {
			deltas = append(deltas, group...)
		}
	}
	return deltas, nil
}
func sumEvents(events []Event, unknown []string) (Usage, error) {
	out := Usage{Unknown: append([]string(nil), unknown...)}
	for _, e := range events {
		out.Unknown = append(out.Unknown, e.Usage.Unknown...)
		dst := []**int64{&out.InputUncached, &out.InputCached, &out.Output}
		for j, n := range []*int64{e.Usage.InputUncached, e.Usage.InputCached, e.Usage.Output} {
			if n == nil {
				out.Unknown = append(out.Unknown, fmt.Sprintf("%s: missing category %d", e.EventID, j))
				continue
			}
			if *dst[j] == nil {
				v := int64(0)
				*dst[j] = &v
			}
			if **dst[j] > math.MaxInt64-*n {
				return out, fmt.Errorf("usage overflow")
			}
			**dst[j] += *n
		}
	}
	return out, nil
}
func RunUsage(r Run) (map[string]Usage, error) {
	var events []Event
	owners := map[string]string{}
	for _, a := range r.Attempts {
		for _, e := range a.Usage {
			events = append(events, e)
			owners[e.SessionID+"/"+e.EventID] = a.AttemptID
		}
	}
	deltas, err := eventDeltas(events)
	if err != nil {
		return nil, err
	}
	byAttempt := map[string][]Event{}
	for _, e := range deltas {
		owner := owners[e.SessionID+"/"+e.EventID]
		byAttempt[owner] = append(byAttempt[owner], e)
	}
	out := map[string]Usage{}
	for _, a := range r.Attempts {
		u, err := sumEvents(byAttempt[a.AttemptID], nil)
		if err != nil {
			return nil, err
		}
		out[a.AttemptID] = u
	}
	return out, nil
}
