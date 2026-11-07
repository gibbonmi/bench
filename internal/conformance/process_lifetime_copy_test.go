package conformance

import (
	"slices"
	"testing"
)

// processLifetimeCopyHeader opens each synthetic copy file. Its function opens at line 10.
const processLifetimeCopyHeader = "package example\n\nimport (\n\t\"os\"\n\n\t\"github.com/gibbonmi/bench/internal/subprocess\"\n\t\"github.com/gibbonmi/bench/internal/subprocess/protection\"\n)\n\n"

// TestProcessLifetimeInheritedRegistrationCopy is RS18. It grades the hand-built registration
// rule on synthetic trees: each counted reader with a Registration literal, in the reading
// function or in a same-package helper, turns red; a read of another value or a read with no
// Registration stays green. The copy file spells no owner call, so only the reader tokens let
// the prefilter parse it.
func TestProcessLifetimeInheritedRegistrationCopy(t *testing.T) {
	const rel = "internal/example/copy.go"
	grade := func(function string) []string {
		t.Helper()
		files := map[string]string{"internal/example/run.go": processLifetimeOwned, rel: processLifetimeCopyHeader + function}
		return processLifetimeClosureDiags(throwawayRoot{files: files}.build(t), processLifetimeFixtureCensus)
	}
	red := []struct{ name, function, want string }{
		{"inherited", "func register() subprocess.Registration {\n\tinherited, _ := protection.Inherited()\n\treturn subprocess.Registration{Resources: inherited}\n}\n",
			rel + ":11 register reads the process vector through protection.Inherited and builds a subprocess.Registration by hand outside internal/subprocess"},
		{"inherited from the process environment", "func register() subprocess.Registration {\n\tinherited, _ := protection.InheritedFrom(os.Environ())\n\treturn subprocess.Registration{Resources: inherited}\n}\n",
			rel + ":11 register reads the process vector through protection.InheritedFrom and builds a subprocess.Registration by hand outside internal/subprocess"},
		{"parsed process variable", "func register() subprocess.Registration {\n\tinherited, _ := protection.ParseVector(os.Getenv(protection.Env))\n\treturn subprocess.Registration{Resources: inherited}\n}\n",
			rel + ":11 register reads the process vector through protection.ParseVector and builds a subprocess.Registration by hand outside internal/subprocess"},
		{"registration in a same-package helper", "func register() subprocess.Registration {\n\tinherited, _ := protection.Inherited()\n\treturn wrap(inherited)\n}\n\nfunc wrap(resources []protection.Descriptor) subprocess.Registration {\n\treturn subprocess.Registration{Resources: resources}\n}\n",
			rel + ":11 register reads the process vector through protection.Inherited and builds a subprocess.Registration by hand outside internal/subprocess"},
	}
	for _, c := range red {
		if diags := grade(c.function); !slices.Equal(diags, []string{c.want}) {
			t.Errorf("%s: want %q, got %q", c.name, c.want, diags)
		}
	}
	green := []struct{ name, function string }{
		{"inherited from another environment", "func register(environ []string) subprocess.Registration {\n\tinherited, _ := protection.InheritedFrom(environ)\n\treturn subprocess.Registration{Resources: inherited}\n}\n"},
		{"parsed other value", "func register(value string) subprocess.Registration {\n\tinherited, _ := protection.ParseVector(os.Getenv(value))\n\treturn subprocess.Registration{Resources: inherited}\n}\n"},
		{"reader with no registration", "func members() int {\n\tinherited, _ := protection.Inherited()\n\treturn len(inherited)\n}\n"},
	}
	for _, c := range green {
		if diags := grade(c.function); len(diags) != 0 {
			t.Errorf("%s: want no diagnostics, got %q", c.name, diags)
		}
	}
}
