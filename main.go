package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	files, err := getHashedFiles("/home/solomizi/Downloads/wetransfer_imgl8632-jpg_2021-08-06_1930")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		fmt.Printf("MD5 File HASH: %x\n", file.fileHash)
		fmt.Println(file.path)
	}
}

// walk the file dir
// get files exclude dirs
// hash files
func getHashedFiles(root string) ([]File, error) {
	var fs []File
	var files []File
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			var f = File{
				path: path,
				file: info,
			}
			fs = append(fs, f)
		}
		return nil
	})
	for i := 0; i < len(fs); i++ {
		f, err := getFile(fs[i].path)
		if err != nil {
			return nil, err
		}
		var pFile = File{
			path:     fs[i].path,
			file:     fs[i].file,
			fileHash: md5.Sum(f),
		}
		files = append(files, pFile)
	}
	return files, err
}

// get file bytes from path
func getFile(filepath string) ([]byte, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

type File struct {
	path     string
	file     os.FileInfo
	fileHash [16]byte
}
