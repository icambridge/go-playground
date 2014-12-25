package main

import (
	"fmt"
	"time"
	"sync"
)

var wg sync.WaitGroup

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		if j == -1 {
			wg.Done()
			return
		}
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}

func printer(results <-chan int) {

	for r := range results {
		fmt.Println(r)
	}
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results)
	}

	go printer(results)
	// Here we send 9 `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	for w := 1; w <= 3; w++ {
		jobs <- -1
	}
	wg.Wait()
}
