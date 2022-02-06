// windows: go build -o gsm-coordinator.exe
// linux: go build -o gsm-coordinator

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/laper32/regsm-console/src/app/coordinator/conf"
	"github.com/laper32/regsm-console/src/lib/log"
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
	message    chan []byte
}

var (
	upgrader = websocket.Upgrader{}
	hub      = &Hub{
		actors:     make(map[uint]*Actor),
		register:   make(chan *Actor),
		unregister: make(chan *Actor),
		message:    make(chan []byte),
	}
)

func wsHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	// Read the first connect message, and resolve it.
	// We need to know what role the connection is.
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Info("Connection lost:", conn.RemoteAddr().String())
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(msg, &data)
	if err != nil {
		log.Error(err)
		return
	}
	role := data["role"].(string)

	actor := &Actor{
		role:     role,
		conn:     conn,
		identity: make(map[string]interface{}),
	}

	if role == "server" {
		serverID := uint(data["server_id"].(float64))
		testActor := hub.actors[serverID]
		if testActor != nil {
			for k := range data {
				delete(data, k)
			}
			detail := make(map[string]interface{})
			detail["server_started"] = true
			data["level"] = "error"
			data["role"] = "coordinator"
			data["message"] = detail
			conn.WriteJSON(&data)
			return
		} else {
			actor.identity["server_id"] = serverID
			actor.identity["daemon_pid"] = int(data["daemon_pid"].(float64))
		}
	} else if role == "coordinator" {
		fmt.Println("This connection comes from an another coordinator.")
		return
	} else if role == "cli" {
		whatToDo := data["command"].(string)
		switch whatToDo {
		case "send":
			fmt.Println("Sending command to the specific server.")
			message := data["message"].(map[string]interface{})
			serverID := uint(message["server_id"].(float64))
			toExecute := message["message"].(string)
			thisActor := hub.actors[serverID]
			if thisActor == nil {
				log.Info("This server currently offline. ID:", serverID)
				return
			}
			sendJSON := make(map[string]interface{})
			sendJSON["command"] = "send"
			sendJSON["message"] = toExecute
			err = thisActor.conn.WriteJSON(&sendJSON)
			if err != nil {
				log.Error(err)
				return
			}

			msg := <-hub.message

			err = conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Info("Connection lost:", conn.RemoteAddr().String())
				return
			}
		case "attach":
			break
		case "stop":
			break
		case "restart":
			break
		case "update":
			break
		default:
			log.Error("Unknown command:", whatToDo)
		}

		return
	} else {
		fmt.Println("Unknown role. Terminate this connection.")
		return
	}

	hub.register <- actor

	for k := range data {
		delete(data, k)
	}
	detail := make(map[string]interface{})
	detail["connected"] = true
	data["level"] = "info"
	data["role"] = "coordinator"
	data["message"] = detail
	actor.conn.WriteJSON(&data)

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
		if err == nil {

			// message := data["message"].(map[string]interface{})
			// exited := message["exited"].(bool)
			// if exited {
			// 	log.Info("Server exited. ID:", actor.identity["server_id"])
			// 	hub.unregister <- actor
			// 	break
			// }
		}

		hub.message <- msg
	}
}

func (h *Hub) run() {
	for {
		select {
		case actor := <-h.register:
			serverID := actor.identity["server_id"].(uint)
			h.actors[serverID] = actor
		case actor := <-h.unregister:
			serverID := actor.identity["server_id"].(uint)
			delete(h.actors, serverID)
		}
	}
}

func main() {
	cfg := conf.Init()
	log.Init(cfg.Log)

	go hub.run()
	http.HandleFunc("/", wsHandle)
	http.ListenAndServe(fmt.Sprintf("%v:%v", os.Args[1], os.Args[2]), nil)
}
