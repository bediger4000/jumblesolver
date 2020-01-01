package srvr

import (
	"jumble/dictionary"
	"net/http"
)

type Srvr struct {
	// Configuration values or whatever every
	// handler function needs to know - common info
	FindWords dictionary.Dictionary
	Router    *http.ServeMux
}
