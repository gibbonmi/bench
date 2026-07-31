package specbuild

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/spec"
)

type record struct {
	Version      int                   `json:"version"`
	Slug         string                `json:"slug"`
	Spec         string                `json:"spec"`
	Run          string                `json:"run"`
	Branch       string                `json:"branch"`
	Base         string                `json:"base"`
	Candidate    string                `json:"candidate"`
	CandidateTip string                `json:"candidate_tip"`
	Terminal     bool                  `json:"terminal,omitempty"`
	Assignments  map[string]assignment `json:"assignments"`
}

const zeroObjectID = "0000000000000000000000000000000000000000"

type assignment struct {
	ID, Path, Base, Request, Ticket, TicketDigest, Created   string
	Rows, Fence, Assumptions                                 []string
	Checkpoint, CheckpointRef, CheckpointTree, ReceiptDigest string
	CheckpointPatch, Integrated                              string
	DelegatePending                                          bool
}

func (a assignment) public() Assignment {
	return Assignment{ID: a.ID, Path: a.Path, Base: a.Base, Rows: append([]string(nil), a.Rows...), Fence: append([]string(nil), a.Fence...), Assumptions: append([]string(nil), a.Assumptions...)}
}

func (r record) status() Status {
	state, next := "active", "bench spec build assign "+r.Slug
	if r.Terminal {
		state, next = "terminal", ""
	}
	pending := ""
	for _, assigned := range r.Assignments {
		if assigned.DelegatePending && (pending == "" || assigned.ID < pending) {
			pending = assigned.ID
		}
	}
	if pending != "" {
		next = "delegate assignment " + pending
	}
	return Status{Slug: r.Slug, State: state, Subject: r.CandidateTip, Next: next}
}

func (s *Service) resolve(slug string) (string, error) {
	if strings.TrimSpace(slug) == "" {
		return "", errors.New("spec build slug is required")
	}
	_, resolved, _, ok, err := spec.Resolve(s.root, slug)
	if err != nil {
		return "", fmt.Errorf("resolve spec: %w", err)
	}
	if !ok {
		return "", errors.New("spec build spec does not exist")
	}
	return resolved, nil
}

func (s *Service) lock(slug string) (func(), error) {
	path, err := s.statePath(slug)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("create spec build state directory: %w", err)
	}
	f, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open spec build lock: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("lock spec build: %w", err)
	}
	return func() { _ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN); _ = f.Close() }, nil
}

func (s *Service) load(slug string) (record, bool, error) {
	path, err := s.statePath(slug)
	if err != nil {
		return record{}, false, err
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return record{}, false, nil
	}
	if err != nil {
		return record{}, false, fmt.Errorf("read spec build state: %w", err)
	}
	var run record
	if err := jsonfile.Decode(b, &run); err != nil || !run.valid(slug) {
		return record{}, false, errors.New("spec build has incomplete prior state")
	}
	return run, true, nil
}

func (s *Service) save(run record) error {
	path, err := s.statePath(run.Slug)
	if err != nil {
		return err
	}
	b, err := json.Marshal(run)
	if err != nil {
		return fmt.Errorf("encode spec build state: %w", err)
	}
	return replaceState(path, append(b, '\n'))
}

func (r record) valid(slug string) bool {
	return r.Version == 1 && r.Slug == slug && r.Spec != "" && r.Run == digest(r.Spec) && r.Branch != "" && r.Base != "" && r.Candidate == "refs/bench/specbuild/candidate/"+digest(r.Spec) && r.CandidateTip != "" && r.Assignments != nil
}

func replaceState(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".specbuild-*")
	if err != nil {
		return fmt.Errorf("create spec build state replacement: %w", err)
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err := os.Chmod(name, 0o600); err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil || tmp.Sync() != nil || tmp.Close() != nil {
		_ = tmp.Close()
		return errors.New("write spec build state")
	}
	if err := os.Rename(name, path); err != nil {
		return fmt.Errorf("replace spec build state: %w", err)
	}
	installed, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open installed spec build state: %w", err)
	}
	if err := installed.Sync(); err != nil {
		_ = installed.Close()
		return fmt.Errorf("sync installed spec build state: %w", err)
	}
	if err := installed.Close(); err != nil {
		return fmt.Errorf("close installed spec build state: %w", err)
	}
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("open spec build state directory: %w", err)
	}
	defer dir.Close()
	if err := dir.Sync(); err != nil {
		return fmt.Errorf("sync spec build state directory: %w", err)
	}
	return nil
}

func (s *Service) statePath(slug string) (string, error) {
	common, err := benchgit.Output("-C", s.root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("resolve common Git directory: %w", err)
	}
	return filepath.Join(common, "bench", "specbuild", digest(slug)+".json"), nil
}

func workingSubject(root string) (string, string, error) {
	branch, err := benchgit.Output("-C", root, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil || branch == "" {
		return "", "", errors.New("spec build start requires a checked-out working branch")
	}
	dirty, err := benchgit.Output("-C", root, "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return "", "", err
	}
	if dirty != "" {
		return "", "", errors.New("spec build start requires a clean working checkout")
	}
	tip, err := benchgit.Output("-C", root, "rev-parse", "HEAD^{commit}")
	if err != nil {
		return "", "", err
	}
	return branch, tip, nil
}

func updateRef(root, ref, new, old string) error {
	args := []string{"-C", root, "update-ref", ref, new}
	if old != "" {
		args = append(args, old)
	}
	if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func refAt(root, ref, want string) bool {
	got, err := benchgit.Output("-C", root, "rev-parse", "--verify", ref+"^{commit}")
	return err == nil && got == want
}

func refAbsent(root, ref string) (bool, error) {
	err := exec.Command("git", "-C", root, "show-ref", "--verify", "--quiet", ref).Run()
	if err == nil {
		return false, nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 {
		return true, nil
	}
	return false, fmt.Errorf("inspect candidate identity: %w", err)
}

func digest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
