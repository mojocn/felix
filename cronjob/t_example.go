package cronjob

import "fmt"

//defining schedule task function here
//then add the function in manger.go
func task() {
	fmt.Println("task one is called")
}
func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}
