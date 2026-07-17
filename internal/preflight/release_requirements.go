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
	if requirements.SchemaVersion != 1 || len(requirements.Records) == 0 {
		return errors.New("invalid requirement registry version or records")
	}
	seen := map[string]bool{}
	public, bank := map[string]bool{}, map[string]bool{}
	for _, record := range requirements.Records {
		if seen[record.Key] || record.Key == "" || record.Owner == "" || record.Schema == "" || !safeRegistryPath(record.Path) || len(record.Profiles) == 0 || record.Requiredness != "required" && record.Requiredness != "conditional" {
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
	var record map[string]json.RawMessage
	if err := decodeStrict(data, &record); err != nil {
		return "", fmt.Errorf("rollback policy is malformed: %w", err)
	}
	var target string
	if err := json.Unmarshal(record["rollback_target"], &target); err != nil || target == "" || hasControlBytes(target) {
		return "", errors.New("rollback policy has no rollback target")
	}
	return target, nil
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
		var document map[string]json.RawMessage
		if err := decodeStrict(data, &document); err != nil {
			return fmt.Errorf("governance record %s has malformed SPDX JSON", record.Key)
		}
		var version string
		if err := json.Unmarshal(document["SPDXVersion"], &version); err != nil || version != "SPDX-2.3" {
			return fmt.Errorf("governance record %s has unsupported SPDX schema version", record.Key)
		}
		return nil
	}
	var value map[string]json.RawMessage
	if err := decodeStrict(data, &value); err != nil {
		return fmt.Errorf("governance record %s is malformed: %w", record.Key, err)
	}
	var version int
	if err := json.Unmarshal(value["schema_version"], &version); err != nil || version != 1 {
		return fmt.Errorf("governance record %s has unsupported schema version", record.Key)
	}
	return nil
}

func isProducerOwner(owner string) bool {
	return owner == "FT71" || owner == "FT87" || owner == "FT88"
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
