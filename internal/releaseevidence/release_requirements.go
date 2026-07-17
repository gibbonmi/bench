package releaseevidence

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type RequirementStatus struct {
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

func hasControlBytes(value string) bool {
	for _, byteValue := range []byte(value) {
		if byteValue < 0x20 || byteValue == 0x7f {
			return true
		}
	}
	return false
}

func validateRequirementRegistry(registry Registry) error {
	if registry.SchemaVersion != 1 || len(registry.Toolchains) == 0 || len(registry.Records) == 0 {
		return errors.New("invalid requirement registry version or records")
	}
	manifest := registry.ComponentManifest
	if manifest.SchemaVersion != 1 || !safeRegistryPath(manifest.Path) {
		return errors.New("invalid component manifest schema")
	}
	manifestFields := []string{manifest.RootFields.SchemaVersion, manifest.RootFields.Component, manifest.RootFields.Files, manifest.ComponentFields.Name, manifest.ComponentFields.Version, manifest.ComponentFields.Target, manifest.TargetFields.OS, manifest.TargetFields.Arch, manifest.FileFields.Path, manifest.FileFields.Mode, manifest.FileFields.Size, manifest.FileFields.SHA256}
	seenManifestField := map[string]bool{}
	for _, field := range manifestFields {
		if field == "" || hasControlBytes(field) || seenManifestField[field] {
			return errors.New("invalid component manifest schema field")
		}
		seenManifestField[field] = true
	}
	seen := map[string]bool{}
	packagePaths := map[string]bool{}
	public, bank := map[string]bool{}, map[string]bool{}
	for _, record := range registry.Records {
		coreSchema := record.Schema == "license/v1" || record.Schema == "notices/v1" || record.Schema == "spdx-json/2.3" || record.Schema == "governance-policy/v1"
		if seen[record.Key] || record.Key == "" || record.Owner == "" || record.Schema == "" || (!record.Producer && !coreSchema) || !safeRegistryPath(record.Path) || len(record.Profiles) == 0 || record.Requiredness != "required" && record.Requiredness != "conditional" {
			return fmt.Errorf("invalid requirement registry record %q", record.Key)
		}
		for _, value := range []string{record.Key, record.Owner, record.Schema, record.Path} {
			if hasControlBytes(value) {
				return fmt.Errorf("requirement registry record %q contains control bytes", record.Key)
			}
		}
		seen[record.Key] = true
		if record.PackageMode != "" {
			if record.Producer || record.PackageMode != "0644" || packagePaths[record.Path] {
				return fmt.Errorf("invalid packaged requirement %q", record.Key)
			}
			packagePaths[record.Path] = true
		}
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
	toolchains := map[string]bool{}
	for _, toolchain := range registry.Toolchains {
		if toolchains[toolchain.Name] || toolchain.Name == "" || len(toolchain.VersionArgv) == 0 || toolchain.VersionArgv[0] != toolchain.Name {
			return fmt.Errorf("invalid toolchain requirement %q", toolchain.Name)
		}
		toolchains[toolchain.Name] = true
		for operation, argv := range toolchain.Operations {
			if operation == "" || len(argv) == 0 {
				return fmt.Errorf("invalid toolchain operation %q for %s", operation, toolchain.Name)
			}
			for _, arg := range argv {
				if arg == "" || hasControlBytes(arg) {
					return fmt.Errorf("invalid toolchain operation %q for %s", operation, toolchain.Name)
				}
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
	data, err := ReadRegular(filepath.Join(root, "governance", "policies", "recovery-rollback.json"))
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

func validateRequirementBytes(record Requirement, data []byte, identity Identity) error {
	if !bytes.HasSuffix(data, []byte("\n")) {
		return fmt.Errorf("requirement %s is missing a final newline", record.Key)
	}
	if record.Producer {
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
	if record.Schema == "license/v1" {
		if len(bytes.TrimSpace(data)) == 0 {
			return fmt.Errorf("license record %s is empty", record.Key)
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
		if err := validateSPDXDocument(data, "", "__PACKAGE_VERSION__"); err != nil {
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

func containsProfile(profiles []Profile, want Profile) bool {
	for _, profile := range profiles {
		if profile == want {
			return true
		}
	}
	return false
}
