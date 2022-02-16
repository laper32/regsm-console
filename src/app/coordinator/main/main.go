// windows: go build -o gsm-coordinator.exe
// linux: go build -o gsm-coordinator

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/laper32/regsm-console/src/lib/status"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 5 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Actor struct {
	role     string
	conn     *websocket.Conn
	identity map[string]interface{}
	io       struct {
		input  chan []byte
		output chan []byte
	}
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
		// 要写到文件里面的哈
		log.Error(err)
		return
	}
	// Read the first connect message, and resolve it.
	// We need to know what role the connection is.
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Error(err)
		return
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(msg, &data)
	if err != nil {
		log.Error(err)
		return
	}
	role := data["role"].(string)
	detail := data["detail"].(map[string]interface{})
	var actor *Actor
	statusCode := int(data["code"].(float64))
	// 设计失误
	// 连接红蓝字应该是全局的 但我当时没考虑到
	// 2.0会重新设计
	isLoggingIn := status.ServerConnectedCoordinatorAndLoggingIn.ToInt() == statusCode
	if isLoggingIn {
		if role == "server" || role == "coordinator" {
			actor = &Actor{
				role:     role,
				conn:     conn,
				identity: make(map[string]interface{}),
				io: struct {
					input  chan []byte
					output chan []byte
				}{input: make(chan []byte), output: make(chan []byte)},
			}
		}
		switch role {
		case "cli":
			whatToDo := detail["command"].(string)
			serverID := uint(detail["server_id"].(float64))
			thisActor := hub.actors[serverID]
			if thisActor == nil {
				log.Info(status.CoordinatorServerOffline.WriteDetail(""))
				conn.Close()
				return
			}
			thisActor.io.input <- msg
			switch whatToDo {
			case "attach":

				// If we are now attaching the server session
				// What we need to do is: retrieve the output message
				// and send to the CLI's console
				for {
					_, input, err := conn.ReadMessage()
					if err != nil {
						return
					}
					thisActor.io.input <- input

					output := <-thisActor.io.output
					conn.WriteMessage(websocket.TextMessage, output)
				}
			case "stop":
				hub.unregister <- thisActor
				return
			default:
				conn.Close()
				return
			}
		case "coordinator":
			return
		case "server":
			// Forcely close repeated connection...
			serverID := uint(detail["server_id"].(float64))
			hasActor := hub.actors[serverID]
			if hasActor != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(status.CoordinatorServerAlreadyExist.WriteDetail("")))
				return
			}
			actor.identity["server_id"] = serverID
		default:
			log.Error(fmt.Sprintf("Terminating this connection due to the unknown role: %v", role))
			conn.Close()
			return
		}
		for k := range data {
			delete(data, k)
		}

		data["role"] = "coordinator"
		data["code"] = status.OK
		data["message"] = status.OK.Message()
		fmt.Println(data)
		actor.conn.WriteJSON(&data)
		// ...
		// 我草
		// 为啥啊
		// 为啥在struct里面的conn就没办法读消息 裸conn就可以正常拿到消息?
		// =========================这段是可以取到消息的
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		fmt.Println(string(msg))
		// 下面的都不行

		// for k := range data {
		// 	delete(data, k)
		// }
		// _, msg, _ := actor.conn.ReadMessage()
		// fmt.Println(string(msg))
		// if data["code"] == nil {
		// 	fmt.Println("Why...the code is null?")
		// 	return
		// }
		// responseStatus := int(data["code"].(float64))
		// // responseStatus := status.ToCode(int(data["code"].(float64)))
		// if responseStatus == status.OK.ToInt() {
		// 	detail = data["detail"].(map[string]interface{})
		// 	if serverID := uint(detail["server_id"].(float64)); serverID != actor.identity["server_id"].(uint) {
		// 		// This should not happen.
		// 		return
		// 	}
		// 	hub.register <- actor
		// 	log.Info(fmt.Sprintf("Actor: %v (%v) connected.", actor.identity["server_id"], actor.role))
		// } else {
		// 	fmt.Println("Aborting this connection since the message received is incorrect. IP:", actor.conn.RemoteAddr().String())
		// 	fmt.Println("Message:", data)
		// 	hub.unregister <- actor
		// 	return
		// }
		// go actor.read()
		// go actor.write()
	} else {
		// Unknown.
		// Terminate this connection
		return
	}
}

func (h *Hub) run() {
	for {
		select {
		case actor := <-h.register:
			serverID := actor.identity["server_id"].(uint)
			h.actors[serverID] = actor
			fmt.Printf("Actor: %v (Role: %v) connected.\n", serverID, actor.role)
		case actor := <-h.unregister:
			serverID := actor.identity["server_id"].(uint)
			fmt.Printf("Actor: %v (Role: %v) disconnected.\n", serverID, actor.role)
			delete(h.actors, serverID)
		}
	}
}

func (actor *Actor) read() {
	defer func() {
		hub.unregister <- actor
		actor.conn.Close()
	}()
	actor.conn.SetReadDeadline(time.Now().Add(pongWait))
	actor.conn.SetPongHandler(func(string) error { actor.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := actor.conn.ReadMessage()
		if err != nil {
			hub.unregister <- actor
			log.Info(fmt.Sprintf("Actor: %v (%v) disconnected.", actor.identity["server_id"], actor.role))
			break
		}

		actor.io.output <- msg
	}
}

func (actor *Actor) write() {
	tick := time.NewTicker(pingPeriod)
	defer func() {
		tick.Stop()
		hub.unregister <- actor
		actor.conn.Close()
	}()
	for {
		select {
		case message, ok := <-actor.io.input:
			if !ok {
				return
			}
			actor.conn.WriteMessage(websocket.TextMessage, message)
		case <-tick.C:
			actor.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := actor.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
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
