package input

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func readRequest(path string) map[string]any {
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

func readBcDescription(path string) map[int]map[string]string {
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
	var bcDescription map[int]map[string]string

	// Десериализация в интерфейс
	err = json.Unmarshal(bytes, &bcDescription)
	if err != nil {
		log.Fatalf("Ошибка разбора JSON: %v", err)
		return nil
	}

	return bcDescription
}
