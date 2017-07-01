package utils

type TermError struct {
	S string
}
func (t TermError) Error() string {
	return t.S
}

type NotExistsError struct {
	S string
}
func (t NotExistsError) Error() string {
	return t.S
}