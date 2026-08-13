package canary

import (
	"reflect"
	"testing"
)

func TestSelectCanaryOwnersFromImmutableFixtures(t *testing.T) {
	accepted := []Fixture{
		{Dir: "/fixtures/a", Check: "line-routing"},
		{Dir: "/fixtures/b", Check: "package-core-guard"},
	}
	for _, tc := range []struct {
		name     string
		fixtures []Fixture
		accept   bool
	}{
		{name: "accepted", fixtures: accepted, accept: true},
		{name: "empty", fixtures: nil},
		{name: "missing owner", fixtures: []Fixture{{Dir: "/fixtures/unowned"}}},
		{name: "duplicate", fixtures: []Fixture{{Dir: "/first/repeated", Family: "test"}, {Dir: "/second/repeated", Family: "test"}}},
		{name: "hostile", fixtures: []Fixture{{Dir: "/fixtures/hostile\x1b", Family: "test"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			first := Select(tc.fixtures)
			second := Select(tc.fixtures)
			if first.Accepted != tc.accept || !reflect.DeepEqual(first, second) {
				t.Fatalf("Select() = %#v then %#v, accept=%v", first, second, tc.accept)
			}
		})
	}
}
