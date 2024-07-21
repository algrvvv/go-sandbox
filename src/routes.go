package src

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

var (
	sessions = make(map[string]*Session)
	upgrader = websocket.Upgrader{}
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmp.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func offlineHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/offline.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmp.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func offlineRunHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Code string `json:"code"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	}

	res := executeUserCode(data.Code)
	json, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json)
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	res := executeUserCode(code)

	tmp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = tmp.Execute(w, res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newSessionHandler(w http.ResponseWriter, r *http.Request) {}

func wsHandler(w http.ResponseWriter, r *http.Request) {}
