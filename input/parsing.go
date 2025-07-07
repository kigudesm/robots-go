package input

import (
	"math"
	"robots-go/constants"
	"robots-go/gametypes"
	"robots-go/structures"
	"strconv"
)

func convertEventToStruct(ev map[string]any) structures.EventStruct {

	var event structures.EventStruct

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

func parsingEventsFun(request map[string]any) []structures.EventStruct {

	evs := request["events"].([]any)

	var events []structures.EventStruct

	for _, item := range evs {
		if ev, ok := item.(map[string]any); ok {
			events = append(events, convertEventToStruct(ev))
		}
	}

	return events
}

func getProviders(set map[string]any) map[string]structures.ProviderStruct {
	result := make(map[string]structures.ProviderStruct)
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

func parsingSettingsFun(request map[string]any) structures.SettingsStruct {

	set := request["settings"].(map[string]any)
	var settings structures.SettingsStruct

	// Парсим targetEventKind
	tek := set["targetEventKind"].([]any)
	var targetEventKind []string
	for _, item := range tek {
		if ek, ok := item.(string); ok {
			// Убираем необрабатываемые targetEventKind
			if _, ok := constants.EventKinds[ek]; ok {
				targetEventKind = append(targetEventKind, ek)
			}
		}
	}
	settings.TargetEventKind = targetEventKind

	// Парсим ident - продолжительность таймов, товарищеский или нет
	ident := set["sportRules"].(map[string]any)["object"].(map[string]any)["ident"].(string)
	settings.MatchType = ident

	// Продолжительности матча, и времена начала и окончания таймов, компенсированное время
	matchType := gametypes.MatchTypes[ident]
	settings.HalfDuration = int64(matchType.HalfDuration)
	settings.MatchDuration = 2 * int64(matchType.HalfDuration)
	settings.PartTimes = map[int]structures.PartBeginEnd{
		1: {Begin: 0, End: settings.HalfDuration},
		2: {Begin: settings.HalfDuration, End: settings.MatchDuration},
		3: {Begin: settings.MatchDuration, End: settings.MatchDuration + 900},
		4: {Begin: settings.MatchDuration + 900, End: settings.MatchDuration + 1800},
		5: {Begin: settings.MatchDuration + 1800, End: settings.MatchDuration + 1800},
	}
	settings.InjuryDefault[0] = matchType.TFirstHalf - matchType.HalfDuration
	settings.InjuryDefault[1] = matchType.TSecondHalf - matchType.HalfDuration

	// Парсим EventGameTypeIdent - обычный или с добавочным временем и пенальти
	var eGTI string
	if set["eventGameTypeIdent"] == nil {
		eGTI = "regular"
	} else {
		eGTI = set["eventGameTypeIdent"].(string)
	}
	settings.EventGameTypeIdent = eGTI

	// Парсим условие реверсности трансляции
	settings.SportscastReverseTeams = set["sportscastReverseTeams"].(bool)

	// Серверное время
	num, _ := strconv.ParseInt(set["serverTime"].(string), 10, 64)
	settings.ServerTime = num / 1000

	// Время начала матча
	num, _ = strconv.ParseInt(set["startTime"].(string), 10, 64)
	settings.StartTime = num / 1000

	// Переводим в вероятность число из switchToTwoWayBetsProbability
	prStr := set["autoLiveScheme"].(map[string]any)["object"].(map[string]any)["switchToTwoWayBetsProbability"].(string)
	pr, _ := strconv.ParseFloat(prStr, 64)
	settings.SwitchToTwoWayBetsProbability = pr * math.Pow10(-8)

	// Следование за снятиями провайдера
	settings.FollowProviderCancels = set["autoLiveScheme"].(map[string]any)["object"].(map[string]any)["followProviderCancels"].(bool)

	// Находим провайдеров по каждому eventKind
	settings.Providers = getProviders(set)

	return settings
}
