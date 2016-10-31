package parsehub_go

// ParseHub Run Wrapper
type Run struct {
	parsehub *ParseHub

	data     *RunResponse
}

// Get run data
func (r *Run) GetData() *RunResponse {
	return r.data
}
