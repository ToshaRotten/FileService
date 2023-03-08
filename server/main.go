package main

import (
	"flag"
	"github.com/ToshaRotten/fileService/APIServer"
	"github.com/ToshaRotten/fileService/APIServer/config"
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
)

var (
	configPath string
	logger     = logrus.New()
)

func init() {
	flag.StringVar(&configPath, "configPath", "configs/config.yaml", "Path to server config file")
}

func main() {
	logo := figure.NewFigure("FileService - Server", "", true)
	logo.Print()
	logger.Info("Starting ...")
	flag.Parse()
	conf := config.New()

	err := conf.ParseFile(configPath)
	if err != nil {
		logger.Error(err)
	}
	server := APIServer.New()
	if err = server.Start(conf); err != nil {
		logger.Error(err)
	}
}
