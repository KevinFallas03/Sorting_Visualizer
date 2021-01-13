package algorithms

//SelectionSort ...
// func SelectionSort(data []int, c chan [][]int, stopCh chan struct{}, msgCh chan string) {
// 	swaps := 0
// 	comparations := 0
// 	loops := 0
// 	//t := time.Now()

// 	for i := 0; i < len(data); i++ {
// 		loops++
// 		maxIndex := 0
// 		for j := 1; j < len(data)-i; j++ {
// 			loops++
// 			comparations++
// 			if data[j] > data[maxIndex] {
// 				maxIndex = j
// 				data[len(data)-i-1], data[maxIndex] = data[maxIndex], data[len(data)-i-1]
// 				select {
// 				case <-stopCh:
// 					close(c)
// 					return
// 				// case c <- data:
// 				case c <- [][]int{[]int{data[len(data)-i-1], len(data) - i - 1}, []int{data[maxIndex], maxIndex}}:
// 				}

// 			}
// 		}
// 		swaps++

// 	}

// 	// hi, mi, si := t.Clock()
// 	// hf, mf, sf := time.Now().Clock()
// 	// msgCh <- "\nSelectionSort:" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparations) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)
// 	close(c)
// }

// SelectionSort ...
func SelectionSort(arr []int, c chan [][]int, stopCh chan struct{}, msgCh chan string) {
	len := len(arr)
	for i := 0; i < len-1; i++ {
		minIndex := i
		for j := i + 1; j < len; j++ {
			if arr[j] <= arr[minIndex] {
				select {
				case <-stopCh:
					close(c)
					return
				// case c <- data:
				case c <- [][]int{[]int{arr[j], j}, []int{arr[minIndex], minIndex}}:
				}
				arr[j], arr[minIndex] = arr[minIndex], arr[j]

			}
		}
	}
	close(c)
}
