package main

import (
	"fmt"
	"net/http"
	"sync"
)

type result struct {
	url    string
	status string
}

func check(url string, wg *sync.WaitGroup, results chan result) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		results <- result{url: url, status: "ERROR"}
		return
	}
	defer resp.Body.Close()
	results <- result{url: url, status: resp.Status}
}

func main() {
	urls := []string{"http://google.com", "http://facebook.com", "http://twitter.com"}

	var wg sync.WaitGroup

	results := make(chan result)

	for _, url := range urls {
		wg.Add(1)
		go check(url, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Printf("Results chan type : %T, results chan value %v \n", results, results)

	for result := range results {
		fmt.Println(result)
	}
}
