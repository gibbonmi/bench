package protection_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestDirectoryRefusesAForeignBinding hands each directory adapter shape the descriptor of its
// own directory carrying the identity of a sibling of the same shape. Only the binding differs,
// so only the owner identity comparison refuses it.
func TestDirectoryRefusesAForeignBinding(t *testing.T) {
	for name, adapter := range map[string]protection.Directory{
		"plain":  {Kind: protection.KindKitTest, Store: "store", Noun: "run"},
		"prefix": {Kind: protection.KindBinary, Store: "store", Noun: "run", Prefix: "owned-", Class: "owned directory"},
		"lock":   {Kind: protection.KindCache, Store: "store", Noun: "run", Lock: "lock"},
	} {
		t.Run(name, func(t *testing.T) {
			parent := t.TempDir()
			own, other := adapterDirectory(t, adapter, parent, "own"), adapterDirectory(t, adapter, parent, "other")
			if err := adapter.Validate(own); err != nil {
				t.Fatalf("the adapter refused its own descriptor %+v: %v", own, err)
			}
			foreign := own
			foreign.Binding = other.Binding
			if err := adapter.Validate(foreign); err == nil {
				t.Fatalf("the adapter accepted the foreign binding %+v", foreign)
			}
		})
	}
}

// adapterDirectory creates one directory that adapter accepts and describes it.
func adapterDirectory(t *testing.T, adapter protection.Directory, parent, name string) protection.Descriptor {
	t.Helper()
	dir := filepath.Join(parent, adapter.Prefix+name)
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if adapter.Lock != "" {
		if err := os.WriteFile(filepath.Join(dir, adapter.Lock), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	d, err := adapter.Describe(dir)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
