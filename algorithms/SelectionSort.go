package algorithms

import (
	"fmt"
	"time"
)

//SelectionSort ...
func SelectionSort(data []int, c chan []int) {
	t := time.Now()
	length := len(data)
	for i := 0; i < length; i++ {
		maxIndex := 0
		for j := 1; j < length-i; j++ {
			if data[j] > data[maxIndex] {
				maxIndex = j
			}
		}
		data[length-i-1], data[maxIndex] = data[maxIndex], data[length-i-1]
		c <- data
	}
	fmt.Println("SelectionSort: ", time.Since(t))
	close(c)
}
