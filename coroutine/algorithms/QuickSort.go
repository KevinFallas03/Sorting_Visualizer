package algorithms

import (
	"fmt"
	"math/rand"
	"time"
)

// QUICKSORT CON PIVOTE ALEATORIO
func QuickSort(data []int, c chan []int ){
	t := time.Now()
    if len(data) < 2 {
        c <- data
        return
    }
      
    left, right := 0, len(data)-1
      
    pivot := rand.Int() % len(data)
      
    data[pivot], data[right] = data[right], data[pivot]

    for i, _ := range data {
        if data[i] < data[right] {
            data[left], data[i] = data[i], data[left]
            left++
        }
    }
      
    data[left], data[right] = data[right], data[left]
      
    QuickSort(data[:left], c)
    QuickSort(data[left+1:], c)
    fmt.Println("QuickSort: ", time.Since(t))

    c <- data
    close(c)
    return

}

func QuickSortAux(){

}