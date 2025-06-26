package input

import (
	"time"
)

type Event struct {
	ID      int64     `json:"id"`
	RegTime time.Time `json:"regtime"` // или используйте time.Time с кастомным парсингом
	Type    int       `json:"type"`
	I1      *int      `json:"i1,omitempty"` // указатель для optional-полей
	I2      *int      `json:"i2,omitempty"`
	I3      *int      `json:"i3,omitempty"`
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
