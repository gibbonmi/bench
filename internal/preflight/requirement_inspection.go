package preflight

import (
	"fmt"
	"os"
	"path/filepath"
)

func inspectRequirements(root string, run RunEvidence, profile Profile) ([]requirementStatus, []evidenceDigest, string, error) {
	statuses := make([]requirementStatus, 0, len(requirements.Records))
	inputs := make([]evidenceDigest, 0, len(requirements.Records)+1)
	for _, record := range requirements.Records {
		status := requirementStatus{Key: record.Key, Owner: record.Owner, Schema: record.Schema, Requiredness: record.Requiredness, Status: "not_applicable"}
		if !containsProfile(record.Profiles, profile) {
			status.Reason = "requirement is not applicable to selected profile"
			statuses = append(statuses, status)
			continue
		}
		status.Applicable = true
		data, err := readRegular(filepath.Join(root, filepath.FromSlash(record.Path)))
		if os.IsNotExist(err) && record.Producer {
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
		if record.Requiredness == "conditional" && record.Producer {
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
	inputs = append(inputs, evidenceDigest{Path: "internal/preflight/requirements.json", SHA256: digest(requirementsJSON)})
	for _, rel := range releaseInputPaths {
		data, err := readRegular(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return nil, nil, "", fmt.Errorf("release input %s is unreadable: %w", rel, err)
		}
		inputs = append(inputs, evidenceDigest{Path: rel, SHA256: digest(data)})
	}
	return statuses, inputs, unsatisfied, nil
}
