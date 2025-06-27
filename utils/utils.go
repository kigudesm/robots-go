package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func ReadRequest(path string) map[string]any {
	// Открываем JSON файл
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	// Читаем содержимое файла
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
		return nil
	}

	// Объявляем переменную для неизвестной структуры
	var request map[string]any

	// Десериализация в интерфейс
	err = json.Unmarshal(bytes, &request)
	if err != nil {
		log.Fatalf("Ошибка разбора JSON: %v", err)
		return nil
	}

	return request
}

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
