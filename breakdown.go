package onfido

const (
	BreakdownClear    BreakdownResult = "clear"
	BreakdownConsider BreakdownResult = "consider"

	SubBreakdownClear        BreakdownSubResult = "clear"
	SubBreakdownConsider     BreakdownSubResult = "consider"
	SubBreakdownUnidentified BreakdownSubResult = "unidentified"
)

// BreakdownResult represents a report's breakdown result
type BreakdownResult string

// BreakdownSubResult represents a report's sub-breakdown result
type BreakdownSubResult string

type Breakdowns map[string]Breakdown

type Breakdown struct {
	Result        *BreakdownResult `json:"result"`
	SubBreakdowns SubBreakdowns    `json:"breakdown"`
}

type SubBreakdowns map[string]SubBreakdown

type SubBreakdown struct {
	Result     *BreakdownSubResult `json:"result"`
	Properties Properties          `json:"properties"`
}
