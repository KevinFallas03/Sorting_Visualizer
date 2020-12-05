package algorithms

import (
	//"fmt"
	"time"
)

var closed bool

// QuickSort con pivote aleatorio
func QuickSort(data []int, c chan []int, stopCh chan struct{}, msgCh chan string) {
	t := time.Now()
	closed = false
	m := QuickSortAux(data, c, stopCh, msgCh)
	if !closed {
		c <- m
		close(c)
		msgCh <- "QuickSort: " + time.Since(t).String()
	}

	//fmt.Println("QuickSort: ", time.Since(t))
}

// QuickSortAux ...
func QuickSortAux(data []int, c chan []int, stopCh chan struct{}, msgCh chan string) []int {

	var less []int
	var equals []int
	var greater []int
	var slice6 []int

	if len(data) > 1 {
		pivot := data[0] //rand.Int() % len(data)
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
		slice6 = append(slice5, QuickSortAux(greater, c, stopCh, msgCh)...)
		if !closed {
			select {
			case <-stopCh:
				close(c)
				closed = true
				return data
			case c <- slice6:
			}
		}
		c <- slice6
		return slice6
	} else {
		return slice6
	}

}
