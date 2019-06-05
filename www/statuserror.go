package www

// StatusError stores the values of a failed http.Request
type StatusError struct {
	Status     string
	StatusCode int
	Body       []byte
}

// Error returns the Status value of the failed  request
func (se StatusError) Error() string {
	return se.Status
}
