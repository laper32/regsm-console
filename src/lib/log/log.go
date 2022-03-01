package log

import (
	"fmt"
	"net/url"
	"os"
	"runtime"

	"go.uber.org/zap"
)

type Config struct {
	Debug      bool
	OutputPath []string // stdout and stderr
}

func defaultConfig() *Config {
	return &Config{
		Debug:      false,
		OutputPath: []string{"stdout"}, // stdout, and named by today.
	}
}

var logger *zap.Logger

// https://github.com/marcus-crane/october/blob/main/app.go#L33
// https://github.com/uber-go/zap/issues/621
// We need to manipulate file sink
// Since on windows, the file URI is file:/// => Perhaps uber is using file:// => ERROR
func registerWinFileSink(u *url.URL) (zap.Sink, error) {
	// Remove leading slash left by url.Parse()
	return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}

func Init(conf *Config) {
	if conf == nil {
		conf = defaultConfig()
	}
	var err error
	var cfg zap.Config
	if conf.Debug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	// Fuck off
	// Why then dont manipulate file path?
	// Just bcs they dont need windows s t they don't fix it?
	// wtf is it.
	if runtime.GOOS == "windows" {
		err := zap.RegisterSink("winfile", registerWinFileSink)
		if err != nil {
			panic("Failed to register windows file sink")
		}
		for i := 0; i < len(conf.OutputPath); i++ {
			if conf.OutputPath[i] == "stdout" || conf.OutputPath[i] == "stderr" {
				continue
			}
			conf.OutputPath[i] = "winfile:///" + conf.OutputPath[i]
		}
	}

	cfg.ErrorOutputPaths = conf.OutputPath
	cfg.OutputPaths = conf.OutputPath
	logger, err = cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
}
func Close() {
	err := logger.Sync()
	if err != nil {
		fmt.Println("zap error:", err.Error())
	}
}

func Info(args ...interface{}) {
	logger.Sugar().Info(args)
}

func Error(args ...interface{}) {
	logger.Sugar().Error(args)
}

func Warn(args ...interface{}) {
	logger.Sugar().Warn(args)
}

func Debug(args ...interface{}) {
	logger.Sugar().Debug(args)
}

func Panic(args ...interface{}) {
	logger.Sugar().Panic(args)
}

func CheckErr(err error) {
	if err != nil {
		logger.Sugar().Info(err)
	}
}

func GetLogger() *zap.Logger {
	return logger
}
