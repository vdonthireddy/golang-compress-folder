package main

import (
	"strings"
	"os"
)

func main() {
    codePath := strings.Split(os.Args[1], "=")[1] + "/"
    compress(codePath)
}