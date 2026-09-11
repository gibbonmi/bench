package assessment

import "time"

type Measure struct {
	Value     *float64  `json:"value,omitempty"`
	Reference Reference `json:"reference"`
}

type Reference struct {
	Producer string `json:"producer"`
	Native   string `json:"native"`
}

type Event struct {
	EventID   string    `json:"event_id"`
	SessionID string    `json:"session_id"`
	Epoch     int       `json:"epoch"`
	Sequence  int       `json:"sequence"`
	Mode      string    `json:"mode"`
	Counter   string    `json:"counter"`
	Reference Reference `json:"reference"`
	Usage     Usage     `json:"usage"`
}

type Rates struct {
	InputUncached *float64 `json:"input_uncached,omitempty"`
	InputCached   *float64 `json:"input_cached,omitempty"`
	Output        *float64 `json:"output,omitempty"`
	UnitScale     float64  `json:"unit_scale"`
	Currency      string   `json:"currency"`
	Source        string   `json:"source"`
	Date          string   `json:"date"`
	Conditions    string   `json:"conditions"`
}

type Charge struct {
	Kind      string    `json:"kind"`
	Amount    *float64  `json:"amount,omitempty"`
	Currency  string    `json:"currency"`
	Reference Reference `json:"reference"`
}

type Cost struct {
	Estimated *Rates   `json:"estimated,omitempty"`
	Other     []Charge `json:"other,omitempty"`
	Actual    []Charge `json:"actual,omitempty"`
}

type Attempt struct {
	AttemptID     string      `json:"attempt_id"`
	ChunkID       string      `json:"chunk_id"`
	Role          string      `json:"role"`
	SessionID     string      `json:"session_id"`
	Model         string      `json:"model"`
	Effort        string      `json:"effort"`
	State         string      `json:"state"`
	TimeReference *Reference  `json:"time_reference,omitempty"`
	StartedAt     *time.Time  `json:"started_at,omitempty"`
	EndedAt       *time.Time  `json:"ended_at,omitempty"`
	Usage         []Event     `json:"usage"`
	Cost          Cost        `json:"cost"`
	Evidence      []Reference `json:"evidence"`
}

type Run struct {
	Version       int                `json:"version"`
	RunID         string             `json:"run_id"`
	RepoKey       string             `json:"repo_key"`
	Source        string             `json:"source"`
	Condition     string             `json:"condition"`
	TaskID        string             `json:"task_id"`
	Holdout       bool               `json:"holdout"`
	TimeReference *Reference         `json:"time_reference,omitempty"`
	StartedAt     *time.Time         `json:"started_at,omitempty"`
	EndedAt       *time.Time         `json:"ended_at,omitempty"`
	State         string             `json:"state"`
	Attempts      []Attempt          `json:"attempts"`
	Evidence      []Reference        `json:"evidence"`
	Quality       map[string]Measure `json:"quality"`
}
