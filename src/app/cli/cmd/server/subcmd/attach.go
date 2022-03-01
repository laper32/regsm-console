package subcmd

import (
	"fmt"

	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/spf13/cobra"
)

type RetGram struct {
	Role    string                 `json:"role"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Detail  map[string]interface{} `json:"detail"`
}

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
			panic("Work in progress")
			// cfg, err := cliconf.CoordinatorConfiguration()
			// if err != nil {
			// 	log.Error(err)
			// 	return
			// }

			// url := &url.URL{
			// 	Scheme: "ws",
			// 	Host:   fmt.Sprintf("%v:%v", cfg.GetString("coordinator.ip"), cfg.GetUint("coordinator.port")),
			// }
			// var retGram *RetGram
			// client := clientws.New(url.String())
			// client.OnConnected(func() {
			// 	// 设计失误
			// 	// 连接红蓝字应该是全局的 但我当时没考虑到
			// 	// 2.0会重新设计
			// 	retGram = &RetGram{
			// 		Role:    misc.Role,
			// 		Code:    status.ServerConnectedCoordinatorAndLoggingIn.ToInt(),
			// 		Message: status.ServerConnectedCoordinatorAndLoggingIn.Message(),
			// 		Detail: map[string]interface{}{
			// 			"server_id": serverID,
			// 			"command":   "attach",
			// 			"message":   strings.Join(args, " "),
			// 		},
			// 	}
			// 	msg, _ := json.Marshal(&retGram)
			// 	err := client.SendTextMessage(string(msg))
			// 	if err != nil {
			// 		log.Error("Failed to establish connection since the error occured. Message:", err)
			// 		client.Close()
			// 		return
			// 	}

			// 	go func() {
			// 		scanner := bufio.NewScanner(os.Stdin)
			// 		for scanner.Scan() {
			// 			text := scanner.Text()
			// 			text = strings.TrimSpace(text)
			// 			if len(text) == 0 {
			// 				continue
			// 			}
			// 			err = client.SendTextMessage(text)
			// 			if err != nil {
			// 				break
			// 			}
			// 		}
			// 	}()
			// })
			// client.OnTextMessageReceived(func(msg string) {
			// 	fmt.Println(msg)
			// })
			// client.Connect()
			// c := make(chan os.Signal, 1)
			// signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			// for {
			// 	s := <-c
			// 	switch s {
			// 	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			// 		// Client want to shutting down the connection.
			// 		retGram = &RetGram{
			// 			Role:    misc.Role,
			// 			Code:    status.ServerTerminateAttachConsole.ToInt(),
			// 			Message: status.ServerTerminateAttachConsole.Message(),
			// 			Detail: map[string]interface{}{
			// 				"server_id": serverID,
			// 			},
			// 		}
			// 		msg, _ := json.Marshal(&retGram)
			// 		err = client.SendTextMessage(string(msg))
			// 		if err != nil {
			// 			log.Warn("Coordinator offline. Message:", err)
			// 		}
			// 		return
			// 	default:
			// 		return
			// 	}
			// }
		},
	}
	attach.Flags().UintVar(&serverID, "server-id", 0, "Server ID")
	attach.MarkFlagRequired("server-id")
	return attach
}
