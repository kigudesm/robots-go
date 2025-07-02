package utils

func IsinSet(item int, set map[int]struct{}) bool {
	_, exists := set[item]
	return exists
}

// Ordered ограничивающий тип для чисел и строк
type Ordered interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64 | ~string
}

// Max возвращает максимальный элемент в слайсе
func Max[T Ordered](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero // возвращает нулевое значение, если слайс пустой
	}

	maxVal := slice[0]
	for _, v := range slice[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}
