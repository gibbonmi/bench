package compatibility

import "strings"

// Session binds live observations to one active chat and lifecycle interval.
type Session struct {
	ID, Epoch string
	Context   Context
}

// Observation retains the provenance of an actual capability result.
type Observation struct {
	Capability, Operation, Route, Provenance, Permission string
	Session                                              Session
	Success                                              bool
}

// SessionRequest supplies current context and in-memory observations only.
type SessionRequest struct {
	Session      Session
	Operation    string
	Observations []Observation
}

// SessionReport derives pending checks without publishing a reusable certificate.
func SessionReport(request SessionRequest) Report {
	report := Report{}
	operation := request.Operation
	known := false
	for _, item := range capabilities() {
		if required(item, operation) {
			known = true
		}
	}
	if !known {
		report.Checks = append(report.Checks, CheckRow{Check: "operation", State: StateUnknown, Action: "select an operation before dependent work: " + strings.Join(operations(), ", ")})
		return report
	}
	for _, item := range capabilities() {
		row := CheckRow{Check: item.name, State: StateNotRequired, Action: "none"}
		if required(item, operation) {
			row.State, row.Action = StateUnknown, item.action
			for _, observed := range request.Observations {
				if observed.Capability != item.name || observed.Operation != operation {
					continue
				}
				if !sameSession(observed.Session, request.Session) {
					continue
				}
				if observed.Provenance != "actual-tool" || observed.Permission != "normal" || observed.Route == "" {
					continue
				}
				row.State, row.Action = StateFailed, item.action
				if observed.Success {
					row.State, row.Action = StateOK, "none"
				}
			}
			if row.State != StateOK {
				report.Live = append(report.Live, LiveRow{Capability: item.name, Action: item.action})
			}
		}
		report.Checks = append(report.Checks, row)
	}
	environment := strings.ToLower(request.Session.Context.Environment.Value)
	if strings.HasPrefix(environment, "windows") {
		report.Checks = append(report.Checks, CheckRow{Check: "supported-environment", State: StateFailed, Action: "native Windows Bench execution is outside this claim; the Windows desktop agent runs in WSL2"})
	}
	return report
}

func sameSession(previous, current Session) bool {
	if current.ID == "" || current.Epoch == "" || current.ID != previous.ID || current.Epoch != previous.Epoch {
		return false
	}
	if !comparableContext(current.Context) || !comparableContext(previous.Context) {
		return false
	}
	return Fingerprint(previous.Context) == Fingerprint(current.Context)
}

func comparableContext(context Context) bool {
	if _, ok := ParseInterface(string(context.Interface)); !ok {
		return false
	}
	for _, fact := range []Fact{context.Repository, context.Environment, context.ConfigurationHome, context.ActiveRuntime, context.PolicyProvenance} {
		if fact.Value == "" || fact.Value == "unknown" || fact.Source == "" || fact.Source == "unknown" {
			return false
		}
	}
	return true
}
