package utils

type TermError struct {
	s string
}
func (t TermError) Error() string {
	return t.s
}