package examples

import (
	"fmt"
	"sync"
)

func ParallelPrimes() {
	fmt.Println("Enter N")
	var n int
	fmt.Scanln(&n)

	jobs := make(chan int, n)
	results := make(chan int, n)

	for i := 2; i < n; i++ {
		jobs <- i
	}

	var wg sync.WaitGroup
	for i := range 3 {
		wg.Add(1)
		go worker(i, results, jobs, &wg)
	}

	close(jobs)
	wg.Wait()

	for i := 0; i < n; i++ {
		fmt.Println(<-results)
	}
}

func worker(workerId int, results chan<- int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("WorkerId: %d, Number: %d\n", workerId, j)
		if isPrime(j) {
			results <- j
		}
	}
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}

	for i := 2; i*i <= n; i++ {
		if n % i == 0 {
			return false
		}
	}
	return true
}