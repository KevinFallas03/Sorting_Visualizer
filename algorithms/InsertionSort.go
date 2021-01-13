package algorithms

//InsertionSort ...
func InsertionSort(items []int, c chan [][]int, stopCh chan struct{}, msgCh chan string) {
	var n = len(items)
	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			if items[j-1] > items[j] {
				items[j-1], items[j] = items[j], items[j-1]
				select {
				case <-stopCh:
					close(c)
					return
				case c <- [][]int{[]int{items[j-1], j - 1}, []int{items[j], j}}:
				}
			}
			j = j - 1
		}
	}
	close(c)
}
