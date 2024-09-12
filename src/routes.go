package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/algrvvv/go-sandbox/src/logger"
)

var (
	sessions = make(map[SessionID][]*websocket.Conn)
	mu       sync.Mutex
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		logger.Error("failed to parse template: "+err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmp.Execute(w, nil); err != nil {
		logger.Error("failed to execute template: "+err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func offlineHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/offline.html")
	if err != nil {
		logger.Error("failed to parse template: "+err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmp.Execute(w, nil); err != nil {
		logger.Error("failed to execute template: "+err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func offlineRunHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Code string `json:"code"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logger.Error("failed to parse request: "+err.Error(), err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	}

	res := executeUserCode(data.Code)
	jsonData, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		logger.Error("failed to marshal response: "+err.Error(), err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}

func onlineHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/online.html")
	if err != nil {
		logger.Error("failed to parse template: "+err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmp.Execute(w, nil); err != nil {
		logger.Error("failed to execute template: "+err.Error(), err)
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

func newOnlineHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := uuid.New().String()
	uid := uuid.New().String()
	url := fmt.Sprintf("/online?s=%s&u=%s", sessionID, uid)
	http.Redirect(w, r, url, http.StatusFound)
}

func connectOnlineHandler(w http.ResponseWriter, r *http.Request) {
	s := SessionID(r.URL.Query().Get("s"))
	if _, exists := sessions[s]; exists {
		uid := uuid.New().String()
		//ss := append(*currentSessions, Session{Uid: uid})
		//sessions[s] = &ss
		url := fmt.Sprintf("/online?s=%s&u=%s", s, uid)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		http.Error(w, "Session not found", http.StatusNotFound)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("online err: %v", err)
		return
	}
	defer conn.Close()

	var (
		sessionID = SessionID(r.URL.Query().Get("s"))
		uid       = r.URL.Query().Get("u")
	)

	if sessionID == "" || uid == "" {
		http.Error(w, "Missing session id or uid", http.StatusBadRequest)
		return
	}

	//if _, ok := sessions[sessionID]; !ok {
	//	sessions[sessionID] = make([]*websocket.Conn, 0)
	//}

	mu.Lock()
	sessions[sessionID] = append(sessions[sessionID], conn)
	mu.Unlock()

	defer func() {
		mu.Lock()
		conns := sessions[sessionID]
		for i, c := range conns {
			if c == conn {
				sessions[sessionID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}

		if len(sessions[sessionID]) == 0 {
			delete(sessions, sessionID)
		}
		mu.Unlock()
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ws read err: %v", err)
			break
		}

		mu.Lock()
		for _, c := range sessions[sessionID] {
			if c != conn {
				err = c.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("ws write err: %v", err)
				}
			}
		}
		mu.Unlock()
	}
}

func handleConnection(conn *websocket.Conn) {

}
