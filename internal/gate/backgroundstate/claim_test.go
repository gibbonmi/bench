package backgroundstate_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// key is one grade key of the fixtures.
const key = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

// currentBoot answers the boot id of this boot.
func currentBoot(t *testing.T) string {
	t.Helper()
	boot, err := protection.BootID()
	if err != nil {
		t.Fatal(err)
	}
	return boot
}

// fields answers the top-level fields of the JSON record at path.
func fields(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	return got
}

// TestClaimNamesKeyWorktreeOwnerAndBoot is BG32. The claim of a running grade names its key,
// its worktree, its owner pid, and the boot id, and it reads back unchanged.
func TestClaimNamesKeyWorktreeOwnerAndBoot(t *testing.T) {
	gitdir, worktree := t.TempDir(), "/work tree/grade"
	claim, err := backgroundstate.ClaimFor(key, worktree)
	if err != nil {
		t.Fatal(err)
	}
	if err := backgroundstate.WriteClaim(durablefile.Native(), gitdir, claim); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"key": key, "worktree": worktree, "owner_pid": float64(os.Getpid()), "boot_id": currentBoot(t)}
	if got := fields(t, backgroundstate.ClaimPath(gitdir, key)); !reflect.DeepEqual(got, want) {
		t.Fatalf("claim fields = %v, want %v", got, want)
	}
	read, err := backgroundstate.ReadClaim(gitdir, key)
	if err != nil || read != claim {
		t.Fatalf("ReadClaim = %+v, %v; want %+v", read, err, claim)
	}
}

// TestClaimOfADeadOwnerIsNotLive is BG33. The claim of this live process is live, and the same
// claim with the pid of an exited process is not.
func TestClaimOfADeadOwnerIsNotLive(t *testing.T) {
	exited := processtest.Helper(t, "exit")
	if err := exited.Run(); err != nil {
		t.Fatal(err)
	}
	claim, err := backgroundstate.ClaimFor(key, "worktree")
	if err != nil {
		t.Fatal(err)
	}
	if live, err := claim.Live(); err != nil || !live {
		t.Fatalf("the claim of this process: live = %v, %v; want live", live, err)
	}
	claim.Owner = exited.Process.Pid
	if live, err := claim.Live(); err != nil || live {
		t.Fatalf("the claim of exited pid %d: live = %v, %v; want not live", claim.Owner, live, err)
	}
}

// TestClaimOfAnotherBootIsNotLive is BG34. A live pid under another boot id can name a new
// process, so the claim is not live.
func TestClaimOfAnotherBootIsNotLive(t *testing.T) {
	claim, err := backgroundstate.ClaimFor(key, "worktree")
	if err != nil {
		t.Fatal(err)
	}
	if live, err := claim.Live(); err != nil || !live {
		t.Fatalf("the claim of this boot: live = %v, %v; want live", live, err)
	}
	claim.Boot = "another-boot"
	if live, err := claim.Live(); err != nil || live {
		t.Fatalf("the claim of another boot: live = %v, %v; want not live", live, err)
	}
}

// TestReleaseClaimRemovesOnlyItsOwnClaim guards the release of a claim. A grade releases the
// claim that it holds, and it keeps a claim of the same key that another owner holds.
func TestReleaseClaimRemovesOnlyItsOwnClaim(t *testing.T) {
	gitdir := t.TempDir()
	claim, err := backgroundstate.ClaimFor(key, "worktree")
	if err != nil {
		t.Fatal(err)
	}
	other := claim
	other.Owner = claim.Owner + 1
	if err := backgroundstate.WriteClaim(durablefile.Native(), gitdir, other); err != nil {
		t.Fatal(err)
	}
	if err := backgroundstate.ReleaseClaim(gitdir, claim); err != nil {
		t.Fatal(err)
	}
	if read, err := backgroundstate.ReadClaim(gitdir, key); err != nil || read != other {
		t.Fatalf("after the release of another owner's claim: ReadClaim = %+v, %v; want %+v", read, err, other)
	}
	if err := backgroundstate.WriteClaim(durablefile.Native(), gitdir, claim); err != nil {
		t.Fatal(err)
	}
	if err := backgroundstate.ReleaseClaim(gitdir, claim); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(backgroundstate.ClaimPath(gitdir, key)); !os.IsNotExist(err) {
		t.Fatalf("after the release of its own claim: Lstat = %v; want no claim", err)
	}
}

// TestWriteClaimRefusesADirectoryThatIsNotPrivate guards the record directory. A claim
// directory that another user can read, or a file at its path, is not private, so the claim
// is not written.
func TestWriteClaimRefusesADirectoryThatIsNotPrivate(t *testing.T) {
	claim, err := backgroundstate.ClaimFor(key, "worktree")
	if err != nil {
		t.Fatal(err)
	}
	for name, create := range map[string]func(path string) error{
		"open directory": func(path string) error {
			if err := os.Mkdir(path, 0o700); err != nil {
				return err
			}
			return os.Chmod(path, 0o755)
		},
		"file": func(path string) error { return os.WriteFile(path, nil, 0o700) },
	} {
		gitdir := t.TempDir()
		if err := create(filepath.Dir(backgroundstate.ClaimPath(gitdir, key))); err != nil {
			t.Fatal(err)
		}
		err := backgroundstate.WriteClaim(durablefile.Native(), gitdir, claim)
		if err == nil || !strings.Contains(err.Error(), "is not a private directory") {
			t.Errorf("%s: WriteClaim = %v, want the refusal of a directory that is not private", name, err)
		}
	}
}

// TestReadClaimRefusesAClaimOfAnotherKey guards the key of a claim. A record at the path of
// one key that names another key does not read as the claim of either key.
func TestReadClaimRefusesAClaimOfAnotherKey(t *testing.T) {
	gitdir, other := t.TempDir(), strings.Repeat("f", len(key))
	claim, err := backgroundstate.ClaimFor(other, "worktree")
	if err != nil {
		t.Fatal(err)
	}
	if err := backgroundstate.WriteClaim(durablefile.Native(), gitdir, claim); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(backgroundstate.ClaimPath(gitdir, other), backgroundstate.ClaimPath(gitdir, key)); err != nil {
		t.Fatal(err)
	}
	read, err := backgroundstate.ReadClaim(gitdir, key)
	if err == nil || !strings.Contains(err.Error(), "the claim names another key") {
		t.Fatalf("ReadClaim of a claim of another key = %+v, %v; want the refusal", read, err)
	}
	if read != (backgroundstate.Claim{}) {
		t.Fatalf("ReadClaim answered %+v with its refusal, want the empty claim", read)
	}
}
