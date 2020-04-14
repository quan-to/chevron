package pages

import (
	"github.com/gorilla/mux"
	"github.com/quan-to/chevron/internal/models"
	"net/http"
	"os"
	"path"
)

func displayFile(filename string, w http.ResponseWriter, r *http.Request) {
	fileData, err := Asset(filename[1:]) // Remove leading slash
	if err != nil {
		w.Header().Set("Content-Type", models.MimeText)
		w.WriteHeader(404)
		_, _ = w.Write([]byte(os.ErrNotExist.Error()))
		return
	}

	w.WriteHeader(200)
	_, _ = w.Write([]byte(fileData))
}

// AddHandlers attach handlers for all page assets into the specified router
func AddHandlers(r *mux.Router) {
	files := AssetNames()

	for _, v := range files {
		filePath := path.Join("/", v)
		r.HandleFunc(filePath, func(w http.ResponseWriter, r *http.Request) {
			displayFile(filePath, w, r)
		})
	}
}
