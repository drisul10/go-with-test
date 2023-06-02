package arrays_slices

func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

func SumAllTails(numbersToSum ...[]int) (sums []int) {
	for _, numbers := range numbersToSum {
		tailNumbers := numbers[1:]
		sums = append(sums, Sum(tailNumbers))
	}

	return sums
}
