package main

import (
	"fmt"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	goev "github.com/MuggleWei/goev"
	demo "github.com/MuggleWei/goev/examples/demo"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ClientConfig struct {
	// service
	host string
	port uint

	// log
	logLevel         string
	logFile          string
	logEnableConsole bool
}

func initConfig(filePath string) (*ClientConfig, error) {
	pflag.StringP("config", "f", "./config/client.yml", "config file")

	pflag.StringP("host", "H", "0.0.0.0", "bind host")
	pflag.UintP("port", "P", 8080, "listen port")

	pflag.String("log.level", "info", "log level")
	pflag.String("log.file", "./log/client.log", "log file path")
	pflag.Bool("log.console", false, "enable/disable log console output")

	pflag.Parse()

	// config
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			fmt.Fprintf(os.Stdout, "config file not found: %v\n", filePath)
		} else {
			panic(fmt.Errorf("error config file: %v", err))
		}
	}

	// viper bind command line
	viper.BindPFlags(pflag.CommandLine)

	return &ClientConfig{
		host: viper.GetString("host"),
		port: viper.GetUint("port"),

		logLevel:         viper.GetString("log.level"),
		logFile:          viper.GetString("log.file"),
		logEnableConsole: viper.GetBool("log.console"),
	}, nil
}

func PrintConfig(cfg *ClientConfig) {
	log.Info("--------------------")
	log.Info("auth config:")
	log.Infof("host=%v, port=%v", cfg.host, cfg.port)
	log.Infof("log.level=%v, log.file=%v, log.console=%v",
		cfg.logLevel, cfg.logFile, cfg.logEnableConsole)
}

func main() {
	// try get config file path
	configFilePath := "config/client.yml"
	if len(os.Args) == 3 && os.Args[1] == "-f" {
		configFilePath = os.Args[2]
	}

	// init config
	cfg, err := initConfig(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed init config: %v", configFilePath)
		os.Exit(1)
	}

	// init log
	demo.InitLog(cfg.logLevel, cfg.logFile, cfg.logEnableConsole)
	PrintConfig(cfg)

	// init handle
	handle := NewClientHandle()

	// init evloop
	evloop := goev.NewEvloop()
	evloop.SetTimerTick(3 * time.Second)
	evloop.SetCallbackOnAddConn(handle.OnAddSession)
	evloop.SetCallbackOnClose(handle.OnClose)
	evloop.SetCallbackOnMessage(handle.onMessage)
	evloop.SetCallbackOnTimer(handle.onTimer)

	// handle set evloop
	handle.evloop = evloop

	// listen
	servAddr := fmt.Sprintf("%v:%v", cfg.host, cfg.port)
	conn, err := net.Dial("tcp", servAddr)
	if err != nil {
		log.Errorf("failed dial addr: %v", servAddr)
		os.Exit(1)
	}
	log.Infof("success dial: %v", servAddr)

	session := &ClientSession{
		conn: conn,
		codec: &demo.BytesCodec{
			MaxPayloadLimit: 512 * 1024,
		},
		userData: &ClientUserData{
			servAddr: servAddr,
		},
	}
	evloop.AddSession(session)

	evloop.Run()
}
