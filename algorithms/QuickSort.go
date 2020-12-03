package algorithms

import (
	"fmt"
	"time"
)
var closed bool

// QUICKSORT CON PIVOTE ALEATORIO
func QuickSort(data []int, c chan []int,stopCh chan struct{}) {
	t := time.Now()
	closed = false
	m := QuickSortAux(data, c, stopCh)
	if !closed{
		c <- m
		close(c)
	}
	fmt.Println("QuickSort: ", time.Since(t))
}

// QuickSortAux ...
func QuickSortAux(data []int, c chan []int, stopCh chan struct{}) []int {

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
		slice5 := append(QuickSortAux(less, c, stopCh), equals...)
		slice6 = append(slice5, QuickSortAux(greater, c, stopCh)...)
		if !closed{
			select {
				case <-stopCh:
					close(c)
					closed = true
					return data
				case c <- slice6:
			}
		}
		return slice6
	} else {
		return slice6
	}

}
