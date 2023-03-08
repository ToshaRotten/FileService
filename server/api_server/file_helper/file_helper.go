package file_helper

import (
	"crypto/sha1"
	"fmt"
	"github.com/fsnotify/fsnotify"
	_ "github.com/fsnotify/fsnotify"
	"io/fs"
	"io/ioutil"
	"os"
)

type File struct {
	Hash [20]byte `json:"hash"`
	Name string   `json:"name"`
}

type FileHelper struct {
	traceDirectory string
	Files          []File
}

func New() *FileHelper {
	return &FileHelper{
		Files:          nil,
		traceDirectory: "/tmp",
	}
}

func (f *FileHelper) CountOfFiles() int {
	var counter int
	for counter, _ = range f.Files {
		counter++
	}
	return counter
}

func (f *FileHelper) Inotify() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() error {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return err
				}
				if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					err = f.UpdateFiles()
					if err != nil {
						return err
					}
				}
			case err := <-watcher.Errors:
				if err != nil {
					return err
				}
			}
		}
	}()
	err = watcher.Add("." + f.traceDirectory)
	if err != nil {
		return err
	}
	<-done
	return nil
}

// SetTraceDirectory ..
func (f *FileHelper) SetTraceDirectory(path string) {
	f.traceDirectory = path
}

func (f *FileHelper) GetAllNames() []string {
	var names []string
	for _, file := range f.Files {
		names = append(names, file.Name)
	}
	fmt.Println(names)
	return nil
}

func (f *FileHelper) GetAllHashes() [][20]byte {
	var hashes [][20]byte
	for _, file := range f.Files {
		hashes = append(hashes, file.Hash)
	}
	return nil
}

func (f *FileHelper) CheckFileByName(fileName string) bool {
	for _, file := range f.Files {
		if file.Name == fileName {
			return true
		}
	}
	return false
}

func (f *FileHelper) CheckFileByHash(fileHash [20]byte) bool {
	for _, file := range f.Files {
		if file.Hash == fileHash {
			return true
		}
	}
	return false
}

func (f *FileHelper) ClearFiles() {
	f.Files = nil
}

func (f *FileHelper) UpdateFiles() error {
	f.ClearFiles()
	var files []fs.FileInfo
	var temp File
	files, err := ioutil.ReadDir("." + f.traceDirectory)
	if err != nil {
		return err
	}
	for _, file := range files {
		temp.Name = file.Name()
		err, temp.Hash = f.GetFileHash(file.Name())
		if err != nil {
			return err
		}
		if !f.CheckFileByName(temp.Name) {
			f.Files = append(f.Files, temp)
		}
	}
	return nil
}

// GetFileHash ..
func (f *FileHelper) GetFileHash(fileName string) (error, [20]byte) {
	var hash [20]byte
	data, err := os.ReadFile("." + f.traceDirectory + "/" + fileName)
	if err != nil {
		return err, hash
	}
	hash = sha1.Sum(data)
	return nil, hash
}

// WriteFile ..
func (f *FileHelper) WriteFile(data []byte, fileName string) error {
	file, err := os.Create("." + f.traceDirectory + "/" + fileName)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	f.UpdateFiles()
	return nil
}

// GetFileData ..
func (f *FileHelper) GetFileData(fileName string) (error, []byte) {
	data, err := os.ReadFile("." + f.traceDirectory + "/" + fileName)
	if err != nil {
		return err, nil
	}
	return nil, data
}

// RemoveFile ..
func (f *FileHelper) RemoveFile(fileName string) error {
	err := os.Remove("." + f.traceDirectory + "/" + fileName)
	if err != nil {
		return err
	}
	f.UpdateFiles()
	return nil
}
