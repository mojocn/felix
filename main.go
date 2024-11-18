package main

import (
	"github.com/mojocn/felix/cmd"
	"log"
)

var buildTime, gitHash string

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd.Execute(buildTime, gitHash)
}
