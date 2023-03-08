package api_client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ToshaRotten/fileService/api_client/config"
	"github.com/ToshaRotten/fileService/api_client/file_helper"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var (
	nilHashConst [20]byte = [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	client                = &http.Client{}
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
	Config     *config.Config
	Logger     *logrus.Logger
	FileHelper *file_helper.FileHelper
}

// New ..
func New(config *config.Config) *APIClient {
	c := &APIClient{
		Config:     config,
		Logger:     logrus.New(),
		FileHelper: file_helper.New(),
	}
	c.configureLogger()
	return c
}

func (c *APIClient) configureLogger() {
	c.Logger.Level = logrus.DebugLevel
}

// GetFileList ..
func (c *APIClient) GetFileList() {
	c.FileHelper.UpdateFiles()
	var temp []Foo
	req, err := http.NewRequest("GET", "http://"+c.Config.Host+c.Config.Port+"/file/get/file_list", nil)
	if err != nil {
		c.Logger.Error(err)
	}
	response, err := client.Do(req)
	if err != nil {
		c.Logger.Error(err)
	}
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&temp)
	if err != nil {
		c.Logger.Error(err)
	}
	if len(temp) < 1 {
		c.Logger.Info("No files")
	}
	for _, foo := range temp {
		c.Logger.Debug("file name: ", foo.Bar.Name, "; hash: ", foo.Bar.Hash)
	}
	defer response.Body.Close()
}

// DeleteFileByName ..
func (c *APIClient) DeleteFileByName(fileName string) {
	c.FileHelper.UpdateFiles()
	var temp Bar
	temp.Name = fileName
	body, err := json.Marshal(temp)
	if err != nil {
		c.Logger.Error(err)
	}
	fmt.Println(temp)
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

// GetFileByName ..
func (c *APIClient) GetFileByName(fileName string) {
	c.FileHelper.UpdateFiles()
	var temp Bar
	var file Bar
	file.Name = fileName
	body, err := json.Marshal(file)
	if err != nil {
		c.Logger.Error(err)
	}
	req, err := http.NewRequest("GET", "http://"+c.Config.Host+c.Config.Port+"/file/get", bytes.NewBuffer(body))
	if err != nil {
		c.Logger.Error(err)
	}
	response, err := client.Do(req)
	if err != nil {
		c.Logger.Error(err)
	}
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.Logger.Error(err)
	}
	err = json.Unmarshal(respBody, &temp)
	if err != nil {
		c.Logger.Error(err)
	}
	c.Logger.Info(response.Status)

	dst := make([]byte, base64.StdEncoding.EncodedLen(len(temp.Data)))
	_, err = base64.StdEncoding.Decode(dst, []byte(temp.Data))
	if err != nil {
		c.Logger.Error(err)
	}
	err = c.FileHelper.WriteFile(dst, temp.Name)
	if err != nil {
		c.Logger.Error(err)
	}
}

// PutFile ..
func (c *APIClient) PutFile(fileName string) {
	c.FileHelper.UpdateFiles()
	temp, err := c.FileHelper.AllFileData(fileName)
	body, err := json.Marshal(temp)
	if err != nil {
		c.Logger.Error(err)
	}
	req, err := http.NewRequest("PUT", "http://"+c.Config.Host+c.Config.Port+"/file/put", bytes.NewBuffer(body))
	if err != nil {
		c.Logger.Error(err)
	}
	response, err := client.Do(req)
	if err != nil {
		c.Logger.Error(err)
	}
	c.Logger.Info(response.Status)
}

func (c *APIClient) UpdateFile(fileName string) {
	c.FileHelper.UpdateFiles()
	temp, err := c.FileHelper.AllFileData(fileName)
	if err != nil {
		c.Logger.Error(err)
	}
	body, err := json.Marshal(temp)
	if err != nil {
		c.Logger.Error(err)
	}
	req, err := http.NewRequest("POST", "http://"+c.Config.Host+c.Config.Port+"/file/update", bytes.NewBuffer(body))
	if err != nil {
		c.Logger.Error(err)
	}
	response, err := client.Do(req)
	if err != nil {
		c.Logger.Error(err)
	}
	if response.StatusCode == 200 {
		c.Logger.Info(response.Status)
	} else {
		c.Logger.Warn(response.Status)
	}

}
