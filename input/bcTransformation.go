package input

import (
	"math"
	"robots-go/constants"
	"robots-go/structures"
	"sort"
)

func bcExcludeEvents(events []structures.EventStruct) []structures.EventStruct {
	// Создаем map для исключаемых ID (используем map для быстрого поиска)
	excludeMap := make(map[int64]bool)

	for _, ev := range events {
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
	var result []structures.EventStruct
	for _, ev := range events {
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
func bcExcludeMistakes(events []structures.EventStruct, settings structures.SettingsStruct) (
	[]structures.EventStruct, structures.SettingsStruct) {

	var timer int64
	var part structures.PartStruct
	excludeMap := make(map[int]bool)

	for i, event := range events {
		switch event.Type {
		case 1103:
			{
				_, timer = partTimer(events[i+1:], bcTimeToTimestamp(event.RegTime), settings)
				if timer < settings.MatchDuration-120 {
					excludeMap[i] = true
					flag := true
					for _, ev := range events[:i] {
						if _, ok := constants.Unblocks[ev.Type]; ok {
							flag = false
							break
						}
					}
					if flag {
						settings.Block = settings.TargetEventKind
					}
				}
			}
		case 1102:
			{
				part, timer = partTimer(events[i+1:], bcTimeToTimestamp(event.RegTime), settings)
				if timer < settings.PartTimes[part.Nmb].End-120 {
					excludeMap[i] = true
					flag := true
					for _, ev := range events[:i] {
						if _, ok := constants.Unblocks[ev.Type]; ok {
							flag = false
							break
						}
					}
					if flag {
						settings.Block = settings.TargetEventKind
					}
				}
			}
		}
	}

	var result []structures.EventStruct
	for i, v := range events {
		if !excludeMap[i] {
			result = append(result, v)
		}
	}

	return result, settings
}

func moveUp1102(events []structures.EventStruct) []structures.EventStruct {

	// Находим индекс первого события таймера первого тайма
	var timerIdx int = len(events)
	for i, ev := range events {
		if _, ok := constants.BcTimer[ev.Type]; ok && *ev.I2 == 1 {
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
			events = append(events[:timerIdx], append([]structures.EventStruct{ev}, events[timerIdx:]...)...)
			break
		}
	}

	return events
}

func bcReverse(events []structures.EventStruct) []structures.EventStruct {
	for _, event := range events {
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
	return events
}

func bcGetInjury(events []structures.EventStruct, settings structures.SettingsStruct) structures.SettingsStruct {
	settings.Injury = settings.InjuryDefault
	for i := range 2 {
		for _, event := range events {
			if event.Type == 1104 && *event.I1 != 0 && *event.I2 == i+1 {
				settings.Injury[i] = *event.I1 * 60
				break
			}
		}
	}
	return settings
}

func bcTransformation(request map[string]any, settings structures.SettingsStruct) (
	[]structures.EventStruct, structures.SettingsStruct) {
	events := parsingEventsFun(request)                    // parse events
	events = bcExcludeEvents(events)                       // exclude 1020 and statistics
	events, settings = bcExcludeMistakes(events, settings) // exclude ends 1102 and 1103 with mistakes
	events = moveUp1102(events)                            // move up 1102 if mistake
	settings = bcGetInjury(events, settings)               // get injury
	if settings.SportscastReverseTeams {                   //reverse broadcast if SportscastReverseTeams == true
		events = bcReverse(events)
	}

	return events, settings
}
