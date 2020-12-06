package main

import (
	"fmt"
	"time"
)

func main() {
	numList := [20]int{29, 20, 27, 14, 8, 25, 5, 20, 18, 1, 14, 26, 28, 6, 25, 1, 7, 11, 3, 7}
	fmt.Println("Lista Inicial")
	fmt.Println(numList)
	//InsertionSort(numList)
	quickSort(numList)

}

func InsertionSort(data [20]int) {
	fmt.Println("Insertion sort")
	start := time.Now()
	for i := 1; i < len(data); i++ {
		fmt.Println(data)
		if data[i] < data[i-1] {
			j := i - 1
			temp := data[i]
			for j >= 0 && data[j] > temp {
				data[j+1] = data[j]
				j--
			}
			data[j+1] = temp
		}
	}
	elapsed := time.Since(start)
	fmt.Println("Ordenada: ")
	fmt.Println(data)
	fmt.Println("Tiempo transcurrido: ", elapsed)
}

func quickSort(nums [20]int) {
	start := time.Now()
	recursionSort(nums, 0, len(nums)-1)
	elapsed := time.Since(start)
	fmt.Println("Ordenada: ")
	fmt.Println(nums)
	fmt.Println("Tiempo transcurrido: ", elapsed)
}

func recursionSort(data [20]int, left int, right int) {
	if left < right {
		pivot := partition(data, left, right)
		recursionSort(data, left, pivot-1)
		recursionSort(data, pivot+1, right)
	}
}

func partition(data [20]int, left int, right int) int {
	for left < right {
		for left < right && data[left] <= data[right] {
			right--
		}
		if left < right {
			data[left], data[right] = data[right], data[left]
			left++
		}

		for left < right && data[left] <= data[right] {
			left++
		}
		if left < right {
			data[left], data[right] = data[right], data[left]
			right--
		}
	}
	fmt.Println(data)
	return left
}
