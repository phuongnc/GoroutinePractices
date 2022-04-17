package main

import (
	"strings"
	"sync"
)

func (m *CharactersMap) countCharacter(wg *sync.WaitGroup, url string) {
	defer wg.Done()

	strBody, err := readBody(url)
	if err == nil {
		for _, c := range strBody {
			cString := string(c)
			index := strings.Index(lstChar, strings.ToLower(cString))
			if index != -1 {
				m.lock.Lock()
				if _, exist := m.resultMap[string(lstChar[index])]; exist {
					m.resultMap[string(lstChar[index])] += 1
				} else {
					m.resultMap[string(lstChar[index])] = 1
				}
				m.lock.Unlock()
			}
		}
	}
}
