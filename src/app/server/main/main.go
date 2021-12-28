package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

func pumpStdin(ws *websocket.Conn, w io.Writer) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message, '\n')
		if _, err := w.Write(message); err != nil {
			break
		}
	}
}

func pumpStdout(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	defer func() {
	}()
	s := bufio.NewScanner(r)
	for s.Scan() {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer ws.Close()

	outr, outw, err := os.Pipe()
	if err != nil {
		internalError(ws, "stdout: ", err)
		return
	}
	defer outr.Close()
	defer outw.Close()

	inr, inw, err := os.Pipe()
	if err != nil {
		internalError(ws, "stdout: ", err)
		return
	}
	defer inr.Close()
	defer inw.Close()

	proc, err := os.StartProcess("D:/regsm/server/1/hlds.exe", []string{" -console -game cstrike -ip 0.0.0.0 -port 23333 +map de_dust2 0"}, &os.ProcAttr{
		Files: []*os.File{inr, outw, outw},
		Dir:   "D:/regsm/server/1",
	})
	if err != nil {
		internalError(ws, "start: ", err)
		return
	}

	inr.Close()
	outw.Close()

	stdoutDone := make(chan struct{})
	go pumpStdout(ws, outr, stdoutDone)
	go ping(ws, stdoutDone)

	pumpStdin(ws, inw)

	// Some commands will exit when stdin is closed.
	inw.Close()

	// Other commands need a bonk on the head.
	if err := proc.Signal(os.Interrupt); err != nil {
		log.Println("inter:", err)
		os.Exit(0)
	}

	select {
	case <-stdoutDone:
	case <-time.After(time.Second):
		// A bigger bonk on the head.
		if err := proc.Signal(os.Kill); err != nil {
			log.Println("term:", err)
		}
		<-stdoutDone
	}

	if _, err := proc.Wait(); err != nil {
		log.Println("wait:", err)
	}
}
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	// D:/regsm/server/1/hlds.exe -console -game cstrike -ip 0.0.0.0 -port 23333 +map de_dust2

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWS)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

// func main() {

// 	// Process
// 	// 	1. Connect to coordinator
// 	// 	2. Redirect stdin to socket connection receive
// 	// 	3. Start this server
// 	// 	4. Write start info (PID, etc)
// 	// 	5. Send to coordinator to make a presistent storage.
// 	// 	6. DONE.

// 	type Application struct {
// 		ID         uint     `json:"server_id"`
// 		Dir        string   `json:"dir"`
// 		Executable string   `json:"executable"`
// 		Args       []string `json:"args"`
// 	}
// 	var app Application
// 	json.Unmarshal([]byte(os.Args[1]), &app)

// 	// We cannot pass in []strings{} directly.
// 	// Because for example: if we pass the slice directly, we will found that 'de_dust2"' not found.
// 	// Yeah, wtf, why occur '"'?
// 	// Not very sure how golang pass in the executable params, and don't know the implementation detail...
// 	// Overall, stupid, very stupid
// 	// Based on this, we need to construct a parameter string
// 	// And we need also add up a space to make sure it can be executed correctly.
// 	// Test passed on CS1.6
// 	serverParamStr := strings.Join(app.Args, " ")
// 	serverEXE := exec.Command(app.Executable, serverParamStr+" ")
// 	serverEXE.Stdin = os.Stdin
// 	serverEXE.Stdout = os.Stdout
// 	serverEXE.Stderr = os.Stderr
// 	serverEXE.Dir = app.Dir

// 	err := serverEXE.Start()
// 	if err != nil {
// 		fmt.Println("ERROR:", err)
// 		return
// 	}
// 	fmt.Printf("Server started. Server ID: %v(Process ID: %v)", app.ID, serverEXE.Process.Pid)
// 	type Status struct {
// 		ServerID  uint
// 		ProcessID int
// 		Done      bool
// 	}
// 	status := &Status{
// 		ServerID:  app.ID,
// 		ProcessID: serverEXE.Process.Pid,
// 	}
// 	json.Marshal(status)
// }
