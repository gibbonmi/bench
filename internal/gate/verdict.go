package gate

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"syscall"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

const (
	verdictSchema = 1
	cacheLimit    = 16 * 1024
	manifestLimit = 16 * 1024
	freshness     = 10 * time.Minute
	policyVersion = "oracle-v1/freshness-v1"
)

type State string

const (
	Absent      State = "absent"
	Ready       State = "ready"
	Pending     State = "pending"
	Invalid     State = "invalid"
	Unavailable State = "unavailable"
)

type Inspection struct {
	State         State
	Status        string
	PendingStatus string
	CachedTree    string
	CurrentTree   string
	RecordedAt    time.Time
	Reason        string
	CacheBytes    int
	ReusableGreen bool
}

type verdictRecord struct {
	Schema     int    `json:"schema"`
	State      State  `json:"state"`
	Status     string `json:"status,omitempty"`
	Tree       string `json:"tree"`
	Oracle     string `json:"oracle"`
	RecordedAt string `json:"recorded_at,omitempty"`
	StartedAt  string `json:"started_at,omitempty"`
	OwnerPID   int    `json:"owner_pid,omitempty"`
}

type manifest struct {
	Schema      int      `json:"schema"`
	Closure     string   `json:"closure"`
	Environment []string `json:"environment"`
	Paths       []string `json:"paths"`
	Tools       []string `json:"tools"`
}

type subject struct {
	Tree, Oracle string
	Closed       bool
	Reason       string
	Resolution   Resolution
	Env          []string
}

type gateFile interface {
	Name() string
	Chmod(os.FileMode) error
	Write([]byte) (int, error)
	Sync() error
	Close() error
	Fd() uintptr
}

func durableReplace(gitdir string, rec verdictRecord) error {
	return durableReplaceWithEngine(productionGateEngine{}, gitdir, rec)
}

func persistInterruptedIfGreen(engine gateEngine, root, gitdir string, plan subject) {
	if !inspectAt(root, engine.Now()).ReusableGreen {
		return
	}
	pending := interruptedRecord(plan, engine.Now())
	if err := durableReplaceWithEngine(engine, gitdir, pending); err != nil {
		_ = durableReplaceWithEngine(engine, gitdir, pending)
	}
}

func durableReplaceWithEngine(engine gateEngine, gitdir string, rec verdictRecord) error {
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := engine.CreateTemp(gitdir, ".bench-last-gate-")
	if err != nil {
		return err
	}
	name, installed := tmp.Name(), false
	defer func() {
		_ = tmp.Close()
		if !installed {
			_ = os.Remove(name)
		}
	}()
	if err := tmp.Chmod(0o600); err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := engine.Rename(name, filepath.Join(gitdir, benchgit.GateCacheFile)); err != nil {
		return err
	}
	dir, err := engine.OpenDir(gitdir)
	if err != nil {
		return err
	}
	syncErr, closeErr := dir.Sync(), dir.Close()
	if syncErr != nil {
		return syncErr
	}
	if closeErr != nil {
		return closeErr
	}
	installed = true
	return nil
}

func Inspect(root string) Inspection { return inspectWithEngine(root, productionGateEngine{}) }

func inspectWithEngine(root string, engine gateEngine) Inspection {
	return inspectAt(root, engine.Now())
}

func inspectAt(root string, now time.Time) Inspection {
	s, err := buildSubject(root)
	if err != nil {
		return Inspection{State: Unavailable, Reason: "subject unavailable"}
	}
	gi := Inspection{CurrentTree: s.Tree, Reason: s.Reason}
	gitdir, err := benchgit.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		gi.State = Unavailable
		gi.Reason = "git directory unavailable"
		return gi
	}
	path := filepath.Join(gitdir, benchgit.GateCacheFile)
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		gi.State = Absent
		return gi
	}
	if err != nil {
		gi.State = Unavailable
		gi.Reason = "cache metadata unavailable"
		return gi
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		gi.State = Invalid
		gi.Reason = "invalid cache metadata"
		return gi
	}
	gi.CacheBytes = int(info.Size())
	f, err := os.Open(path)
	if err != nil {
		gi.State = Unavailable
		gi.Reason = "cache unavailable"
		return gi
	}
	data, readErr := io.ReadAll(io.LimitReader(f, cacheLimit+1))
	closeErr := f.Close()
	if readErr != nil || closeErr != nil {
		gi.State = Unavailable
		gi.Reason = "cache unavailable"
		return gi
	}
	if len(data) == 0 || len(data) > cacheLimit {
		gi.State = Invalid
		gi.Reason = "invalid cache record"
		return gi
	}
	if data[0] != '{' || (data[len(data)-1] != '}' && (data[len(data)-1] != '\n' || len(data) < 2 || data[len(data)-2] != '}')) {
		gi.State = Invalid
		gi.Reason = "invalid cache framing"
		return gi
	}
	var rec verdictRecord
	if err := strictJSON(data, &rec); err != nil || validateRecordBytes(data, rec, now) != nil {
		gi.State = Invalid
		gi.Reason = "invalid cache record"
		return gi
	}
	gi.State, gi.Status, gi.CachedTree = rec.State, rec.Status, rec.Tree
	if rec.State == Pending {
		held, err := lockHeld(gitdir)
		if err != nil {
			gi.State = Unavailable
			gi.Reason = "gate lock unavailable"
			return gi
		}
		if held {
			gi.PendingStatus = "locked-pending"
		} else {
			gi.PendingStatus = "interrupted-pending"
		}
		return gi
	}
	tm, _ := time.Parse(time.RFC3339, rec.RecordedAt)
	gi.RecordedAt = tm
	if !s.Closed {
		return gi
	}
	if rec.Status != "green" {
		gi.Reason = "recorded " + rec.Status
		return gi
	}
	if rec.Tree != s.Tree {
		gi.Reason = "working tree changed"
		return gi
	}
	if rec.Oracle != s.Oracle {
		gi.Reason = "oracle changed"
		return gi
	}
	if now.Sub(tm) >= freshness {
		gi.Reason = "verdict expired"
		return gi
	}
	gi.ReusableGreen = true
	return gi
}

func strictJSON(data []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	var trailing any
	if err := dec.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("trailing JSON")
	}
	// encoding/json accepts duplicate keys; a token walk rejects them recursively.
	dec = json.NewDecoder(bytes.NewReader(data))
	return rejectDuplicateNames(dec)
}

func rejectDuplicateNames(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	d, ok := tok.(json.Delim)
	if !ok {
		return nil
	}
	if d == '{' {
		seen := map[string]bool{}
		for dec.More() {
			k, err := dec.Token()
			if err != nil {
				return err
			}
			key := k.(string)
			if seen[key] {
				return errors.New("duplicate name")
			}
			seen[key] = true
			if err := rejectDuplicateNames(dec); err != nil {
				return err
			}
		}
	} else if d == '[' {
		for dec.More() {
			if err := rejectDuplicateNames(dec); err != nil {
				return err
			}
		}
	}
	_, err = dec.Token()
	return err
}

func validateRecordBytes(data []byte, r verdictRecord, now time.Time) error {
	if r.Schema != verdictSchema || !treeHashRE.MatchString(r.Tree) || len(r.Oracle) != 64 {
		return errors.New("invalid record")
	}
	if _, err := hex.DecodeString(r.Oracle); err != nil || strings.ToLower(r.Oracle) != r.Oracle {
		return errors.New("invalid oracle")
	}
	switch r.State {
	case Ready:
		if err := requireObjectFields(data, []string{"oracle", "recorded_at", "schema", "state", "status", "tree"}); err != nil {
			return err
		}
		if (r.Status != "green" && r.Status != "red" && r.Status != "timeout") || r.RecordedAt == "" || r.StartedAt != "" || r.OwnerPID != 0 {
			return errors.New("invalid ready")
		}
		if tm, err := strictRecordTime(r.RecordedAt); err != nil || tm.After(now) {
			return errors.New("invalid ready time")
		}
	case Pending:
		if err := requireObjectFields(data, []string{"oracle", "owner_pid", "schema", "started_at", "state", "tree"}); err != nil {
			return err
		}
		if r.Status != "" || r.RecordedAt != "" || r.StartedAt == "" || r.OwnerPID <= 0 {
			return errors.New("invalid pending")
		}
		if tm, err := strictRecordTime(r.StartedAt); err != nil || tm.After(now) {
			return errors.New("invalid pending time")
		}
	default:
		return errors.New("invalid state")
	}
	return nil
}

func strictRecordTime(value string) (time.Time, error) {
	tm, err := time.Parse(time.RFC3339, value)
	if err != nil || tm.Nanosecond() != 0 || !strings.HasSuffix(value, "Z") || tm.Format(time.RFC3339) != value {
		return time.Time{}, errors.New("invalid record time")
	}
	return tm, nil
}

func requireObjectFields(data []byte, want []string) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	got := make([]string, 0, len(fields))
	for name := range fields {
		got = append(got, name)
	}
	sort.Strings(got)
	if !slices.Equal(got, want) {
		return errors.New("invalid record fields")
	}
	return nil
}

func lockHeld(gitdir string) (bool, error) {
	path := filepath.Join(gitdir, "bench-gate.lock")
	executionLockOwners.Lock()
	defer executionLockOwners.Unlock()
	if executionLockOwners.paths[path] {
		return true, nil
	}
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer f.Close()
	lock := recordLock(syscall.F_RDLCK)
	if err := syscall.FcntlFlock(f.Fd(), syscall.F_GETLK, &lock); err != nil {
		return false, err
	}
	return lock.Type != syscall.F_UNLCK, nil
}
