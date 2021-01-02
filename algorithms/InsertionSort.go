package algorithms

import (
	"strconv"
	"time"
)

//InsertionSort ...
func InsertionSort(data []int, c chan []int, stopCh chan struct{}, msgCh chan string) {
	swaps := 0
	comparations := 0
	loops := 0
	t := time.Now()

	for i := 1; i < len(data); i++ {
		loops++
		comparations++
		if data[i] < data[i-1] {
			j := i - 1
			temp := data[i]
			for j >= 0 && data[j] > temp {
				loops++
				swaps++
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

	hi, mi, si := t.Clock()
	hf, mf, sf := time.Now().Clock()
	msgCh <- "InsertionSort:" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparations) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)
	close(c)
}
