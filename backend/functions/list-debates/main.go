package listdebates

import "net/http"

// ListDebates is the Cloud Functions entry point
func ListDebates(w http.ResponseWriter, r *http.Request) {
	ListDebatesHandler(w, r)
}
