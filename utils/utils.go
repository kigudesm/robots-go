package utils

// Создает множество map[T]bool из среза []T
func SliceToSet[T comparable](slice []T) map[T]bool {
	set := make(map[T]bool, len(slice))
	for _, item := range slice {
		set[item] = true
	}
	return set
}

// Объединяет два множества типа map[T]struct{}
func UnionSets[T comparable](a, b map[T]bool) map[T]bool {
	result := make(map[T]bool)
	for k := range a {
		result[k] = true
	}
	for k := range b {
		result[k] = true
	}
	return result
}

// // Ordered ограничивающий тип для чисел и строк
// type Ordered interface {
// 	~int | ~int32 | ~int64 | ~float32 | ~float64 | ~string
// }

// // Max возвращает максимальный элемент в слайсе
// func Max[T Ordered](slice []T) T {
// 	if len(slice) == 0 {
// 		var zero T
// 		return zero // возвращает нулевое значение, если слайс пустой
// 	}

// 	maxVal := slice[0]
// 	for _, v := range slice[1:] {
// 		if v > maxVal {
// 			maxVal = v
// 		}
// 	}
// 	return maxVal
// }
