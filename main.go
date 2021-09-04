package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a file path: ")
	for scanner.Scan() {
		input = scanner.Text()
		break
	}
	_, err := GetHashedFiles(input)
	if err != nil {
		fmt.Println(err)
	}
}

// GetHashedFiles walk the file dir
// get files exclude dirs
// hash files
func GetHashedFiles(root string) ([]File, error) {
	var fs []File
	var files []File
	var err error
	_, err = os.Stat(root)

	if err != nil {
		return nil, err
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			var f = File{
				path: path,
				file: info,
			}
			fs = append(fs, f)
		}
		if err != nil {
			return err
		}
		return nil
	})
	for i := 0; i < len(fs); i++ {
		f, err := GetFile(fs[i].path)
		if err != nil {
			return nil, err
		}
		var pFile = File{
			path:     fs[i].path,
			file:     fs[i].file,
			fileHash: md5.Sum(f),
		}
		colorReset := "\033[0m"
		colorRed := "\033[31m"

		for j := 0; j < len(files); j++ {
			if fmt.Sprintf("%x", pFile.fileHash) == fmt.Sprintf("%x", files[j].fileHash) {
				fmt.Print(colorRed)
				fmt.Printf("%v\n", files[j].path)
				fmt.Printf("%v\n", pFile.path)
				fmt.Print(colorReset)
				break
			}
		}
		files = append(files, pFile)
	}
	return files, err
}

// GetFile get file bytes from path
func GetFile(filepath string) ([]byte, error) {
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
