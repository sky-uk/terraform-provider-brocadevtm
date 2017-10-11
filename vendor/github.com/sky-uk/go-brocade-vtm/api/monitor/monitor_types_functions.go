package monitor

// NewArgumentIssue - Returns the Membernodes object
func NewArgumentIssue(name string, description string, value string) ArgumentIssue {
	argumentIssue := ArgumentIssue{
		Name:        name,
		Description: description,
		Value:       value,
	}
	return argumentIssue

}
