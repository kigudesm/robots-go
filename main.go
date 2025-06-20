package main

import (

    // "fmt"        
    "robots-go/utils"
)

func main() {
    // data := utils.Read()
    utils.Read()

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
