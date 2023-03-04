package main

import (
	"fileService/api_server"
	"fileService/api_server/config"
	"flag"
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

	//fmt.Println(file_helper.GetFileNames("/tmp"))
	//fmt.Println(file_helper.GetFileHashes("/tmp"))
	//fmt.Println(file_helper.GetFileHash("/tmp/file.txt"))
	//fmt.Println(file_helper.GetFileData("/tmp/file.txt"))

	err := conf.ParseFile(configPath)
	if err != nil {
		logger.Error(err)
	}
	server := api_server.New()
	if err = server.Start(conf); err != nil {
		logger.Error(err)
	}
}
