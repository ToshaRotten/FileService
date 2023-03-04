package api_client

import (
	"log"
	"main/config"
	"net/http"
)

var (
	client = &http.Client{}
)

type APIClient struct {
	Config *config.Config
}

func New(config *config.Config) *APIClient {
	c := &APIClient{
		Config: config,
	}
	return c
}

func (c *APIClient) GetFileList() {
	req, err := http.NewRequest("GET", "http://"+c.Config.Host+c.Config.Port+"/file/get/file_list", nil)
	if err != nil {
		log.Fatalln(err)
	}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	//json.NewDecoder(response.Body).Decode(&result)
	log.Println(result)
	defer response.Body.Close()
}
