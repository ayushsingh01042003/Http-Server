package examples

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

func ParallelUrl() {
	var n int
	fmt.Println("Number of URLs")
	fmt.Scanln(&n)

	urlSlice := make([]string, n)
	for i := range len(urlSlice) {
		fmt.Scanln(&urlSlice[i])
	}

	var wg sync.WaitGroup
	for _, url := range urlSlice {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			data, err := hitURL(url)
			if err != nil {
				fmt.Println(err)
				return
			}
	
			title, err := getTitle(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(title)
		}(url)
	}
	wg.Wait()
}

func hitURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getTitle(body string) (string, error) {
	idx := strings.Index(body, "<title>")
	if idx == -1 {
		return "", errors.New("no title found")
	}
	idx += 7

	endIdx := -1
	for i := idx; i < len(body); i++ {
		if body[i] == '<' {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return "", errors.New("some problem when finding title")
	}

	title := body[idx : endIdx]
	return title, nil
}
