package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func countCharacters(lstUrl []string) map[string]int {
	returnMap := make(map[string]int)

	for _, url := range lstUrl {

		strBody, err := readBody(url)
		if err != nil {
			continue
		}

		for _, c := range strBody {
			index := strings.Index(lstChar, strings.ToLower(string(c)))
			if index != -1 {
				if _, exist := returnMap[string(lstChar[index])]; exist {
					returnMap[string(lstChar[index])] += 1
				} else {
					returnMap[string(lstChar[index])] = 1
				}
			}
		}

	}
	return returnMap
}

func readBody(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Get data from url %s error: %v", url, err)
		return "", err
	}
	defer resp.Body.Close()

	bytesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read body data from %s error: %v", url, err)
		return "", err
	}
	return string(bytesBody), nil
}
