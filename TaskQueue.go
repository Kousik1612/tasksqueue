package main

import (
	"fmt"
	"net/http"
)

// Task ...
type Task struct {
	ID           int
	SiteAddress  string
	IsCompleted  bool
	Status       string //completed, failed, timeout
	RetryCount   int
	FailedReason string
}

func main() {
	siteAddress := []string{"https://google.com", "https://fb.com", "https://alivecorr.com"}
	messages := make(chan string, 3)
	statusList := make([]Task, 3)
	TaskQueue(siteAddress, messages, statusList)
	go TaskCleanup(siteAddress, statusList, messages)
}

// TaskQueue ...
func TaskQueue(queue []string, messages chan string, statusList []Task) {
	var taskCount int
	for _, val := range queue {
		select {
		case messages <- val:
			fmt.Println("sent message", val)
			res, err := http.Get(val)
			if err != nil {
				statusList[taskCount] = Task{ID: taskCount + 1, SiteAddress: val, IsCompleted: false, Status: "failed"}
			} else if res.StatusCode == http.StatusGatewayTimeout {
				statusList[taskCount] = Task{ID: taskCount + 1, SiteAddress: val, IsCompleted: false, Status: "timeout"}
			} else {
				statusList[taskCount] = Task{ID: taskCount + 1, SiteAddress: val, IsCompleted: true, Status: "completed"}
			}
		default:
			fmt.Println("no message sent")
		}
		taskCount++
	}
}

// TaskCleanup ...
func TaskCleanup(siteAddress []string, queue []Task, messages chan string) {
	var i int
	for i < len(siteAddress) {
		select {
		case msg := <-messages:
			fmt.Println("received message", msg)
			if queue[i].Status == "failed" && queue[i].RetryCount < 3 {
				queue[i].RetryCount++
				queue = append(queue, queue[i])
				queue[i] = queue[len(queue)-1]
				queue = queue[:len(queue)-1]
				continue
			}
		default:
			fmt.Println("no message received")
		}
		i++
	}
}
