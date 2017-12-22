package controllers

import (
	"net/http"
	"strings"

	"github.com/euskadi31/go-server"
)

// UIController struct
type UIController struct {
}

// NewUIController constructor
func NewUIController() *UIController {
	return &UIController{}
}

// Mount implements server.Controller
func (c UIController) Mount(r *server.Router) {
	r.AddPrefixRouteFunc("/ui/", c.StaticFileHandler)
}

// StaticFileHandler endpoint
func (c UIController) StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.Replace(r.URL.Path, "/ui/", "/", 1)

	extensions := []string{".js", ".css", ".map", ".ico"}
	for _, ext := range extensions {
		if strings.HasSuffix(r.URL.Path, ext) {
			http.ServeFile(w, r, "/opt/cryptotrader/ui/"+filename)

			return
		}
	}

	http.ServeFile(w, r, "/opt/cryptotrader/ui/index.html")
}
