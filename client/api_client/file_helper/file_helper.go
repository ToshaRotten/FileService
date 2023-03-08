package file_helper

import (
	"crypto/sha1"
	"encoding/base64"
	"io/fs"
	"io/ioutil"
	"os"
)

type File struct {
	Hash [20]byte `json:"hash"`
	Name string   `json:"name"`
	Data string   `json:"data"`
}

type FileHelper struct {
	traceDirectory string
	Files          []File
}

//New ..
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

func (f *FileHelper) CheckFileByName(fileName string) bool {
	for _, file := range f.Files {
		if file.Name == fileName {
			return true
		}
	}
	return false
}

func (f *FileHelper) ClearFiles() {
	f.Files = nil
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
	return nil
}

func (f *FileHelper) AllFileData(fileName string) (File, error) {
	var file File
	var err error
	file.Name = fileName
	err, file.Hash = f.GetFileHash(fileName)
	if err != nil {
		return File{}, err
	}
	err, data := f.GetFileData(fileName)
	if err != nil {
		return File{}, err
	}
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)
	if err != nil {
		return File{}, err
	}
	file.Data = string(dst)
	return file, nil
}
