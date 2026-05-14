package models

// Action represents the action that was performed given the request to the API
type Action string

// Outcome represents the outcome of the action given the request to the API
type Outcome string

const (
	ActionCreate Action = "CREATE"
	ActionRead   Action = "READ"
	ActionUpdate Action = "UPDATE"
	ActionDelete Action = "DELETE"

	OutcomeSuccess Outcome = "success"
	OutcomeFailure Outcome = "failure"
)
