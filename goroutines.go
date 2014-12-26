package main

import (
	"fmt"
	"time"
	"sync"
)

var wg sync.WaitGroup

type Task interface {
	Name() string
	Exec() string
}

type STask struct {
	Word string
}

func (t STask) Name() string {
	return t.Word
}

func (t STask) Exec() string {
	return t.Word
}


func worker(id int, jobs <-chan Task, results chan<- string) {
	for  {
		j, ok := <-jobs
		if !ok {
			wg.Done()
			break
		}
		fmt.Println("worker", id, "processing job", j.Name())

		time.Sleep(time.Second)
		results <- j.Exec()
	}
}

func printer(results <-chan string) {

	for r := range results {
		fmt.Println(r)
	}
}

func main() {
	jobs := make(chan Task, 100)
	results := make(chan string, 100)
	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results)
	}

	go printer(results)

	jobs <- STask{"Iain"}
	jobs <- STask{"Ian"}
	jobs <- STask{"John"}
	jobs <- STask{"Sally"}
	jobs <- STask{"James"}
	jobs <- STask{"Adrian"}
	
	close(jobs)
	wg.Wait()
}
