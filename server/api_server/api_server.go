package api_server

//• GET, получить список файлов и их hash по содержимому на сервере в папке /tmp
//• GET, получить по имени файла из папки /tmp файл, если файла нет - возвращать ошибку
//• PUT, положить файл в папку /tmp, если уже файл есть - возвращать ошибку
//• POST, обновить файл в папке /tmp новым файлом, если файла нет - возвращать ошибку, если файл есть и по hash совпадает - возвращать что не требуется обновление
//• DELETE, удалить по имени файла из папки /tmp файл, если файла нет - возвращать ошибку

import (
	"fileService/api_server/config"
	"fileService/api_server/file_helper"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

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
	s.configureLogger()
	s.configureRouter()
	s.configureFileHelper()
	s.Logger.Info("Server is started ...")
	s.Logger.Info("Bind addr: http://", s.Config.Host+s.Config.Port)
	err := http.ListenAndServe(s.Config.Host+s.Config.Port, s.Router)
	if err != nil {
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
	s.Logger.Trace("FileHelper ...")
	s.FileHelper.SetTraceDirectory(s.Config.TraceDirectory)
	err := s.FileHelper.UpdateFiles()
	if err != nil {
		s.Logger.Error(err)
	}
}

// configureRouter ..
func (s *APIServer) configureRouter() {
	s.Logger.Trace("Router ...")
	s.Router.HandleFunc("/file/get/file_list", s.getFileList())
	s.Router.HandleFunc("/file/get/", s.getFile())
	s.Router.HandleFunc("/file/put/", s.putFile())
	s.Router.HandleFunc("/file/update/", s.updateFile())
	s.Router.HandleFunc("/file/delete/", s.deleteFile())
}

func (s *APIServer) getFileList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (s *APIServer) getFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (s *APIServer) putFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (s *APIServer) updateFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (s *APIServer) deleteFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
