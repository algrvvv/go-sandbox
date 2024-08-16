package src

import (
	"net/http"
)

func GetNewServer() *http.ServeMux {
	serv := http.NewServeMux()

	serv.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	serv.HandleFunc("GET /", indexHandler)
	serv.HandleFunc("GET /offline", offlineHandler)
	serv.HandleFunc("POST /offline/run", offlineRunHandler)
	serv.HandleFunc("GET /new", newOnlineHandler)
	serv.HandleFunc("GET /online", onlineHandler)
	serv.HandleFunc("GET /connect", connectOnlineHandler)
	serv.HandleFunc("GET /ws", wsHandler)

	return serv
}
