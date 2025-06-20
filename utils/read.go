package utils

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
)

func Read() {
    // Открываем JSON файл
    file, err := os.Open("Input_samples/regular_match.json")
    if err != nil {
        log.Fatalf("Ошибка открытия файла: %v", err)
    }
    defer file.Close()

    // Читаем содержимое файла
    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatalf("Ошибка чтения файла: %v", err)
    }

    // Объявляем переменную для неизвестной структуры
    var data interface{}

    // Десериализация в интерфейс
    err = json.Unmarshal(bytes, &data)
    if err != nil {
        log.Fatalf("Ошибка разбора JSON: %v", err)
    }

    // Теперь, data содержит разобранную структуру
    switch v := data.(type) {
    case map[string]interface{}:
        fmt.Println("Объект JSON (словарь):")
        fmt.Printf("%#vn", v)
    case []interface{}:
        fmt.Println("Массив JSON:")
        fmt.Printf("%#vn", v)
    default:
        fmt.Printf("Неопознанный тип данных: %Tn", v)
    }
}
