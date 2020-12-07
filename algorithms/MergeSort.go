package algorithms

import (
	"time"
)

//MergeSort initialize the sorting
func MergeSort(data []int, c chan []int, stopCh chan struct{}, msgCh chan string) {
	t := time.Now()
	closed = false
	m := MergeSortAux(data, c, stopCh, msgCh)
	if !closed {
		c <- m
		close(c)
		msgCh <- "MergeSort: " + time.Since(t).String()
	}
}

//MergeSortAux do the recursive part of the sort
func MergeSortAux(data []int, c chan []int, stopCh chan struct{}, msgCh chan string) []int {
	var num = len(data)

	if num == 1 {
		return data
	}

	middle := int(num / 2)
	var (
		left  = make([]int, middle)
		right = make([]int, num-middle)
	)
	for i := 0; i < num; i++ {
		if i < middle {
			left[i] = data[i]
		} else {
			right[i-middle] = data[i]
		}
	}
	result := Merge(MergeSortAux(left, c, stopCh, msgCh), MergeSortAux(right, c, stopCh, msgCh))
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
}

//Merge it merge the left and right list
func Merge(left, right []int) (result []int) {
	result = make([]int, len(left)+len(right))

	i := 0
	for len(left) > 0 && len(right) > 0 {
		if left[0] < right[0] {
			result[i] = left[0]
			left = left[1:]
		} else {
			result[i] = right[0]
			right = right[1:]
		}
		i++
	}

	for j := 0; j < len(left); j++ {
		result[i] = left[j]
		i++
	}
	for j := 0; j < len(right); j++ {
		result[i] = right[j]
		i++
	}

	return
}
