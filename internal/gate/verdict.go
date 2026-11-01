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
	"time"

	"github.com/gibbonmi/bench/internal/durablefile"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/jsonfile"
)

const (
	verdictSchema = 1
	cacheLimit    = 16 * 1024
	manifestLimit = 16 * 1024
	freshness     = 60 * time.Minute
	policyVersion = "oracle-v2/freshness-v1/prospective-v2"
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
	// Drifted marks a record whose tree or oracle is not this subject's, whatever verdict
	// it carries. It travels with the inspection so a consumer never decides for itself
	// whether a record describes the tree its reader is on. A red graded elsewhere is not
	// this tree's red.
	Drifted bool
	// Only Inspect fills the two reader facts, so a gate run pays no history read for them.
	// ReaderRecorded marks a current tree that holds a reader record, on a HEAD that descends
	// from the green marker. ReaderLandings is the count of reader-record commits after the
	// green marker up to HEAD, and it is counted only when ReaderRecorded holds.
	ReaderRecorded bool
	ReaderLandings int
}

type verdictRecord struct {
	Schema     int    `json:"schema"`
	State      State  `json:"state,omitempty"`
	Status     string `json:"status,omitempty"`
	Tree       string `json:"tree"`
	Oracle     string `json:"oracle,omitempty"`
	RecordedAt string `json:"recorded_at,omitempty"`
	StartedAt  string `json:"started_at,omitempty"`
	OwnerPID   int    `json:"owner_pid,omitempty"`

	// An evidence record carries the grade key in place of the oracle digest, because the
	// oracle digest frames the checkout path and the grade key does not.
	Grade string `json:"grade,omitempty"`

	// The lane class carries these three and none of the verdict fields above. A lane
	// grades a declared check list rather than the oracle, so it has no state, no
	// status, and no oracle identity to record.
	Lane      string `json:"lane,omitempty"`
	Outcome   string `json:"outcome,omitempty"`
	RunBinary string `json:"run_binary,omitempty"`

	// The delta class carries these two and the run binary. A delta grades the delta set of
	// a test-only change from the base tree, so it has no oracle identity and no status.
	Delta    string `json:"delta,omitempty"`
	BaseTree string `json:"base_tree,omitempty"`

	// The reader class carries these three and the run binary. A reader set grades the readers
	// of a landing's changed families, so it names its schedule, the destination tree, and the
	// green marker commit, and it has no oracle identity and no status.
	ReaderSet       string `json:"reader_set,omitempty"`
	DestinationTree string `json:"destination_tree,omitempty"`
	GreenMarker     string `json:"green_marker,omitempty"`
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
	// Identity is the tree-free oracle identity: every oracle frame but the tree, and then
	// the capability posture. A landing compares it across trees.
	Identity   string
	Closed     bool
	Reason     string
	Resolution Resolution
	Env        []string
}

// durableReplace installs rec as the verdict cache of gitdir through replacer.
func durableReplace(replacer durablefile.Replacer, gitdir string, rec verdictRecord) error {
	return durableReplaceAt(replacer, gitdir, benchgit.GateCacheFile, rec)
}

// storeRecord answers the fields that every record of the store shares: the schema and the
// tree. Each record constructor starts from it.
func storeRecord(tree string) verdictRecord {
	return verdictRecord{Schema: verdictSchema, Tree: tree}
}

// recordTime answers the store form of a record time: UTC, whole seconds, RFC 3339.
func recordTime(at time.Time) string {
	return at.UTC().Truncate(time.Second).Format(time.RFC3339)
}

// readyVerdict answers the gate cache record of a terminal status of plan, recorded at at.
func readyVerdict(plan subject, status string, at time.Time) verdictRecord {
	rec := storeRecord(plan.Tree)
	rec.State, rec.Status, rec.Oracle, rec.RecordedAt = Ready, status, plan.Oracle, recordTime(at)
	return rec
}

func Inspect(root string) Inspection {
	return withReaderFacts(root, inspectAt(root, time.Now().UTC()))
}

func inspectAt(root string, now time.Time) Inspection {
	s, err := buildSubject(root)
	if err != nil {
		return Inspection{State: Unavailable, Reason: "subject unavailable"}
	}
	return inspectSubjectAt(root, s, now)
}

func inspectSubjectAt(root string, s subject, now time.Time) Inspection {
	gi := Inspection{CurrentTree: s.Tree, Reason: s.Reason}
	gitdir, err := benchgit.AdminDir(root)
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
	// A record that is not a gate verdict answers here and never reaches the reuse
	// path below. It carries no oracle to compare and no status to read, so the reader
	// names its class instead of grading it as a verdict of some other shape.
	if !loaded.class.isGateVerdict() {
		gi.CachedTree = loaded.record.Tree
		gi.RecordedAt, _ = time.Parse(time.RFC3339, loaded.record.RecordedAt)
		gi.Drifted = loaded.record.Tree != s.Tree
		gi.Reason = loaded.class.reuseRefusal
		return gi
	}
	rec := loaded.record
	gi.Status, gi.CachedTree = rec.Status, rec.Tree
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
	// Drift is graded ahead of the record's own status, and ahead of the closure check
	// that refuses reuse. A record naming a tree or an oracle this subject does not have
	// is evidence about some other run. A red of another tree read as this one's sends a
	// reader to fix work that is no longer here.
	drift := driftReason(rec, s)
	gi.Drifted = drift != ""
	if !s.Closed {
		return gi
	}
	if gi.Drifted {
		gi.Reason = drift
		return gi
	}
	if rec.Status != "green" {
		gi.Reason = "recorded " + rec.Status
		return gi
	}
	if now.Sub(tm) >= freshness {
		gi.Reason = "verdict expired"
		return gi
	}
	// This runs after drift and expiry: those retire a narrow record exactly as they
	// retire a full one. Naming the narrowness of an expired record would dress
	// retired evidence as current.
	if reason := narrowVerdictReason(loaded.class); reason != "" {
		gi.Reason = reason
		return gi
	}
	gi.ReusableGreen = true
	return gi
}

// driftReason names how a record fails to describe the subject, and "" when it
// describes it. One derivation answers both the reuse refusal and the Drifted flag
// consumers read. So no surface can call a retired record current by comparing tree
// hashes for itself.
func driftReason(rec verdictRecord, s subject) string {
	switch {
	case rec.Tree != s.Tree:
		return "working tree changed"
	case rec.Oracle != s.Oracle:
		return "oracle changed"
	}
	return ""
}

// narrowVerdictReason returns the reuse reason for the class the loader selected, and
// "" for a class whose green is reusable. The answer is the class's own declaration.
// Deriving it from the class name instead would silently readmit any class whose name
// happens to miss the spelling the derivation matched on.
func narrowVerdictReason(class verdictRecordClass) string {
	return class.reuseRefusal
}

type loadedVerdict struct {
	record verdictRecord
	class  verdictRecordClass
	state  State
	reason string
	bytes  int
}

// storeRecordBytes is one record file read from the evidence store. data is non-nil
// only when the file cleared every check that holds whatever class the bytes turn out
// to name. Otherwise state and reason say why the store has nothing readable there.
type storeRecordBytes struct {
	data   []byte
	bytes  int
	state  State
	reason string
}

// readStoreRecord applies the file discipline every record class in the store
// shares. It must be a regular 0600 file, within the size cap, framed as a single
// JSON object. What the bytes mean is the reading class's question.
// readStoreRecord answers only whether there are bytes worth asking about. So a
// class added to the store cannot be given a laxer file than the verdict cache gets.
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
	if err := strictJSON(read.data, &loaded.record); err != nil {
		loaded.state, loaded.reason = Invalid, "invalid cache record"
		return loaded
	}
	class, err := validateRecordBytes(read.data, loaded.record, now)
	if err != nil {
		loaded.state, loaded.reason = Invalid, "invalid cache record"
		return loaded
	}
	loaded.class = class
	loaded.state = loaded.record.State
	if !class.isGateVerdict() {
		// A non-verdict class declares no state field, so the record's own State is
		// empty. It is nonetheless a well-formed record the reader can name, which is
		// what Ready means to every consumer of this state.
		loaded.state = Ready
	}
	return loaded
}

// The two rejections strictJSON adds on top of encoding/json are sentinels so a caller
// can class them without matching on message text.
var (
	errTrailingJSON      = errors.New("trailing JSON")
	errDuplicateJSONName = errors.New("duplicate name")
)

// strictJSON takes its structural grammar from jsonfile and keeps the gate's diagnostic
// precedence. A typed refusal of the first value comes first, then trailing data, then a
// duplicate key. The typed decode stays here, so the manifest diagnostic keeps the
// encoding/json error text without a jsonfile prefix.
func strictJSON(data []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if jsonfile.RejectTrailing(dec) != nil {
		return errTrailingJSON
	}
	err := jsonfile.DecodeDocument(data, new(json.RawMessage))
	if errors.Is(err, jsonfile.ErrDuplicateKey) {
		return errDuplicateJSONName
	}
	return err
}

func isContentAddress(value string) bool {
	if len(value) != 64 || strings.ToLower(value) != value {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

// validateRecordBytes grades one record file against the class its exact field set
// names. The class is selected first, because the oracle identity is a gate verdict's
// field and not every class in the store is a gate verdict. Only the schema and the
// tree are common to all of them.
func validateRecordBytes(data []byte, r verdictRecord, now time.Time) (verdictRecordClass, error) {
	class, err := selectVerdictRecordClass(data)
	if err != nil {
		return verdictRecordClass{}, err
	}
	if err := validateStoreRecordCommon(r); err != nil {
		return verdictRecordClass{}, err
	}
	if class.isGateVerdict() {
		if len(r.Oracle) != 64 {
			return verdictRecordClass{}, errors.New("invalid record")
		}
		if _, err := hex.DecodeString(r.Oracle); err != nil || strings.ToLower(r.Oracle) != r.Oracle {
			return verdictRecordClass{}, errors.New("invalid oracle")
		}
	}
	return class, class.validate(data, r, now)
}

// validateStoreRecordCommon grades the fields that every record of the store shares: the
// schema and the tree. The gate cache classes and the evidence store both read through it.
func validateStoreRecordCommon(r verdictRecord) error {
	if r.Schema != verdictSchema || !treeHashRE.MatchString(r.Tree) {
		return errors.New("invalid record")
	}
	return nil
}

func strictRecordTime(value string) (time.Time, error) {
	tm, err := time.Parse(time.RFC3339, value)
	if err != nil || recordTime(tm) != value {
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
