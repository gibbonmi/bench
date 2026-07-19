package publication

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	if err := VerifyPublishAuthority(root, profile); err != nil {
		return Record{}, err
	}
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

// Exact next_action values the staged path hands back to the human running
// 2FA approval, and the terminal value once the complete set is live.
const (
	nextActionApprovePlatforms = "approve-platform-packages"
	nextActionApproveWrapper   = "approve-wrapper"
	nextActionPromote          = "promote"
)

// alreadyStaged reports the stage id a prior run recorded for pkg@version, if
// any — the record is the one source of truth for resume, never re-derived
// from the registry (the fixture has no "list pending stages" query).
func alreadyStaged(record Record, pkg ApprovedPackage) (stageID string, found bool) {
	for _, t := range record.Transitions {
		if t.Package == pkg.Name && t.Version == pkg.Version && t.Action == "stage" && t.Result == "success" {
			stageID = t.StageID
			found = true
		}
	}
	return stageID, found
}

// stageAndVerify stages pkg (or resumes a prior stage id from the record) and,
// when it is already live at the registry, verifies its integrity against the
// approved local tarball and makes sure tag points at it. verified is true
// only once the package is confirmed live with matching integrity — the one
// signal callers use to decide whether a step is still pending 2FA approval.
// A live package whose integrity does not match is a terminal mismatch,
// exactly like the first-publication path.
func stageAndVerify(ctx context.Context, root string, record *Record, registry Registry, pkg ApprovedPackage, tag string) (verified bool, err error) {
	registryIntegrity, live, err := registry.Integrity(ctx, pkg.Name, pkg.Version)
	if err != nil {
		return false, err
	}
	if live {
		if registryIntegrity != pkg.Integrity {
			record.Transitions = append(record.Transitions, Transition{
				Package: pkg.Name, Version: pkg.Version, Action: "verify", AuthMode: "approval",
				LocalIntegrity: pkg.Integrity, RegistryIntegrity: registryIntegrity, TagState: tag,
				Result: "mismatch", Timestamp: clock().Format(time.RFC3339),
			})
			record.Result = "failed"
			if saveErr := SaveRecord(root, *record); saveErr != nil {
				return false, saveErr
			}
			return false, fmt.Errorf("registry integrity for %s@%s does not match the approved local tarball; publication stopped", pkg.Name, pkg.Version)
		}
		if err := registry.TagAdd(ctx, pkg.Name, tag, pkg.Version); err != nil {
			return false, fmt.Errorf("could not point %s at %s: %w", tag, pkg.Name, err)
		}
		record.Transitions = append(record.Transitions, Transition{
			Package: pkg.Name, Version: pkg.Version, Action: "verify", AuthMode: "approval",
			LocalIntegrity: pkg.Integrity, RegistryIntegrity: registryIntegrity, TagState: tag,
			Result: "success", Timestamp: clock().Format(time.RFC3339),
		})
		return true, SaveRecord(root, *record)
	}

	if _, found := alreadyStaged(*record, pkg); found {
		return false, nil
	}
	tarball, err := os.ReadFile(pkg.FilePath)
	if err != nil {
		return false, fmt.Errorf("approved artifact %s is unreadable: %w", pkg.Name, err)
	}
	stageID, err := registry.StageSubmit(ctx, pkg.Name, pkg.Version, tarball)
	record.Transitions = append(record.Transitions, Transition{
		Package: pkg.Name, Version: pkg.Version, Action: "stage", AuthMode: "oidc-stage", StageID: stageID,
		LocalIntegrity: pkg.Integrity, TagState: tag,
		Result: resultOf(err), Timestamp: clock().Format(time.RFC3339),
	})
	if err != nil {
		record.Result = "failed"
		_ = SaveRecord(root, *record)
		return false, fmt.Errorf("stage-submit failed for %s@%s: %w", pkg.Name, pkg.Version, err)
	}
	return false, SaveRecord(root, *record)
}

// RunStagedPublication drives the resumable subsequent-release (staged) path:
// every package — platforms and the wrapper — is stage-submitted under
// auth_mode oidc-stage, which never makes anything live (the fixture's stage
// endpoint only stores pending bytes). Nothing advances past staging on its
// own; a human runs the registry's out-of-band 2FA approval, and a rerun of
// this function observes the resulting live state and advances. The wrapper
// is never treated as approved until every platform package is verified live
// first — if the registry reports the wrapper live while a platform is still
// pending, that is an ordering violation and a terminal failure, never a
// silent skip. Once every package is live and verified, next_action is
// "promote"; the actual dist-tag promotion is a separate, explicit step
// (RunPromotion) so promoting is never an accidental side effect of resuming
// submit/status.
func RunStagedPublication(ctx context.Context, root, version, profile string, registry Registry) (Record, string, error) {
	if err := VerifyPublishAuthority(root, profile); err != nil {
		return Record{}, "", err
	}
	releaseIndexSHA256, approved, err := VerifyApprovedSet(root, version)
	if err != nil {
		return Record{}, "", err
	}
	record, err := LoadRecord(root)
	if err != nil {
		return Record{}, "", err
	}
	if record.SchemaVersion == 0 {
		record = Record{Path: "public", Profile: profile, Result: "in_progress"}
		for _, pkg := range approved {
			record.Provenance = append(record.Provenance, Provenance{Package: pkg.Name, SHA256: pkg.SHA256})
		}
	} else if record.ReleaseIndexSHA256 != releaseIndexSHA256 {
		return Record{}, "", fmt.Errorf("publication record was built against a different approved release index; resume is unsafe")
	}
	record.ReleaseIndexSHA256 = releaseIndexSHA256

	if record.Result == "success" {
		return record, "release-complete", nil
	}

	tag := CandidateTag(version)
	ordered := orderedForFirstPublication(approved)
	if len(ordered) == 0 {
		return record, "", fmt.Errorf("release plan named no publication artifacts for version %s", version)
	}
	wrapper := ordered[len(ordered)-1]
	platforms := ordered[:len(ordered)-1]

	platformsVerified := true
	for _, pkg := range platforms {
		verified, err := stageAndVerify(ctx, root, &record, registry, pkg, tag)
		if err != nil {
			return record, "", err
		}
		if !verified {
			platformsVerified = false
		}
	}

	if !platformsVerified {
		// The wrapper must never be approved ahead of every platform package.
		// Detect the registry already reporting it live in that state as a
		// hard ordering violation rather than silently accepting it.
		_, wrapperLive, err := registry.Integrity(ctx, wrapper.Name, wrapper.Version)
		if err != nil {
			return record, "", err
		}
		if wrapperLive {
			record.Transitions = append(record.Transitions, Transition{
				Package: wrapper.Name, Version: wrapper.Version, Action: "verify", AuthMode: "approval",
				TagState: tag, Result: "order-violation", Timestamp: clock().Format(time.RFC3339),
			})
			record.Result = "failed"
			if err := SaveRecord(root, record); err != nil {
				return record, "", err
			}
			return record, "", fmt.Errorf("wrapper %s was approved before every platform package; publication stopped", wrapper.Name)
		}
		if _, found := alreadyStaged(record, wrapper); !found {
			tarball, err := os.ReadFile(wrapper.FilePath)
			if err != nil {
				return record, "", fmt.Errorf("approved artifact %s is unreadable: %w", wrapper.Name, err)
			}
			stageID, err := registry.StageSubmit(ctx, wrapper.Name, wrapper.Version, tarball)
			record.Transitions = append(record.Transitions, Transition{
				Package: wrapper.Name, Version: wrapper.Version, Action: "stage", AuthMode: "oidc-stage", StageID: stageID,
				LocalIntegrity: wrapper.Integrity, TagState: tag,
				Result: resultOf(err), Timestamp: clock().Format(time.RFC3339),
			})
			if err != nil {
				record.Result = "failed"
				_ = SaveRecord(root, record)
				return record, "", fmt.Errorf("stage-submit failed for %s@%s: %w", wrapper.Name, wrapper.Version, err)
			}
			if err := SaveRecord(root, record); err != nil {
				return record, "", err
			}
		}
		return record, nextActionApprovePlatforms, nil
	}

	wrapperVerified, err := stageAndVerify(ctx, root, &record, registry, wrapper, tag)
	if err != nil {
		return record, "", err
	}
	if !wrapperVerified {
		return record, nextActionApproveWrapper, nil
	}

	return record, nextActionPromote, nil
}

// RunPromotion moves the "latest" dist-tag onto version, platform packages
// first and the wrapper strictly last, but only once the complete approved
// set reverifies live at the registry with matching integrity — a fresh
// check, never trusted from an earlier transition. It refuses to run at all
// unless every package is confirmed, so promoting a partial or unverified set
// (wrapper included) is not reachable through this function.
func RunPromotion(ctx context.Context, root, version, profile string, registry Registry) (Record, error) {
	releaseIndexSHA256, approved, err := VerifyApprovedSet(root, version)
	if err != nil {
		return Record{}, err
	}
	record, err := LoadRecord(root)
	if err != nil {
		return Record{}, err
	}
	if record.SchemaVersion == 0 {
		return Record{}, fmt.Errorf("no publication record for %s; run release submit first", version)
	}
	if record.ReleaseIndexSHA256 != releaseIndexSHA256 {
		return Record{}, fmt.Errorf("publication record was built against a different approved release index; promote is unsafe")
	}
	if record.Result == "success" {
		return record, nil
	}
	if record.Result == "failed" {
		return record, fmt.Errorf("publication record is in a failed state; roll back before promoting")
	}

	ordered := orderedForFirstPublication(approved)
	for _, pkg := range ordered {
		registryIntegrity, live, err := registry.Integrity(ctx, pkg.Name, pkg.Version)
		if err != nil {
			return record, err
		}
		if !live || registryIntegrity != pkg.Integrity {
			return record, fmt.Errorf("release is not fully verified live for %s; cannot promote", pkg.Name)
		}
	}

	for _, pkg := range ordered {
		err := registry.TagAdd(ctx, pkg.Name, "latest", pkg.Version)
		record.Transitions = append(record.Transitions, Transition{
			Package: pkg.Name, Version: pkg.Version, Action: "promote", AuthMode: "approval",
			LocalIntegrity: pkg.Integrity, TagState: "latest",
			Result: resultOf(err), Timestamp: clock().Format(time.RFC3339),
		})
		if err != nil {
			record.Result = "failed"
			_ = SaveRecord(root, record)
			return record, fmt.Errorf("promote failed for %s@%s: %w", pkg.Name, pkg.Version, err)
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

// RunRollback recovers from a partial or failed release: it removes the
// candidate tag from every approved package (never touching "latest", which
// it never even reads), deprecates the bad version with message, and never
// issues an unpublish call — there is no such operation in the Registry port.
// It tolerates a tag that was never set (a package that failed before it was
// ever tagged) and a version that never went live (nothing to deprecate).
func RunRollback(ctx context.Context, root, version, profile, message string, registry Registry) (Record, error) {
	_, approved, err := VerifyApprovedSet(root, version)
	if err != nil {
		return Record{}, err
	}
	record, err := LoadRecord(root)
	if err != nil {
		return Record{}, err
	}
	if record.SchemaVersion == 0 {
		return Record{}, fmt.Errorf("no publication record for %s; nothing to roll back", version)
	}

	tag := CandidateTag(version)
	for _, pkg := range approved {
		err := registry.TagRemove(ctx, pkg.Name, tag)
		result := resultOf(err)
		if err != nil && strings.Contains(err.Error(), "tag not set") {
			result = "absent"
			err = nil
		}
		record.Transitions = append(record.Transitions, Transition{
			Package: pkg.Name, Version: pkg.Version, Action: "tag-remove", AuthMode: "approval",
			TagState: tag, Result: result, Timestamp: clock().Format(time.RFC3339),
		})
		if err != nil {
			record.Result = "failed"
			_ = SaveRecord(root, record)
			return record, fmt.Errorf("tag-remove failed for %s: %w", pkg.Name, err)
		}

		_, live, integrityErr := registry.Integrity(ctx, pkg.Name, pkg.Version)
		if integrityErr != nil {
			return record, integrityErr
		}
		if !live {
			continue
		}
		deprecateErr := registry.Deprecate(ctx, pkg.Name, pkg.Version, message)
		record.Transitions = append(record.Transitions, Transition{
			Package: pkg.Name, Version: pkg.Version, Action: "deprecate", AuthMode: "approval",
			Result: resultOf(deprecateErr), Timestamp: clock().Format(time.RFC3339),
		})
		if deprecateErr != nil {
			record.Result = "failed"
			_ = SaveRecord(root, record)
			return record, fmt.Errorf("deprecate failed for %s@%s: %w", pkg.Name, pkg.Version, deprecateErr)
		}
		if err := SaveRecord(root, record); err != nil {
			return record, err
		}
	}

	record.Result = "rolled_back"
	if err := SaveRecord(root, record); err != nil {
		return record, err
	}
	return record, nil
}
