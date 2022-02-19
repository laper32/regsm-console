package subcmd

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/app/cli/misc"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/laper32/regsm-console/src/lib/status"
	"github.com/spf13/cobra"
)

func InitAttachCMD() *cobra.Command {
	var serverID uint
	attach := &cobra.Command{
		Use: "attach",
		Run: func(cmd *cobra.Command, args []string) {
			sv := dpkg.FindIdentifiedServer(serverID)
			if sv == nil {
				fmt.Println("Server NOT found")
				return
			}
			if sv.Deleted {
				fmt.Println("Server has been deleted")
				return
			}
			detail := make(map[string]interface{})
			detail["server_id"] = serverID
			detail["command"] = "attach"
			detail["message"] = strings.Join(args, " ")

			message := make(map[string]interface{})
			message["role"] = misc.Role
			// 设计失误
			// 连接红蓝字应该是全局的 但我当时没考虑到
			// 2.0会重新设计
			message["code"] = status.ServerConnectedCoordinatorAndLoggingIn
			message["message"] = status.ServerConnectedCoordinatorAndLoggingIn.Message()
			message["detail"] = detail

			cfg, err := cliconf.CoordinatorConfiguration()
			if err != nil {
				log.Error(err)
				return
			}
			url := &url.URL{
				Scheme: "ws",
				Host:   fmt.Sprintf("%v:%v", cfg.GetString("coordinator.ip"), cfg.GetUint("coordinator.port")),
			}
			fmt.Printf("[%v] Connecting to the coordinator...", url.String())
			conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
			if err != nil {
				fmt.Println()
				log.Error(err)
				return
			}
			fmt.Println("OK")
			defer conn.Close()
			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
			scanner := bufio.NewScanner(os.Stdin)
			for {
				signal := <-interrupt
				switch signal {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					detail := make(map[string]interface{})
					detail["server_id"] = serverID

					data := make(map[string]interface{})
					data["role"] = "cli"
					data["code"] = status.ServerTerminateAttachConsole
					data["message"] = status.ServerTerminateAttachConsole.Message()
					data["detail"] = detail
					err := conn.WriteJSON(&data)
					if err != nil {
						fmt.Println("ERROR:", err)
						return
					}
					conn.Close()
				default:
					return
				}
				scanner.Scan()
				text := scanner.Text()
				err := conn.WriteMessage(websocket.TextMessage, []byte(text))
				if err != nil {
					fmt.Println("ERROR:", err)
					continue
				}
			}
		},
	}
	attach.Flags().UintVar(&serverID, "server-id", 0, "Server ID")
	attach.MarkFlagRequired("server-id")
	return attach
}
