package main

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

func main() {

    files, err := listFiles("./folder")
    if err != nil {
        panic(err)
    }

    zipMe(files, "test.zip")

    for _, f := range files {
        fmt.Println(f)
    }
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

func zipMe(filepaths []string, target string) error {

    flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
    file, err := os.OpenFile(target, flags, 0644)

    if err != nil {
        return fmt.Errorf("Failed to open zip for writing: %s", err)
    }
    defer file.Close()

    zipw := zip.NewWriter(file)
    defer zipw.Close()

    for _, filename := range filepaths {
        if err := addFileToZip(filename, zipw); err != nil {
            return fmt.Errorf("Failed to add file %s to zip: %s", filename, err)
        }
    }
    return nil

}

func addFileToZip(filename string, zipw *zip.Writer) error {
    file, err := os.Open(filename)

    if err != nil {
        return fmt.Errorf("Error opening file %s: %s", filename, err)
    }
    defer file.Close()

    wr, err := zipw.Create(filename)
    if err != nil {

        return fmt.Errorf("Error adding file; '%s' to zip : %s", filename, err)
    }

    if _, err := io.Copy(wr, file); err != nil {
        return fmt.Errorf("Error writing %s to zip: %s", filename, err)
    }

    return nil
}