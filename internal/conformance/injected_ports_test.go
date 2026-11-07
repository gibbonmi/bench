package conformance

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/conformance/injectedports"
	"github.com/gibbonmi/bench/internal/gate"
)

// The injected-interface composition audit found the same gap at five seams at once. A
// port injected into a consumer whose tests only ever drive a fake lets a softened real
// producer leave the gate green. Junction tests close the seams that were open. This
// check keeps them closed. It refuses a derived port that names neither a test driving
// its real producer nor a written reason it does not.
//
// The two halves are deliberately independent. The inventory is derived from the Go source
// of every production package under cmd and internal, so a port added tomorrow appears
// without anyone remembering to say so. The registry file is the advertisement, and disagreement between them is the red.
// A single hand-maintained list would satisfy itself.
//
// The named-test-exists half is a tripwire, not a behavior check. It catches a junction
// test deleted or renamed, but it cannot catch one gutted back to driving a fake. The
// behavior half of the defense is the junction tests themselves and the canary fixture
// under tests/canary/injected-ports. That fixture keeps the unregistered-port red alive.

// checkInjectedPortRegistry grades the registry against a fresh derivation over root.
func checkInjectedPortRegistry(root string) []string {
	return injectedports.Check(root, injectedports.Policy{Rows: injectedPortRegistry, Required: auditedPortPackages, Families: gradedRootFamilies(root)})
}

// gradedRootFamilies binds graded-root data to the local names that construct a hook.
// The release phase names of root's own preflight registry bind the production reader's
// phase parameter and the native probe's override loop. The environment names that
// root's gate input manifest declares bind the gate's declared-environment reader. An
// unusable source supplies no family, so its reader stays unresolved instead of grading
// silence; a valid manifest that declares no names resolves to no hooks.
func gradedRootFamilies(root string) map[string][]string {
	families := map[string][]string{}
	if registries := readReleaseRegistries(root); registries.preflightUsable && len(registries.preflight.Verify) > 0 {
		phases := registries.preflight.Verify
		families["internal/releasepreflight.runner.runExternal.name"] = phases
		families["internal/conformance.preflightPhaseOverrides.phase"] = phases
	}
	if names, declared := gate.DeclaredEnvironment(root); declared {
		families["internal/gate.suppliedEnvironment.name"] = append([]string{}, names...)
	}
	return families
}

// preflightPhaseOverrides returns one phase override entry per phase, each naming path.
func preflightPhaseOverrides(phases []string, path string) []string {
	var env []string
	for _, phase := range phases {
		env = append(env, "BENCH_PREFLIGHT_"+strings.ToUpper(phase)+"="+path)
	}
	return env
}

// runAtEvidenceReadyPoint runs a preflight that holds at the release-evidence readiness
// point. The observer sees the ready file that the held run writes, then removes it to
// release the run. A run that exits without writing it never reached that point.
func runAtEvidenceReadyPoint(cmd *exec.Cmd, ready string) ([]byte, error) {
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	for {
		select {
		case err := <-exited:
			return output.Bytes(), errors.Join(err, errors.New("the run never reached the evidence readiness point"))
		default:
		}
		if _, err := os.Stat(ready); err == nil {
			break
		}
		time.Sleep(bounds.FixedWindow(10 * time.Millisecond))
	}
	if err := os.Remove(ready); err != nil {
		return nil, errors.Join(err, cmd.Process.Kill(), <-exited)
	}
	return output.Bytes(), <-exited
}

// TestInjectedPortPhaseFamilyReadsTheGradedRegistry proves that the wrapper's phase family
// comes from the graded root's own release registry. Each case writes that registry, a
// phase reader, and a test that sets two phase hooks. The registry alone decides which of
// them are hooks, and an absent registry leaves the reader unresolved.
func TestInjectedPortPhaseFamilyReadsTheGradedRegistry(t *testing.T) {
	const hook = "environment hook has no registry row: BENCH_PREFLIGHT_"
	for _, test := range []struct {
		name, registry string
		heads          []string
	}{
		{"the registry's vulnerability phase", `["gate", "vulnerability"]`, []string{hook + "VULNERABILITY "}},
		{"a newly registered phase", `["gate", "vulnerability", "newphase"]`, []string{hook + "NEWPHASE ", hook + "VULNERABILITY "}},
		{"an absent registry", "", []string{injectedports.UnresolvedHookMessage + ": internal/releasepreflight.runner.runExternal "}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			files := map[string]string{
				"internal/releasepreflight/command.go":      "package releasepreflight\n\nimport (\n\t\"os\"\n\t\"strings\"\n)\n\ntype runner struct{}\n\nfunc (r *runner) runExternal(name string) string { return os.Getenv(\"BENCH_PREFLIGHT_\" + strings.ToUpper(name)) }\n",
				"internal/releasepreflight/command_test.go": "package releasepreflight\n\nimport \"testing\"\n\nfunc TestPhases(t *testing.T) {\n\tt.Setenv(\"BENCH_PREFLIGHT_VULNERABILITY\", \"x\")\n\tt.Setenv(\"BENCH_PREFLIGHT_NEWPHASE\", \"x\")\n\t_ = (&runner{}).runExternal(\"x\")\n}\n",
			}
			if test.registry != "" {
				files["internal/releaseevidence/registry.json"] = `{"verify": ` + test.registry + `, "publish_before": [], "publish_only": []}`
			}
			for rel, body := range files {
				writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(rel)), body)
			}
			diags := injectedports.Check(root, injectedports.Policy{Families: gradedRootFamilies(root)})
			if len(diags) != len(test.heads) {
				t.Fatalf("diagnostics %q, want heads %q", diags, test.heads)
			}
			for i, head := range test.heads {
				if !strings.HasPrefix(diags[i], head) {
					t.Fatalf("diagnostic %d = %q, want head %q", i, diags[i], head)
				}
			}
		})
	}
}

// TestInjectedPortRegistryCheckBites is the recorded bite proof for
// checkInjectedPortRegistry. The unregistered-port red is retained permanently by the
// tests/canary/injected-ports fixture. The other four modes have no fixture of their own
// and are proved here, each against a synthetic tree.
func TestInjectedPortRegistryCheckBites(t *testing.T) {
	// This tree holds one package, one port-shaped type, and one injection site. It is the
	// smallest tree the derivation reports anything at all for.
	const portSource = "package sample\n\n// Port is the fixture's injected port.\ntype Port interface {\n\tDo() error\n}\n\nfunc use(p Port) error { return p.Do() }\n"
	const testSource = "package sample\n\nimport \"testing\"\n\nfunc TestRealProducer(t *testing.T) {}\n"

	write := func(t *testing.T, files map[string]string) string {
		t.Helper()
		root := t.TempDir()
		for rel, body := range files {
			writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(rel)), body)
		}
		return root
	}
	check := func(root string, packages []string, rows []injectedports.Row) []string {
		return injectedports.Check(root, injectedports.Policy{Rows: rows, Required: packages})
	}
	sampleTree := map[string]string{
		"internal/sample/sample.go":      portSource,
		"internal/sample/sample_test.go": testSource,
	}
	honestRow := injectedports.Row{
		Package: "internal/sample", Port: "Port",
		TestFile: "internal/sample/sample_test.go", TestName: "TestRealProducer",
	}
	packages := []string{"internal/sample"}

	t.Run("an honest row is silent", func(t *testing.T) {
		diags := check(write(t, sampleTree), packages, []injectedports.Row{honestRow})
		if len(diags) != 0 {
			t.Fatalf("a derived port with a real-producer test the tree declares: want no diagnostics, got %v", diags)
		}
	})

	t.Run("a package the tree does not hold is not graded", func(t *testing.T) {
		if diags := check(t.TempDir(), packages, nil); len(diags) != 0 {
			t.Fatalf("an absent audited package: want no diagnostics, got %v", diags)
		}
	})

	t.Run("an unregistered port fires", func(t *testing.T) {
		diags := check(write(t, sampleTree), packages, nil)
		if len(diags) != 1 || !strings.Contains(diags[0], injectedports.UnregisteredPortMessage) ||
			!strings.Contains(diags[0], "internal/sample.Port") || !strings.Contains(diags[0], "a parameter of use in internal/sample/sample.go") {
			t.Fatalf("a derived port with no registry row: want one unregistered diagnostic naming it and its site, got %v", diags)
		}
	})

	t.Run("a row naming a test the tree lacks fires", func(t *testing.T) {
		stale := honestRow
		stale.TestName = "TestRealProducerRenamedAway"
		diags := check(write(t, sampleTree), packages, []injectedports.Row{stale})
		if len(diags) != 1 || !strings.Contains(diags[0], injectedports.MissingPortTestMessage) ||
			!strings.Contains(diags[0], "TestRealProducerRenamedAway") {
			t.Fatalf("a row advertising an absent test: want one missing-test diagnostic naming it, got %v", diags)
		}
	})

	t.Run("a test named only in a comment does not satisfy the tripwire", func(t *testing.T) {
		commented := map[string]string{
			"internal/sample/sample.go":      portSource,
			"internal/sample/sample_test.go": "package sample\n\n// TestRealProducer used to live here.\nconst gone = \"func TestRealProducer(t *testing.T) {}\"\n",
		}
		diags := check(write(t, commented), packages, []injectedports.Row{honestRow})
		if len(diags) != 1 || !strings.Contains(diags[0], injectedports.MissingPortTestMessage) {
			t.Fatalf("a test surviving only as prose: want one missing-test diagnostic, got %v", diags)
		}
	})

	t.Run("an exemption with an empty reason fires", func(t *testing.T) {
		blank := injectedports.Row{Package: "internal/sample", Port: "Port", Exempt: "   "}
		diags := check(write(t, sampleTree), packages, []injectedports.Row{blank})
		if len(diags) != 1 || !strings.Contains(diags[0], injectedports.EmptyPortExemptionMessage) ||
			!strings.Contains(diags[0], "internal/sample.Port") {
			t.Fatalf("an exemption with a blank reason: want one empty-exemption diagnostic, got %v", diags)
		}
	})

	t.Run("a written exemption is silent", func(t *testing.T) {
		reasoned := injectedports.Row{Package: "internal/sample", Port: "Port", Exempt: "covered by the runtime contract suite"}
		if diags := check(write(t, sampleTree), packages, []injectedports.Row{reasoned}); len(diags) != 0 {
			t.Fatalf("an exemption carrying a reason: want no diagnostics, got %v", diags)
		}
	})

	t.Run("a package present but deriving nothing fails closed", func(t *testing.T) {
		empty := map[string]string{"internal/sample/sample.go": "package sample\n\n// Name carries no port.\nfunc Name() string { return \"sample\" }\n"}
		diags := check(write(t, empty), packages, []injectedports.Row{honestRow})
		if len(diags) != 1 || !strings.Contains(diags[0], injectedports.ZeroPortInventoryMessage) ||
			!strings.Contains(diags[0], "internal/sample") {
			t.Fatalf("an audited package deriving zero ports: want one zero-inventory diagnostic naming it, got %v", diags)
		}
	})

	t.Run("a retired live-tree row is absent", func(t *testing.T) {
		for _, row := range injectedPortRegistry {
			if row.Package == "internal/canary" && row.Port == "Runner" {
				t.Fatal("the retired internal/canary.Runner registry row returned")
			}
		}
	})

	t.Run("a retired live-tree audit fires", func(t *testing.T) {
		packages := append(append([]string(nil), auditedPortPackages...), "internal/canary")
		kit := NewHarness(t).KitRoot
		diags := injectedports.Check(kit, injectedports.Policy{Rows: injectedPortRegistry, Required: packages, Families: gradedRootFamilies(kit)})
		if len(diags) != 1 || !strings.Contains(diags[0], injectedports.ZeroPortInventoryMessage) || !strings.Contains(diags[0], "internal/canary") {
			t.Fatalf("a reintroduced canary audit: want its zero-port diagnostic, got %v", diags)
		}
	})

	t.Run("unparsable audited source is reported", func(t *testing.T) {
		broken := map[string]string{"internal/sample/sample.go": "package sample\n\nfunc Port broken {\n"}
		diags := check(write(t, broken), packages, nil)
		if len(diags) != 2 || !anyContains(diags, "cannot parse internal/sample/sample.go") ||
			!anyContains(diags, injectedports.ZeroPortInventoryMessage) {
			t.Fatalf("unparsable audited source: want a parse diagnostic and the fail-closed inventory red, got %v", diags)
		}
	})

	// twoPortTree carries two ports, PortA and PortB, both derived and honestly rowed. The
	// fixtures below pair them with orphanRow, a registry row naming PortC, a port the tree
	// never declares. The two honest ports sit beside the orphan, so the check must fire
	// only for the underived row, not collaterally for the derived ones.
	twoPortTree := map[string]string{
		"internal/sample2/sample2.go":      "package sample2\n\ntype PortA interface {\n\tDo() error\n}\n\ntype PortB interface {\n\tDo() error\n}\n\nfunc use(a PortA, b PortB) error {\n\tif err := a.Do(); err != nil {\n\t\treturn err\n\t}\n\treturn b.Do()\n}\n",
		"internal/sample2/sample2_test.go": "package sample2\n\nimport \"testing\"\n\nfunc TestRealProducerA(t *testing.T) {}\nfunc TestRealProducerB(t *testing.T) {}\n",
	}
	twoPortPackages := []string{"internal/sample2"}
	portARow := injectedports.Row{
		Package: "internal/sample2", Port: "PortA",
		TestFile: "internal/sample2/sample2_test.go", TestName: "TestRealProducerA",
	}
	portBRow := injectedports.Row{
		Package: "internal/sample2", Port: "PortB",
		TestFile: "internal/sample2/sample2_test.go", TestName: "TestRealProducerB",
	}
	orphanRow := injectedports.Row{
		Package: "internal/sample2", Port: "PortC",
		TestFile: "internal/sample2/sample2_test.go", TestName: "TestRealProducerC",
	}

	t.Run("a row naming an underived port fires the orphan diagnostic alone", func(t *testing.T) {
		diags := check(write(t, twoPortTree), twoPortPackages, []injectedports.Row{portARow, portBRow, orphanRow})
		if len(diags) != 1 || !strings.Contains(diags[0], injectedports.OrphanPortRowMessage) ||
			!strings.Contains(diags[0], "internal/sample2.PortC") {
			t.Fatalf("PortA and PortB derived and honestly rowed, PortC registered but never derived: want exactly one orphan-row diagnostic naming PortC, got %v", diags)
		}
	})

	t.Run("the orphan and missing-test diagnostics coexist without masking each other", func(t *testing.T) {
		// PortB's injection arm and its named test both survive here. The row points at a test
		// name the tree does not declare, so it takes the missing-test path. A row whose arm
		// was deleted along with its test would instead fall to the orphan path above. Both
		// cases run together to prove neither masks the other.
		staleTestRow := injectedports.Row{
			Package: "internal/sample2", Port: "PortB",
			TestFile: "internal/sample2/sample2_test.go", TestName: "TestRealProducerBRenamedAway",
		}
		diags := check(write(t, twoPortTree), twoPortPackages, []injectedports.Row{portARow, staleTestRow, orphanRow})
		if len(diags) != 2 ||
			!anyContains(diags, injectedports.MissingPortTestMessage) || !anyContains(diags, "internal/sample2.PortB") ||
			!anyContains(diags, injectedports.OrphanPortRowMessage) || !anyContains(diags, "internal/sample2.PortC") {
			t.Fatalf("a derived port with a stale test name beside an orphaned row: want one missing-test diagnostic for PortB and one orphan-row diagnostic for PortC, got %v", diags)
		}
	})
}
