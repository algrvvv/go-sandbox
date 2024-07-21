package src

import (
	"net/http"
)

func GetNewServer() *http.ServeMux {
	serv := http.NewServeMux()

	serv.HandleFunc("GET /", indexHandler)
	serv.HandleFunc("GET /offline", offlineHandler)
	serv.HandleFunc("GET /online", indexHandler)
	serv.HandleFunc("POST /offline/run", offlineRunHandler)
	serv.HandleFunc("GET /new", newSessionHandler)
	serv.HandleFunc("GET /ws", wsHandler)

	return serv
}
