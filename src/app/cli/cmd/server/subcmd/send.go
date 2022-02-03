package subcmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	cliconf "github.com/laper32/regsm-console/src/app/cli/conf"
	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/laper32/regsm-console/src/app/cli/misc"
	"github.com/spf13/cobra"
)

func InitSendCMD() *cobra.Command {
	var (
		serverID uint
	)
	send := &cobra.Command{
		Use:  "send",
		Args: cobra.MinimumNArgs(1),
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
			detail["message"] = strings.Join(args, " ")

			message := make(map[string]interface{})
			message["role"] = misc.Role
			message["message"] = detail

			// now, read the coordinator address
			coordinator_cfg, err := cliconf.CoordinatorConfiguration()
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}
			url := url.URL{
				Scheme: "ws",
				Host:   fmt.Sprintf("%v:%v", coordinator_cfg.GetString("coordinator.ip"), coordinator_cfg.GetUint("coordinator.port")),
			}
			fmt.Println("Establishing a connection to the coordinator:", url.String())
			c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}
			fmt.Println("DONE")

			done := make(chan struct{})
			go func() {
				defer close(done)
				for {
					_, message, err := c.ReadMessage()
					if err != nil {
						fmt.Println("ERROR:", err)
						return
					}
					fmt.Printf("%s\n", message)
				}
			}()

			var send bool = false
			for {

				select {
				case <-done:
					c.Close()
					return
				default:
					if !send {
						err = c.WriteJSON(message)
						if err != nil {
							fmt.Println("ERROR:", err)
							break
						}
						send = true
					}
				}
			}
		},
	}
	send.Flags().UintVar(&serverID, "server-id", 0, "Server ID")
	send.MarkFlagRequired("server-id")
	return send
}
