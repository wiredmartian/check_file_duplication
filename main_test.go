package main

import (
	"testing"
)

//TestGetFile test get file
// expecting file bytes
func TestGetFile(t *testing.T) {
	file, err := GetFile("./test/test/image1.jpg")

	if err != nil {
		t.Errorf("test failed: input: %v, expceted: %v, received: %v", "./test/test/image1.jpg", "[]byte", err.Error())
	}
	if len(file) == 0 {
		t.Errorf("test failed: expected []byte length > 0, received: %v", len(file))
	}
}

// TestEmptyDir test get file with directory
// expecting path is a dir error
func TestEmptyDir(t *testing.T) {
	file, err := GetFile("./test/test1")
	if err == nil || len(file) > 0 {
		t.Errorf("test failed: expected not a directory err, received %v file bytes instead", len(file))
	}
}

func TestGetHashedFiles(t *testing.T) {
	fs, err := GetHashedFiles("./test")
	if err != nil {
		t.Fatalf("test failed: expected []File, received: %v", err.Error())
	}
	if fs == nil || len(fs) != 8 {
		t.Errorf("test failed: expected 8 []File, received: %v", len(fs))
	}
}
