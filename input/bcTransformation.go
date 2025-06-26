package input

import (
	"math"
)

// // Исключение событий с типом 1020 и тех на которые они ссылаются
// func bcExclude1020(events []Event) []Event {
// 	var result []Event
// 	for _, ev := range events {
// 		if condition(v) {
// 			result = append(result, ev)
// 		}
// 	}
// 	return result
// }

func bcExclude1020(events []Event) []Event {
	// Создаем map для исключаемых ID (используем map для быстрого поиска)
	excludeMap := make(map[int64]bool)

	for _, ev := range events {
		if ev.Type == 1020 {
			// Проверяем наличие optional-полей перед вычислением
			if ev.I1 != nil && ev.I2 != nil {
				// Вычисляем значение для исключения: i1 * 2^32 + i2
				excludedID := int64(*ev.I1)*int64(math.Pow(2, 32)) + int64(*ev.I2)
				excludeMap[excludedID] = true
			}
		}
	}

	// Фильтруем исходный слайс
	var result []Event
	for _, ev := range events {
		// Исключаем события типа 1020 и те, чьи ID есть в excludeMap
		if ev.Type != 1020 && !excludeMap[ev.ID] {
			result = append(result, ev)
		}
	}

	return result
}
