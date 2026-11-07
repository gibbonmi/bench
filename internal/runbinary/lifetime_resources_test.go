package runbinary

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestMain removes an outer run's resource vector, so this package owns the only vector. A
// builder then registers only the selection's resources, under the protocol that its test names.
func TestMain(m *testing.M) {
	if err := os.Unsetenv(protection.Env); err != nil {
		os.Exit(3)
	}
	os.Exit(m.Run())
}

// binaryProtocol admits the selection kind alone, as the composition root's binary row does.
func binaryProtocol() *protection.Protocol {
	return &protection.Protocol{Validators: protection.Validators{protection.KindBinary: ValidateProtection}, ReadLimit: bounds.ControlRecordLimit}
}

// selectedFactory builds a selection that writes a stub executable under the binary protocol.
func selectedFactory(t *testing.T) Factory {
	return Factory{
		TempRoot:   t.TempDir(),
		Protection: binaryProtocol(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	}
}

// TestSelectionValidatorRefusesAForeignBinding pins the selection validator that the composition
// root installs: it accepts a real selection's descriptor and refuses the same descriptor when its
// binding names another directory at the same path.
func TestSelectionValidatorRefusesAForeignBinding(t *testing.T) {
	t.Parallel()
	factory := selectedFactory(t)
	selection, err := factory.Own(context.Background(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := selection.Close(); err != nil {
			t.Error(err)
		}
	})
	resources := selection.Resources()
	if len(resources) != 1 {
		t.Fatalf("selection resources = %+v, want one binary generation", resources)
	}
	d := resources[0]
	if err := ValidateProtection(d); err != nil {
		t.Fatalf("ValidateProtection(own selection) = %v, want nil", err)
	}
	d.Binding = "0:0"
	if err := ValidateProtection(d); err == nil || !strings.Contains(err.Error(), "binding 0:0 does not name selection") {
		t.Fatalf("ValidateProtection(foreign binding) = %v, want a binding refusal", err)
	}
}

// TestBuildScriptMalformedResources is RS23. A malformed inherited vector refuses the build
// script before it starts, so the script's marker never appears.
func TestBuildScriptMalformedResources(t *testing.T) {
	t.Setenv(protection.Env, "{not a vector")
	source := t.TempDir()
	marker := filepath.Join(t.TempDir(), "started")
	if err := os.MkdirAll(filepath.Join(source, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "scripts", "go-build.sh"), []byte("#!/usr/bin/env bash\ntouch \""+marker+"\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := Factory{TempRoot: t.TempDir(), Protection: binaryProtocol()}.Own(context.Background(), source)
	if err == nil || !strings.Contains(err.Error(), protection.Env) {
		t.Fatalf("Own under a malformed vector = %v, want the vector refusal", err)
	}
	if _, statErr := os.Stat(marker); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("the build script ran under a malformed vector: %v", statErr)
	}
}

// TestBinaryClose is PL36. A selection user exits and leaves a live descendant in its
// registered group, so Close must keep the executable and report the retention.
func TestBinaryClose(t *testing.T) {
	t.Parallel()
	factory := selectedFactory(t)
	selection, err := factory.Own(context.Background(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if resources := selection.Resources(); len(resources) != 1 || ValidateProtection(resources[0]) != nil {
		t.Fatalf("selection resources = %+v, want one valid binary generation", resources)
	}
	pidFile := filepath.Join(t.TempDir(), "descendant")
	user := exec.Command("sh", "-c", `sleep 300 </dev/null >/dev/null 2>&1 & echo $! > "$1.partial" && mv "$1.partial" "$1"`, "sh", pidFile)
	policy := subprocess.Policy{Normal: subprocess.Preserve, Final: bounds.ProcessWindow(bounds.FixedWindow(bounds.ProcessFinalWindow)), Poll: bounds.ProcessPoll}
	lifetime := subprocess.NewLifetime(user, policy, subprocess.Streams{}, subprocess.Registration{Protocol: factory.Protection, Resources: selection.Resources()})
	if err := lifetime.Start(); err != nil {
		t.Fatal(err)
	}
	outcome, _ := lifetime.Wait()
	processtest.WaitPID(t, pidFile)
	if outcome.Process.Code != 0 || outcome.Protection.State != protection.Protected {
		t.Fatalf("user outcome = %+v, want exit 0 with its live group protected", outcome)
	}
	closeErr := selection.Close()
	if closeErr == nil || !strings.Contains(closeErr.Error(), string(protection.Protected)) {
		t.Fatalf("Close with a live registered group = %v, want a protected retention", closeErr)
	}
	if info, err := os.Stat(selection.Path); err != nil || !info.Mode().IsRegular() {
		t.Fatalf("retained executable = %v, %v; want the selected file", info, err)
	}
}
