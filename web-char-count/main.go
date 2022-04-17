package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	lstChar   = "abcdefghiklmnopqrstuvwxyz"
	urlPrefix = "https://datatracker.ietf.org/doc/html/rfc%d"
)

type CharactersMap struct {
	resultMap map[string]int
	lock      sync.Mutex
}

func main() {
	lstUrl := []string{}

	for i := 1; i < 100; i++ {
		url := fmt.Sprintf(urlPrefix, i)
		lstUrl = append(lstUrl, url)
	}

	start := time.Now()

	// NOT USE GOROUTINE
	result := countCharacters(lstUrl)
	elapse := time.Since(start)
	fmt.Println(result)
	fmt.Println("It tooks ", elapse)

	// USE GOROUTINE
	// result := CharactersMap{resultMap: make(map[string]int)}
	// var wg sync.WaitGroup
	// for _, url := range lstUrl {
	// 	wg.Add(1)
	// 	go (&result).countCharacter(&wg, url)
	// }
	// wg.Wait()

	// elapse := time.Since(start)
	// fmt.Println(result.resultMap)
	// fmt.Println("It tooks ", elapse)
}
