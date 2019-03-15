package cronjob

import "fmt"

func RunTasks() {
	// Do jobs with params
	Every(1).Second().Do(taskWithParams, 1, "hello")

	// Do jobs without params
	Every(1).Second().Do(task)
	Every(2).Seconds().Do(task)
	Every(1).Minute().Do(task)
	Every(2).Minutes().Do(task)
	Every(1).Hour().Do(task)
	Every(2).Hours().Do(task)
	Every(1).Day().Do(task)
	Every(2).Days().Do(task)

	// Do jobs on specific weekday
	Every(1).Monday().Do(task)
	Every(1).Thursday().Do(task)

	// function At() take a string like 'hour:min'
	Every(1).Day().At("10:30").Do(task)
	Every(1).Monday().At("18:30").Do(task)

	// remove, clear and next_run
	_, time := NextRun()
	fmt.Println(time)

	// Remove(task)
	// Clear()

	// function Start start all the pending jobs
	<-Start()

	// also , you can create a your new scheduler,
	// to run two scheduler concurrently
	s := NewScheduler()
	s.Every(3).Seconds().Do(task)
	<-s.Start()
}
