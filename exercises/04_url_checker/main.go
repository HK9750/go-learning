package main

import (
	"net/http"
	"sync"
)

type result struct {
	url    string
	status string
}

func check(url string, wg *sync.WaitGroup) result {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		return result{url: url, status: "ERROR"}
	}
	return result{url: url, status: resp.Status}
}

func main() {
	urls := []string{"http://site1.com",
		"http://site2.com",
		"http://site3.com",
		"http://site4.com",
		"http://site5.com",
		"http://site6.com",
		"http://site6.com",
		"http://site7.com",
		"http://site8.com",
		"http://site9.com",
		"http://site10.com"}

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go check(url, &wg)
	}
	wg.Wait()
}
