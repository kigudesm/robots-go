package input

import (
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

func parsingSettingsFun(request map[string]any) structures.SettingsStruct {

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
	var settings structures.SettingsStruct
	settings.TargetEventKind = targetEventKind
	settings.MatchType = ident
	settings.EventGameTypeIdent = eGTI
	settings.HalfDuration = int64(matchType.HalfDuration)
	settings.MatchDuration = 2 * int64(matchType.HalfDuration)
	settings.PartTimes = map[int]structures.PartBeginEnd{
		1: {Begin: 0, End: settings.HalfDuration},
		2: {Begin: settings.HalfDuration, End: settings.MatchDuration},
		3: {Begin: settings.MatchDuration, End: settings.MatchDuration + 900},
		4: {Begin: settings.MatchDuration + 900, End: settings.MatchDuration + 1800},
		5: {Begin: settings.MatchDuration + 1800, End: settings.MatchDuration + 1800},
	}
	settings.SportscastReverseTeams = set["sportscastReverseTeams"].(bool)
	settings.InjuryDefault[0] = matchType.TFirstHalf - matchType.HalfDuration
	settings.InjuryDefault[1] = matchType.TSecondHalf - matchType.HalfDuration
	settings.ServerTime = num / 1000

	return settings
}
