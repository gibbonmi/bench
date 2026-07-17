package preflight

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type requirementStatus struct {
	Key          string `json:"key"`
	Owner        string `json:"owner"`
	Schema       string `json:"schema"`
	Requiredness string `json:"requiredness"`
	Applicable   bool   `json:"applicable"`
	Status       string `json:"status"`
	Reason       string `json:"reason,omitempty"`
	Digest       string `json:"digest,omitempty"`
}

type evidenceDigest struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

var releaseInputPaths = []string{
	"go.mod",
	"go.sum",
	"package.json",
	"scripts/go-build.sh",
	"scripts/build-artifacts.sh",
	"scripts/smoke-artifacts.sh",
	"scripts/platforms.json",
	"scripts/wrapper-assets.json",
}

func hasControlBytes(value string) bool {
	for _, byteValue := range []byte(value) {
		if byteValue < 0x20 || byteValue == 0x7f {
			return true
		}
	}
	return false
}

func validateRequirementRegistry() error {
	if requirements.SchemaVersion != 1 || len(requirements.PackageEvidence) == 0 || len(requirements.Records) == 0 {
		return errors.New("invalid requirement registry version or records")
	}
	packagePaths := map[string]bool{}
	for _, evidence := range requirements.PackageEvidence {
		if packagePaths[evidence.Path] || !safeRegistryPath(evidence.Path) || evidence.Schema == "" || evidence.Mode != "0644" {
			return fmt.Errorf("invalid package evidence registry entry %q", evidence.Path)
		}
		if evidence.Schema != "license/v1" && !knownRequirementSchema(evidence.Schema) {
			return fmt.Errorf("unsupported package evidence schema %q", evidence.Schema)
		}
		if evidence.Path != "LICENSE" {
			record, found := requirementForPath(evidence.Path)
			if !found || record.Schema != evidence.Schema {
				return fmt.Errorf("package evidence registry is not bound to requirement %q", evidence.Path)
			}
		}
		packagePaths[evidence.Path] = true
	}
	seen := map[string]bool{}
	public, bank := map[string]bool{}, map[string]bool{}
	for _, record := range requirements.Records {
		if seen[record.Key] || record.Key == "" || record.Owner == "" || !knownRequirementSchema(record.Schema) || !safeRegistryPath(record.Path) || len(record.Profiles) == 0 || record.Requiredness != "required" && record.Requiredness != "conditional" {
			return fmt.Errorf("invalid requirement registry record %q", record.Key)
		}
		for _, value := range []string{record.Key, record.Owner, record.Schema, record.Path} {
			if hasControlBytes(value) {
				return fmt.Errorf("requirement registry record %q contains control bytes", record.Key)
			}
		}
		seen[record.Key] = true
		profiles := map[Profile]bool{}
		for _, profile := range record.Profiles {
			if profiles[profile] {
				return fmt.Errorf("requirement %s repeats profile %q", record.Key, profile)
			}
			profiles[profile] = true
			switch profile {
			case ProfilePublic:
				public[record.Key] = true
			case ProfileBank:
				bank[record.Key] = true
			default:
				return fmt.Errorf("requirement %s has unknown profile %q", record.Key, profile)
			}
		}
	}
	for key := range public {
		if !bank[key] {
			return fmt.Errorf("bank profile is not a strict public superset: %s", key)
		}
	}
	if len(bank) <= len(public) {
		return errors.New("bank profile is not a strict public superset")
	}
	return nil
}

func safeRegistryPath(value string) bool {
	if value == "" || strings.Contains(value, "\\") || filepath.IsAbs(filepath.FromSlash(value)) || strings.Contains(value, "\x00") {
		return false
	}
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(value)))
	return clean == value && clean != "." && !strings.HasPrefix(clean, "../")
}

func readRollbackTarget(root string) (string, error) {
	data, err := readRegular(filepath.Join(root, "governance", "policies", "recovery-rollback.json"))
	if err != nil {
		return "", fmt.Errorf("rollback policy is unreadable: %w", err)
	}
	var record recoveryRollbackPolicy
	if err := decodeStrict(data, &record); err != nil {
		return "", fmt.Errorf("rollback policy is malformed: %w", err)
	}
	if record.SchemaVersion != 1 || record.Policy != "recovery-rollback" || record.RollbackTarget == "" || hasControlBytes(record.RollbackTarget) {
		return "", errors.New("rollback policy has no rollback target")
	}
	return record.RollbackTarget, nil
}

func inspectRequirements(root string, run RunEvidence, profile Profile) ([]requirementStatus, []evidenceDigest, string, error) {
	statuses := make([]requirementStatus, 0, len(requirements.Records))
	inputs := make([]evidenceDigest, 0, len(requirements.Records)+1)
	for _, record := range requirements.Records {
		status := requirementStatus{Key: record.Key, Owner: record.Owner, Schema: record.Schema, Requiredness: record.Requiredness, Status: "not_applicable"}
		applicable := containsProfile(record.Profiles, profile)
		if !applicable {
			status.Reason = "requirement is not applicable to selected profile"
			statuses = append(statuses, status)
			continue
		}
		status.Applicable = true
		path := filepath.Join(root, filepath.FromSlash(record.Path))
		data, err := readRegular(path)
		if os.IsNotExist(err) && isProducerOwner(record.Owner) {
			if run.Mode == ModeVerify {
				status.Status, status.Reason = "pending", "producer record is not present"
			} else {
				status.Status, status.Reason = "missing", "required producer record is not present"
			}
			statuses = append(statuses, status)
			continue
		}
		if err != nil {
			if os.IsNotExist(err) {
				status.Status, status.Reason = "missing", "required governance record is not present"
				statuses = append(statuses, status)
				continue
			}
			return nil, nil, "", fmt.Errorf("requirement %s is unreadable: %w", record.Key, err)
		}
		if len(data) == 0 {
			status.Status, status.Reason = "invalid", "record is empty"
			statuses = append(statuses, status)
			continue
		}
		if err := validateRequirementBytes(record, data, run.Identity); err != nil {
			status.Status, status.Reason = "invalid", err.Error()
			statuses = append(statuses, status)
			continue
		}
		status.Status = "satisfied"
		if record.Requiredness == "conditional" && isProducerOwner(record.Owner) {
			var envelope producerEnvelope
			if err := decodeStrict(data, &envelope); err == nil && envelope.Status == "not_applicable" {
				status.Status, status.Reason = "not_applicable", envelope.Reason
			}
		}
		status.Digest = digest(data)
		inputs = append(inputs, evidenceDigest{Path: record.Path, SHA256: status.Digest})
		statuses = append(statuses, status)
	}
	unsatisfied := ""
	for _, status := range statuses {
		if status.Status == "missing" || status.Status == "invalid" {
			unsatisfied = fmt.Sprintf("release requirement %s is %s: %s", status.Key, status.Status, status.Reason)
			break
		}
	}
	registryDigest := digest(requirementsJSON)
	inputs = append(inputs, evidenceDigest{Path: "internal/preflight/requirements.json", SHA256: registryDigest})
	for _, rel := range releaseInputPaths {
		data, err := readRegular(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return nil, nil, "", fmt.Errorf("release input %s is unreadable: %w", rel, err)
		}
		inputs = append(inputs, evidenceDigest{Path: rel, SHA256: digest(data)})
	}
	return statuses, inputs, unsatisfied, nil
}

func validateRequirementBytes(record Requirement, data []byte, identity Identity) error {
	if !bytes.HasSuffix(data, []byte("\n")) {
		return fmt.Errorf("requirement %s is missing a final newline", record.Key)
	}
	if isProducerOwner(record.Owner) {
		if record.Schema != producerSchema(record.Owner) {
			return fmt.Errorf("producer record %s has unsupported schema %s", record.Key, record.Schema)
		}
		var envelope producerEnvelope
		if err := decodeStrict(data, &envelope); err != nil {
			return fmt.Errorf("producer record %s is malformed: %w", record.Key, err)
		}
		if envelope.SchemaVersion != 1 || envelope.Key != record.Key || envelope.Owner != record.Owner || envelope.Status != "satisfied" && !(record.Requiredness == "conditional" && envelope.Status == "not_applicable" && envelope.Reason != "") {
			return fmt.Errorf("producer record %s has mismatched schema, key, owner, or status", record.Key)
		}
		if identity.SourceCommit == nil || envelope.Identity.SourceCommit != *identity.SourceCommit || identity.PackageVersion == nil || envelope.Identity.PackageVersion != *identity.PackageVersion {
			return fmt.Errorf("producer record %s identity does not match release", record.Key)
		}
		if len(envelope.Payload) == 0 || string(envelope.Payload) == "null" {
			return fmt.Errorf("producer record %s payload is empty", record.Key)
		}
		if envelope.Digest != digest(mustCanonicalPayload(envelope.Payload)) {
			return fmt.Errorf("producer record %s digest does not match payload", record.Key)
		}
		return nil
	}
	if record.Schema == "notices/v1" {
		if len(bytes.TrimSpace(data)) == 0 {
			return fmt.Errorf("governance record %s is empty", record.Key)
		}
		return nil
	}
	if record.Schema == "spdx-json/2.3" {
		if err := validateSPDXDocument(data, "", ""); err != nil {
			return fmt.Errorf("governance record %s has invalid SPDX JSON: %w", record.Key, err)
		}
		return nil
	}
	if record.Schema == "governance-policy/v1" {
		if err := validateGovernancePolicy(record, data); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("requirement %s has unsupported schema %s", record.Key, record.Schema)
}

func isProducerOwner(owner string) bool {
	return owner == "FT71" || owner == "FT87" || owner == "FT88"
}

func producerSchema(owner string) string {
	switch owner {
	case "FT71":
		return "ft71/local-event/v1"
	case "FT87":
		return "ft87/offline-network-control/v1"
	case "FT88":
		return "ft88/data-handling/v1"
	default:
		return ""
	}
}

type supportedVersionsPolicy struct {
	SchemaVersion     int    `json:"schema_version"`
	Policy            string `json:"policy"`
	Support           string `json:"support"`
	PreviousMinorDays int    `json:"previous_minor_days"`
	EOL               string `json:"eol"`
}

type securityResponsePolicy struct {
	SchemaVersion           int    `json:"schema_version"`
	Policy                  string `json:"policy"`
	PrivateRoute            string `json:"private_route"`
	AcknowledgeBusinessDays int    `json:"acknowledge_business_days"`
	TriageBusinessDays      int    `json:"triage_business_days"`
	MitigationBusinessDays  struct {
		Critical int    `json:"critical"`
		High     int    `json:"high"`
		Medium   int    `json:"medium"`
		Low      string `json:"low"`
	} `json:"mitigation_business_days"`
}

type dependencyLicensePolicy struct {
	SchemaVersion            int    `json:"schema_version"`
	Policy                   string `json:"policy"`
	ReviewRequired           bool   `json:"review_required"`
	LicenseChangeRequires    string `json:"license_change_requires"`
	DependencyChangeRequires string `json:"dependency_change_requires"`
}

type threatModelPolicy struct {
	SchemaVersion  int      `json:"schema_version"`
	Policy         string   `json:"policy"`
	ReleaseInput   bool     `json:"release_input"`
	TrustBoundary  string   `json:"trust_boundary"`
	PrimaryThreats []string `json:"primary_threats"`
}

type recoveryRollbackPolicy struct {
	SchemaVersion           int    `json:"schema_version"`
	Policy                  string `json:"policy"`
	RollbackTarget          string `json:"rollback_target"`
	PreservePriorGeneration bool   `json:"preserve_prior_generation"`
	NPMRollback             string `json:"npm_rollback"`
}

type supportPolicy struct {
	SchemaVersion int    `json:"schema_version"`
	Policy        string `json:"policy"`
	Route         string `json:"route"`
	PersonalEmail bool   `json:"personal_email"`
	NonPersonal   bool   `json:"non_personal"`
}

func validateGovernancePolicy(record Requirement, data []byte) error {
	valid := false
	switch record.Key {
	case "core.policy.supported_versions":
		var value supportedVersionsPolicy
		if err := decodeStrict(data, &value); err != nil {
			return fmt.Errorf("governance record %s is malformed: %w", record.Key, err)
		}
		valid = value.SchemaVersion == 1 && value.Policy == "supported-versions" && value.Support != "" && value.PreviousMinorDays > 0 && value.EOL != ""
	case "core.policy.security_response":
		var value securityResponsePolicy
		if err := decodeStrict(data, &value); err != nil {
			return fmt.Errorf("governance record %s is malformed: %w", record.Key, err)
		}
		valid = value.SchemaVersion == 1 && value.Policy == "security-response" && value.PrivateRoute != "" && value.AcknowledgeBusinessDays > 0 && value.TriageBusinessDays > 0 && value.MitigationBusinessDays.Critical > 0 && value.MitigationBusinessDays.High > 0 && value.MitigationBusinessDays.Medium > 0 && value.MitigationBusinessDays.Low != ""
	case "core.policy.dependency_license_change":
		var value dependencyLicensePolicy
		if err := decodeStrict(data, &value); err != nil {
			return fmt.Errorf("governance record %s is malformed: %w", record.Key, err)
		}
		valid = value.SchemaVersion == 1 && value.Policy == "dependency-license-change" && value.ReviewRequired && value.LicenseChangeRequires != "" && value.DependencyChangeRequires != ""
	case "core.policy.threat_model":
		var value threatModelPolicy
		if err := decodeStrict(data, &value); err != nil {
			return fmt.Errorf("governance record %s is malformed: %w", record.Key, err)
		}
		valid = value.SchemaVersion == 1 && value.Policy == "threat-model" && value.ReleaseInput && value.TrustBoundary != "" && len(value.PrimaryThreats) > 0
	case "core.policy.recovery_rollback":
		var value recoveryRollbackPolicy
		if err := decodeStrict(data, &value); err != nil {
			return fmt.Errorf("governance record %s is malformed: %w", record.Key, err)
		}
		valid = value.SchemaVersion == 1 && value.Policy == "recovery-rollback" && value.RollbackTarget != "" && value.PreservePriorGeneration && value.NPMRollback != ""
	case "core.policy.support":
		var value supportPolicy
		if err := decodeStrict(data, &value); err != nil {
			return fmt.Errorf("governance record %s is malformed: %w", record.Key, err)
		}
		valid = value.SchemaVersion == 1 && value.Policy == "support" && value.Route != "" && !value.PersonalEmail && value.NonPersonal
	default:
		return fmt.Errorf("governance record %s has an unknown policy schema", record.Key)
	}
	if !valid {
		return fmt.Errorf("governance record %s has invalid schema version or policy values", record.Key)
	}
	return nil
}

type spdxDocument struct {
	SPDXID            string             `json:"SPDXID"`
	SPDXVersion       string             `json:"SPDXVersion"`
	CreationInfo      spdxCreationInfo   `json:"creationInfo"`
	DataLicense       string             `json:"dataLicense"`
	Name              string             `json:"name"`
	DocumentNamespace string             `json:"documentNamespace"`
	Packages          []spdxPackage      `json:"packages"`
	Relationships     []spdxRelationship `json:"relationships"`
}

type spdxCreationInfo struct {
	Creators []string `json:"creators"`
}

type spdxPackage struct {
	SPDXID           string `json:"SPDXID"`
	Name             string `json:"name"`
	VersionInfo      string `json:"versionInfo"`
	DownloadLocation string `json:"downloadLocation"`
	LicenseConcluded string `json:"licenseConcluded"`
	LicenseDeclared  string `json:"licenseDeclared"`
}

type spdxRelationship struct {
	SPDXElementID      string `json:"spdxElementId"`
	RelationshipType   string `json:"relationshipType"`
	RelatedSPDXElement string `json:"relatedSpdxElement"`
}

func validateSPDXDocument(data []byte, expectedName, expectedVersion string) error {
	var document spdxDocument
	if err := decodeStrict(data, &document); err != nil {
		return fmt.Errorf("malformed SPDX JSON: %w", err)
	}
	if document.SPDXID != "SPDXRef-DOCUMENT" || document.SPDXVersion != "SPDX-2.3" || document.DataLicense != "CC0-1.0" || document.Name == "" || document.DocumentNamespace == "" || len(document.CreationInfo.Creators) == 0 || len(document.Packages) != 1 || len(document.Relationships) != 1 {
		return errors.New("SPDX document has invalid required fields or schema version")
	}
	pkg := document.Packages[0]
	if pkg.SPDXID == "" || pkg.Name == "" || pkg.VersionInfo == "" || pkg.DownloadLocation == "" || pkg.LicenseConcluded == "" || pkg.LicenseDeclared == "" || document.Relationships[0].SPDXElementID != document.SPDXID || document.Relationships[0].RelationshipType != "DESCRIBES" || document.Relationships[0].RelatedSPDXElement != pkg.SPDXID {
		return errors.New("SPDX document has incomplete package relationship")
	}
	if expectedName != "" && pkg.Name != expectedName {
		return fmt.Errorf("SPDX package name %q does not match %q", pkg.Name, expectedName)
	}
	if expectedVersion != "" && pkg.VersionInfo != expectedVersion {
		return fmt.Errorf("SPDX package version %q does not match %q", pkg.VersionInfo, expectedVersion)
	}
	return nil
}

type producerEnvelope struct {
	SchemaVersion int              `json:"schema_version"`
	Key           string           `json:"key"`
	Owner         string           `json:"owner"`
	Identity      producerIdentity `json:"identity"`
	Status        string           `json:"status"`
	Reason        string           `json:"reason"`
	Payload       json.RawMessage  `json:"payload"`
	Digest        string           `json:"digest"`
}

type producerIdentity struct {
	SourceCommit   string `json:"source_commit"`
	PackageVersion string `json:"package_version"`
}

func inputFingerprint(root string, run RunEvidence) (string, error) {
	h := sha256.New()
	paths := append([]string{"internal/preflight/requirements.json"}, releaseInputPaths...)
	for _, record := range requirements.Records {
		paths = append(paths, record.Path)
	}
	paths = append(paths, "dist/artifacts")
	sort.Strings(paths)
	seen := map[string]bool{}
	for _, rel := range paths {
		if seen[rel] {
			continue
		}
		seen[rel] = true
		if err := fingerprintPath(h, root, rel); err != nil {
			return "", err
		}
	}
	for _, result := range run.Phases {
		_, _ = h.Write([]byte(result.Name))
		_, _ = h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func fingerprintPath(h io.Writer, root, rel string) error {
	path := filepath.Join(root, filepath.FromSlash(rel))
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		_, _ = io.WriteString(h, "absent:\x00"+rel+"\n")
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || (!info.Mode().IsRegular() && !info.IsDir()) {
		return fmt.Errorf("unsafe release evidence input: %s", rel)
	}
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		sort.Strings(names)
		for _, name := range names {
			if err := fingerprintPath(h, root, filepath.ToSlash(filepath.Join(rel, name))); err != nil {
				return err
			}
		}
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, _ = io.WriteString(h, fmt.Sprintf("file:%s:%o:%d:", rel, info.Mode().Perm(), len(data)))
	_, _ = h.Write(data)
	_, _ = h.Write([]byte{0})
	return nil
}

func containsProfile(profiles []Profile, want Profile) bool {
	for _, profile := range profiles {
		if profile == want {
			return true
		}
	}
	return false
}

func decodeStrict(data []byte, value any) error {
	if !json.Valid(data) {
		return errors.New("invalid JSON")
	}
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("trailing JSON")
	}
	return nil
}

func rejectDuplicateJSONKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	var scan func() error
	scan = func() error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		delim, isDelim := token.(json.Delim)
		if !isDelim {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return err
				}
				key, ok := keyToken.(string)
				if !ok {
					return errors.New("JSON object key is not a string")
				}
				if seen[key] {
					return fmt.Errorf("duplicate JSON key %q", key)
				}
				seen[key] = true
				if err := scan(); err != nil {
					return err
				}
			}
		case '[':
			for decoder.More() {
				if err := scan(); err != nil {
					return err
				}
			}
		default:
			return errors.New("unexpected JSON delimiter")
		}
		_, err = decoder.Token()
		return err
	}
	return scan()
}
