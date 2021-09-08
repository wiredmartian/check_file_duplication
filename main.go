package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	fmt.Println(timestamppb.Now().AsTime().Local())
	_, duplicates, err := GetHashedFiles(input)
	if err != nil {
		fmt.Println(err)
	}
	err = WriteResults(duplicates)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(timestamppb.Now().AsTime().Local())
}

// GetHashedFiles walk the file dir
// get files exclude dirs
// hash files
func GetHashedFiles(root string) ([]File, map[int]string, error) {
	var err error
	var fs []File
	var files []File
	duplicates := make(map[int]string)
	_, err = os.Stat(root)

	if err != nil {
		return nil, nil, err
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
	count := 0
	for i := 0; i < len(fs); i++ {
		f, err := GetFile(fs[i].path)
		if err != nil {
			return nil, nil, err
		}
		var pFile = File{
			path:     fs[i].path,
			file:     fs[i].file,
			fileHash: md5.Sum(f),
		}

		for j := 0; j < len(files); j++ {
			newFHashStr := fmt.Sprintf("%x", pFile.fileHash)
			fHashStr := fmt.Sprintf("%x", files[j].fileHash)
			if fHashStr == newFHashStr {
				duplicates[count] = pFile.path
				count++
				duplicates[count] = files[j].path
				fmt.Printf("%v (DUPLICATE)\n", pFile.path)
				fmt.Printf("%v (DUPLICATE)\n", files[j].path)
				count++
				break
			}
		}
		files = append(files, pFile)
	}
	return files, duplicates, err
}

// GetFile get file bytes from path
func GetFile(fpath string) ([]byte, error) {
	var err error
	info, err := os.Stat(fpath)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory")
	}
	file, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// WriteResults write report to a text file
func WriteResults(files map[int]string) error {
	if len(files) != 0 {
		var err error
		_, err = os.Stat("./outputs/")
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir("./outputs/", os.ModePerm)
			} else {
				// don't know
				return err
			}
		}
		txtPath := fmt.Sprintf("./outputs/%v.txt", timestamppb.Now().Seconds)
		textF, err := os.OpenFile(txtPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		for _, f := range files {
			_, err := textF.WriteString(f + "\n")
			if err != nil {
				return err
			}
		}
		_ = textF.Close()
	}
	return nil
}

type File struct {
	path     string
	file     os.FileInfo
	fileHash [16]byte
}
