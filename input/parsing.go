package input

import (
	"robots-go/gametypes"
	"strconv"
)

// Трансляция
type EventStruct struct {
	ID      int64  // id события
	RegTime string // время
	Type    int    // тип
	// Параметры события (могут отсутствовать)
	I1 *int
	I2 *int
	I3 *int
}

func convertEventToStruct(ev map[string]any) EventStruct {

	var event EventStruct

	// Обязательные поля
	event.ID = int64(ev["id"].(float64))
	event.RegTime = ev["regtime"].(string)
	event.Type = int(ev["type"].(float64))

	// Опциональные поля
	if i1, ok := ev["i1"].(float64); ok {
		val := int(i1)
		event.I1 = &val
	}

	if i2, ok := ev["i2"].(float64); ok {
		val := int(i2)
		event.I2 = &val
	}

	if i3, ok := ev["i3"].(float64); ok {
		val := int(i3)
		event.I3 = &val
	}

	return event
}

func parsingEventsFun(request map[string]any) []EventStruct {

	evs := request["events"].([]any)

	var events []EventStruct

	for _, item := range evs {
		if ev, ok := item.(map[string]any); ok {
			events = append(events, convertEventToStruct(ev))
		}
	}

	return events

}

// Настройки матча
type SettingsStruct struct {
	MatchType          string   // Тип матча (обычный, товарищеский, 2 по 40 и т.п.)
	EventGameTypeIdent string   // Тип матча (с добаленным временем, серией пенальти и т.п.)
	HalfDuration       int64    // Длительность тайма
	MatchDuration      int64    // Длительность матча
	InjuryDefault      [2]int   // Компенсированное время по умолчанию
	ServerTime         int64    // Время сервера
	TargetEventKind    []string // Целевые рынки
	Block              []string // Заблокированные eventKind
}

func parsingSettingsFun(request map[string]any) SettingsStruct {

	set := request["settings"].(map[string]any)

	// Парсим targetEventKind
	tek := set["targetEventKind"].([]any)
	var targetEventKind []string
	for _, item := range tek {
		if ek, ok := item.(string); ok {
			targetEventKind = append(targetEventKind, ek)
		}
	}

	// Парсим ident
	ident, _ := set["sportRules"].(map[string]any)["object"].(map[string]any)["ident"].(string)
	matchType := gametypes.MatchTypes[ident]

	// Парсим EventGameTypeIdent
	var eGTI string
	if set["eventGameTypeIdent"] == nil {
		eGTI = "regular"
	} else {
		eGTI = set["eventGameTypeIdent"].(string)
	}

	//Парсим serverTime
	num, _ := strconv.ParseInt(set["serverTime"].(string), 10, 64)

	// Заполняем Settings
	var settings SettingsStruct
	settings.TargetEventKind = targetEventKind
	settings.MatchType = ident
	settings.EventGameTypeIdent = eGTI
	settings.HalfDuration = int64(matchType.HalfDuration)
	settings.MatchDuration = 2 * int64(matchType.HalfDuration)
	settings.InjuryDefault[0] = matchType.TFirstHalf - matchType.HalfDuration
	settings.InjuryDefault[1] = matchType.TSecondHalf - matchType.HalfDuration
	settings.ServerTime = num / 1000

	return settings
}
