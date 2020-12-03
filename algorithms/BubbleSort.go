package algorithms

import (
	"fmt"
	"time"
)

//BubbleSort ...
func BubbleSort(data []int, c chan []int, stopCh chan struct{}) {
	t := time.Now()
	for i := 0; i < len(data); i++ {
		for j := 1; j < len(data)-i; j++ {
			if data[j] < data[j-1] {
				data[j], data[j-1] = data[j-1], data[j]
				select {
					case <-stopCh:
						close(c)
						return
					case c <- data:
				}

			}
		}
	}
	fmt.Println("BubbleSort: ", time.Since(t))
	close(c)
}
