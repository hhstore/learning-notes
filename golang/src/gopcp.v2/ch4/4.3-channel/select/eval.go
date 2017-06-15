package main

import "fmt"

var intChan1, intChan2 chan int
var chans = []chan int{intChan1, intChan2}
var numbers = []int{1, 2, 3, 4, 5}

func main() {

	fmt.Println("len(intChan1):", len(intChan1))

	select {
	case getChan(0) <- getNumber(0):
		fmt.Println("the 1th case is selected.")
	case getChan(1) <- getNumber(1):
		fmt.Println("the 2nd case is selected.")
	default:
		fmt.Println("Default!")
	}
}

func getNumber(i int) int {
	fmt.Printf("numbers[%d]\n", i)
	return numbers[i]
}

func getChan(i int) chan int {
	fmt.Printf("channels[%d]\n", i)
	return chans[i]

}
