package algorithms

import (
	//"fmt"
	"time"
)

//InsertionSort ...
func InsertionSort(data []int, c chan []int,stopCh chan struct{}, msgCh chan string) {
	t := time.Now()
	for i := 1; i < len(data); i++ {
		if data[i] < data[i-1] {
			j := i - 1
			temp := data[i]
			for j >= 0 && data[j] > temp {
				data[j+1] = data[j]
				j--
			}
			data[j+1] = temp
			select {
				case <-stopCh:
					close(c)
					return
				case c <- data:
			}
		}
	}
	msgCh <- "InsertionSort: "+time.Since(t).String()
	//fmt.Println("InsertionSort: ", time.Since(t))
	close(c)
}
