package publication

import (
	"context"
	"fmt"
	"os"
	"time"
)

// clock is overridable for deterministic tests; production uses time.Now.
var clock = func() time.Time { return time.Now().UTC() }

// CandidateTag names the version-specific non-default dist-tag a first
// publication uses: candidate-<version>. Direct-publishing under this tag
// (never "latest") keeps every platform package invisible to a plain
// `npm install` until an explicit later promotion step moves the tag —
// callers never need to know this naming, only RunFirstPublication does.
func CandidateTag(version string) string {
	return "candidate-" + version
}

// orderedForFirstPublication imposes the one publish order the first-
// publication path requires: every platform package (in whatever order
// release-plan.mjs names them), then the wrapper last. release-plan.mjs's own
// artifact-records order sorts the wrapper file name *before* the platform
// file names (a leading digit sorts before a letter), so this reordering is
// the actual publication-order policy — it must not be left to file-name sort.
func orderedForFirstPublication(packages []ApprovedPackage) []ApprovedPackage {
	ordered := make([]ApprovedPackage, 0, len(packages))
	var wrapper *ApprovedPackage
	for i := range packages {
		if packages[i].Kind == "wrapper" {
			w := packages[i]
			wrapper = &w
			continue
		}
		ordered = append(ordered, packages[i])
	}
	if wrapper != nil {
		ordered = append(ordered, *wrapper)
	}
	return ordered
}

// RunFirstPublication drives the resumable first-publication state machine: it
// locally verifies the complete approved set, then direct-publishes every
// platform package under CandidateTag(version), verifying each live registry
// integrity against the approved local tarball before advancing, and finally
// publishes and verifies the wrapper last. A retry (record already on disk)
// resumes: a package already live with matching registry integrity is treated
// as complete and is never republished (idempotent); any live package whose
// registry integrity does not exactly match the approved local tarball is a
// terminal failure that stops the whole release without touching the rest of
// the build set. Callers never see this ordering or resume logic — they only
// ever get back the durable Record.
func RunFirstPublication(ctx context.Context, root, version, profile string, registry Registry) (Record, error) {
	releaseIndexSHA256, approved, err := VerifyApprovedSet(root, version)
	if err != nil {
		return Record{}, err
	}
	record, err := LoadRecord(root)
	if err != nil {
		return Record{}, err
	}
	if record.SchemaVersion == 0 {
		record = Record{Path: "first", Profile: profile, Result: "in_progress"}
		for _, pkg := range approved {
			record.Provenance = append(record.Provenance, Provenance{Package: pkg.Name, SHA256: pkg.SHA256})
		}
	} else if record.ReleaseIndexSHA256 != releaseIndexSHA256 {
		return Record{}, fmt.Errorf("publication record was built against a different approved release index; resume is unsafe")
	}
	record.ReleaseIndexSHA256 = releaseIndexSHA256

	if record.Result == "success" {
		// Already complete: resuming a finished release is a no-op.
		return record, nil
	}

	tag := CandidateTag(version)

	for _, pkg := range orderedForFirstPublication(approved) {
		registryIntegrity, live, err := registry.Integrity(ctx, pkg.Name, pkg.Version)
		if err != nil {
			return record, err
		}
		if live {
			if registryIntegrity != pkg.Integrity {
				record.Transitions = append(record.Transitions, Transition{
					Package: pkg.Name, Version: pkg.Version, Action: "verify", AuthMode: "direct",
					LocalIntegrity: pkg.Integrity, RegistryIntegrity: registryIntegrity, TagState: tag,
					Result: "mismatch", Timestamp: clock().Format(time.RFC3339),
				})
				record.Result = "failed"
				if err := SaveRecord(root, record); err != nil {
					return record, err
				}
				return record, fmt.Errorf("registry integrity for %s@%s does not match the approved local tarball; publication stopped", pkg.Name, pkg.Version)
			}
			record.Transitions = append(record.Transitions, Transition{
				Package: pkg.Name, Version: pkg.Version, Action: "verify", AuthMode: "direct",
				LocalIntegrity: pkg.Integrity, RegistryIntegrity: registryIntegrity, TagState: tag,
				Result: "resumed", Timestamp: clock().Format(time.RFC3339),
			})
			if err := SaveRecord(root, record); err != nil {
				return record, err
			}
			continue
		}

		tarball, err := os.ReadFile(pkg.FilePath)
		if err != nil {
			return record, fmt.Errorf("approved artifact %s is unreadable: %w", pkg.Name, err)
		}
		publishIntegrity, err := registry.Publish(ctx, pkg.Name, pkg.Version, tag, tarball)
		record.Transitions = append(record.Transitions, Transition{
			Package: pkg.Name, Version: pkg.Version, Action: "publish", AuthMode: "direct",
			LocalIntegrity: pkg.Integrity, TagState: tag,
			Result: resultOf(err), Timestamp: clock().Format(time.RFC3339),
		})
		if err != nil {
			record.Result = "failed"
			_ = SaveRecord(root, record)
			return record, fmt.Errorf("publish failed for %s@%s: %w", pkg.Name, pkg.Version, err)
		}
		if err := SaveRecord(root, record); err != nil {
			return record, err
		}

		verifyIntegrity, live, err := registry.Integrity(ctx, pkg.Name, pkg.Version)
		if err != nil {
			return record, err
		}
		matched := live && verifyIntegrity == pkg.Integrity && verifyIntegrity == publishIntegrity
		result := "success"
		if !matched {
			result = "mismatch"
		}
		record.Transitions = append(record.Transitions, Transition{
			Package: pkg.Name, Version: pkg.Version, Action: "verify", AuthMode: "direct",
			LocalIntegrity: pkg.Integrity, RegistryIntegrity: verifyIntegrity, TagState: tag,
			Result: result, Timestamp: clock().Format(time.RFC3339),
		})
		if !matched {
			record.Result = "failed"
			if err := SaveRecord(root, record); err != nil {
				return record, err
			}
			return record, fmt.Errorf("registry integrity for %s@%s did not verify after publish; publication stopped", pkg.Name, pkg.Version)
		}
		if err := SaveRecord(root, record); err != nil {
			return record, err
		}
	}

	record.Result = "success"
	if err := SaveRecord(root, record); err != nil {
		return record, err
	}
	return record, nil
}

func resultOf(err error) string {
	if err != nil {
		return "failed"
	}
	return "success"
}
