package rsver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type App struct {
	mu       sync.Mutex
	allUsers map[string]int64 // uuid -> timestamp unix
}

func NewApp() *App {
	return &App{
		allUsers: make(map[string]int64),
		mu:       sync.Mutex{},
	}
}
func (a *App) RunTimer() {
	tk := time.NewTicker(time.Minute)
	defer tk.Stop()
	for {
		select {
		case <-tk.C:
			a.registerNodePullAuthedUsers("") //todo
		}
	}
}

func (a *App) registerNodePullAuthedUsers(url string) {
	if url == "" {
		return
	}
	body := `{"name":"aw"}`
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(body))
	if err != nil {
		log.Println("Error registering:", err)
		return
	}
	defer resp.Body.Close()
	users := make(map[string]int64)
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		log.Println("Error decoding response:", err)
		return
	}
	a.mu.Lock()
	a.allUsers = users
	a.mu.Unlock()
}

func (a *App) IsUserNotAllowed(uuid string) (isNotAllowed bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	ts, ok := a.allUsers[uuid]
	if !ok {
		log.Println("Unauthorized user:", uuid)
		return true
	}
	if ts < time.Now().Unix() {
		log.Println("User expired:", uuid)
		return true
	}
	return false
}

func (a *App) ping(w http.ResponseWriter, r *http.Request) {
	//json response hello world
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ping": "pong"}`))
}
func Run() {
	app := NewApp()
	go app.RunTimer()

	fmt.Println("WebSocket server starting on 127.0.0.1:80")
	http.HandleFunc("/ws-s5", app.wsS5)
	http.HandleFunc("/ws-vless", app.wsVless)
	http.HandleFunc("/", app.ping)
	http.HandleFunc("/ping", app.ping)
	log.Fatal(http.ListenAndServe(":80", nil))
}
