package src

import (
	"net/http"
)

func GetNewServer() *http.ServeMux {
	serv := http.NewServeMux()

	serv.HandleFunc("GET /", indexHandler)
	serv.HandleFunc("GET /offline", indexHandler)
	serv.HandleFunc("GET /online", indexHandler)
	serv.HandleFunc("POST /offline/run", runHandler)
	serv.HandleFunc("GET /new", newSessionHandler)
	serv.HandleFunc("GET /ws", wsHandler)

	return serv
}
