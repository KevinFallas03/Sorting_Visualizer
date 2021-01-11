package algorithms

//SelectionSort ...
func SelectionSort(data []int, c chan []int, stopCh chan struct{}, msgCh chan string) {
	swaps := 0
	comparations := 0
	loops := 0
	//t := time.Now()
	length := len(data)

	for i := 0; i < length; i++ {
		loops++
		maxIndex := 0
		for j := 1; j < length-i; j++ {
			loops++
			comparations++
			if data[j] > data[maxIndex] {
				maxIndex = j
			}
		}
		swaps++
		data[length-i-1], data[maxIndex] = data[maxIndex], data[length-i-1]
		select {
		case <-stopCh:
			close(c)
			return
		case c <- data:
		}
	}

	// hi, mi, si := t.Clock()
	// hf, mf, sf := time.Now().Clock()
	// msgCh <- "\nSelectionSort:" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparations) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)
	close(c)
}
