package srvr

// Routes sets the URIs and functions that work them.
// "pattern" string documented here: https://golang.org/pkg/net/http/#ServeMux
func (s *Srvr) Routes() {
	s.Router.HandleFunc("/jumble", s.handleJumble())
	s.Router.HandleFunc("/solve", s.handleSolve())
	s.Router.HandleFunc("/", s.handleJumble())
}
