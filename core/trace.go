package core

// VTrace determines the type of visual tracing
type VTrace int

const (
	// VTraceGrid shows the page grid
	VTraceGrid VTrace = iota
)

// SetVTrace sets the visual tracing
func (xd *xtsDocument) SetVTrace(t VTrace) {
	xd.tracing |= 1 << t
}

// IsTrace returns true if tracing t is set
func (xd *xtsDocument) IsTrace(t VTrace) bool {
	return (xd.tracing>>t)&1 == 1
}
