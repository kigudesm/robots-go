package input

import (
	"math"
	"robots-go/utils"
	"sort"
)

var BCStatistics = map[int]struct{}{
	1006: {}, 1052: {}, 1145: {}, 1171: {}, 1172: {}, 1173: {}, 1621: {}, 1200: {}, 1201: {}, 1202: {},
	1203: {}, 1204: {}, 1205: {}, 1206: {}, 1207: {}, 1208: {}, 1209: {}, 1210: {}, 1211: {}, 1212: {},
	1213: {}, 1214: {}, 1215: {}, 1216: {}, 1217: {}, 1218: {}, 1219: {}, 1220: {}, 1221: {}, 1222: {},
	1224: {}, 1225: {}, 1226: {}, 1227: {}, 1233: {}, 1250: {}, 1251: {}, 1411: {}, 1418: {}, 1538: {},
	1550: {}, 1551: {}, 1552: {}, 1554: {}, 1560: {}, 1561: {}, 1564: {}, 1565: {}, 1576: {}, 1586: {},
	1597: {}, 1620: {}, 1622: {}, 1720: {}, 1721: {}, 1857: {}, 1853: {}, 1854: {}, 1855: {}, 1858: {},
	2001: {}, 2002: {}, 2003: {}, 2011: {}, 2021: {}, 2022: {}, 2031: {}, 2041: {}, 2061: {}, 2062: {},
	2063: {}, 2064: {}, 2065: {}, 2066: {}, 2067: {}, 2681: {}, 2682: {}, 3020: {}, 3021: {}, 3100: {},
	3101: {}, 3102: {}, 3103: {}, 3105: {}, 3201: {}, 1813: {}, 1859: {}, 1863: {}, 2070: {}, 1477: {},
	1867: {}, 1242: {}, 1866: {}, 1161: {}, 3104: {},
}

func bcExcludeEvents(events []EventStruct) []EventStruct {
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
	var result []EventStruct
	for _, ev := range events {
		// Исключаем события типа 1020 и в BCStatistics и те, чьи ID есть в excludeMap
		if ev.Type != 1020 && !excludeMap[ev.ID] && !utils.IsinSet(ev.Type, BCStatistics) {
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
func bcExcludeMistakes(events []EventStruct, settings SettingsStruct) ([]EventStruct, SettingsStruct) {

	var timer int64
	var part Part
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
						if utils.IsinSet(ev.Type, utils.Unblocks) {
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
						if utils.IsinSet(ev.Type, utils.Unblocks) {
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

	var result []EventStruct
	for i, v := range events {
		if !excludeMap[i] {
			result = append(result, v)
		}
	}

	return result, settings
}

func moveUp1102(events []EventStruct) []EventStruct {

	// Находим индекс первого события таймера первого тайма
	var timerIdx int = len(events)
	for i, ev := range events {
		if utils.IsinSet(ev.Type, utils.BcTimer) && *ev.I2 == 1 {
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
			events = append(events[:timerIdx], append([]EventStruct{ev}, events[timerIdx:]...)...)
			break
		}
	}

	return events
}

func bcReverse(events []EventStruct) []EventStruct {
	for _, event := range events {
		if value, ok := eventsWithTeam[event.Type]; ok {
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

func bcTransformation(request map[string]any, settings SettingsStruct) ([]EventStruct, SettingsStruct) {
	events := parsingEventsFun(request)                    // parse events
	events = bcExcludeEvents(events)                       // exclude 1020 and statistics
	events, settings = bcExcludeMistakes(events, settings) // exclude ends 1102 and 1103 with mistakes
	events = moveUp1102(events)                            // move up 1102 if mistake
	if settings.SportscastReverseTeams {                   //reverse broadcast if SportscastReverseTeams == true
		events = bcReverse(events)
	}

	return events, settings
}
