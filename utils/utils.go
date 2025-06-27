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
