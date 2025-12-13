package listdebates

import (
	"net/http"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

// ListDebates is the Cloud Functions entry point
func ListDebates(w http.ResponseWriter, r *http.Request) {
	ListDebatesHandler(w, r)
}
