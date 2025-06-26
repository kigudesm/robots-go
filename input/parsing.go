package input

import (
	"time"
)

// Настройки матча
type Settings struct {
	TargetEventKind []string // Целевые рынки
	// MatchType       string // Тип матча (обычный, товарищеский, 2 по 40 и т.п.)
	// HalfDuration    int // Длительность тайма
	// MatchDuration   int // Длительность матча
	// InjuryDefault   [2]int // Компенсированное время по умолчанию
}

// Трансляция
type Event struct {
	ID      int64     // id события
	RegTime time.Time // время
	Type    int       // тип
	// Параметры события (могут отсутствовать)
	I1 *int
	I2 *int
	I3 *int
}

func convertEventToStruct(ev map[string]any) Event {

	var event Event

	// Обязательные поля
	event.ID = int64(ev["id"].(float64))

	event.RegTime, _ = time.Parse("02.01.2006 15:04:05", ev["regtime"].(string))

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

func parsingEventsFun(request map[string]any) []Event {

	evs := request["events"].([]any)

	var events []Event

	for _, item := range evs {
		if ev, ok := item.(map[string]any); ok {
			events = append(events, convertEventToStruct(ev))
		}
	}

	return events

}

func parsingSettingsFun(request map[string]any) Settings {

	set := request["settings"].(map[string]any)

	tek := set["targetEventKind"].([]any)

	var targetEventKind []string

	for _, item := range tek {
		if ek, ok := item.(string); ok {
			targetEventKind = append(targetEventKind, ek)
		}
	}

	var settings Settings

	settings.TargetEventKind = targetEventKind

	return settings

}
