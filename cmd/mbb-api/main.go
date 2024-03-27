package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	flag "github.com/spf13/pflag"

	"github.com/mmart-pro/mart-brute-blocker/internal/app"
	"github.com/mmart-pro/mart-brute-blocker/internal/config"
)

func main() {
	var (
		configFlag string
		hostFlag   string
		portFlag   int
		help       bool
		version    bool
	)

	flag.BoolVarP(&help, "help", "?", false, "help")
	flag.BoolVarP(&version, "version", "v", false, "версия приложения")
	flag.StringVarP(&configFlag, "config", "c", "config.json", "json-файл конфигурации")
	flag.StringVarP(&hostFlag, "host", "h", "", "адрес сервера")
	flag.IntVarP(&portFlag, "port", "p", 0, "порт сервера")

	flag.Parse()

	if help {
		flag.Usage()
		return
	}
	if version {
		printVersion()
		return
	}

	// config
	runtimeConfig := config.APIConfig{}
	if configFlag != "" {
		cfg, err := config.NewAPIConfig(configFlag)
		if err != nil {
			log.Fatal(fmt.Errorf("error read config from %s: %w", configFlag, err))
		}
		runtimeConfig = cfg
	}

	// параметр из командной строки перекрывает значение конфига
	if hostFlag != "" {
		runtimeConfig.GrpcConfig.GrpcHost = hostFlag
	}
	if portFlag != 0 {
		runtimeConfig.GrpcConfig.GrpcPort = strconv.Itoa(portFlag)
	}

	err := app.NewApp(runtimeConfig).
		Startup(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
