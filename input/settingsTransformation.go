package input

import (
	"math"
	"robots-go/constants"
	"robots-go/gametypes"
	"robots-go/structures"
	"strconv"
)

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

	// Следование за блокировками провайдера
	settings.FollowProviderBlocks = set["autoLiveScheme"].(map[string]any)["object"].(map[string]any)["followProviderBlocks"].(bool)

	// Находим провайдеров по каждому eventKind
	settings.Providers = getProviders(set)

	return settings
}
