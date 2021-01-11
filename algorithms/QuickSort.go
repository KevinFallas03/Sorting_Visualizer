package algorithms

import (
	"math/rand"
)

var closed bool

// QuickSort ...
func QuickSort(data []int, c chan []int, stopCh chan struct{}, msgCh chan string) {
	//t := time.Now()
	swaps := 0
	comparations := 0
	loops := 0

	closed = false
	m := QuickSortAux(data, c, stopCh, msgCh, &swaps, &comparations, &loops)
	if !closed {
		c <- m
		// hi, mi, si := t.Clock()
		// hf, mf, sf := time.Now().Clock()
		// msgCh <- "\nQuickSort" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparations) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)
		close(c)
	}
}

// QuickSortAux ...
// func QuickSortAux(data []int, c chan []int, stopCh chan struct{}, msgCh chan string, swaps, comparations, loops *int) []int {

// 	var less []int
// 	var equals []int
// 	var greater []int
// 	var result []int

// 	if len(data) > 1 {
// 		pivot := data[0]
// 		for i := range data {
// 			if data[i] < pivot {
// 				*comparations++
// 				less = append(less, data[i])
// 			} else if data[i] == pivot {
// 				equals = append(equals, data[i])
// 			} else if data[i] > pivot {
// 				greater = append(greater, data[i])
// 			}
// 		}
// 		slice5 := append(QuickSortAux(less, c, stopCh, msgCh, swaps, comparations, loops), equals...)
//  	result = append(slice5, QuickSortAux(greater, c, stopCh, msgCh, swaps, comparations, loops)...)
// 		if !closed {
// 			select {
// 			case <-stopCh:
// 				close(c)
// 				closed = true
// 				return data
// 			case c <- result:
// 			}
// 		}
// 		c <- result
// 		return result
// 	} else {
// 		return result
// 	}

// }

// QuickSortAux ...
func QuickSortAux(a []int, c chan []int, stopCh chan struct{}, msgCh chan string, swaps, comparations, loops *int) []int {
	*loops++
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	pivot := rand.Int() % len(a)

	a[pivot], a[right] = a[right], a[pivot]

	for i := range a {
		*comparations++
		if a[i] < a[right] {
			*swaps++
			a[left], a[i] = a[i], a[left]
			left++
		}
	}

	*swaps++
	a[left], a[right] = a[right], a[left]

	QuickSortAux(a[:left], c, stopCh, msgCh, swaps, comparations, loops)
	QuickSortAux(a[left+1:], c, stopCh, msgCh, swaps, comparations, loops)
	if !closed {
		select {
		case <-stopCh:
			close(c)
			closed = true
			return a
		case c <- a:
		}
	}
	c <- a
	return a
}
