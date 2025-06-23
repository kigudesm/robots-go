package main

import (

    // "fmt"        
    "github.com/kigudesm/robots-go/data_preparation/utils"
)

func main() {
    path := "Input_samples/regular_match.json"
    // data := utils.Read(path)
    utils.readRequest(path)

    // // Простое извлечение
    // value, exists := data["events"]
    // if exists {
    //     fmt.Println("Найдено:", value)
    // } else {
    //     fmt.Println("Ключ не найден")
    // }

    // Теперь, data содержит разобранную структуру
    // switch v := data.(type) {
    // case map[string]interface{}:
    //     fmt.Println("Объект JSON (словарь):")
    //     fmt.Printf("%#vn", v)
    // case []interface{}:
    //     fmt.Println("Массив JSON:")
    //     fmt.Printf("%#vn", v)
    // default:
    //     fmt.Printf("Неопознанный тип данных: %Tn", v)
    // }
}
