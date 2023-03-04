package main

import (
	"flag"
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"main/api_client"
	"main/config"
)

var (
	configPath string
	port       string
	host       string
	logger     = logrus.New()
)

func init() {
	flag.StringVar(&configPath, "path", "configs/config.yaml", "Set host")
	flag.StringVar(&host, "host", "", "Set host")
	flag.StringVar(&port, "port", "", "Set port")
}

func main() {
	logo := figure.NewFigure("FileService - Client", "", true)
	logo.Print()

	conf := config.New()
	conf.ParseFile(configPath)
	if host != "" {
		conf.Host = host
	}
	if port != "" {
		conf.Port = port
	}
	client := api_client.New(conf)
	client.GetFileList()
}
