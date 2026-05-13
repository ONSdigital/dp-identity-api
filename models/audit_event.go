package models

// Action represents the action that was performed on the policy
type Action string

// Outcome represents the outcome of the action
type Outcome string

const (
	ActionCreate Action = "CREATE"
	ActionRead   Action = "READ"
	ActionUpdate Action = "UPDATE"
	ActionDelete Action = "DELETE"

	OutcomeSuccess Outcome = "success"
	OutcomeFailure Outcome = "failure"
)
