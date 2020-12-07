package algorithms

import (
	"time"
)

var closed bool

// QuickSort ...
func QuickSort(data []int, c chan []int, stopCh chan struct{}, msgCh chan string) {
	t := time.Now()
	closed = false
	m := QuickSortAux(data, c, stopCh, msgCh)
	if !closed {
		c <- m
		close(c)
		msgCh <- "QuickSort: " + time.Since(t).String()
	}
}

// QuickSortAux ...
func QuickSortAux(data []int, c chan []int, stopCh chan struct{}, msgCh chan string) []int {

	var less []int
	var equals []int
	var greater []int
	var result []int

	if len(data) > 1 {
		pivot := data[0]
		for i := range data {
			if data[i] < pivot {
				less = append(less, data[i])
			} else if data[i] == pivot {
				equals = append(equals, data[i])
			} else if data[i] > pivot {
				greater = append(greater, data[i])
			}
		}
		slice5 := append(QuickSortAux(less, c, stopCh, msgCh), equals...)
		result = append(slice5, QuickSortAux(greater, c, stopCh, msgCh)...)
		if !closed {
			select {
			case <-stopCh:
				close(c)
				closed = true
				return data
			case c <- result:
			}
		}
		c <- result
		return result
	} else {
		return result
	}

}
