package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func onlineRunHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Session string `json:"session"`
		Uid     string `json:"uid"`
		Code    string `json:"code"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logger.Error("failed to parse request: "+err.Error(), err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	s := SessionID(data.Session)
	if _, exists := sessions[s]; !exists {
		logger.Warn("session not found")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("session not found"))
		return
	}

	if data.Uid == "" {
		logger.Warn("u is empty")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("u is empty"))
		return
	}

	res := executeUserCode(data.Code)
	jsonData, err := json.Marshal(map[string]interface{}{
		"type": "console",
		// TODO планировалось, что первоначальному пользователю будет сразу же ответом приходить,
		// но пока что будет отправляться тоже по вебсокетам.
		"exceptedUser": data.Uid,
		"data":         res,
	})
	if err != nil {
		logger.Error("failed to marshal response: "+err.Error(), err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// TODO понаблюдать за поведением, наличием дедлоков
	mu.Lock()
	for _, c := range sessions[s] {
		err = c.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			logger.Error("failed to write message: "+err.Error(), err)
		}
	}
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Success execute"))
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
		// ss := append(*currentSessions, Session{Uid: uid})
		// sessions[s] = &ss
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
		logger.Error("failed to upgrade connection: "+err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	var (
		sessionID = SessionID(r.URL.Query().Get("s"))
		uid       = r.URL.Query().Get("u")
	)

	if sessionID == "" || uid == "" {
		logger.Error("Missing session id or uid", nil)
		http.Error(w, "Missing session id or uid", http.StatusBadRequest)
		return
	}

	// if _, ok := sessions[sessionID]; !ok {
	//	sessions[sessionID] = make([]*websocket.Conn, 0)
	// }

	mu.Lock()
	sessions[sessionID] = append(sessions[sessionID], conn)
	SendCountActiveUsers(sessionID)
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
		SendCountActiveUsers(sessionID)
		mu.Unlock()
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error("failed to read message: "+err.Error(), err)
			log.Printf("ws read err: %v", err)
			break
		}

		mu.Lock()
		for _, c := range sessions[sessionID] {
			if c != conn {
				err = c.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					logger.Error("failed to write message: "+err.Error(), err)
				}
			}
		}
		mu.Unlock()
	}
}

func goFmtHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.Error("failed to parse request: "+err.Error(), err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	tmp, err := os.CreateTemp("", "gofmt_*.go")
	if err != nil {
		logger.Error("failed to create temp file: "+err.Error(), err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer func() { _ = os.Remove(tmp.Name()) }()

	_, _ = tmp.Write([]byte(data.Code))

	cmd := exec.Command("gofmt", "-w", tmp.Name())
	if err = cmd.Run(); err != nil {
		logger.Error("failed to run gofmt: "+err.Error(), err)
	}

	cmd = exec.Command("goimports", "-w", tmp.Name())
	if err = cmd.Run(); err != nil {
		logger.Error("failed to run goimports: "+err.Error(), err)
	}

	final, err := os.ReadFile(tmp.Name())
	if err != nil {
		logger.Error("failed to read final result file: "+err.Error(), err)
		return
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"status": true,
		"code":   string(final),
	})
	if err != nil {
		logger.Error("failed to marshal response: "+err.Error(), err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}
