package http

import (
	"fmt"
	"net/http"
)

// homeHandler shows the server name for rooted requests. Otherwise, a 404 is
// returned.
func homeHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "matchbox\n")
	}
	return http.HandlerFunc(fn)
}
