package preflight

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

type proposalAuthoritySnapshot struct {
	headTree, index, status, headBytes, worktreeBytes, assignment []byte
}

func commandBytes(t *testing.T, name string, args ...string) []byte {
	t.Helper()
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
	return out
}

func proposalAuthorityState(t *testing.T, root string) proposalAuthoritySnapshot {
	t.Helper()
	status := commandBytes(t, "git", "-C", root, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	paths := bytes.Split(commandBytes(t, "git", "-C", root, "ls-files", "-z"), []byte{0})
	var head, worktree bytes.Buffer
	for _, raw := range paths {
		if len(raw) == 0 {
			continue
		}
		path := string(raw)
		committed := commandBytes(t, "git", "-C", root, "show", "HEAD:"+path)
		live, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatalf("read tracked path %s: %v", path, err)
		}
		fmt.Fprintf(&head, "%d:%s%d:", len(path), path, len(committed))
		head.Write(committed)
		fmt.Fprintf(&worktree, "%d:%s%d:", len(path), path, len(live))
		worktree.Write(live)
	}
	indexPath := strings.TrimSpace(string(commandBytes(t, "git", "-C", root, "rev-parse", "--git-path", "index")))
	index, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	assignmentPath, err := intent.Address(root)
	if err != nil {
		t.Fatalf("assignment address: %v", err)
	}
	assignment, err := os.ReadFile(assignmentPath)
	if err != nil {
		t.Fatalf("read assignment authority: %v", err)
	}
	return proposalAuthoritySnapshot{
		headTree:      commandBytes(t, "git", "-C", root, "ls-tree", "-r", "--full-tree", "HEAD"),
		index:         index,
		status:        status,
		headBytes:     head.Bytes(),
		worktreeBytes: worktree.Bytes(),
		assignment:    assignment,
	}
}

// TestWritesProposalReadOnly covers DP11 through the public command.
func TestWritesProposalReadOnly(t *testing.T) {
	root, slug := seedProposal(t)
	args := proposalArgs(t, root, slug)
	warm := proposalAuthorityState(t, root)
	before := proposalAuthorityState(t, root)
	if !reflect.DeepEqual(warm, before) {
		t.Fatal("authority observer changed the state it measures")
	}
	first, firstCode := Command(args)
	afterFirst := proposalAuthorityState(t, root)
	second, secondCode := Command(args)
	afterSecond := proposalAuthorityState(t, root)
	if firstCode != 0 || secondCode != 0 || first != second {
		t.Fatalf("repeat proposal = (%d, %d, equal=%t)", firstCode, secondCode, first == second)
	}
	if !reflect.DeepEqual(before, afterFirst) || !reflect.DeepEqual(before, afterSecond) {
		t.Fatalf("proposal changed tracked, index, status, or assignment authority state")
	}
}
