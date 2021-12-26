// windows: go build -o gsm-coordinator.exe
// linux: go build -o gsm-coordinator
package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/laper32/regsm-console/src/app/coordinator/conf"
	"github.com/laper32/regsm-console/src/lib/log"
)

func main() {
	conf := conf.Init()
	log.Init(conf.Log)

	log.Info("Initalizing coordinator...")
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(int(conf.Coordinator.Port)))
	checkError(err)
	log.Info("Done.")

	// TBD:
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	request := make([]byte, 128)                          // set maxium request length to 128B to prevent flood attack
	// close connection before exit
	defer func() {
		fmt.Println("Exiting...")
		conn.Close()
	}()

	for {
		read_len, err := conn.Read(request)

		if err != nil {
			fmt.Println(err)
			break
		}

		if read_len == 0 {
			break // connection already closed by client
		} else if strings.TrimSpace(string(request[:read_len])) == "timestamp" {
			daytime := strconv.FormatInt(time.Now().Unix(), 10)
			conn.Write([]byte(daytime))
		} else {
			daytime := time.Now().String()
			conn.Write([]byte(daytime))
		}

		request = make([]byte, 128) // clear last read content
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
