package gate

// Ordinary conformance evidence is one atomically replaced ledger rather than one file per
// check. A green updates only entries it earned, a red deletes only entries it contradicted,
// and interruption never publishes a half-authored set. The ledger's map key preserves
// attribution when one entry is malformed; damage to the envelope widens every check.

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

const conformanceCheckSlotSchema = 1

const conformanceCheckSlotStoreName = "conformance-check-slots"

var conformanceCheckSlotStoreFields = []string{"check_slots", "schema"}

type conformanceCheckSlotStore struct {
	Schema int                        `json:"schema"`
	Slots  map[string]json.RawMessage `json:"check_slots"`
}

type conformanceCheckSlot struct {
	Schema               int           `json:"schema"`
	Check                string        `json:"check"`
	Tier                 registry.Tier `json:"tier"`
	Identity             string        `json:"identity"`
	AuthoredAt           string        `json:"authored_at"`
	CanaryImplementation string        `json:"canary_implementation"`
	CanaryShared         string        `json:"canary_shared"`
}

type conformanceCheckSkip struct {
	Check      string
	Identity   string
	AuthoredAt time.Time
}

type conformanceCheckPartition struct {
	Tier           registry.Tier
	Identities     map[string]string
	Executed       []string
	Inherited      []conformanceCheckSkip
	Canary         conformanceCanaryIdentities
	CanaryFull     bool
	CanaryFamilies []string
}

func (p conformanceCheckPartition) verdictExecuted() []string {
	executed := make(map[string]bool, len(p.Executed))
	for _, name := range p.Executed {
		executed[name] = true
	}
	names := make([]string, 0, len(p.Executed))
	for _, check := range registry.Checks {
		if check.RunsAt(p.Tier) && (check.Meta || executed[check.Name]) {
			names = append(names, check.Name)
		}
	}
	return names
}

func (p conformanceCheckPartition) verdictInherited() []string {
	names := make([]string, 0, len(p.Inherited))
	for _, inherited := range p.Inherited {
		names = append(names, inherited.Check)
	}
	return names
}

func (p conformanceCheckPartition) verdictEvidence() map[string]skipEvidence {
	evidence := make(map[string]skipEvidence, len(p.Inherited))
	for _, inherited := range p.Inherited {
		evidence[inherited.Check] = skipEvidence{
			Identity:   inherited.Identity,
			AuthoredAt: inherited.AuthoredAt.UTC().Truncate(time.Second).Format(time.RFC3339),
		}
	}
	return evidence
}

type conformanceCheckOutcome int

const (
	checkRunGreen conformanceCheckOutcome = iota
	checkRunRed
)

func partitionConformanceChecks(root string, tier registry.Tier, identities map[string]string, identityErr error, now time.Time) conformanceCheckPartition {
	partition := conformanceCheckPartition{Tier: tier, Identities: identities}
	checks := ordinaryConformanceChecks(tier)
	if identityErr != nil {
		for _, check := range checks {
			partition.Executed = append(partition.Executed, check.Name)
		}
		return partition
	}
	slots, valid := loadConformanceCheckSlots(root)
	if !valid {
		for _, check := range checks {
			partition.Executed = append(partition.Executed, check.Name)
		}
		return partition
	}
	for _, check := range checks {
		identity := identities[check.Name]
		if !isContentAddress(identity) {
			partition.Executed = append(partition.Executed, check.Name)
			continue
		}
		authoredAt, err := validateConformanceCheckSlot(slots[check.Name], check.Name, check.Tier, identity, now)
		if err != nil {
			partition.Executed = append(partition.Executed, check.Name)
			continue
		}
		partition.Inherited = append(partition.Inherited, conformanceCheckSkip{Check: check.Name, Identity: identity, AuthoredAt: authoredAt})
	}
	return partition
}

func executeAllConformanceChecks(tier registry.Tier, identities map[string]string) conformanceCheckPartition {
	partition := conformanceCheckPartition{Tier: tier, Identities: identities}
	for _, check := range ordinaryConformanceChecks(tier) {
		partition.Executed = append(partition.Executed, check.Name)
	}
	return partition
}

func applyConformanceCheckOutcome(root string, partition conformanceCheckPartition, outcome conformanceCheckOutcome, at time.Time) error {
	slots, valid := loadConformanceCheckSlots(root)
	if !valid {
		if outcome == checkRunRed {
			return removeConformanceCheckSlotStore(root)
		}
		slots = map[string]json.RawMessage{}
	}
	if err := validateCheckPartitionForPersistence(partition, slots, at, outcome == checkRunGreen); err != nil {
		return err
	}
	for _, name := range partition.Executed {
		if outcome == checkRunRed {
			delete(slots, name)
			continue
		}
		if !isContentAddress(partition.Identities[name]) {
			continue
		}
		check, _ := registry.Find(name)
		record := conformanceCheckSlot{
			Schema:               conformanceCheckSlotSchema,
			Check:                name,
			Tier:                 check.Tier,
			Identity:             partition.Identities[name],
			AuthoredAt:           at.UTC().Truncate(time.Second).Format(time.RFC3339),
			CanaryImplementation: partition.Canary.Bound[name],
			CanaryShared:         partition.Canary.Shared,
		}
		data, err := json.Marshal(record)
		if err != nil {
			return err
		}
		if _, err := validateConformanceCheckSlot(data, name, check.Tier, record.Identity, at); err != nil {
			return err
		}
		slots[name] = data
	}
	return replaceConformanceCheckSlotStore(root, slots)
}

func decorateConformanceCanarySelection(root string, partition *conformanceCheckPartition, current conformanceCanaryIdentities, resolutionErr error) {
	partition.Canary = current
	slots, valid := loadConformanceCheckSlots(root)
	if !valid {
		partition.CanaryFull = true
		return
	}
	if resolutionErr != nil || !isContentAddress(current.Shared) {
		partition.CanaryFull = true
		return
	}
	families := map[string]bool{}
	for _, check := range ordinaryConformanceChecks(partition.Tier) {
		owned := registry.CanaryFamilies(check.Name)
		var prior conformanceCheckSlot
		if strictJSON(slots[check.Name], &prior) != nil || prior.CanaryShared != current.Shared {
			partition.CanaryFull = true
			partition.CanaryFamilies = nil
			return
		}
		if len(owned) == 0 || prior.CanaryImplementation == current.Bound[check.Name] {
			continue
		}
		if !isContentAddress(prior.CanaryImplementation) || !isContentAddress(current.Bound[check.Name]) {
			partition.CanaryFull = true
			partition.CanaryFamilies = nil
			return
		}
		for _, family := range owned {
			families[family] = true
		}
	}
	for family := range families {
		partition.CanaryFamilies = append(partition.CanaryFamilies, family)
	}
	sort.Strings(partition.CanaryFamilies)
}

func validateCheckPartitionForPersistence(partition conformanceCheckPartition, slots map[string]json.RawMessage, now time.Time, requireInheritedEvidence bool) error {
	seen := map[string]bool{}
	for _, name := range partition.Executed {
		check, found := registry.Find(name)
		if !found || check.Meta || !check.RunsAt(partition.Tier) || seen[name] {
			return fmt.Errorf("invalid executed conformance check %q", name)
		}
		seen[name] = true
	}
	for _, inherited := range partition.Inherited {
		check, found := registry.Find(inherited.Check)
		if !found || check.Meta || !check.RunsAt(partition.Tier) || seen[inherited.Check] {
			return fmt.Errorf("invalid inherited conformance check %q", inherited.Check)
		}
		seen[inherited.Check] = true
		if requireInheritedEvidence {
			if _, err := validateConformanceCheckSlot(slots[inherited.Check], inherited.Check, check.Tier, partition.Identities[inherited.Check], now); err != nil {
				return fmt.Errorf("invalid inherited conformance check %q: %w", inherited.Check, err)
			}
		}
	}
	return nil
}

func loadConformanceCheckSlots(root string) (map[string]json.RawMessage, bool) {
	dir, err := componentSlotDir(root)
	if err != nil {
		return nil, false
	}
	read := readStoreRecord(filepath.Join(dir, conformanceCheckSlotStoreName))
	if read.data == nil {
		if read.state == Absent {
			return map[string]json.RawMessage{}, true
		}
		return nil, false
	}
	var store conformanceCheckSlotStore
	if strictJSON(read.data, &store) != nil || store.Schema != conformanceCheckSlotSchema || store.Slots == nil {
		return nil, false
	}
	for name := range store.Slots {
		check, found := registry.Find(name)
		if !found || check.Meta {
			return nil, false
		}
	}
	return store.Slots, true
}

func validateConformanceCheckSlot(data []byte, name string, tier registry.Tier, identity string, now time.Time) (time.Time, error) {
	if len(data) == 0 {
		return time.Time{}, errors.New("slot absent")
	}
	var record conformanceCheckSlot
	if strictJSON(data, &record) != nil || record.Schema != conformanceCheckSlotSchema || record.Check == "" || !isContentAddress(record.Identity) {
		return time.Time{}, errors.New("invalid check slot record")
	}
	if record.Check != name {
		return time.Time{}, errors.New("slot answers for another check")
	}
	if record.Tier != tier {
		return time.Time{}, errors.New("slot answers for another tier")
	}
	if record.Identity != identity {
		return time.Time{}, errors.New("slot answers for another identity")
	}
	check, found := registry.Find(name)
	if !found {
		return time.Time{}, errors.New("slot answers for an unknown check")
	}
	legacyCanary := record.CanaryShared == "" && record.CanaryImplementation == ""
	if !legacyCanary {
		if !isContentAddress(record.CanaryShared) {
			return time.Time{}, errors.New("invalid check slot canary identity")
		}
		if len(registry.CanaryFamilies(check.Name)) > 0 {
			if !isContentAddress(record.CanaryImplementation) {
				return time.Time{}, errors.New("invalid check slot canary implementation")
			}
		} else if record.CanaryImplementation != "" {
			return time.Time{}, errors.New("unexpected check slot canary implementation")
		}
	}
	authoredAt, err := strictRecordTime(record.AuthoredAt)
	if err != nil || authoredAt.After(now) {
		return time.Time{}, errors.New("invalid check slot time")
	}
	return authoredAt, nil
}

func replaceConformanceCheckSlotStore(root string, slots map[string]json.RawMessage) error {
	if len(slots) == 0 {
		return removeConformanceCheckSlotStore(root)
	}
	dir, err := componentSlotDir(root)
	if err != nil {
		return err
	}
	if err := ensureEvidenceDir(filepath.Dir(dir), dir); err != nil {
		return err
	}
	data, err := json.Marshal(conformanceCheckSlotStore{Schema: conformanceCheckSlotSchema, Slots: slots})
	if err != nil {
		return err
	}
	return durableReplaceRecordAt(dir, conformanceCheckSlotStoreName, data)
}

func removeConformanceCheckSlotStore(root string) error {
	dir, err := componentSlotDir(root)
	if err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(dir, conformanceCheckSlotStoreName)); err != nil && !errors.Is(err, os.ErrNotExist) {
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

func ordinaryConformanceChecks(tier registry.Tier) []registry.Check {
	var checks []registry.Check
	for _, check := range registry.Checks {
		if !check.Meta && check.RunsAt(tier) {
			checks = append(checks, check)
		}
	}
	return checks
}
