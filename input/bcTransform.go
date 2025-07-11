package input

import (
	"math"
	"robots-go/constants"
	"robots-go/structures"
	"sort"
)

func convertEventToStruct(ev map[string]any) structures.EventInfo {

	var event structures.EventInfo

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

	if i4, ok := ev["i4"].(float64); ok {
		val := int(i4)
		event.I4 = &val
	}

	if i5, ok := ev["i5"].(float64); ok {
		val := int(i5)
		event.I5 = &val
	}

	return event
}

func parsingEventsFun(request map[string]any) []structures.EventInfo {

	evs := request["events"].([]any)

	var events []structures.EventInfo

	for _, item := range evs {
		if ev, ok := item.(map[string]any); ok {
			events = append(events, convertEventToStruct(ev))
		}
	}

	return events
}

func getProviders(set map[string]any) map[string]structures.ProviderInfo {
	result := make(map[string]structures.ProviderInfo)
	pEKs, _ := set["betScannerSourcesSettingsByEventKinds"].([]any)
	for _, item := range pEKs {
		peK, _ := item.(map[string]any)
		eK := peK["eventKindId"].(string)
		source := peK["sourcesSettings"].([]any)
		weigth := 0
		tmp := result[eK]
		for _, provider := range source {
			pr, _ := provider.(map[string]any)
			w := int(pr["weight"].(float64))
			if w > weigth {
				weigth = w
				tmp.ID = pr["providerLayerId"].(string)
			}
		}
		sourcesData := set["betScannerSourcesData"].([]any)
		for _, sources := range sourcesData {
			elem := sources.(map[string]any)
			id := elem["providerLayerId"].(string)
			if tmp.ID == id {
				tmp.MatchClosed = elem["matchClosed"].(bool)
			}
		}
		result[eK] = tmp
	}
	return result
}

func bcExcludeEvents(eventsPtr *[]structures.EventInfo) []structures.EventInfo {
	// Создаем map для исключаемых ID (используем map для быстрого поиска)
	excludeMap := make(map[int64]bool)

	for _, ev := range *eventsPtr {
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
	var result []structures.EventInfo
	for _, ev := range *eventsPtr {
		// Исключаем события типа 1020 и в BCStatistics и те, чьи ID есть в excludeMap
		if _, ok := constants.BcStatistics[ev.Type]; !ok && ev.Type != 1020 && !excludeMap[ev.ID] {
			result = append(result, ev)
		}
	}

	// Сортировка по убыванию ID
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID > result[j].ID
	})

	return result
}

// Исключение из трансляции ошибочных событий 1102 и 1103
func bcExcludeMistakes(eventsPtr *[]structures.EventInfo, settPtr *structures.MatchSettings) []structures.EventInfo {

	var timer int64
	var part structures.Part
	excludeMap := make(map[int]bool)

	for i, event := range *eventsPtr {
		switch event.Type {
		case 1103:
			{
				eventsRemaining := (*eventsPtr)[i+1:]
				_, timer = partTimer(&eventsRemaining, bcTimeToTimestamp(event.RegTime), settPtr)
				if timer < settPtr.MatchDuration-120 {
					excludeMap[i] = true
					settPtr.BlockAll = true
					for _, ev := range (*eventsPtr)[:i] {
						if _, ok := constants.Unblocks[ev.Type]; ok {
							settPtr.BlockAll = false
							break
						}
					}
				}
			}
		case 1102:
			{
				eventsRemaining := (*eventsPtr)[i+1:]
				part, timer = partTimer(&eventsRemaining, bcTimeToTimestamp(event.RegTime), settPtr)
				if timer < settPtr.PartTimes[part.Nmb].End-120 {
					excludeMap[i] = true
					settPtr.BlockAll = true
					for _, ev := range (*eventsPtr)[:i] {
						if _, ok := constants.Unblocks[ev.Type]; ok {
							settPtr.BlockAll = false
							break
						}
					}
				}
			}
		}
	}

	var result []structures.EventInfo
	for i, v := range *eventsPtr {
		if !excludeMap[i] {
			result = append(result, v)
		}
	}
	return result
}

func moveUp1102(eventsPtr *[]structures.EventInfo) []structures.EventInfo {

	events := *eventsPtr
	// Находим индекс первого события таймера первого тайма
	var timerIdx int = len(events)
	for i, ev := range events {
		if constants.BcTimer[ev.Type] && *ev.I2 == 1 {
			timerIdx = i
			break
		}
	}

	// Ищем событие 1102 первого тайма после найденного индекса
	for i := timerIdx; i < len(events); i++ {
		if events[i].Type == 1102 && *events[i].I1 == 1 {
			// Удаляем и вставляем на новую позицию
			ev := events[i]
			events = append(events[:i], events[i+1:]...)
			events = append(events[:timerIdx], append([]structures.EventInfo{ev}, events[timerIdx:]...)...)
			break
		}
	}

	return events
}

func bcReverse(eventsPtr *[]structures.EventInfo) []structures.EventInfo {
	for _, event := range *eventsPtr {
		if value, ok := constants.EventsWithTeam[event.Type]; ok {
			switch value {
			case "i1":
				{
					*event.I1 = 3 - *event.I1
				}
			case "i2":
				{
					*event.I2 = 3 - *event.I2
				}
			case "i3":
				{
					*event.I3 = 3 - *event.I3
				}
			case "i4":
				{
					*event.I4 = 3 - *event.I4
				}
			case "i5":
				{
					*event.I5 = 3 - *event.I5
				}
			}
		}
	}
	return *eventsPtr
}

func bcTransformation(request map[string]any, settPtr *structures.MatchSettings) []structures.EventInfo {
	events := parsingEventsFun(request)          // parse events
	events = bcExcludeEvents(&events)            // exclude 1020 and statistics
	events = bcExcludeMistakes(&events, settPtr) // exclude ends 1102 and 1103 with mistakes
	events = moveUp1102(&events)                 // move up 1102 if mistake
	if settPtr.SportscastReverseTeams {          //reverse broadcast if SportscastReverseTeams == true
		events = bcReverse(&events)
	}

	return events
}
