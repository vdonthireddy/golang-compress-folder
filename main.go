package main

import (
	"strings"
	"path"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
    codePath := strings.Split(os.Args[1], "=")[1] + "/"
    compress(codePath)
}

func compress(codePath string) {
	files, err := listFiles(codePath)
	if err != nil {
		panic(err)
	}

    zipFileName := path.Base(codePath) + ".zip"

    zipMe(files, path.Dir(codePath), zipFileName)
    
	// for _, f := range files {
	// 	fmt.Println(f)
	// }
	fmt.Println("Done!")
}

func listFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil

}

func zipMe(filepaths []string,  removePath string, target string) error {

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(target, flags, 0644)

	if err != nil {
		return fmt.Errorf("Failed to open zip for writing: %s", err)
	}
	defer file.Close()

	zipw := zip.NewWriter(file)
	defer zipw.Close()

	for _, filename := range filepaths {
		if err := addFileToZip(filename, removePath, zipw); err != nil {
			return fmt.Errorf("Failed to add file %s to zip: %s", filename, err)
		}
	}
	return nil
}

func addFileToZip(filename string, removePath string, zipw *zip.Writer) error {
	file, err := os.Open(filename)

	if err != nil {
		return fmt.Errorf("Error opening file %s: %s", filename, err)
	}
	defer file.Close()

    filename = strings.Replace(filename, removePath, "", -1)

	wr, err := zipw.Create(filename)
	if err != nil {

		return fmt.Errorf("Error adding file; '%s' to zip : %s", filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Error writing %s to zip: %s", filename, err)
	}

	return nil
}
