package srvr

/*
 * "pattern" string documented here: https://golang.org/pkg/net/http/#ServeMux
 */
func (s *Srvr) Routes() {
	s.Router.HandleFunc("/index.html", s.handleIndex())
	s.Router.HandleFunc("/form", s.handleForm())
	s.Router.HandleFunc("/jumble", s.handleJumble())
	s.Router.HandleFunc("/solve", s.handleSolve())
	s.Router.HandleFunc("/", s.handleIndex())
}
