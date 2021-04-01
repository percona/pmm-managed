package models

//go:generate reform

type Interval string

const (
	Standard Interval = "standard"
	Frequent Interval = "frequent"
	Rare     Interval = "rare"
)

// ChecksState represents any changes to an STT check loaded in pmm-managed.
//reform:checks_state
type ChecksState struct {
	Name     string   `reform:"name,pk"`
	Interval Interval `reform:"interval"`
}
