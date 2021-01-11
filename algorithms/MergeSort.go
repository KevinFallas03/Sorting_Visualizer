package algorithms

//MergeSort initialize the sorting
func MergeSort(data []int, c chan []int, stopCh chan struct{}, msgCh chan string) {
	//t := time.Now()
	swaps := 0
	comparations := 0
	loops := 0

	closed = false
	m := MergeSortAux(data, c, stopCh, msgCh, &swaps, &comparations, &loops)
	if !closed {
		c <- m
		// hi, mi, si := t.Clock()
		// hf, mf, sf := time.Now().Clock()
		// msgCh <- "\nMergeSort:" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparations) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)
		close(c)
	}
}

//MergeSortAux do the recursive part of the sort
func MergeSortAux(data []int, c chan []int, stopCh chan struct{}, msgCh chan string, swaps, comparations, loops *int) []int {
	*loops++

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
		*comparations++
		if i < middle {
			left[i] = data[i]
		} else {
			right[i-middle] = data[i]
		}
		*swaps++
	}
	result := Merge(MergeSortAux(left, c, stopCh, msgCh, swaps, comparations, loops), MergeSortAux(right, c, stopCh, msgCh, swaps, comparations, loops), swaps, comparations, loops)
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
func Merge(left, right []int, swaps, comparations, loops *int) (result []int) {
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
