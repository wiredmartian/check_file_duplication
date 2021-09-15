package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	var input string
	c := make(chan string)
	duplicates := make(map[int]string)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a file path: ")
	for scanner.Scan() {
		input = scanner.Text()
		break
	}
	fmt.Println(timestamppb.Now().AsTime().Local())
	go ProcessDuplicates(input, c)
	i := 0
	for path := range c {
		duplicates[i] = path
		i++
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		WriteResults(duplicates)
	}()
	wg.Wait()
	fmt.Println(timestamppb.Now().AsTime().Local())
}

func ProcessDuplicates(root string, c chan string) {
	var err error
	var fs []File
	var files []File
	_, err = os.Stat(root)
	if err != nil {
		log.Fatalf(err.Error())
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
	if err != nil {
		log.Fatalf(err.Error())
	}
	for i := 0; i < len(fs); i++ {
		f, err := GetFile(fs[i].path)
		if err != nil {
			fmt.Println(err.Error())
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
				c <- files[j].path
				c <- pFile.path // maybe??
				fmt.Println(files[j].path)
				break
			}
		}
		files = append(files, pFile)
	}
	close(c)
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
func WriteResults(files map[int]string) {
	if len(files) != 0 {
		var err error
		_, err = os.Stat("./outputs/")
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir("./outputs/", os.ModePerm)
				if err != nil {
					log.Fatal(err.Error())
				}
			} else {
				log.Fatal(err.Error())
			}
		}
		txtPath := fmt.Sprintf("./outputs/%v.txt", timestamppb.Now().Seconds)
		textF, err := os.OpenFile(txtPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, f := range files {
			_, err := textF.WriteString(f + "\n")
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		_ = textF.Close()
	}
}

type File struct {
	path     string
	file     os.FileInfo
	fileHash [16]byte
}
