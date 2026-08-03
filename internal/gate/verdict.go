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

	"github.com/gibbonmi/bench/internal/conformance/registry"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

const (
	verdictSchema = 1
	cacheLimit    = 16 * 1024
	manifestLimit = 16 * 1024
	freshness     = 60 * time.Minute
	policyVersion = "oracle-v2/freshness-v1/prospective-v1"
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
	State          State
	Status         string
	PendingStatus  string
	CachedTree     string
	CurrentTree    string
	RecordedAt     time.Time
	Reason         string
	CacheBytes     int
	ReusableGreen  bool
	Partition      *Partition
	CheckPartition *CheckPartition
}

// ComponentSkip is one entry of a partition: a component that did not run, and the evidence
// that covered it. The two evidence forms are alternatives — an ancestor slot's
// content-address identity with the time that slot was authored, or the source digest of the
// seal a reused build artifact was published under — so exactly one of them is populated.
type ComponentSkip struct {
	Component  string
	Identity   string
	AuthoredAt time.Time
	Seal       string
}

// Partition is what a partial verdict graded: the components the run executed for itself,
// and the components it skipped with the evidence that stood in for each. A consumer reads
// it to render the run's narrowness or to refuse a release for it, so the evidence travels
// with the names rather than being looked up again against a store that has since moved.
type Partition struct {
	Executed []string
	Skipped  []ComponentSkip
}

// CheckPartition is the conformance work a mixed run executed and inherited.
type CheckPartition struct {
	Executed  []string
	Inherited []ComponentSkip
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

	Executed     []string                `json:"executed,omitempty"`
	Skipped      []string                `json:"skipped,omitempty"`
	SkipEvidence map[string]skipEvidence `json:"skip_evidence,omitempty"`

	CheckExecuted  []string                `json:"check_executed,omitempty"`
	CheckInherited []string                `json:"check_inherited,omitempty"`
	CheckEvidence  map[string]skipEvidence `json:"check_evidence,omitempty"`
}

// skipEvidence is one component's recorded skip evidence. Every field is optional in the
// struct and exact in the bytes: an entry is graded against one of two exact field sets, so
// an absent field is always a refusal here rather than a zero value some reader defaults.
type skipEvidence struct {
	Identity   string `json:"identity,omitempty"`
	AuthoredAt string `json:"authored_at,omitempty"`
	Seal       string `json:"seal,omitempty"`
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
	loaded := loadVerdict(filepath.Join(gitdir, benchgit.GateCacheFile), now)
	gi.State, gi.CacheBytes = loaded.state, loaded.bytes
	if loaded.reason != "" {
		gi.Reason = loaded.reason
	}
	if gi.State != Ready && gi.State != Pending {
		return gi
	}
	rec := loaded.record
	gi.Status, gi.CachedTree = rec.Status, rec.Tree
	gi.Partition = rec.partition()
	gi.CheckPartition = rec.checkPartition()
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
	// Checked after drift and expiry: those retire a narrow record exactly as they retire a
	// full one, and naming the narrowness of an expired record would dress retired evidence
	// as current. The Partition on the inspection carries the narrowness either way.
	if reason := narrowVerdictReason(rec); reason != "" {
		gi.Reason = reason
		return gi
	}
	gi.ReusableGreen = true
	return gi
}

// narrowVerdictReason names the class of a verdict that graded less than the whole tree, and
// returns "" for one that graded all of it. A partial verdict ran only the components whose
// inputs moved. It is evidence about its own tree and never the whole-tree green a reuse
// credits, so this is the single place a reuse asks how wide the grading was.
func narrowVerdictReason(r verdictRecord) string {
	if r.partitions() || r.checkPartitions() {
		return "partial verdict"
	}
	return ""
}

type loadedVerdict struct {
	record verdictRecord
	state  State
	reason string
	bytes  int
}

// storeRecordBytes is one record file read from the evidence store. data is non-nil only
// when the file cleared every check that holds whatever class the bytes turn out to name;
// otherwise state and reason say why the store has nothing readable there.
type storeRecordBytes struct {
	data   []byte
	bytes  int
	state  State
	reason string
}

// readStoreRecord applies the file discipline every record class in the store shares: a
// regular 0600 file, within the size cap, framed as a single JSON object. What the bytes
// mean is the reading class's question — this answers only whether there are bytes worth
// asking about, so a class added to the store cannot be given a laxer file than the
// verdict cache gets.
func readStoreRecord(path string) storeRecordBytes {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return storeRecordBytes{state: Absent}
	}
	if err != nil {
		return storeRecordBytes{state: Unavailable, reason: "cache metadata unavailable"}
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		return storeRecordBytes{state: Invalid, reason: "invalid cache metadata"}
	}
	read := storeRecordBytes{bytes: int(info.Size())}
	f, err := os.Open(path)
	if err != nil {
		read.state, read.reason = Unavailable, "cache unavailable"
		return read
	}
	data, readErr := io.ReadAll(io.LimitReader(f, cacheLimit+1))
	closeErr := f.Close()
	if readErr != nil || closeErr != nil {
		read.state, read.reason = Unavailable, "cache unavailable"
		return read
	}
	if len(data) == 0 || len(data) > cacheLimit {
		read.state, read.reason = Invalid, "invalid cache record"
		return read
	}
	if data[0] != '{' || (data[len(data)-1] != '}' && (data[len(data)-1] != '\n' || len(data) < 2 || data[len(data)-2] != '}')) {
		read.state, read.reason = Invalid, "invalid cache framing"
		return read
	}
	read.data = data
	return read
}

func loadVerdict(path string, now time.Time) loadedVerdict {
	read := readStoreRecord(path)
	loaded := loadedVerdict{state: read.state, reason: read.reason, bytes: read.bytes}
	if read.data == nil {
		return loaded
	}
	if err := strictJSON(read.data, &loaded.record); err != nil || validateRecordBytes(read.data, loaded.record, now) != nil {
		loaded.state, loaded.reason = Invalid, "invalid cache record"
		return loaded
	}
	loaded.state = loaded.record.State
	return loaded
}

// The two rejections strictJSON adds on top of encoding/json are sentinels so a caller
// can class them without matching on message text.
var (
	errTrailingJSON      = errors.New("trailing JSON")
	errDuplicateJSONName = errors.New("duplicate name")
)

func strictJSON(data []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	var trailing any
	if err := dec.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errTrailingJSON
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
				return errDuplicateJSONName
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

// The two exact field sets a ready record may carry. They are alternatives, never a
// spectrum: a record holding part of one and part of another names no class the loader
// can resolve, and resolving it by guess would credit work that nobody ran. The narrow
// sets are derived — the full set plus that class's own fields — so a field added to the
// full record joins them without a second edit; restated, that addition would make the
// narrow classes silently reject every valid record.
var (
	fullReadyFields            = []string{"oracle", "recorded_at", "schema", "state", "status", "tree"}
	partialReadyFields         = sortedFieldSet(fullReadyFields, "executed", "skip_evidence", "skipped")
	checkPartialReadyFields    = sortedFieldSet(fullReadyFields, "check_evidence", "check_executed", "check_inherited")
	combinedPartialReadyFields = sortedFieldSet(partialReadyFields, "check_evidence", "check_executed", "check_inherited")
)

// readyFieldClasses is the one place every ready verdict class is enumerated, keyed by the
// name storeRecordClasses (record_classes.go) reports it under. A *ReadyFields variable above
// that never joins this map is caught by TestVerdictReadyFieldsAreAllRegistered, which parses
// this file's declarations and fails for any it does not find registered here.
var readyFieldClasses = map[string][]string{
	"full verdict":             fullReadyFields,
	"partial verdict":          partialReadyFields,
	"check-partial verdict":    checkPartialReadyFields,
	"combined-partial verdict": combinedPartialReadyFields,
}

// The two exact field sets one skip-evidence entry may carry, under the same alternatives
// discipline as the record classes: an ancestor slot's identity with the time it was
// authored, or a reused build's seal digest. An entry holding parts of both describes no
// evidence, and reading it as either would credit a component on a proof it never named.
var (
	ancestorEvidenceFields = []string{"authored_at", "identity"}
	sealEvidenceFields     = []string{"seal"}
)

// sortedFieldSet joins a base field set with extras in the sorted order
// requireObjectFields compares against.
func sortedFieldSet(base []string, extra ...string) []string {
	fields := append(slices.Clone(base), extra...)
	sort.Strings(fields)
	return fields
}

// partitions reports whether the record reaches for the partial class at all: any single
// partial field measures the record against the whole partial set, so a fragment is
// refused for what it is missing rather than read as a wider class carrying something
// extra. A retired class's fields (the reduced verdict's reduced/phases/ancestor family)
// are unknown fields now, so a legacy reduced record fails the exact-field-set validation
// below and reads as non-reusable rather than as a full green.
func (r verdictRecord) partitions() bool {
	return r.Executed != nil || r.Skipped != nil || r.SkipEvidence != nil
}

func (r verdictRecord) checkPartitions() bool {
	return r.CheckExecuted != nil || r.CheckInherited != nil || r.CheckEvidence != nil
}

// partition projects the partial class onto the shape consumers read, and nil for every
// other class — so "this verdict skipped something" is one nil check rather than a marker a
// consumer could read without the evidence that explains it. It reads a record the loader
// has already validated, where every skipped component has exactly one evidence entry.
func (r verdictRecord) partition() *Partition {
	if !r.partitions() {
		return nil
	}
	partition := &Partition{Executed: slices.Clone(r.Executed), Skipped: make([]ComponentSkip, 0, len(r.Skipped))}
	for _, component := range r.Skipped {
		evidence := r.SkipEvidence[component]
		authoredAt, _ := time.Parse(time.RFC3339, evidence.AuthoredAt)
		partition.Skipped = append(partition.Skipped, ComponentSkip{
			Component:  component,
			Identity:   evidence.Identity,
			AuthoredAt: authoredAt,
			Seal:       evidence.Seal,
		})
	}
	return partition
}

func (r verdictRecord) checkPartition() *CheckPartition {
	if !r.checkPartitions() {
		return nil
	}
	partition := &CheckPartition{Executed: slices.Clone(r.CheckExecuted), Inherited: make([]ComponentSkip, 0, len(r.CheckInherited))}
	for _, check := range r.CheckInherited {
		evidence := r.CheckEvidence[check]
		authoredAt, _ := time.Parse(time.RFC3339, evidence.AuthoredAt)
		partition.Inherited = append(partition.Inherited, ComponentSkip{
			Component: check, Identity: evidence.Identity, AuthoredAt: authoredAt,
		})
	}
	return partition
}

// validatePartition grades what a partial record claims: which components the run executed
// for itself, and which evidence covered each one it skipped. The two collections are
// cross-checked in both directions — a skip with no evidence credits a component nobody
// graded, and evidence naming a component the record did not skip is a proof of nothing the
// record claims.
//
// Both lists are sorted and duplicate-free so one partition has exactly one spelling: reuse
// compares records byte for byte, and two orderings of one partition would read as two
// different verdicts of one run.
func validatePartition(data []byte, r verdictRecord, recordedAt time.Time) error {
	if err := requireComponentSet(r.Executed); err != nil {
		return err
	}
	if err := requireComponentSet(r.Skipped); err != nil {
		return err
	}
	for _, component := range r.Skipped {
		if slices.Contains(r.Executed, component) {
			return errors.New("invalid partition")
		}
	}
	entries, err := rawSkipEvidence(data)
	if err != nil {
		return err
	}
	// Equal cardinality closes the second direction: every skipped component is looked up
	// below, so a map no larger than the skipped set holds evidence for nothing else.
	if len(entries) != len(r.Skipped) {
		return errors.New("invalid skip evidence")
	}
	for _, component := range r.Skipped {
		entry, ok := entries[component]
		if !ok {
			return errors.New("invalid skip evidence")
		}
		if err := validateSkipEvidence(entry, r.SkipEvidence[component], recordedAt); err != nil {
			return err
		}
	}
	return nil
}

// requireComponentSet holds a partition's halves to a non-empty, strictly ascending list of
// named components. Strict ascent carries the duplicate refusal: a component named twice in
// one half says nothing a reader can act on and would let two partitions share one spelling.
func requireComponentSet(components []string) error {
	if len(components) == 0 {
		return errors.New("invalid component set")
	}
	for i, component := range components {
		if component == "" || (i > 0 && components[i-1] >= component) {
			return errors.New("invalid component set")
		}
	}
	return nil
}

// validateSkipEvidence grades one entry against exactly one of the two evidence forms. The
// seal form is dispatched on its own field so an entry reaching for both is measured against
// the seal set and refused for the ancestor fields it also carries, rather than read as
// whichever form the reader checked first.
func validateSkipEvidence(entry json.RawMessage, evidence skipEvidence, recordedAt time.Time) error {
	if evidence.Seal != "" {
		if requireObjectFields(entry, sealEvidenceFields) != nil || !isContentAddress(evidence.Seal) {
			return errors.New("invalid skip evidence")
		}
		return nil
	}
	if requireObjectFields(entry, ancestorEvidenceFields) != nil || !isContentAddress(evidence.Identity) {
		return errors.New("invalid skip evidence")
	}
	// A slot is authored when its component runs green, so its authorship precedes every run
	// that reads it. A later time is evidence written after the run it is claimed to cover.
	if tm, err := strictRecordTime(evidence.AuthoredAt); err != nil || tm.After(recordedAt) {
		return errors.New("invalid skip evidence")
	}
	return nil
}

// rawSkipEvidence returns the record's evidence entries as unparsed objects, which is what
// grading an entry against an exact field set needs: the decoded struct cannot tell a field
// that was absent from one that was present and empty.
func rawSkipEvidence(data []byte) (map[string]json.RawMessage, error) {
	return rawEvidence(data, "skip_evidence")
}

func rawEvidence(data []byte, field string) (map[string]json.RawMessage, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	var entries map[string]json.RawMessage
	if err := json.Unmarshal(fields[field], &entries); err != nil {
		return nil, errors.New("invalid skip evidence")
	}
	return entries, nil
}

func validateCheckPartition(data []byte, r verdictRecord, recordedAt time.Time) error {
	wantExecuted := make(map[string]bool, len(r.CheckExecuted))
	for _, name := range r.CheckExecuted {
		wantExecuted[name] = true
	}
	wantInherited := make(map[string]bool, len(r.CheckInherited))
	for _, name := range r.CheckInherited {
		wantInherited[name] = true
	}
	var executed, inherited []string
	for _, check := range registry.Checks {
		if !check.RunsAt(registry.Dev) {
			continue
		}
		if check.Meta || wantExecuted[check.Name] {
			executed = append(executed, check.Name)
		}
		if !check.Meta && wantInherited[check.Name] {
			inherited = append(inherited, check.Name)
		}
	}
	if !slices.Equal(r.CheckExecuted, executed) || !slices.Equal(r.CheckInherited, inherited) || len(r.CheckInherited) == 0 {
		return errors.New("invalid conformance check partition")
	}
	seen := make(map[string]bool, len(r.CheckExecuted)+len(r.CheckInherited))
	for _, name := range r.CheckExecuted {
		if seen[name] {
			return errors.New("invalid conformance check partition")
		}
		seen[name] = true
	}
	for _, name := range r.CheckInherited {
		if seen[name] {
			return errors.New("invalid conformance check partition")
		}
		seen[name] = true
	}
	if len(seen) != len(registry.Names(registry.Dev)) {
		return errors.New("invalid conformance check partition")
	}
	entries, err := rawEvidence(data, "check_evidence")
	if err != nil || len(entries) != len(r.CheckInherited) {
		return errors.New("invalid conformance check evidence")
	}
	for _, name := range r.CheckInherited {
		entry, ok := entries[name]
		if !ok || validateSkipEvidence(entry, r.CheckEvidence[name], recordedAt) != nil || r.CheckEvidence[name].Seal != "" {
			return errors.New("invalid conformance check evidence")
		}
	}
	return nil
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
		want := fullReadyFields
		switch {
		case r.partitions() && r.checkPartitions():
			want = combinedPartialReadyFields
		case r.partitions():
			want = partialReadyFields
		case r.checkPartitions():
			want = checkPartialReadyFields
		}
		if err := requireObjectFields(data, want); err != nil {
			return err
		}
		if (r.Status != "green" && r.Status != "red" && r.Status != "timeout") || r.RecordedAt == "" || r.StartedAt != "" || r.OwnerPID != 0 {
			return errors.New("invalid ready")
		}
		tm, err := strictRecordTime(r.RecordedAt)
		if err != nil || tm.After(now) {
			return errors.New("invalid ready time")
		}
		if r.partitions() {
			if err := validatePartition(data, r, tm); err != nil {
				return err
			}
		}
		if r.checkPartitions() {
			return validateCheckPartition(data, r, tm)
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
