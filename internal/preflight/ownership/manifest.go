package ownership

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/jsonfile"
)

// Manifest is the authored planning-claim inventory of one spec. Authority and
// Completeness state what its author declared; neither authenticates review, approval, or
// semantic completeness. A present manifest with a diagnostic keeps no item.
type Manifest struct {
	Present                           bool
	Version                           int
	Authority, Completeness, Baseline string
	Scope                             []ScopeItem
	Changes                           []Change
	Premises                          []PremiseItem
	Seams                             []Seam
	Decisions                         []Decision
	ManualClaims                      []ManualClaim
	Diagnostics                       []Diagnostic
}

// Location is one exact place: a path, its positive line, the needle bytes that line holds,
// and the source it is read at. An item location reads the spec; an owner location reads
// the pinned base or tip.
type Location struct {
	Path   string `json:"path"`
	Line   int    `json:"line"`
	Needle string `json:"needle"`
	Source string `json:"source"`
}

// ScopeItem declares the exact paths or segment-bounded prefixes an inventory covers.
type ScopeItem struct {
	ID       string   `json:"id"`
	Location Location `json:"location"`
	Paths    []string `json:"paths"`
}

// Owner is one exact declaration identity.
type Owner struct {
	Path      string `json:"path"`
	Package   string `json:"package"`
	Symbol    string `json:"symbol"`
	Kind      string `json:"kind"`
	Signature string `json:"signature"`
}

// Disposition is the action a change takes for one retirement consumer.
type Disposition struct {
	Symbol string `json:"symbol"`
	Action string `json:"action"`
}

// Change declares one changed owner and the callers its change reaches.
type Change struct {
	ID           string        `json:"id"`
	Location     Location      `json:"location"`
	Kind         string        `json:"kind"`
	Owner        Owner         `json:"owner"`
	CallerScope  string        `json:"caller_scope"`
	State        string        `json:"state"`
	Destination  *string       `json:"destination"`
	Dispositions []Disposition `json:"dispositions"`
}

// PremiseItem declares one premise. Its subject and expected values are tagged objects
// whose per-kind grammar belongs to the premise's own collector.
type PremiseItem struct {
	ID       string          `json:"id"`
	Location Location        `json:"location"`
	Kind     string          `json:"kind"`
	State    string          `json:"state"`
	Subject  json.RawMessage `json:"subject"`
	Expected json.RawMessage `json:"expected"`
}

// Seam declares one test seam and the rows it supports.
type Seam struct {
	ID       string   `json:"id"`
	Location Location `json:"location"`
	Path     string   `json:"path"`
	Package  string   `json:"package"`
	Function string   `json:"function"`
	State    string   `json:"state"`
	Rows     []string `json:"rows"`
}

// Decision binds one source predicate's exact bytes to one coverage row.
type Decision struct {
	ID        string   `json:"id"`
	Location  Location `json:"location"`
	Source    Location `json:"source"`
	Predicate string   `json:"predicate"`
	Row       string   `json:"row"`
}

// ManualApplicability is the explicit form of one manual claim.
type ManualApplicability string

const (
	// CurrentPremise is a required current-source premise that no collector supports.
	CurrentPremise ManualApplicability = "current-premise"
	// NativeFuture is future native acceptance evidence that its own implementation owns.
	NativeFuture ManualApplicability = "native-future"
)

// ManualClaim is one unsupported claim. NativeWitness is present exactly on a
// native-future claim.
type ManualClaim struct {
	ID            string
	Location      Location
	Reason        string
	Applicability ManualApplicability
	NativeWitness *NativeWitness
}

// NativeWitness names the one native acceptance row a native-future claim waits for, and
// the plan route that later verifies it. Observation is always unobserved here.
type NativeWitness struct {
	Row          string   `json:"row"`
	Chunk        string   `json:"chunk"`
	Ticket       string   `json:"ticket"`
	Verification string   `json:"verification"`
	PhaseOwner   Location `json:"phase_owner"`
	Observation  string   `json:"observation"`
}

const completenessDeclared = "declared-complete"

var (
	locationKeys = []string{"path", "line", "needle", "source"}
	ownerKeys    = []string{"path", "package", "symbol", "kind", "signature"}
	manualKeys   = []string{"id", "location", "reason", "applicability"}
	nativeKeys   = []string{"row", "chunk", "ticket", "verification", "phase_owner", "observation"}
	// reusedTestKeys is the reused-test subject, which a posture premise's writer probe also
	// takes.
	reusedTestKeys = []string{"path", "function", "record_spec", "source_commit", "evidence_commit", "chunk", "verification", "mutation"}
	premiseKinds   = []string{"lines", "glob-members", "headroom", "symbol", "import", "anchor-kind", "sentinel", "error", "bound",
		"reused-test", "reused-owner", "retirement", "posture", "fence-exclusion", "printed-field"}
)

// decodeManifest grades one payload against the version 1 schema. Every top-level array
// is required, every object carries exactly its declared keys, and every enum takes its
// exact spelling. rows is the spec's coverage map, which each row reference must name. The
// strict shared decoder reads the payload first, because a plain decode keeps the last
// value of a repeated key and loses the first promise.
func decodeManifest(payload []byte, rows []coverage.MapRow) (Manifest, error) {
	var top map[string]json.RawMessage
	if err := jsonfile.DecodeDocument(payload, &top); err != nil {
		return Manifest{}, fmt.Errorf("%s payload: %w", MetadataLabel, err)
	}
	if top == nil {
		return Manifest{}, fmt.Errorf("%s payload is not a JSON object", MetadataLabel)
	}
	var m Manifest
	if raw, ok := top["version"]; ok {
		if err := json.Unmarshal(raw, &m.Version); err != nil || m.Version != 1 {
			return Manifest{}, fmt.Errorf("unsupported metadata version %s", raw)
		}
	}
	object, err := fields(payload, "", "version", "authority", "completeness", "baseline", "scope", "changes", "premises", "seams", "decisions", "manual_claims")
	if err != nil {
		return Manifest{}, err
	}
	for _, field := range []struct {
		key    string
		target *string
	}{{"authority", &m.Authority}, {"completeness", &m.Completeness}, {"baseline", &m.Baseline}} {
		if err := decodeInto(object, "", field.key, field.target); err != nil {
			return Manifest{}, err
		}
	}
	if err := firstError(enum("", "authority", m.Authority, "authored-spec"), enum("", "completeness", m.Completeness, completenessDeclared, "partial"), enum("", "baseline", m.Baseline, "source-base")); err != nil {
		return Manifest{}, err
	}
	decoder := itemDecoder{object: object, rows: rows}
	m.Scope = items(&decoder, "scope", decodeScope)
	m.Changes = items(&decoder, "changes", decodeChange)
	m.Premises = items(&decoder, "premises", decodePremise)
	m.Seams = items(&decoder, "seams", decodeSeam)
	m.Decisions = items(&decoder, "decisions", decodeDecision)
	m.ManualClaims = items(&decoder, "manual_claims", decodeManual)
	if decoder.err != nil {
		return Manifest{}, decoder.err
	}
	return m, nil
}

// itemDecoder decodes each item array and keeps the first fault, the item IDs seen across
// every array, and the coverage rows a row reference must name.
type itemDecoder struct {
	object map[string]json.RawMessage
	rows   []coverage.MapRow
	ids    []string
	err    error
}

// items decodes one top-level item array through decode.
func items[T any](d *itemDecoder, key string, decode func(*itemDecoder, json.RawMessage, string) (T, string, error)) []T {
	if d.err != nil {
		return nil
	}
	var raws []json.RawMessage
	if err := json.Unmarshal(d.object[key], &raws); err != nil || raws == nil {
		d.err = fmt.Errorf("%s is not a JSON array", key)
		return nil
	}
	decoded := make([]T, 0, len(raws))
	for i, raw := range raws {
		where := fmt.Sprintf("%s[%d]: ", key, i)
		item, id, err := decode(d, raw, where)
		if err == nil && (id == "" || slices.Contains(d.ids, id)) {
			err = fmt.Errorf("duplicate item id %q", id)
			if id == "" {
				err = fmt.Errorf("%sempty id", where)
			}
		}
		if err != nil {
			d.err = err
			return nil
		}
		d.ids = append(d.ids, id)
		decoded = append(decoded, item)
	}
	return decoded
}

// typed checks raw's exact keys, then decodes it into target.
func typed(raw json.RawMessage, where string, target any, keys ...string) (map[string]json.RawMessage, error) {
	object, err := fields(raw, where, keys...)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return nil, fmt.Errorf("%sinvalid value: %v", where, err)
	}
	return object, nil
}

// location grades one nested location. An item location reads the spec; an owner location
// reads the pinned base or tip.
func location(object map[string]json.RawMessage, key, where string, owner bool) error {
	var loc Location
	if _, err := typed(object[key], where+key+": ", &loc, locationKeys...); err != nil {
		return err
	}
	if loc.Path == "" || loc.Line < 1 || loc.Needle == "" {
		return fmt.Errorf("%s%s needs a path, a positive line, and a nonempty needle", where, key)
	}
	if owner {
		return enum(where+key+": ", "source", loc.Source, "base", "tip")
	}
	return enum(where+key+": ", "source", loc.Source, "spec")
}

// row refuses a row reference that names no coverage-map row.
func (d *itemDecoder) row(where, id string) error {
	if !slices.ContainsFunc(d.rows, func(row coverage.MapRow) bool { return row.ID == id }) {
		return fmt.Errorf("%sunknown row %q", where, id)
	}
	return nil
}

func decodeScope(_ *itemDecoder, raw json.RawMessage, where string) (ScopeItem, string, error) {
	var item ScopeItem
	object, err := typed(raw, where, &item, "id", "location", "paths")
	if err == nil {
		err = location(object, "location", where, false)
	}
	if err == nil && len(item.Paths) == 0 {
		err = fmt.Errorf("%sscope names no path", where)
	}
	return item, item.ID, err
}

func decodeChange(_ *itemDecoder, raw json.RawMessage, where string) (Change, string, error) {
	var item Change
	object, err := typed(raw, where, &item, "id", "location", "kind", "owner", "caller_scope", "state", "destination", "dispositions")
	if err == nil {
		_, err = identity(object["owner"], where+"owner: ")
	}
	var dispositions []json.RawMessage
	if err == nil {
		err = json.Unmarshal(object["dispositions"], &dispositions)
	}
	for i := 0; err == nil && i < len(dispositions); i++ {
		_, err = disposition(dispositions[i], fmt.Sprintf("%sdispositions[%d]: ", where, i))
	}
	if err == nil {
		err = firstError(location(object, "location", where, false),
			enum(where, "kind", item.Kind, "signature", "fixture-helper", "grammar", "posture", "liveness"),
			enum(where, "caller_scope", item.CallerScope, "production-and-tests"),
			enum(where, "state", item.State, "existing", "planned"))
	}
	if err == nil && item.Destination != nil && *item.Destination == "" {
		err = fmt.Errorf("%sdestination is empty", where)
	}
	return item, item.ID, err
}

func decodePremise(_ *itemDecoder, raw json.RawMessage, where string) (PremiseItem, string, error) {
	var item PremiseItem
	object, err := typed(raw, where, &item, "id", "location", "kind", "state", "subject", "expected")
	if err == nil {
		err = firstError(location(object, "location", where, false),
			enum(where, "kind", item.Kind, premiseKinds...),
			enum(where, "state", item.State, "current", "planned"))
	}
	for _, key := range []string{"subject", "expected"} {
		var tagged map[string]json.RawMessage
		if err == nil && (json.Unmarshal(object[key], &tagged) != nil || tagged == nil) {
			err = fmt.Errorf("%s%s is not a tagged object", where, key)
		}
	}
	return item, item.ID, err
}

func decodeSeam(d *itemDecoder, raw json.RawMessage, where string) (Seam, string, error) {
	var item Seam
	object, err := typed(raw, where, &item, "id", "location", "path", "package", "function", "state", "rows")
	if err == nil {
		err = firstError(location(object, "location", where, false), enum(where, "state", item.State, "existing", "planned"))
	}
	if err == nil && len(item.Rows) == 0 {
		err = fmt.Errorf("%sseam supports no row", where)
	}
	for i := 0; err == nil && i < len(item.Rows); i++ {
		err = d.row(where, item.Rows[i])
	}
	return item, item.ID, err
}

func decodeDecision(d *itemDecoder, raw json.RawMessage, where string) (Decision, string, error) {
	var item Decision
	object, err := typed(raw, where, &item, "id", "location", "source", "predicate", "row")
	if err == nil {
		err = firstError(location(object, "location", where, false), location(object, "source", where, true))
	}
	if err == nil && item.Predicate == "" {
		err = fmt.Errorf("%spredicate is empty", where)
	}
	if err == nil {
		err = d.row(where, item.Row)
	}
	return item, item.ID, err
}

// decodeManual grades one manual claim. The applicability field selects the exact key set:
// a current-premise claim carries no native field, and a native-future claim carries every
// one, with an unobserved observation.
func decodeManual(d *itemDecoder, raw json.RawMessage, where string) (ManualClaim, string, error) {
	var head struct {
		ID            string              `json:"id"`
		Reason        string              `json:"reason"`
		Applicability ManualApplicability `json:"applicability"`
		Location      Location            `json:"location"`
	}
	if err := json.Unmarshal(raw, &head); err != nil {
		return ManualClaim{}, "", fmt.Errorf("%sinvalid value: %v", where, err)
	}
	keys := manualKeys
	if head.Applicability == NativeFuture {
		keys = append(append([]string(nil), manualKeys...), nativeKeys...)
	}
	object, err := fields(raw, where, keys...)
	if err == nil {
		err = firstError(enum(where, "applicability", string(head.Applicability), string(CurrentPremise), string(NativeFuture)), location(object, "location", where, false))
	}
	if err == nil && head.Reason == "" {
		err = fmt.Errorf("%sreason is empty", where)
	}
	claim := ManualClaim{ID: head.ID, Location: head.Location, Reason: head.Reason, Applicability: head.Applicability}
	if err == nil && head.Applicability == NativeFuture {
		var native NativeWitness
		if err = json.Unmarshal(raw, &native); err != nil {
			err = fmt.Errorf("%sinvalid value: %v", where, err)
		}
		if err == nil {
			err = firstError(location(object, "phase_owner", where, true), enum(where, "observation", native.Observation, "unobserved"), d.row(where, native.Row))
		}
		claim.NativeWitness = &native
	}
	return claim, head.ID, err
}

// firstError is the first fault of an ordered check list.
func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
