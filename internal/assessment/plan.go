package assessment

// Plan pins comparison inputs; its approval reference does not authenticate consent.
type Plan struct {
	Version          int         `json:"version"`
	ID               string      `json:"id"`
	Purpose          string      `json:"purpose"`
	Approval         Reference   `json:"approval"`
	Budget           Budget      `json:"budget"`
	Tasks            []PlanTask  `json:"tasks"`
	Conditions       []Condition `json:"conditions"`
	Variable         string      `json:"variable"`
	QualityTolerance *Tolerance  `json:"quality_tolerance"`
}
type Budget struct {
	Amount   *float64 `json:"amount"`
	Currency string   `json:"currency"`
}
type PlanTask struct {
	ID          string `json:"id"`
	Holdout     bool   `json:"holdout"`
	Repetitions int    `json:"repetitions"`
}
type Line struct {
	Model  string `json:"model"`
	Effort string `json:"effort"`
}
type Condition struct {
	ID           string            `json:"id"`
	Revision     string            `json:"revision"`
	Lines        map[string]Line   `json:"lines"`
	Harness      string            `json:"harness"`
	Limits       map[string]string `json:"limits"`
	Capabilities map[string]string `json:"capabilities"`
	Acceptance   []string          `json:"acceptance"`
	ReviewAxes   []string          `json:"review_axes"`
}
type Tolerance struct {
	MaxFailureRate *float64         `json:"max_failure_rate"`
	Measures       map[string]Range `json:"measures"`
}
type Range struct {
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
}

// Trial carries observations that the ordinary run schema cannot derive.
type Trial struct {
	PlanID       string            `json:"plan_id"`
	Repetition   *int              `json:"repetition,omitempty"`
	Harness      string            `json:"harness"`
	Limits       map[string]string `json:"limits"`
	Capabilities map[string]string `json:"capabilities"`
	Acceptance   []string          `json:"acceptance"`
	ReviewAxes   []string          `json:"review_axes"`
	Reference    Reference         `json:"reference"`
}
