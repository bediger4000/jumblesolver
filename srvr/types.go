package srvr

import (
	"jumble/dictionary"
	"net/http"
)

type Srvr struct {
	FindWords dictionary.Dictionary
	Router    *http.ServeMux
	Debug     bool
}
