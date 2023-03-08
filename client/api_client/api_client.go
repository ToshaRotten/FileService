package api_client

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"log"
	"main/config"
	"net/http"
)

var (
	client = &http.Client{}
)

type Bar struct {
	Hash [20]byte `json:"hash"`
	Name string   `json:"name"`
	Data string   `json:"data"`
}

type Foo struct {
	Bar Bar `json:"bar"`
}

type APIClient struct {
	Config *config.Config
	Logger *logrus.Logger
}

func New(config *config.Config) *APIClient {
	c := &APIClient{
		Config: config,
		Logger: logrus.New(),
	}
	return c
}

func (c *APIClient) GetFileList() {
	req, err := http.NewRequest("GET", "http://"+c.Config.Host+c.Config.Port+"/file/get/file_list", nil)
	if err != nil {
		c.Logger.Error(err)
	}
	response, err := client.Do(req)
	if err != nil {
		c.Logger.Error(err)
	}
	var result []Foo
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&result)
	if err != nil {
		c.Logger.Error(err)
	}
	log.Println(result)
	defer response.Body.Close()
}

func (c *APIClient) DeleteFileByName(fileName string) {
	hash := [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	foo := Foo{
		Bar: Bar{
			Hash: hash,
			Name: "",
			Data: "",
		},
	}
	body, _ := json.Marshal(foo)
	req, err := http.NewRequest("DELETE", "http://"+c.Config.Host+c.Config.Port+"/file/delete", bytes.NewBuffer(body))
	if err != nil {
		c.Logger.Error(err)
	}
	response, err := client.Do(req)
	if err != nil {
		c.Logger.Error(err)
	}
	c.Logger.Info(response.Status)
}
