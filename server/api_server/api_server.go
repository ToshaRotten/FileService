package api_server

//• GET, получить список файлов и их hash по содержимому на сервере в папке /tmp
//• GET, получить по имени файла из папки /tmp файл, если файла нет - возвращать ошибку
//• PUT, положить файл в папку /tmp, если уже файл есть - возвращать ошибку
//• POST, обновить файл в папке /tmp новым файлом, если файла нет - возвращать ошибку, если файл есть и по hash совпадает - возвращать что не требуется обновление
//• DELETE, удалить по имени файла из папки /tmp файл, если файла нет - возвращать ошибку

import (
	"encoding/base64"
	"encoding/json"
	"fileService/api_server/config"
	"fileService/api_server/file_helper"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Bar struct {
	Hash [20]byte `json:"hash"`
	Name string   `json:"name"`
	Data string   `json:"data"`
}

type Foo struct {
	Bar Bar `json:"bar"`
}

// APIServer ..
type APIServer struct {
	Config     *config.Config
	Logger     *logrus.Logger
	Router     *mux.Router
	FileHelper *file_helper.FileHelper
}

// New ..
func New() *APIServer {
	s := &APIServer{
		Config:     config.New(),
		Logger:     logrus.New(),
		Router:     mux.NewRouter(),
		FileHelper: file_helper.New(),
	}
	return s
}

// Start ..
func (s *APIServer) Start(config *config.Config) error {
	s.Config = config
	err := s.configureLogger()
	if err != nil {
		s.Logger.Error(err)
		return err
	}
	s.configureRouter()
	s.configureFileHelper()
	s.Logger.Info("Server is started ...")
	s.Logger.Info("Bind addr: http://", s.Config.Host+s.Config.Port)
	err = http.ListenAndServe(s.Config.Host+s.Config.Port, s.Router)
	if err != nil {
		s.Logger.Error(err)
		return err
	}
	return nil
}

func (s *APIServer) configureLogger() error {
	s.Logger.Trace("Logger ...")
	err := s.Logger.Level.UnmarshalText([]byte(s.Config.LogLevel))
	if err != nil {
		return err
	}
	return nil
}

func (s *APIServer) configureFileHelper() {
	s.FileHelper.SetTraceDirectory(s.Config.TraceDirectory)
	s.Logger.Trace("FileHelper ...")
	go func() {
		err := s.FileHelper.Inotify()
		if err != nil {
			s.Logger.Error(err)
		}
	}()
}

// configureRouter ..
func (s *APIServer) configureRouter() {
	s.Logger.Trace("Router ...")
	s.Router.HandleFunc("/file/get/file_list", s.getFileList())
	s.Router.HandleFunc("/file/get", s.getFile())
	s.Router.HandleFunc("/file/put", s.putFile())
	s.Router.HandleFunc("/file/update", s.updateFile())
	s.Router.HandleFunc("/file/delete", s.deleteFile())
}

func (s *APIServer) getFileList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		array := make([]Foo, s.FileHelper.CountOfFiles())
		for i, file := range s.FileHelper.Files {
			array[i].Bar.Hash = file.Hash
			array[i].Bar.Name = file.Name
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		data, err := json.Marshal(&array)
		if err != nil {
			s.Logger.Error(err)
		}
		s.Logger.Trace("GET FILE LIST, FILE LIST:", array)
		s.Logger.Trace("COUNT OF FILES:", s.FileHelper.CountOfFiles())
		_, err = w.Write(data)
		if err != nil {
			s.Logger.Error(err)
		}
	})
}

func (s *APIServer) getFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var temp Bar
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(reqBody, &temp)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusNoContent)
		}
		err, data := s.FileHelper.GetFileData(temp.Name)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusNotFound)
		}
		dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(dst, data)
		w.Header().Add("data", string(dst))
	})
}

func (s *APIServer) putFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var temp Bar
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(reqBody, &temp)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		s.Logger.Trace("PUT FILE REQUEST", temp)
		if s.FileHelper.CheckFileByName(temp.Name) {
			w.WriteHeader(http.StatusAlreadyReported)
		} else {
			dst := make([]byte, base64.StdEncoding.EncodedLen(len(temp.Data)))
			_, err = base64.StdEncoding.Decode(dst, []byte(temp.Data))
			if err != nil {
				s.Logger.Error(err)
			}
			err = s.FileHelper.WriteFile(dst, temp.Name)
			if err != nil {
				s.Logger.Error(err)
				w.WriteHeader(http.StatusNotModified)
			}
		}
	})
}

func (s *APIServer) updateFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var temp Bar
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(reqBody, &temp)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
		}

		if !s.FileHelper.CheckFileByName(temp.Name) {
			w.WriteHeader(http.StatusNoContent)
			s.Logger.Error("File not found")
		}
		err, hash := s.FileHelper.GetFileHash(temp.Name)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusNoContent)
		}
		if hash == temp.Hash {
			s.Logger.Trace("UPDATE FILE REQUEST", temp)
			err = s.FileHelper.RemoveFile(temp.Name)
			if err != nil {
				s.Logger.Error(err)
			}
			dst := make([]byte, base64.StdEncoding.EncodedLen(len(temp.Data)))
			_, err = base64.StdEncoding.Decode(dst, []byte(temp.Data))
			if err != nil {
				s.Logger.Error(err)
			}
			err = s.FileHelper.WriteFile(dst, temp.Name)
			if err != nil {
				s.Logger.Error(err)
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	})
}

func (s *APIServer) deleteFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var temp Bar
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(reqBody, &temp)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		s.Logger.Trace("DELETE FILE", temp)
		err = s.FileHelper.RemoveFile(temp.Name)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusNotModified)
		}
	})
}
