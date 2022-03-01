// windows: go build -o gsm-coordinator.exe
// linux: go build -o gsm-coordinator

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	cliconf "github.com/laper32/regsm-console/src/app/coordinator/conf"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/laper32/regsm-console/src/lib/os/shutil"
	"github.com/laper32/regsm-console/src/lib/status"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 5 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 1) / 10
	Role       = "coordinator"
)

type RetGram struct {
	Role    string                 `json:"role"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Detail  map[string]interface{} `json:"detail"`
}

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
	var retGram *RetGram
	err = conn.ReadJSON(&retGram)
	if err != nil {
		log.Info(err)
		return
	}
	if retGram.Code == status.CLICoordinatorSendStopSignal.ToInt() {
		os.Exit(0)
	}
	var actor *Actor
	// 设计失误
	// 连接红蓝字应该是全局的 但我当时没考虑到
	// 2.0会重新设计
	isLoggingIn := status.ServerConnectedCoordinatorAndLoggingIn.ToInt() == retGram.Code
	if isLoggingIn {
		if retGram.Role == "server" || retGram.Role == "coordinator" {
			actor = &Actor{
				role:     retGram.Role,
				conn:     conn,
				identity: make(map[string]interface{}),
				io: struct {
					input  chan []byte
					output chan []byte
				}{input: make(chan []byte), output: make(chan []byte)},
			}
		}
		switch retGram.Role {
		case "cli":
			whatToDo := retGram.Detail["command"].(string)
			serverID := uint(retGram.Detail["server_id"].(float64))
			thisActor := hub.actors[serverID]
			if thisActor == nil {
				detail := make(map[string]interface{})
				detail["server_id"] = serverID
				ret := status.CoordinatorServerOffline.WriteDetail(detail)
				log.Info(ret)
				conn.WriteMessage(websocket.TextMessage, []byte(ret))
				conn.Close()
				return
			}
			retGram.Code = status.CLISendingCommand.ToInt()
			retGram.Message = status.CLISendingCommand.Message()
			msg, _ := json.Marshal(&retGram)
			thisActor.io.input <- msg
			switch whatToDo {
			case "attach":
				log.Info(fmt.Sprintf("Attaching server \"%v\"", serverID))
				go func() {
					for {
						_, input, err := conn.ReadMessage()
						if err != nil {
							retGram.Role = Role
							retGram.Code = status.ServerTerminateAttachConsole.ToInt()
							retGram.Message = status.ServerTerminateAttachConsole.Message()
							thisActor.conn.WriteJSON(&retGram)
							fmt.Println(retGram)
							conn.Close()
							return
						}
						thisActor.io.input <- input
					}
				}()
			default:
				conn.Close()
				return
			}
		case "server":
			// Forcely close repeated connection...
			serverID := uint(retGram.Detail["server_id"].(float64))
			hasActor := hub.actors[serverID]
			if hasActor != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(status.CoordinatorServerAlreadyExist.WriteDetail(retGram.Detail)))
				return
			}
			actor.identity["server_id"] = serverID
			actor.identity["status"] = status.ServerConnectedCoordinatorAndLoggingIn.ToInt()
		case "coordinator":
			return
		default:
			log.Error(fmt.Sprintf("Terminating this connection due to the unknown role: %v", retGram.Role))
			conn.Close()
			return
		}
		hub.register <- actor
		log.Info(fmt.Sprintf("Actor: %v (%v) connected.", actor.identity["server_id"], actor.role))
		go actor.read()
		go actor.write()
	} else {
		conn.Close()
		return
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
			role := actor.role
			actor.conn.Close()
			delete(h.actors, serverID)
			log.Info(fmt.Sprintf("Actor: %v (%v) disconnected.", serverID, role))
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
			break
		}
		msg = bytes.TrimSpace(bytes.ReplaceAll(msg, []byte{'\n'}, []byte{' '}))

		// We unmarshal it, no error => it's JSON indeed.
		// Then we should do further actions.
		var retGram *RetGram
		if err = json.Unmarshal(msg, &retGram); err == nil {
			switch retGram.Role {
			case "server":
				serverID := uint(retGram.Detail["server_id"].(float64))
				thisActor := hub.actors[serverID]
				switch retGram.Code {
				case status.ServerSendingProcessInfo.ToInt():
					serverPID := int(retGram.Detail["server_pid"].(float64))
					daemonPID := int(retGram.Detail["daemon_pid"].(float64))
					thisActor.identity["server_pid"] = serverPID
					thisActor.identity["daemon_pid"] = daemonPID
					continue
				case status.ServerStopping.ToInt():
					hub.unregister <- thisActor
					return
				default:
					continue
				}
			case "coordinator":
			default:
				continue
			}
		} else {
			actor.io.output <- msg
		}
	}
}

func (actor *Actor) write() {
	tick := time.NewTicker(pingPeriod)
	defer func() {
		tick.Stop()
		actor.conn.Close()
	}()
	for {
		select {
		case message, ok := <-actor.io.input:
			if !ok {
				actor.conn.WriteMessage(websocket.CloseMessage, []byte{})
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
	shutil.ClearTerminalScreen()
	cfg := cliconf.Init()
	log.Init(cfg.Log)
	log.Info("Coordinator initialized.")

	go hub.run()
	http.HandleFunc("/", wsHandle)
	http.ListenAndServe(fmt.Sprintf("%v:%v", cfg.Param.IP, cfg.Param.Port), nil)
}
