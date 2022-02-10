package subcmd

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/app/cli/misc"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/laper32/regsm-console/src/lib/status"
	"github.com/spf13/cobra"
)

func InitStopCMD() *cobra.Command {
	var serverID uint
	stop := &cobra.Command{
		Use: "stop",
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
			detail["command"] = "stop"

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
			fmt.Printf("Sending message...")
			err = conn.WriteJSON(&message)
			if err != nil {
				fmt.Println()
				log.Error(err)
				return
			}
			fmt.Println("OK")

		},
	}
	stop.Flags().UintVar(&serverID, "server-id", 0, "")
	stop.MarkFlagRequired("server-id")
	return stop
}
