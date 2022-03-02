package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/laper32/regsm-console/src/app/server/conf"
	"github.com/laper32/regsm-console/src/app/server/entity"
	"github.com/laper32/regsm-console/src/app/server/misc"
	"github.com/laper32/regsm-console/src/app/server/util"
	"github.com/laper32/regsm-console/src/lib/clientws"
	"github.com/laper32/regsm-console/src/lib/container/queue"
	"github.com/laper32/regsm-console/src/lib/log"
	"github.com/laper32/regsm-console/src/lib/status"
	"github.com/laper32/regsm-console/src/lib/sys/windows"
)

type RetGram struct {
	Role    string                 `json:"role"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Detail  map[string]interface{} `json:"detail"`
}

var (
	retGram     *RetGram
	thisCommand string
	outputQueue *queue.Queue
)

func handleClientWSConnection(cfg *conf.Config) {
	url := url.URL{Scheme: "ws", Host: fmt.Sprintf("%v:%v", cfg.Coordinator.IP, cfg.Coordinator.Port)}
	entity.Conn = clientws.New(url.String())
	entity.Conn.OnConnected(func() {
		detail := make(map[string]interface{})
		detail["server_id"] = cfg.Server.ID
		retGram = &RetGram{
			Role:    misc.Role,
			Code:    status.ServerConnectedCoordinatorAndLoggingIn.ToInt(),
			Message: status.ServerConnectedCoordinatorAndLoggingIn.Message(),
			Detail:  detail,
		}
		msg, _ := json.Marshal(&retGram)
		raw := string(msg)
		fmt.Println(raw)
		err := entity.Conn.SendTextMessage(string(msg))
		if err != nil {
			log.Error("Failed to establish connection since the error occured. Message:", err)
			entity.Conn.Close()
			return
		}
	})
	entity.Conn.OnTextMessageReceived(func(msg string) {
		fmt.Println(msg)
		err := json.Unmarshal([]byte(msg), &retGram)
		if err == nil {
			switch retGram.Code {
			case status.CoordinatorServerAlreadyExist.ToInt():
				log.Info("Repeated server exists.")
				entity.Conn.Close()
				os.Exit(0)
			case status.CLISendingCommand.ToInt():
				if entity.Proc.EXE.ProcessState != nil {
					return
				}
				command := retGram.Detail["command"].(string)
				thisCommand = command
				detail := make(map[string]interface{})
				switch command {
				case "attach":
					// go handleAttach(msg)
					return
				case "restart":
					util.ForceStopServer(cfg)
					return
				case "send":
					if entity.Proc.EXE.ProcessState != nil && entity.Proc.EXE.ProcessState.ExitCode() != 0 {
						log.Info("Server crashed.")
						retGram.Code = status.ServerCrashed.ToInt()
						retGram.Message = status.ServerCrashed.Message()
						detail["server_id"] = cfg.Server.ID
						retGram.Detail = detail
						msg, _ := json.Marshal(&retGram)
						err = entity.Conn.SendTextMessage(string(msg))
						if err != nil {
							log.Warn(err)
						}
						return
					}
					util.SendToConsole(cfg.Server.Game, retGram.Detail["message"].(string))
					thisCommand = ""
					return
				case "stop":
					if entity.Proc.EXE.ProcessState != nil {
						log.Info("Server is terminated. Stop the daemon.")

						detail["server_id"] = cfg.Server.ID
						retGram.Code = status.ServerExited.ToInt()
						retGram.Message = status.ServerExited.Message()
						retGram.Detail = detail
						msg, _ := json.Marshal(&retGram)
						err = entity.Conn.SendTextMessage(string(msg))
						if err != nil {
							log.Warn(err)
						}
						os.Exit(0)
					}
					detail := make(map[string]interface{})
					detail["server_id"] = cfg.Server.ID
					retGram.Code = status.ServerStopping.ToInt()
					retGram.Message = status.ServerStopping.Message()
					retGram.Detail = detail
					msg, _ := json.Marshal(&retGram)
					err := entity.Conn.SendTextMessage(string(msg))
					if err != nil {
						log.Warn(err)
					}
					util.ForceStopServer(cfg)
					retGram.Code = status.ServerExited.ToInt()
					retGram.Message = status.ServerExited.Message()
					msg, _ = json.Marshal(&retGram)
					err = entity.Conn.SendTextMessage(string(msg))
					if err != nil {
						log.Warn(err)
					}
					os.Exit(0)
				case "update":
					// TODO: In the future can be done via only daemon.
					if entity.Proc.EXE.ProcessState != nil {
						log.Info("Server is terminated. Stop the daemon.")

						detail["server_id"] = cfg.Server.ID
						retGram.Code = status.ServerExited.ToInt()
						retGram.Message = status.ServerExited.Message()
						retGram.Detail = detail
						msg, _ := json.Marshal(&retGram)
						err = entity.Conn.SendTextMessage(string(msg))
						if err != nil {
							log.Warn(err)
						}
						os.Exit(0)
					}
					detail := make(map[string]interface{})
					detail["server_id"] = cfg.Server.ID
					retGram.Code = status.ServerStopping.ToInt()
					retGram.Message = status.ServerStopping.Message()
					retGram.Detail = detail
					msg, _ := json.Marshal(&retGram)
					err := entity.Conn.SendTextMessage(string(msg))
					if err != nil {
						log.Warn(err)
					}
					util.ForceStopServer(cfg)
					retGram.Code = status.ServerExited.ToInt()
					retGram.Message = status.ServerExited.Message()
					msg, _ = json.Marshal(&retGram)
					err = entity.Conn.SendTextMessage(string(msg))
					if err != nil {
						log.Warn(err)
					}
					os.Exit(0)
				case "start", "backup", "install", "remove", "search", "validate":
					log.Warn(fmt.Sprintf("Invalid command \"%v\"", command))
					thisCommand = ""
					return
				default:
					log.Warn(fmt.Sprintf("Unknown command \"%v\"", command))
					thisCommand = ""
					return
				}
			default:
				fmt.Println(msg)
			}
		}
	})
	entity.Conn.Connect()
}

func listenSignal(cfg *conf.Config) {
	c := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		s := <-c
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Info("Server is now shutting down...")
			p, _ := os.FindProcess(entity.Proc.EXE.Process.Pid)
			detail := make(map[string]interface{})
			detail["server_id"] = cfg.Server.ID
			retGram.Code = status.ServerStopping.ToInt()
			retGram.Message = status.ServerStopping.Message()
			retGram.Detail = detail
			msg, _ := json.Marshal(&retGram)
			err := entity.Conn.SendTextMessage(string(msg))
			if err != nil {
				log.Warn(err)
			}
			if p == nil {
				retGram.Code = status.ServerExited.ToInt()
				retGram.Message = status.ServerExited.Message()
				msg, _ := json.Marshal(&retGram)
				err = entity.Conn.SendTextMessage(string(msg))
				if err != nil {
					log.Warn(err)
				}
			} else {
				util.ForceStopServer(cfg)
				retGram.Code = status.ServerExited.ToInt()
				retGram.Message = status.ServerExited.Message()
				msg, _ := json.Marshal(&retGram)
				err = entity.Conn.SendTextMessage(string(msg))
				if err != nil {
					log.Warn(err)
				}
			}
			log.Info("Server exited.")
			return
		default:
			return
		}
	}
}

func initQueue() {
	outputQueue = queue.New()
	// If message > 500, then pop
	go func() {
		for {
			if outputQueue.Len() > 500 {
				outputQueue.Dequeue()
			}
		}
	}()
}

func startServer(cfg *conf.Config) {
	exeDir, exeName := util.CombineSeverPath(cfg)
	detail := make(map[string]interface{})
	go func() {
		count := 0
		for {
			entity.Proc.EXE = &exec.Cmd{
				Path:  fmt.Sprintf("%v/%v", exeDir, exeName),
				Dir:   exeDir,
				Args:  append([]string{fmt.Sprintf("%v/%v", exeDir, exeName)}, cfg.Server.Args...),
				Env:   os.Environ(),
				Stdin: os.Stdin, // need to redirect stdin to ensure that we can truly write it to the console.
			}

			entity.Proc.EXE.Stderr = entity.Proc.EXE.Stdout
			// stdout导出到两个地方
			// 1. 日志文件
			// 2. 队列
			// 原因:
			// io.Reader会阻塞命令执行(mlgb)，但是我们所期望的是：
			// 1. 命令行执行后正常进行
			// 2. 当用户要求attach服务器的时候，将控制台信息输出给用户
			// 因此，为了避免阻塞现象，将文件输出转向给文件似乎是一个比较正确的选择
			// (毕竟我们又不需要看控制台输出了什么。。对吧)
			// 至于为什么要队列：为了实现那种所谓的，比如说：
			// 你连接到了控制台，你可以得到最近500条控制台消息输出，云云
			// 虽然说这里可能会导致说编码后的JSON会不会过大?但是这个方案某种程度上来说确实是可行的
			// 更何况队列是FIFO，然后500条消息以后自行FIFO更新队列，怎么看都好像没啥大问题吧，对吧
			// attach是最难的部分，这部分做完了，应该就可以拿出来见人了。
			// 剩下的都是小鱼小虾，好搞得很
			o, _ := entity.Proc.EXE.StdoutPipe()
			go func() {
				logPath := fmt.Sprintf("%v/log/server/%v/L%v.log", os.Getenv("GSM_ROOT"), os.Getenv("GSM_SERVER_ID"), time.Now().Format("20060102"))
				f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					log.Panic(err)
					return
				}
				defer f.Close()
				for {
					data := make([]byte, 4096)
					_, err := o.Read(data)
					// Noting that we need to write RAW message to ensure that everything is correct.
					// Erasing all NUL character.
					data = bytes.Trim(data, "\x00")
					out := string(data)
					f.WriteString(out)
					// 将消息塞入队列中
					// 只存储最近500条消息
					outputQueue.Enqueue(out)
					if err != nil {
						break
					}
				}
			}()
			err := entity.Proc.Start()
			if err != nil {
				log.Error("Failed to startup the server.")
				detail["server_id"] = cfg.Server.ID
				detail["error_message"] = err.Error()
				retGram = &RetGram{
					Role:    misc.Role,
					Code:    status.ServerFailedToStart.ToInt(),
					Message: status.ServerFailedToStart.Message(),
					Detail:  detail,
				}
				msg, _ := json.Marshal(&retGram)
				entity.Conn.SendTextMessage(string(msg))
				if count > cfg.Server.MaxRestartCount {
					entity.Conn.Close()
					os.Exit(0)
				}
				time.Sleep(1 * time.Second)
				continue
			}
			count = 0
			time.Sleep(200 * time.Millisecond)
			windows.ShowWindow(entity.Proc.MainWindowHandle, windows.SW_HIDE)
			detail["server_id"] = cfg.Server.ID
			detail["server_pid"] = entity.Proc.EXE.Process.Pid
			detail["daemon_pid"] = os.Getpid()
			retGram = &RetGram{
				Role:    misc.Role,
				Code:    status.ServerStarted.ToInt(),
				Message: status.ServerStarted.Message(),
				Detail:  detail,
			}
			msg, _ := json.Marshal(&retGram)
			entity.Conn.SendTextMessage(string(msg))

			err = entity.Proc.EXE.Wait()
			if err != nil {
				log.Warn("Server crashed. Message:", err.Error())
				for k := range detail {
					delete(detail, k)
				}
				detail["server_id"] = cfg.Server.ID
				detail["error_message"] = err.Error()
				retGram = &RetGram{
					Role:    misc.Role,
					Code:    status.ServerCrashed.ToInt(),
					Message: status.ServerCrashed.Message(),
					Detail:  detail,
				}
				msg, _ := json.Marshal(&retGram)
				entity.Conn.SendTextMessage(string(msg))
				performCountdown(cfg)
				continue
			}
			if thisCommand == "restart" {
				fmt.Println("Restarting the server...")
				thisCommand = ""
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}

		for k := range detail {
			delete(detail, k)
		}
		detail["server_id"] = cfg.Server.ID
		retGram = &RetGram{
			Role:    misc.Role,
			Code:    status.ServerExited.ToInt(),
			Message: status.ServerExited.Message(),
			Detail:  detail,
		}
		msg, _ := json.Marshal(&retGram)
		entity.Conn.SendTextMessage(string(msg))
		entity.Conn.Close()
		os.Exit(0)
	}()
}

func performCountdown(cfg *conf.Config) {
	if cfg.Server.RestartAfterDelay > 0 {
		log.Info(fmt.Sprintf("Server will be restarted after %v seconds.", cfg.Server.RestartAfterDelay))
		detail := make(map[string]interface{})
		detail["server_id"] = cfg.Server.ID
		retGram = &RetGram{
			Role:    misc.Role,
			Code:    status.ServerRestartCountingDown.ToInt(),
			Message: status.ServerRestartCountingDown.Message(),
			Detail:  detail,
		}
		msg, _ := json.Marshal(&retGram)
		entity.Conn.SendTextMessage(string(msg))

		time.Sleep(time.Second * time.Duration(cfg.Server.RestartAfterDelay))
	}
}

func main() {
	cfg := conf.Init()
	log.Init(cfg.Log)
	initQueue()
	handleClientWSConnection(cfg)
	startServer(cfg)
	listenSignal(cfg)
}
