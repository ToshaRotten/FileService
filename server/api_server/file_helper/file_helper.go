package file_helper

import (
	"crypto/sha1"
	"io/ioutil"
	"os"
)

type File struct {
	Hash [20]byte
	Name string
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

// SetTraceDirectory ..
func (f *FileHelper) SetTraceDirectory(path string) {
	f.traceDirectory = path
}

func (f *FileHelper) UpdateFiles() error {
	var temp File
	files, err := ioutil.ReadDir("." + f.traceDirectory)
	if err != nil {
		return err
	}
	for _, file := range files {
		temp.Name = file.Name()
		err, temp.Hash = GetFileHash(file.Name())
		f.Files = append(f.Files, temp)
	}
	return nil
}

// GetFileHash ..
func GetFileHash(path string) (error, [20]byte) {
	var hash [20]byte
	data, err := os.ReadFile("." + path)
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
	return nil
}
