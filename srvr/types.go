package srvr

import (
	"jumblesolver/dictionary"
	"net/http"
)

// Srvr holds all info that's used across instances
// of this web application.
type Srvr struct {
	FindWords dictionary.Dictionary
	Router    *http.ServeMux
	Debug     bool
}
