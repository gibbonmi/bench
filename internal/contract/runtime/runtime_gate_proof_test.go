package runtime

import (
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/gate"
)

const ft78ProofLedgerFailure = "FT78 proof ledger completeness contract failed"

type ft78ProofCase struct {
	id     string
	driver func(*testing.T)
}

var ft78Story2Proofs = []ft78ProofCase{
	manifestProof("R1/manifest-absent", "", "gate input manifest absent", 0),
	manifestProof("R1/manifest-empty", "", "gate input manifest invalid", 1),
	manifestProof("R1/manifest-malformed", `{`, "gate input manifest invalid", 1),
	manifestProof("R1/manifest-wrong-schema", `{"schema":2,"closure":"local","environment":[],"paths":[],"tools":[]}`, "gate input manifest invalid", 1),
	manifestProof("R1/manifest-remote", `{"schema":1,"closure":"remote","environment":[],"paths":[],"tools":[]}`, "remote oracle", 1),
	manifestProof("R1/manifest-missing-variable", `{"schema":1,"closure":"local","environment":["FT78_MISSING"],"paths":[],"tools":[]}`, "declared environment unavailable", 1),
	manifestProof("R1/manifest-missing-path", `{"schema":1,"closure":"local","environment":[],"paths":["missing"],"tools":[]}`, "declared path unavailable", 1),
	manifestProof("R1/manifest-missing-tool", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":["ft78-missing"]}`, "declared tool unavailable", 1),
	{id: "R1/manifest-escaped-symlink", driver: escapedManifestPathProof},
	{id: "R1/manifest-entry-limit", driver: manifestEntryLimitProof},
	{id: "R1/manifest-byte-limit-16384", driver: func(t *testing.T) { manifestByteLimitProof(t, 16_384, gate.Ready, "") }},
	{id: "R1/manifest-byte-limit-16385", driver: func(t *testing.T) { manifestByteLimitProof(t, 16_385, gate.Ready, "gate input manifest invalid") }},
	manifestProof("R1/manifest-decoded-control-byte", `{"schema":1,"closure":"local","environment":[],"paths":["bad\u0007path"],"tools":[]}`, "gate input manifest invalid", 1),

	mutationProof("R2/gate-script-mutation", mutationGateScript),
	mutationProof("R2/gate-interpreter-mutation", mutationGateInterpreter),
	mutationProof("R2/manifest-mutation", mutationManifest),
	mutationProof("R2/declared-tool-content-mutation", mutationToolContent),
	mutationProof("R2/declared-tool-mode-mutation", mutationToolMode),
	mutationProof("R2/declared-tool-target-mutation", mutationToolTarget),
	mutationProof("R2/declared-file-mutation", mutationDeclaredFile),
	mutationProof("R2/declared-directory-mutation", mutationDeclaredDirectory),
	mutationProof("R2/path-environment-mutation", mutationPATH),
	{id: "R2/repository-root-mutation", driver: repositoryRootMutationProof},
	mutationProof("R2/auto-detected-kind-mutation", mutationAutoKind),
	{id: "R2/ignored-input-rerun", driver: ignoredInputRerunProof},
	{id: "R2/ignored-input-refusal", driver: ignoredInputRefusalProof},

	passlistProof("R3/path-retained", "PATH", "/ft78/literal-path:/usr/bin:/bin", "/ft78/literal-path:/usr/bin:/bin", false),
	passlistProof("R3/bench-kit-absent", "BENCH_KIT", "subject-root", "unset", false),
	passlistProof("R3/bench-wrapper-absent", "BENCH_WRAPPER", "bench-wrapper", "unset", false),
	passlistProof("R3/ci-absent", "CI", "ft78-secret", "unset", false),
	passlistProof("R3/lang-absent", "LANG", "ft78-secret", "unset", false),
	passlistProof("R3/lc-all-absent", "LC_ALL", "C", "unset", false),
	passlistProof("R3/lc-ctype-absent", "LC_CTYPE", "C", "unset", false),
	passlistProof("R3/arbitrary-inherited-name-absent", "FT78_AMBIENT", "ft78-secret", "unset", false),
	passlistProof("R3/declared-absent-variable", "FT78_DECLARED_ABSENT", "", "unset", true),

	resolverPathProof("R4/real-pnpm-path", "pnpm-lock.yaml", "pnpm"),
	resolverPathProof("R4/real-python-path", "pyproject.toml", "mypy"),
	resolverPathProof("R4/real-cargo-path", "Cargo.toml", "cargo"),
	resolverMutationProof("R4/gate-sh-launcher-mutation", resolverGateSh, "gate.sh"),
	resolverMutationProof("R4/gate-sh-interpreter-mutation", resolverGateSh, "gate-interpreter"),
	resolverMutationProof("R4/bench-gate-bash-mutation", resolverBenchGate, "bash"),
	resolverMutationProof("R4/pnpm-bash-mutation", resolverPnpm, "bash"),
	resolverMutationProof("R4/pnpm-tool-mutation", resolverPnpm, "pnpm"),
	resolverMutationProof("R4/npm-bash-mutation", resolverNpm, "bash"),
	resolverMutationProof("R4/npm-tool-mutation", resolverNpm, "npm"),
	resolverMutationProof("R4/python-bash-mutation", resolverPython, "bash"),
	resolverMutationProof("R4/python-mypy-mutation", resolverPython, "mypy"),
	resolverMutationProof("R4/python-pytest-mutation", resolverPython, "pytest"),
	resolverMutationProof("R4/python-ruff-mutation", resolverPython, "ruff"),
	resolverMutationProof("R4/cargo-bash-mutation", resolverCargo, "bash"),
	resolverMutationProof("R4/cargo-tool-mutation", resolverCargo, "cargo"),
	resolverMutationProof("R4/cargo-rustc-mutation", resolverCargo, "rustc"),
	resolverMutationProof("R4/cargo-clippy-driver-mutation", resolverCargo, "clippy-driver"),
	{id: "R4/full-real-wrapper-precedence", driver: fullResolverPrecedenceProof},
	{id: "R4/no-gate-exit-3-no-record", driver: noGateResolverProof},
}

var ft78Story2ExpectedIDs = []string{
	"R1/manifest-absent", "R1/manifest-empty", "R1/manifest-malformed", "R1/manifest-wrong-schema", "R1/manifest-remote", "R1/manifest-missing-variable", "R1/manifest-missing-path", "R1/manifest-missing-tool", "R1/manifest-escaped-symlink", "R1/manifest-entry-limit", "R1/manifest-byte-limit-16384", "R1/manifest-byte-limit-16385", "R1/manifest-decoded-control-byte",
	"R2/gate-script-mutation", "R2/gate-interpreter-mutation", "R2/manifest-mutation", "R2/declared-tool-content-mutation", "R2/declared-tool-mode-mutation", "R2/declared-tool-target-mutation", "R2/declared-file-mutation", "R2/declared-directory-mutation", "R2/path-environment-mutation", "R2/repository-root-mutation", "R2/auto-detected-kind-mutation", "R2/ignored-input-rerun", "R2/ignored-input-refusal",
	"R3/path-retained", "R3/bench-kit-absent", "R3/bench-wrapper-absent", "R3/ci-absent", "R3/lang-absent", "R3/lc-all-absent", "R3/lc-ctype-absent", "R3/arbitrary-inherited-name-absent", "R3/declared-absent-variable",
	"R4/real-pnpm-path", "R4/real-python-path", "R4/real-cargo-path", "R4/gate-sh-launcher-mutation", "R4/gate-sh-interpreter-mutation", "R4/bench-gate-bash-mutation", "R4/pnpm-bash-mutation", "R4/pnpm-tool-mutation", "R4/npm-bash-mutation", "R4/npm-tool-mutation", "R4/python-bash-mutation", "R4/python-mypy-mutation", "R4/python-pytest-mutation", "R4/python-ruff-mutation", "R4/cargo-bash-mutation", "R4/cargo-tool-mutation", "R4/cargo-rustc-mutation", "R4/cargo-clippy-driver-mutation", "R4/full-real-wrapper-precedence", "R4/no-gate-exit-3-no-record",
}

func TestFT78Story2ProofLedgerCompleteness(t *testing.T) {
	contract.NoteContractFailure(t, ft78ProofLedgerFailure)
	seen := map[string]int{}
	for _, proof := range ft78Story2Proofs {
		seen[proof.id]++
		if proof.driver == nil {
			t.Fatalf("%s: nil real-wrapper driver", proof.id)
		}
	}
	if len(seen) != len(ft78Story2ExpectedIDs) {
		t.Fatalf("registered IDs = %d, want %d", len(seen), len(ft78Story2ExpectedIDs))
	}
	for _, id := range ft78Story2ExpectedIDs {
		if seen[id] != 1 {
			t.Fatalf("%s registrations = %d, want 1", id, seen[id])
		}
	}
}

func TestFT78Story2ProofLedger(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	for _, proof := range ft78Story2Proofs {
		proof := proof
		t.Run(proof.id, proof.driver)
	}
}
