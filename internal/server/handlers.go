package server

import (
	"fmt"
	"net/http"
)

func docsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "docs")
}
