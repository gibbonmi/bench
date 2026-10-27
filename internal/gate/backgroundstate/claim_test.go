package backgroundstate_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
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
	exited := exec.Command("true")
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
