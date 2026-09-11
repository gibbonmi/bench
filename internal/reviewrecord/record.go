// Package reviewrecord reads source-bound verification and review evidence.
package reviewrecord

import "fmt"

type NativeRef struct {
	Ref     string `json:"ref"`
	Digest  string `json:"digest"`
	Excerpt string `json:"excerpt"`
}
type Evidence struct {
	ID           string    `json:"id"`
	Performer    string    `json:"performer"`
	Role         string    `json:"role"`
	Model        string    `json:"model"`
	Effort       string    `json:"effort"`
	SourceDigest string    `json:"source_digest"`
	State        string    `json:"state"`
	Outcome      string    `json:"outcome"`
	NativeRef    NativeRef `json:"native_ref"`
}
type Probe struct {
	Mutation  string    `json:"mutation"`
	Outcome   string    `json:"outcome"`
	ExitCode  int       `json:"exit_code"`
	Restore   string    `json:"restore"`
	NativeRef NativeRef `json:"native_ref"`
}
type Verification struct {
	Evidence
	Requirement string `json:"requirement"`
	Command     string `json:"command"`
	ExitCode    *int   `json:"exit_code"`
	Probe       *Probe `json:"probe,omitempty"`
}
type Review struct {
	Evidence
	Axis       string   `json:"axis"`
	Base       string   `json:"base"`
	Tip        string   `json:"tip"`
	FindingIDs []string `json:"finding_ids"`
	Supersedes []string `json:"supersedes"`
}
type Chunk struct {
	ID             string         `json:"id"`
	Base           string         `json:"base"`
	Tip            string         `json:"tip"`
	PlanDigest     string         `json:"plan_digest"`
	SourceDigest   string         `json:"source_digest"`
	AcceptanceRows []string       `json:"acceptance_rows"`
	Verification   []Verification `json:"verification"`
	Reviews        []Review       `json:"reviews"`
}
type Completion struct {
	State          string            `json:"state"`
	SourceDigest   string            `json:"source_digest"`
	Performer      string            `json:"performer"`
	Reconciliation map[string]string `json:"reconciliation"`
	Verification   []Verification    `json:"verification"`
}
type Amendment struct {
	From     string              `json:"from"`
	To       string              `json:"to"`
	ChunkIDs map[string][]string `json:"chunk_ids"`
}
type Record struct {
	Version               int         `json:"version"`
	Spec                  string      `json:"spec"`
	PlanDigest            string      `json:"plan_digest"`
	ImplementationSession string      `json:"implementation_session"`
	Chunks                []Chunk     `json:"chunks"`
	Completion            Completion  `json:"completion"`
	Amendments            []Amendment `json:"amendments,omitempty"`
}

func CheckReviews(chunk Chunk, session string) error {
	for _, axis := range Axes() {
		var current *Review
		for i := range chunk.Reviews {
			if chunk.Reviews[i].Axis == axis {
				current = &chunk.Reviews[i]
			}
		}
		if current == nil {
			return fmt.Errorf("chunk %s: missing %s; record the native review result", chunk.ID, axis)
		}
		if current.Role != "independent-review" || current.Performer == session || current.Performer == "" {
			return fmt.Errorf("chunk %s: invalid %s performer or role; obtain independent review", chunk.ID, axis)
		}
		if current.State != "completed" || current.Outcome != "pass" || len(current.FindingIDs) != 0 {
			return fmt.Errorf("chunk %s: %s %s; resolve findings and obtain a completed result", chunk.ID, current.State, axis)
		}
		if current.SourceDigest != chunk.SourceDigest || current.Base != chunk.Base || current.Tip != chunk.Tip {
			return fmt.Errorf("chunk %s: stale %s; review the current frozen pair", chunk.ID, axis)
		}
	}
	return nil
}

func Axes() []string { return []string{"Standards", "Spec", "Coverage"} }
