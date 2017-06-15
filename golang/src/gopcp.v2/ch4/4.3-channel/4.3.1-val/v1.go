package main

import (
	"fmt"
	"time"
)

var mapChan = make(chan map[string]int, 1)

func main() {
	syncChan := make(chan struct{}, 2)

	//
	go func() {
		for {
			if elem, ok := <-mapChan; ok {
				elem["count"]++
				//fmt.Println("\tresult:", elem)
			} else {
				break
			}
		}
		//
		fmt.Println("Stopped. [receiver]")
		syncChan <- struct{}{}
	}()

	//
	go func() {
		countMap := make(map[string]int)

		for i := 0; i < 5; i++ {
			mapChan <- countMap
			//fmt.Println("\tcountMap:", countMap)
			time.Sleep(time.Millisecond)
			fmt.Printf("The count map: %v.[sender]\n", countMap)
		}

		//
		close(mapChan)
		syncChan <- struct{}{}

	}()

	//
	<-syncChan
	<-syncChan
}
