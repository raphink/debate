package getportrait

import (
	"net/http"
	"os"
	
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("GetPortrait", HandleGetPortrait)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	
	http.HandleFunc("/", HandleGetPortrait)
	http.ListenAndServe(":"+port, nil)
}
