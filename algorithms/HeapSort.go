package algorithms

//HeapSort ...
func HeapSort(data []int, c chan [][]int, stopCh chan struct{}, msgCh chan string) {
	//t := time.Now()
	swaps := 0
	comparations := 0
	loops := 0

	heapify(data, &swaps, &comparations, &loops, c)
	for i := len(data) - 1; i > 0; i-- {
		loops++
		swaps++
		data[0], data[i] = data[i], data[0]
		select {
		case <-stopCh:
			close(c)
			return
		case c <- [][]int{[]int{data[0], 0}, []int{data[i], i}}:
		}
		siftDown(data, 0, i, &swaps, &comparations, &loops, c)

	}

	//hi, mi, si := t.Clock()
	//hf, mf, sf := time.Now().Clock()
	//msgCh <- "\nHeapSort:" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparations) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)
	close(c)
}
func heapify(data []int, swaps, comparations, loops *int, c chan [][]int) {
	for i := (len(data) - 1) / 2; i >= 0; i-- {
		siftDown(data, i, len(data), swaps, comparations, loops, c)
	}
}
func siftDown(heap []int, lo, hi int, swaps, comparations, loops *int, c chan [][]int) {
	root := lo
	for {
		child := root*2 + 1
		*comparations++
		if child >= hi {
			break
		}
		*comparations++
		if child+1 < hi && heap[child] < heap[child+1] {
			child++
		}
		*comparations++
		if heap[root] < heap[child] {
			*swaps++
			heap[root], heap[child] = heap[child], heap[root]
			c <- [][]int{[]int{heap[root], root}, []int{heap[child], child}}
			root = child
		} else {
			break
		}

	}
}
