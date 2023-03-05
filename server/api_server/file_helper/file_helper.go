package file_helper

import (
	"crypto/sha1"
	"fmt"
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
	for _, _ = range f.Files {
		counter++
	}
	return counter
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

func (f *FileHelper) UpdateFiles() error {
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
