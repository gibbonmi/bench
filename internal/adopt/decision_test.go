package adopt

import (
	"reflect"
	"testing"
)

func TestPlanLifecycleFromImmutableInventories(t *testing.T) {
	input := LifecycleInput{
		Current:  []Asset{{Path: "changed", Fingerprint: "old"}, {Path: "removed", Fingerprint: "old"}, {Path: "CLAUDE.md", Fingerprint: "keep"}},
		Desired:  []Asset{{Path: "added", Fingerprint: "new"}, {Path: "changed", Fingerprint: "new"}},
		Preserve: []string{"CLAUDE.md"},
	}
	want := []Operation{{Kind: "add", Path: "added"}, {Kind: "change", Path: "changed"}, {Kind: "remove", Path: "removed"}}
	first := PlanLifecycle(input)
	second := PlanLifecycle(input)
	if first.Refusal != "" || !reflect.DeepEqual(first.Operations, want) || !reflect.DeepEqual(second, first) {
		t.Fatalf("PlanLifecycle() = %#v then %#v, want %#v", first, second, want)
	}
	for _, tc := range []struct {
		name  string
		input LifecycleInput
	}{
		{name: "empty", input: LifecycleInput{}},
		{name: "duplicate", input: LifecycleInput{Current: []Asset{{Path: "x", Fingerprint: "a"}, {Path: "x", Fingerprint: "b"}}}},
		{name: "hostile", input: LifecycleInput{Desired: []Asset{{Path: "../escape", Fingerprint: "x"}}}},
		{name: "missing fingerprint", input: LifecycleInput{Desired: []Asset{{Path: "x"}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := PlanLifecycle(tc.input)
			if tc.name == "empty" {
				if got.Refusal != "" || len(got.Operations) != 0 {
					t.Fatalf("empty plan = %#v", got)
				}
				return
			}
			if got.Refusal == "" {
				t.Fatalf("plan accepted invalid input: %#v", got)
			}
		})
	}
}
