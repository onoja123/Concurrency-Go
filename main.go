package main

import (
	"fmt"
	"math"
	"time"
)

const (
	DoExit = iota
	ExitOK
)

type ControlMessage int

type Job struct {
	value int
}

type Result struct {
	job   Job
	value int
}

func square(jobs <-chan Job, results chan<- Result, controlMsg chan ControlMessage) {
	for {
		select {
		case msg := <-controlMsg:
			switch msg {
			case DoExit:
				fmt.Println("exiting goroutine")
				controlMsg <- ExitOK
			default:
				panic("invalid control message")
			}
		case currentJob := <-jobs:
			res := Result{job: currentJob, value: int(math.Pow(float64(currentJob.value), 2))}
			results <- res
		}
	}
}
func main() {
	jobs := make(chan Job, 64)
	results := make(chan Result, 64)
	control := make(chan ControlMessage)

	go square(jobs, results, control)

	for i := 1; i <= 20; i++ {
		jobs <- Job{i}
	}

	for {
		select {
		case result := <-results:
			fmt.Printf("Job %d = %d\n", result.job.value, result.value)
		case <-time.After(500 * time.Millisecond):
			fmt.Println("sending timeout signal")
			control <- DoExit
			<-control
			fmt.Println("got it! exiting now :)")
			return
		}
	}
}
