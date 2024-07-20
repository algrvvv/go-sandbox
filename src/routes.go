package src

import (
	"github.com/gorilla/websocket"
	"net/http"
	"text/template"
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
