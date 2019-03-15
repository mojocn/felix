package main

import "github.com/dejavuzhou/felix/cmd"

var buildTime, gitHash string

func main() {
	cmd.Execute(buildTime, gitHash)
}
