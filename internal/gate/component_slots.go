package gate

// A component ancestor slot is the on-disk evidence that lets one gate component skip. It
// lives in the same retained-evidence store as the whole-tree green, one slot per scoped
// component, addressed by that component's identity.
//
// Two properties carry the whole design. The address is the content, so a slot stands until
// its component's identity moves and no clock retires it — expiring the proof by time would
// re-charge the component for a change its inputs are blind to. And authorship is at
// execution: only a run that executed the component green writes its slot, so a slot's
// existence is a claim some run graded this exact component at this exact identity.
//
// That makes every reading path a forgery question, and every refusal here answers
// run-the-component. Nothing in this file repairs, rewrites, or re-stamps a slot on the read
// path: a record that cannot be validated as this component's slot at this identity is
// evidence of nothing, and a reader that patched it into shape would be crediting work
// nobody graded.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const componentSlotSchema = 1

// componentSlotDomain separates a slot's address from every other name in the store. The
// whole-tree evidence hashes a tree and an oracle under no domain, so without this a slot
// and a verdict could in principle resolve to one file name.
const componentSlotDomain = componentPolicyVersion + "/slot"

// componentSlotRecord is the slot class. It names the component it answers for and the
// identity it was authored at, and carries no field of either verdict class — a slot and a
// verdict are different claims, and a record that could be read as either would let a
// whole-tree green stand in for a component nobody ran.
type componentSlotRecord struct {
	Schema     int    `json:"schema"`
	Component  string `json:"component"`
	Identity   string `json:"identity"`
	AuthoredAt string `json:"authored_at"`
}

// componentSlotFields is the exact field set a slot record carries. It feeds the record-class
// family registry in record_classes.go, which checks the field set against every other class
// for disjointness — the two verdict classes are alternatives rather than a spectrum for the
// same reason, and a record holding part of one class and part of another names nothing a
// reader can resolve without guessing.
var componentSlotFields = []string{"authored_at", "component", "identity", "schema"}

// recordSharedFields are the names carried by every record class in the store, so they are
// the only names a slot may share with a verdict. Everything else in componentSlotFields
// has to be disjoint from the verdict classes for the two to stay mutually unreadable.
var recordSharedFields = []string{"schema"}

// componentSlotInspection is the answer a skip decision reads. Skippable is true only for a
// record that validated as this component's slot at this identity; otherwise Reason says
// what the store held instead, and the component runs.
type componentSlotInspection struct {
	Skippable  bool
	AuthoredAt time.Time
	Reason     string
}

// authorComponentSlot records that a run executed component green at identity. Only the
// named component's slot is written: a run that skipped some other component graded nothing
// about it, and authoring its slot here would launder one component's green into evidence
// for another's.
func authorComponentSlot(root, component, identity string, authoredAt time.Time) error {
	record := componentSlotRecord{
		Schema:     componentSlotSchema,
		Component:  component,
		Identity:   identity,
		AuthoredAt: authoredAt.UTC().Truncate(time.Second).Format(time.RFC3339),
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	// The bytes are graded before they are published, against the same validator the reader
	// applies. A slot that cannot be read back as this component's is a slot that would make
	// the component run forever, and the author is where that is still fixable.
	if _, err := validateComponentSlotBytes(data, record, authoredAt); err != nil {
		return err
	}
	dir, err := componentSlotDir(root)
	if err != nil {
		return err
	}
	if err := ensureEvidenceDir(filepath.Dir(dir), dir); err != nil {
		return err
	}
	return durableReplaceRecordAt(dir, componentSlotName(component, identity), data)
}

// resolveComponentSlot reports whether component may skip at identity. It only reads.
func resolveComponentSlot(root, component, identity string, now time.Time) componentSlotInspection {
	dir, err := componentSlotDir(root)
	if err != nil {
		return componentSlotInspection{Reason: "evidence unavailable"}
	}
	return inspectComponentSlotAt(filepath.Join(dir, componentSlotName(component, identity)), component, identity, now)
}

// invalidateComponentSlot removes component's slot at identity, the answer to a red
// component. Other components' slots are untouched: a red vet says nothing about test.
func invalidateComponentSlot(root, component, identity string) error {
	dir, err := componentSlotDir(root)
	if err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(dir, componentSlotName(component, identity))); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	file, err := os.Open(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Sync()
}

func inspectComponentSlotAt(path, component, identity string, now time.Time) componentSlotInspection {
	read := readStoreRecord(path)
	if read.data == nil {
		if read.state == Absent {
			return componentSlotInspection{Reason: "slot absent"}
		}
		return componentSlotInspection{Reason: read.reason}
	}
	var record componentSlotRecord
	if err := strictJSON(read.data, &record); err != nil {
		return componentSlotInspection{Reason: "invalid slot record"}
	}
	authoredAt, err := validateComponentSlotBytes(read.data, record, now)
	if err != nil {
		return componentSlotInspection{Reason: err.Error()}
	}
	// The record names what it answers for, and the slot's address says what was asked. They
	// are compared rather than assumed equal: the address is derived from both, so a record
	// disagreeing with the file it was found in is a record something other than an author
	// put there.
	if record.Component != component || record.Identity != identity {
		return componentSlotInspection{Reason: "slot answers for another component"}
	}
	return componentSlotInspection{Skippable: true, AuthoredAt: authoredAt}
}

func validateComponentSlotBytes(data []byte, record componentSlotRecord, now time.Time) (time.Time, error) {
	if record.Schema != componentSlotSchema || record.Component == "" || !isContentAddress(record.Identity) {
		return time.Time{}, errors.New("invalid slot record")
	}
	authoredAt, err := strictRecordTime(record.AuthoredAt)
	if err != nil || authoredAt.After(now) {
		return time.Time{}, errors.New("invalid slot time")
	}
	return authoredAt, nil
}

// componentSlotName is the store file name component's slot at identity occupies. Both are
// framed into it, so no component's slot can be filed where another's is looked for.
func componentSlotName(component, identity string) string {
	h := sha256.New()
	frame(h, componentSlotDomain)
	frame(h, component)
	frame(h, identity)
	return hex.EncodeToString(h.Sum(nil))
}

func componentSlotDir(root string) (string, error) {
	gitdir, err := commonGitDir(root)
	if err != nil {
		return "", err
	}
	return filepath.Join(gitdir, "bench-gate-evidence"), nil
}

// isContentAddress reports whether value is a lowercase sha256 digest. Case is part of it:
// two spellings of one digest would address two slots for one identity.
func isContentAddress(value string) bool {
	if len(value) != 64 || strings.ToLower(value) != value {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

// componentSlotSharesVerdictFields returns the names componentSlotFields shares with either
// verdict class beyond the store-wide shared ones. It is what makes "a slot carries no
// verdict-class field" a checked property of the class rather than a claim in a comment.
func componentSlotSharesVerdictFields() []string {
	return recordClassSharesVerdictFields(componentSlotFields)
}
