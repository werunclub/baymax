package cheatcheck

import (
	"math"
)

const (
	N   = 10
	AMP = 28
)

var (
	// dimenTh []int
	valueTH []float64
)

func init() {
	// dimenTh = []int{1, 3, 6, 9, 12}
	valueTH = []float64{math.Inf(0), 3000, 2600, 2000, 1600, 800, 28}
}

func Max(slice []int) int {
	max := slice[0]
	for index := 1; index < len(slice); index++ {
		if slice[index] > max {
			max = slice[index]
		}
	}
	return max
}

func Min(slice []int) int {
	min := slice[0]
	for index := 1; index < len(slice); index++ {
		if slice[index] < min {
			min = slice[index]
		}
	}
	return min
}

func dynamicCluster(data []int) int {
	cluster := 1
	for i := 1; i < len(data); i++ {
		isSingleCluster := true
		for j := 0; j < i; j++ {
			if math.Abs(float64(data[i]-data[j])) <= N {
				isSingleCluster = false
			}
		}

		if isSingleCluster {
			cluster++
		}
	}
	return cluster
}

func ampDiscriminate(data []int) bool {
	if Max(data)-Min(data) < AMP {
		return true
	}
	return false
}

func cheatDiscriminate(data []int) bool {
	if dynamicCluster(data) == 1 && ampDiscriminate(data) {
		return true
	}
	return false
}

func valueSameDiscriminate(data []int) bool {
	for i := 1; i < len(data); i++ {
		if data[i] != data[0] {
			return false
		}
	}

	return true
}

// index:1,2,3,4,5,6
func valueRangeDiscriminate(data []int, index int) bool {
	for i := 0; i < len(data); i++ {
		if float64(data[i]) >= valueTH[index-1] || float64(data[i]) < valueTH[index] {
			return false
		}
	}

	return true
}

func cheatDiscriminateN(data []int, n int) bool {
	if n == 1 && valueRangeDiscriminate(data, n) {
		return true
	} else if n == 3 {
		if valueRangeDiscriminate(data, n-1) {
			return true
		} else if valueRangeDiscriminate(data, n) && !valueSameDiscriminate(data) && cheatDiscriminate(data) {
			return true
		}
	} else {
		isRange := false
		if n == 6 && valueRangeDiscriminate(data, 4) {
			isRange = true
		} else if n == 9 && valueRangeDiscriminate(data, 5) {
			isRange = true
		} else if n == 12 && valueRangeDiscriminate(data, 6) {
			isRange = true
		}

		if isRange && !valueSameDiscriminate(data) && cheatDiscriminate(data) {
			return true
		}
	}

	return false
}

// StepsCheatCheck 计步作弊检查
// 入参数：data []int 为每10分钟的步数数组
// 返回：有作弊嫌疑的步数索引段数组，如: [][]int{[]int{1, 5}, []int{10, 100}}
// 为 1到5和10到100有作弊嫌疑，注意包含 5 和 100
func StepsCheatCheck(data []int) ([][]int, error) {
	cheatIndexes := make([][]int, 0, 1)

	for i := 0; i < len(data); i++ {
		if i > 10 {
			seg1 := data[i : i+1]
			seg3 := data[i-2 : i+1]
			seg6 := data[i-5 : i+1]
			seg9 := data[i-8 : i+1]
			seg12 := data[i-11 : i+1]

			if cheatDiscriminateN(seg12, 12) {
				cheatIndexes = append(cheatIndexes, []int{i - 11, i})
			} else if cheatDiscriminateN(seg9, 9) {
				cheatIndexes = append(cheatIndexes, []int{i - 8, i})
			} else if cheatDiscriminateN(seg6, 6) {
				cheatIndexes = append(cheatIndexes, []int{i - 5, i})
			} else if cheatDiscriminateN(seg3, 3) {
				cheatIndexes = append(cheatIndexes, []int{i - 2, i})
			} else if cheatDiscriminateN(seg1, 1) {
				cheatIndexes = append(cheatIndexes, []int{i, i})
			}
		} else if i > 7 {
			seg1 := data[i : i+1]
			seg3 := data[i-2 : i+1]
			seg6 := data[i-5 : i+1]
			seg9 := data[i-8 : i+1]

			if cheatDiscriminateN(seg9, 9) {
				cheatIndexes = append(cheatIndexes, []int{i - 8, i})
			} else if cheatDiscriminateN(seg6, 6) {
				cheatIndexes = append(cheatIndexes, []int{i - 5, i})
			} else if cheatDiscriminateN(seg3, 3) {
				cheatIndexes = append(cheatIndexes, []int{i - 2, i})
			} else if cheatDiscriminateN(seg1, 1) {
				cheatIndexes = append(cheatIndexes, []int{i, i})
			}
		} else if i > 4 {
			seg1 := data[i : i+1]
			seg3 := data[i-2 : i+1]
			seg6 := data[i-5 : i+1]

			if cheatDiscriminateN(seg6, 6) {
				cheatIndexes = append(cheatIndexes, []int{i - 5, i})
			} else if cheatDiscriminateN(seg3, 3) {
				cheatIndexes = append(cheatIndexes, []int{i - 2, i})
			} else if cheatDiscriminateN(seg1, 1) {
				cheatIndexes = append(cheatIndexes, []int{i, i})
			}
		} else if i > 1 {
			seg1 := data[i : i+1]
			seg3 := data[i-2 : i+1]

			if cheatDiscriminateN(seg3, 3) {
				cheatIndexes = append(cheatIndexes, []int{i - 2, i})
			} else if cheatDiscriminateN(seg1, 1) {
				cheatIndexes = append(cheatIndexes, []int{i, i})
			}
		} else {
			seg1 := data[i : i+1]
			if cheatDiscriminateN(seg1, 1) {
				cheatIndexes = append(cheatIndexes, []int{i, i})
			}
		}
	}

	indexes := make([][]int, 0, 1)
	for _, idx := range cheatIndexes {
		if len(indexes) == 0 {
			indexes = append(indexes, idx)
		} else if idx[1] < indexes[len(indexes)-1][1] {

		} else if idx[0] <= indexes[len(indexes)-1][1]+1 {
			indexes[len(indexes)-1][1] = idx[1]
		} else {
			indexes = append(indexes, idx)
		}
	}

	return indexes, nil
}
