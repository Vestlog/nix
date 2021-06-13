package api

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/vestlog/nix/pkg/storage"
)

type API struct {
	db *storage.GormDatabase
}

func encode(w http.ResponseWriter, r *http.Request, data interface{}) {
	if r.Header.Get("Accept") == "application/xml" {
		h := w.Header()
		h.Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(data)
		return
	}
	h := w.Header()
	h.Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (api *API) handlePosts(w http.ResponseWriter, r *http.Request) {
	seq := strings.Split(r.URL.Path, "/")
	var data interface{}
	var err error
	if len(seq) > 2 && seq[2] != "" {
		data, err = api.db.GetPost(seq[2])
	} else {
		data, err = api.db.GetPosts()
	}
	if err != nil {
		data = map[string]string{"error": err.Error()}
	}
	encode(w, r, data)
}

func (api *API) handleComments(w http.ResponseWriter, r *http.Request) {
	seq := strings.Split(r.URL.Path, "/")
	var data interface{}
	var err error
	if len(seq) > 2 && seq[2] != "" {
		data, err = api.db.GetComment(seq[2])
	} else {
		data, err = api.db.GetComments()
	}
	if err != nil {
		data = map[string]string{"error": err.Error()}
	}
	encode(w, r, data)
}

func CreateAPIHandler(dsn string) (http.Handler, error) {
	db, err := storage.CreateGormDatabase(dsn)
	if err != nil {
		return nil, err
	}
	api := &API{db}
	mux := http.NewServeMux()
	mux.HandleFunc("/posts/", api.handlePosts)
	mux.HandleFunc("/comments/", api.handleComments)
	return mux, nil
}
