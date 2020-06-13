package main

import "github.com/libragen/felix/cmd"

var buildTime, gitHash string

func main() {
	cmd.Execute(buildTime, gitHash)
}
