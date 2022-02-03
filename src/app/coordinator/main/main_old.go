//go:build ignore
// +build ignore

// windows: go build -o gsm-coordinator.exe
// linux: go build -o gsm-coordinator

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Actor struct {
	role     string
	conn     *websocket.Conn
	identity map[string]interface{}
}

type Hub struct {
	actors     map[uint]*Actor
	register   chan *Actor
	unregister chan *Actor
}

var (
	upgrader = websocket.Upgrader{}
	hub      = &Hub{
		actors:     make(map[uint]*Actor),
		register:   make(chan *Actor),
		unregister: make(chan *Actor),
	}
)

func wsHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	// Read the first connect message, and resolve it.
	// We need to know what role the connection is.
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Connection lost:", conn.RemoteAddr().String())
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(msg, &data)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	role := data["role"].(string)
	actor := &Actor{
		role:     role,
		conn:     conn,
		identity: make(map[string]interface{}),
	}

	if role == "server" {
		actor.identity["server_id"] = uint(data["server_id"].(float64))
		actor.identity["daemon_pid"] = int(data["daemon_pid"].(float64))
	} else if role == "coordinator" {
		fmt.Println("This connection comes from an another coordinator.")
		return
	} else if role == "cli" {
		message := data["message"].(map[string]interface{})
		input := message["message"].(string)
		serverID := uint(message["server_id"].(float64))
		fmt.Println(input)
		fmt.Println(serverID)
		return
	} else {
		fmt.Println("Unknown role. Terminate this connection.")
		return
	}
	hub.register <- actor
	defer func() {
		hub.unregister <- actor
	}()
	fmt.Println("Actor:", conn.RemoteAddr().String(), "connected. Role:", actor.role)
	go read(actor)
}

func read(actor *Actor) {
	for {
		_, msg, err := actor.conn.ReadMessage()
		if err != nil {
			fmt.Println("Actor:", actor.conn.RemoteAddr().String(), "disconnected. Role:", actor.role)
			hub.unregister <- actor
			break
		}
		data := make(map[string]interface{})
		err = json.Unmarshal(msg, &data)
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		if data["level"].(string) == "error" {
			fmt.Println(data["message"])
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case actor := <-h.register:
			serverID := actor.identity["server_id"].(uint)
			fmt.Println("Adding server=>Server ID:", serverID)
			h.actors[serverID] = actor
			fmt.Println(h.actors[serverID])
		case actor := <-h.unregister:
			serverID := actor.identity["server_id"].(uint)
			delete(h.actors, serverID)
		}
	}
}

func main() {
	go hub.run()
	http.HandleFunc("/", wsHandle)
	// If we use gsm coordinator start => GSM_PATH is already set
	// cfg := conf.Load(&conf.Config{
	// 	Name: "coordinator",
	// 	Type: "toml",
	// 	Path: []string{os.Getenv("GSM_PATH")},
	// })
	// http.ListenAndServe(fmt.Sprintf("%v:%v", cfg.GetString("coordinator.ip"), cfg.GetUint("coordinator.port")), nil)
	http.ListenAndServe("localhost:3484", nil)
}
