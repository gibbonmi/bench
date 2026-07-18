package releaseevidence

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type ReleaseIntentError struct{ Message string }

func (e *ReleaseIntentError) Error() string { return e.Message }

func FinalizeEvidence(ctx context.Context, root string, run RunEvidence) error {
	if err := validateRun(root, run); err != nil {
		return err
	}
	if ctx.Err() != nil || TerminalStatus(run.Phases) == StatusInterrupted {
		return context.Canceled
	}
	if run.Scope == ScopeFocused {
		if run.Mode == ModePublish {
			return &ReleaseIntentError{Message: "focused publish runs cannot authorize publication"}
		}
		return nil
	}

	built, err := assembleReleaseEvidence(ctx, root, run)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	current, err := inputFingerprint(root, run)
	if err != nil {
		return fmt.Errorf("release evidence inputs became unsafe: %w", err)
	}
	if current != built.fingerprint {
		return errors.New("release evidence input drift detected before promotion")
	}
	indexBytes, err := indexEncoder(built.index)
	if err != nil {
		return fmt.Errorf("release index encoding failed: %w", err)
	}
	sums := deriveChecksums(built.index)
	manifest := manifestFor(run)
	manifest.Status = built.index.Status
	if err := PromoteEvidenceFiles(root, run.Mode, run.Phases, manifest, map[string][]byte{
		"release-index.json": indexBytes,
		"SHA256SUMS":         sums,
	}); err != nil {
		return fmt.Errorf("could not promote complete release evidence: %w", err)
	}
	if built.unsatisfied != "" {
		return &ReleaseIntentError{Message: built.unsatisfied}
	}
	return nil
}

func manifestFor(run RunEvidence) Manifest {
	return Manifest{SchemaVersion: 1, Mode: run.Mode, Scope: run.Scope, Status: TerminalStatus(run.Phases), Identity: run.Identity, Phases: PhaseSummaries(run.Phases)}
}

type assembledEvidence struct {
	index       Index
	fingerprint string
	unsatisfied string
}

func validateRun(root string, run RunEvidence) error {
	if root == "" || (run.Mode != ModeVerify && run.Mode != ModePublish) || (run.Scope != ScopePreflight && run.Scope != ScopeFocused) {
		return errors.New("invalid release evidence run")
	}
	if run.Mode == ModePublish && run.Scope == ScopePreflight && run.Profile != ProfilePublic && run.Profile != ProfileBank {
		return errors.New("publish requires an explicit profile")
	}
	if run.Profile != "" && run.Profile != ProfilePublic && run.Profile != ProfileBank {
		return fmt.Errorf("unknown release profile %q", run.Profile)
	}
	for label, value := range map[string]*string{"tag": run.Identity.Tag, "package_version": run.Identity.PackageVersion, "source_commit": run.Identity.SourceCommit, "binary_version": run.Identity.BinaryVersion, "changelog_heading": run.Identity.ChangelogHeading, "toolchain": run.Identity.Toolchain} {
		if value != nil && hasControlBytes(*value) {
			return fmt.Errorf("release identity %s contains control bytes", label)
		}
	}
	if err := validateRequirementRegistry(requirements); err != nil {
		return err
	}
	want := PhaseNames(run.Mode)
	if run.Scope == ScopeFocused {
		if len(run.Phases) != 1 || !Contains(want, run.Phases[0].Name) {
			return errors.New("focused release evidence must contain one registered phase")
		}
	} else if len(run.Phases) != len(want) {
		return errors.New("release evidence phases do not match the registered phase set")
	} else {
		for i, result := range run.Phases {
			if result.Name != want[i] {
				return errors.New("release evidence phases are not in registry order")
			}
		}
	}
	for _, result := range run.Phases {
		switch result.Status {
		case StatusGreen, StatusRed, StatusNotRun, StatusInterrupted:
		default:
			return fmt.Errorf("unknown release phase status %q", result.Status)
		}
	}
	return nil
}

func assembleReleaseEvidence(ctx context.Context, root string, run RunEvidence) (assembledEvidence, error) {
	if err := ctx.Err(); err != nil {
		return assembledEvidence{}, err
	}
	effectiveProfile := run.Profile
	if effectiveProfile == "" {
		effectiveProfile = ProfilePublic
	}
	statuses, inputs, unsatisfied, err := inspectRequirements(root, run, effectiveProfile)
	if err != nil {
		return assembledEvidence{}, err
	}
	rollbackTarget, err := readRollbackTarget(root)
	if err != nil {
		return assembledEvidence{}, err
	}
	flags := phaseFlags(run)
	artifacts, targets, artifactFingerprint, err := inspectArtifacts(root)
	if err != nil {
		return assembledEvidence{}, err
	}
	if err := waitForEvidenceProbe(ctx); err != nil {
		return assembledEvidence{}, err
	}
	reproducibility, err := readReproducibility(root, artifacts)
	if err != nil {
		return assembledEvidence{}, err
	}
	nativeProofs, err := inspectNativeProofs(root, targets, artifacts)
	if err != nil {
		return assembledEvidence{}, err
	}
	if run.Mode == ModePublish && len(nativeProofs) != len(targets) {
		if unsatisfied != "" {
			unsatisfied += "; "
		}
		unsatisfied += "native target proof is incomplete"
	}
	currentArtifacts, err := fingerprintArtifactSet(root)
	if err != nil || currentArtifacts != artifactFingerprint {
		return assembledEvidence{}, errors.New("release evidence artifact drift detected during assembly")
	}
	toolchains, err := observeToolchains(ctx, root)
	if err != nil {
		return assembledEvidence{}, err
	}
	phases := make([]phaseEvidence, 0, len(run.Phases))
	for _, result := range run.Phases {
		record := Record{SchemaVersion: 1, Phase: result.Name, Mode: run.Mode, Status: result.Status, ExitCode: result.ExitCode, Error: result.Failure}
		data, err := canonicalJSON(record)
		if err != nil {
			return assembledEvidence{}, err
		}
		phases = append(phases, phaseEvidence{Name: result.Name, Status: result.Status, Digest: digest(data)})
	}
	sort.Slice(inputs, func(i, j int) bool { return inputs[i].Path < inputs[j].Path })
	sort.Slice(artifacts, func(i, j int) bool { return artifacts[i].Name < artifacts[j].Name })
	sort.Slice(targets, func(i, j int) bool {
		if targets[i].OS != targets[j].OS {
			return targets[i].OS < targets[j].OS
		}
		return targets[i].Arch < targets[j].Arch
	})
	status := TerminalStatus(run.Phases)
	if status == StatusGreen && unsatisfied != "" {
		status = StatusRed
	}
	index := Index{
		SchemaVersion:   1,
		Mode:            run.Mode,
		Scope:           run.Scope,
		Profile:         run.Profile,
		Status:          status,
		Identity:        releaseIdentityFrom(run.Identity),
		RollbackTarget:  rollbackTarget,
		Flags:           flags,
		Toolchains:      toolchains,
		Inputs:          inputs,
		Targets:         targets,
		Phases:          phases,
		Requirements:    statuses,
		Artifacts:       artifacts,
		Reproducibility: reproducibility,
		NativeProofs:    nativeProofs,
	}
	fingerprint, err := inputFingerprint(root, run)
	if err != nil {
		return assembledEvidence{}, err
	}
	return assembledEvidence{index: index, fingerprint: fingerprint, unsatisfied: unsatisfied}, nil
}

func observeToolchains(ctx context.Context, root string) ([]toolchainEvidence, error) {
	out := make([]toolchainEvidence, 0, len(requirements.Toolchains))
	for _, requirement := range requirements.Toolchains {
		if len(requirement.VersionArgv) == 0 {
			return nil, fmt.Errorf("toolchain %s has no version command", requirement.Name)
		}
		command := exec.CommandContext(ctx, requirement.VersionArgv[0], requirement.VersionArgv[1:]...)
		command.Dir = root
		version, err := command.Output()
		if err != nil || strings.TrimSpace(string(version)) == "" {
			return nil, fmt.Errorf("toolchain %s version is unavailable", requirement.Name)
		}
		flags := make([]string, 0)
		operations := make([]string, 0, len(requirement.Operations))
		for operation := range requirement.Operations {
			operations = append(operations, operation)
		}
		sort.Strings(operations)
		for _, operation := range operations {
			flags = append(flags, operation+":"+strings.Join(requirement.Operations[operation], " "))
		}
		out = append(out, toolchainEvidence{Name: requirement.Name, Version: strings.TrimSpace(string(version)), Flags: flags})
	}
	return out, nil
}

func phaseFlags(run RunEvidence) []string {
	flags := make([]string, 0, len(run.Phases))
	for _, result := range run.Phases {
		definition, ok := PhaseDefinitionFor(result.Name)
		if !ok {
			continue
		}
		flags = append(flags, result.Name+"="+strings.Join(definition.Argv, " "))
	}
	return flags
}

func waitForEvidenceProbe(ctx context.Context) error {
	path := os.Getenv("BENCH_PREFLIGHT_EVIDENCE_READY_FILE")
	if path == "" {
		return nil
	}
	if err := os.WriteFile(path, []byte("ready\n"), 0o644); err != nil {
		return err
	}
	for {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func TerminalStatus(results []Result) Status {
	status := StatusGreen
	for _, result := range results {
		if result.Status == StatusInterrupted {
			return StatusInterrupted
		}
		if result.Status == StatusRed || result.Status == StatusNotRun {
			status = StatusRed
		}
	}
	return status
}

func Contains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}
