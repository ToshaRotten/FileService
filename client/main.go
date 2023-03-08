package main

import (
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"main/api_client"
	"main/api_client/config"
	"os"
	"time"
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
	logger = logrus.New()
}

func loadAnim() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                   // Start the spinner
	time.Sleep(4 * time.Second)                                 // Run for some time to simulate work
	s.Stop()
}

func Loop(client *api_client.APIClient) {
	done := make(chan bool)
	go func() {
		for {
			fmt.Println("Выбирете действие:")
			fmt.Println("- 1 - Чтобы ПОЛУЧИТЬ список файлов")
			fmt.Println("- 2 - Чтобы ПОЛУЧИТЬ файл по его имени с сервера и сохранить на клиенте(по умлочанию папка /tmp)")
			fmt.Println("- 3 - Чтобы ОТПРАВИТЬ файл на сервер(по умолчанию папка /tmp)")
			fmt.Println("- 4 - Чтобы ОБНОВИТЬ файл на сервере(файлы должны быть разные)")
			fmt.Println("- 5 - Чтобы УДАЛИТЬ файл с сервера")
			fmt.Println("- 6 - Чтобы ВЫЙТИ")
			var selector int
			var fileName string

			fmt.Scanf("%d\n", &selector)

			if selector == 1 {
				logger.Info("Получение списка файлов ...")
				client.GetFileList()
			}
			if selector == 2 {
				fmt.Println("Введите имя файла")
				fmt.Scanf("%s\n", &fileName)
				logger.Info("Получение файла с сервера ...")
				client.GetFileByName(fileName)
			}
			if selector == 3 {
				fmt.Println("Введите имя файла")
				fmt.Scanf("%s\n", &fileName)
				logger.Info("Отправка файла на сервер ...")
				client.PutFile(fileName)
			}
			if selector == 4 {
				fmt.Println("Введите имя файла")
				fmt.Scanf("%s\n", &fileName)
				logger.Info("Обновление файла на сервере ...")
				client.UpdateFile(fileName)
			}
			if selector == 5 {
				fmt.Println("Введите имя файла")
				fmt.Scanf("%s\n", &fileName)
				logger.Info("Удаление файла с сервера ...")
				client.DeleteFileByName(fileName)
			}
			if selector == 6 {
				os.Exit(0)
			}
		}
	}()
	<-done
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
	Loop(client)
}
