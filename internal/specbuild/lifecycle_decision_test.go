package specbuild

import "testing"

func TestTransitionLifecycleFromImmutableValues(t *testing.T) {
	prepared := LifecycleValue{Command: "assign", Request: "ticket-1", Input: "digest", State: "prepared"}
	completed := LifecycleValue{Command: "assign", Request: "ticket-1", Input: "digest", Result: "assignment-1", State: "completed"}
	tests := []struct {
		name     string
		current  LifecycleValue
		event    LifecycleEvent
		accepted bool
		want     LifecycleValue
		refusal  string
	}{
		{name: "prepare", event: LifecycleEvent{Kind: "prepare", Command: "assign", Request: "ticket-1", Input: "digest"}, accepted: true, want: prepared},
		{name: "record", current: prepared, event: LifecycleEvent{Kind: "record", Command: "assign", Request: "ticket-1", Input: "digest", Result: "assignment-1"}, accepted: true, want: LifecycleValue{Command: "assign", Request: "ticket-1", Input: "digest", Result: "assignment-1", State: "prepared"}},
		{name: "complete", current: prepared, event: LifecycleEvent{Kind: "complete", Command: "assign", Request: "ticket-1", Input: "digest", Result: "assignment-1"}, accepted: true, want: completed},
		{name: "rerun", current: completed, event: LifecycleEvent{Kind: "prepare", Command: "assign", Request: "ticket-1", Input: "digest"}, accepted: true, want: completed},
		{name: "opaque request", event: LifecycleEvent{Kind: "prepare", Command: "refresh", Request: "assignment\x00base", Input: "digest"}, accepted: true, want: LifecycleValue{Command: "refresh", Request: "assignment\x00base", Input: "digest", State: "prepared"}},
		{name: "empty", event: LifecycleEvent{}, want: LifecycleValue{}, refusal: "incomplete lifecycle event"},
		{name: "hostile", event: LifecycleEvent{Kind: "prepare", Command: "assign\nnext", Request: "ticket-1", Input: "digest"}, want: LifecycleValue{}, refusal: "invalid lifecycle value"},
		{name: "unknown", current: prepared, event: LifecycleEvent{Kind: "promote", Command: "assign", Request: "ticket-1", Input: "digest"}, want: prepared, refusal: "unknown lifecycle event"},
		{name: "conflict", current: prepared, event: LifecycleEvent{Kind: "complete", Command: "assign", Request: "ticket-2", Input: "digest", Result: "assignment-1"}, want: prepared, refusal: "lifecycle operation is not prepared"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := TransitionLifecycle(test.current, test.event)
			if got.Accepted != test.accepted || got.After != test.want || got.Refusal != test.refusal {
				t.Fatalf("TransitionLifecycle() = %#v, want accepted=%t after=%#v refusal=%q", got, test.accepted, test.want, test.refusal)
			}
		})
	}
}
