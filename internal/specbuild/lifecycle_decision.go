package specbuild

import "strings"

// LifecycleValue is the immutable public state of one durable build operation.
type LifecycleValue struct {
	Command, Request, Input, Result, State string
}

// LifecycleEvent describes one requested durable build-operation transition.
type LifecycleEvent struct {
	Kind, Command, Request, Input, Result string
}

// LifecycleTransition reports whether an event was accepted and the resulting value.
type LifecycleTransition struct {
	Accepted bool
	Before   LifecycleValue
	After    LifecycleValue
	Refusal  string
}

// TransitionLifecycle applies one operation event without reading or writing repository state.
func TransitionLifecycle(current LifecycleValue, event LifecycleEvent) LifecycleTransition {
	result := LifecycleTransition{Before: current, After: current}
	if invalidLifecycleText(event.Kind) || invalidLifecycleText(event.Command) {
		result.Refusal = "invalid lifecycle value"
		return result
	}
	if event.Command == "" || event.Request == "" || event.Input == "" {
		result.Refusal = "incomplete lifecycle event"
		return result
	}
	want := LifecycleValue{Command: event.Command, Request: event.Request, Input: event.Input}
	switch event.Kind {
	case "prepare":
		if current.State == "" {
			want.State = "prepared"
			result.Accepted, result.After = true, want
			return result
		}
		if current.Command == want.Command && current.Request == want.Request && current.Input == want.Input {
			result.Accepted = true
			return result
		}
		result.Refusal = "lifecycle input conflicts with durable state"
	case "record", "complete":
		if current.State != "prepared" || current.Command != want.Command || current.Request != want.Request || current.Input != want.Input {
			result.Refusal = "lifecycle operation is not prepared"
			return result
		}
		want.Result = event.Result
		want.State = "prepared"
		if event.Kind == "complete" {
			want.State = "completed"
		}
		result.Accepted, result.After = true, want
	default:
		result.Refusal = "unknown lifecycle event"
	}
	return result
}

func invalidLifecycleText(value string) bool {
	return strings.TrimSpace(value) != value || strings.ContainsAny(value, "\x00\r\n")
}
