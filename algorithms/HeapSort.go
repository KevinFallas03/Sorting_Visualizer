package algorithms

import (
	//"fmt"
	"time"
)

//HeapSort ...
func HeapSort(data []int, c chan []int,stopCh chan struct{}, msgCh chan string) {
	t := time.Now()
	heapify(data)
	for i := len(data) - 1; i > 0; i-- {
		data[0], data[i] = data[i], data[0]
		siftDown(data, 0, i)
		select {
			case <-stopCh:
				close(c)
				return
			case c <- data:
		}
	}
	msgCh <- "HeapSort: "+time.Since(t).String()
	//fmt.Println("HeapSort: ", time.Since(t))
	close(c)
}
func heapify(data []int) {
	for i := (len(data) - 1) / 2; i >= 0; i-- {
		siftDown(data, i, len(data))
	}
}
func siftDown(heap []int, lo, hi int) {
	root := lo
	for {
		child := root*2 + 1
		if child >= hi {
			break
		}
		if child+1 < hi && heap[child] < heap[child+1] {
			child++
		}
		if heap[root] < heap[child] {
			heap[root], heap[child] = heap[child], heap[root]
			root = child
		} else {
			break
		}

	}
}
