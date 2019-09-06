package main

import (
	"fmt"
	"log"
	"os"
	"plugin"
)

type GoodDoctor interface {
	HealthCheck() error
}

func ini()  {
	log.Println("main package init function called")
}

func main() {
	log.Println("main function stared")


	// load module 插件你也可以使用go http.Request从远程下载到本地,在加载做到动态的执行不同的功能
	// 1. open the so file to load the symbols

	plug, err := plugin.Open("./plugin_doctor.so")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("plugin opened")

	// 2. look up a symbol (an exported function or variable)
	// in this case, variable Greeter
	doc, err := plug.Lookup("Doctor")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 3. Assert that loaded symbol is of a desired type
	// in this case interface type Greeter (defined above)
	doctor, ok := doc.(GoodDoctor)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	// 4. use the module
	if err := doctor.HealthCheck(); err != nil {
		log.Println("use plugin doctor failed, ", err)
	}
}
