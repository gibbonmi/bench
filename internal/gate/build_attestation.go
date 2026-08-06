package gate

// A build attestation is the evidence that the gate itself produced the binary a freshness
// seal describes. A seal is a self-consistent pair — a binary and the digests of that binary
// and of the sources it was built from — and anything able to write beside dist/bench can
// produce one, so a seal on its own proves only that whoever wrote it wrote both halves. The
// attestation is the half no seal writer can forge without the evidence store: a record in
// the same retained-evidence store as the component ancestor slots, authored only by a build
// that ran green inside a gate, naming the digest of the binary that build produced.
//
// Reading is therefore a comparison and never a parse. The seal says which bytes sit beside
// the binary now, the attestation says which bytes a gate built, and only their agreement
// lets the build phase skip. Every refusal here means run the build, which republishes the
// seal and re-authors the attestation together, so a disagreement cannot outlive one build.
//
// The attested digest is hashed from the artifact's own bytes, never read back out of the
// seal beside it. A digest taken from a seal would only restate what some other writer put on
// disk, which is the one claim this record exists to not make.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"path/filepath"
	"time"

	// The package's own freshness constant owns the bare name, so the seal package is
	// reached through an alias.
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
)

const buildAttestationSchema = 1

// buildAttestationDomain separates an attestation's address from every other name in the
// store, exactly as the slot domain does for a slot.
const buildAttestationDomain = componentPolicyVersion + "/build-attestation"

// buildAttestationRecord is the attestation class. Executable is the digest of the binary a
// gate build produced; the artifact it answers for is framed into the record's address
// rather than repeated here, since two artifacts whose bytes agree are the same claim.
type buildAttestationRecord struct {
	Schema     int    `json:"schema"`
	Executable string `json:"executable"`
	AuthoredAt string `json:"authored_at"`
}

// buildAttestationFields is the exact field set an attestation carries. It joins the family
// enumeration in record_classes.go, which spans the verdict, slot, and attestation classes
// equally and so belongs to none of the three files that define a single class.
var buildAttestationFields = []string{"authored_at", "executable", "schema"}

// buildAttestationInspection is the answer a build-skip decision reads. Attested is true only
// for a record that validated as an attestation of the exact bytes the seal names; otherwise
// Reason says what the store held instead, and the build runs.
type buildAttestationInspection struct {
	Attested   bool
	AuthoredAt time.Time
	Reason     string
}

// authorBuildAttestation records that a gate build produced the binary whose bytes hash to
// digest, addressed to the artifact path that binary is published as.
func authorBuildAttestation(root, executable, digest string, authoredAt time.Time) error {
	record := buildAttestationRecord{
		Schema:     buildAttestationSchema,
		Executable: digest,
		AuthoredAt: authoredAt.UTC().Truncate(time.Second).Format(time.RFC3339),
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	// Graded before it is published, against the validator the reader applies: an attestation
	// that cannot be read back would make the build phase run forever, and the author is where
	// that is still fixable.
	if _, err := validateBuildAttestationBytes(data, record, authoredAt); err != nil {
		return err
	}
	dir, err := componentSlotDir(root)
	if err != nil {
		return err
	}
	name, err := buildAttestationName(executable)
	if err != nil {
		return err
	}
	if err := ensureEvidenceDir(filepath.Dir(dir), dir); err != nil {
		return err
	}
	return durableReplaceRecordAt(dir, name, data)
}

// verifyBuildAttestation reports whether executable's seal names bytes a gate build
// authored. It only reads: a seal the store cannot vouch for is evidence of nothing, and
// re-stamping it here would credit a binary nobody built.
func verifyBuildAttestation(root, executable string, now time.Time) buildAttestationInspection {
	_, sealed, err := benchfreshness.SealDigests(executable)
	if err != nil {
		return buildAttestationInspection{Reason: "seal unreadable"}
	}
	dir, err := componentSlotDir(root)
	if err != nil {
		return buildAttestationInspection{Reason: "evidence unavailable"}
	}
	name, err := buildAttestationName(executable)
	if err != nil {
		return buildAttestationInspection{Reason: "evidence unavailable"}
	}
	read := readStoreRecord(filepath.Join(dir, name))
	if read.data == nil {
		if read.state == Absent {
			return buildAttestationInspection{Reason: "attestation absent"}
		}
		return buildAttestationInspection{Reason: read.reason}
	}
	var record buildAttestationRecord
	if err := strictJSON(read.data, &record); err != nil {
		return buildAttestationInspection{Reason: "invalid attestation record"}
	}
	authoredAt, err := validateBuildAttestationBytes(read.data, record, now)
	if err != nil {
		return buildAttestationInspection{Reason: err.Error()}
	}
	// The comparison the class exists for. The seal answers for the bytes on disk and the
	// attestation for the bytes a gate built, so a seal recomputed around a planted binary
	// disagrees with the store no matter how self-consistent it is.
	if record.Executable != sealed {
		return buildAttestationInspection{Reason: "attestation names another binary"}
	}
	return buildAttestationInspection{Attested: true, AuthoredAt: authoredAt}
}

func validateBuildAttestationBytes(data []byte, record buildAttestationRecord, now time.Time) (time.Time, error) {
	if record.Schema != buildAttestationSchema || !isContentAddress(record.Executable) {
		return time.Time{}, errors.New("invalid attestation record")
	}
	authoredAt, err := strictRecordTime(record.AuthoredAt)
	if err != nil || authoredAt.After(now) {
		return time.Time{}, errors.New("invalid attestation time")
	}
	return authoredAt, nil
}

// buildAttestationName is the store file name the attestation for executable occupies. The
// artifact path is framed in its cleaned absolute spelling, so two artifacts under one git
// directory cannot share an attestation and one artifact cannot gain a second. The store sits
// under the common git directory, which every worktree shares while each builds its own
// artifact, so that separation is what keeps one worktree's build from retiring another's.
//
// Cleaning is the whole of the normalization, and it suffices: the seal contract refuses an
// artifact reached through any symlinked path component, so a spelling that resolves a
// symlink cannot reach an artifact this address is ever asked about. Resolving them here
// would instead have to fail on the first build, where the artifact does not yet exist.
func buildAttestationName(executable string) (string, error) {
	absolute, err := filepath.Abs(executable)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	frame(h, buildAttestationDomain)
	frame(h, absolute)
	return hex.EncodeToString(h.Sum(nil)), nil
}
