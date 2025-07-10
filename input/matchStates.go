package input

import (
	"log"
	"robots-go/constants"
	"robots-go/structures"
	"robots-go/utils"
)

func bcGetInjury(events []structures.EventStruct, settings structures.SettingsStruct) [2]int {
	result := settings.InjuryDefault
	for i := range 2 {
		for _, event := range events {
			if event.Type == 1104 && *event.I1 != 0 && *event.I2 == i+1 {
				result[i] = *event.I1 * 60
				break
			}
		}
	}
	return result
}

func cutEvents(events []structures.EventStruct, timestamp int64) []structures.EventStruct {
	for i, event := range events {
		if bcTimeToTimestamp(event.RegTime) <= timestamp {
			return events[i:]
		}
	}
	return events[:0]
}

func isLiveFun(settings structures.SettingsStruct, matchState structures.MatchStateCurrentStruct,
	events []structures.EventStruct, basePoints []structures.BasePointStruct) map[string]bool {

	result := make(map[string]bool)
	for _, eK := range settings.TargetEventKind {
		result[eK] = true
	}

	// Блокировка до начала матча
	if matchState.Timestamp < settings.StartTime {
		return make(map[string]bool)
	}

	// Блокировка через 15 сек. после окончания матча
	if matchState.Part.Nmb == -1 {
		for _, event := range events {
			if (event.Type == 1103) || (event.Type == 1102) {
				if matchState.Timestamp > bcTimeToTimestamp(event.RegTime)+15 {
					return make(map[string]bool)
				}
			}
		}
	}

	// Блокировка таймов, после окончания первого тайма
	if !(matchState.Part.Nmb == 0 || (matchState.Part.IsGoing && matchState.Part.Nmb == 1)) {
		for _, eK := range settings.TargetEventKind {
			if constants.EventKinds[eK].Name != "match" {
				result[eK] = false
			}
		}
	}

	// Следование за снятиями провайдера
	if settings.FollowProviderCancels {
		for _, eK := range settings.TargetEventKind {
			switch {
			case settings.Providers[eK].MatchClosed:
				result[eK] = false
			case (eK == "100201" || eK == "100202") && matchState.Timer <= 2400:
				{
					for _, bP := range basePoints {
						if bP.EventKind == eK && len(bP.Factors) > 0 {
							flag := true
							for _, factor := range bP.Factors {
								if !factor.Disabled || factor.Value != 0 {
									flag = false
									break
								}
							}
							if flag {
								result[eK] = false
							}
							break
						}
					}
				}
			}
		}
	}

	return result
}

// func bPIDFun(eK string, basePoints []structures.BasePointStruct) string {
// 	for _, bP := range basePoints {
// 		if bP.EventKind == eK {
// 			return bP.ID
// 		}
// 	}
// 	return "0"
// }

func isActiveFun(settings structures.SettingsStruct, matchState structures.MatchStateCurrentStruct,
	isLiveSet map[string]bool) map[string]bool {

	result := isLiveSet

	// Матч почти завершен
	if matchState.Part.IsGoing && matchState.Part.Nmb == 2 &&
		matchState.Timer >= settings.MatchDuration+int64(matchState.Injury[1]) {
		log.Println("Match is almost over")
		return make(map[string]bool)
	}

	// Первый тайм почти завершен
	if matchState.Part.IsGoing && matchState.Part.Nmb == 1 &&
		matchState.Timer >= settings.HalfDuration+int64(matchState.Injury[0]) {
		log.Println("First half is almost over")
		for _, eK := range settings.TargetEventKind {
			if constants.EventKinds[eK].Name != "match" {
				result[eK] = false
			}
		}
	}
	return result
}

func blockAll(target []string, alerts []string) (map[string]bool, []string) {
	result := make(map[string]bool)
	for _, eK := range target {
		result[eK] = true
	}
	return result, alerts
}

func isBlockedFun(
	settings structures.SettingsStruct,
	matchState structures.MatchStateCurrentStruct,
	events []structures.EventStruct,
	basePoints []structures.BasePointStruct,
	isActiveSet map[string]bool) (
	map[string]bool, []string) {

	result := make(map[string]bool)
	alerts := make([]string, 0)

	if settings.BlockAll { // Если ранее все было заблокировано
		return blockAll(settings.TargetEventKind, alerts)
	}

	// Блокировка по 1017
	for idx, event := range events {
		if event.Type == 1017 && utils.SliceToSet([]int{1, 4, 5, 6, 8, 10})[*event.I2] {
			flag := true
			setStruct := utils.UnionSets(utils.SliceToSet([]int{1019}), constants.Unblocks)
			for _, ev := range events[:idx] {
				if setStruct[ev.Type] {
					flag = false
					break
				}
			}
			if flag {
				return blockAll(settings.TargetEventKind, alerts)
			}
			break
		}
	}

	// Блокировка по 1006, 1134, 1135
	for idx, event := range events {
		if (event.Type == 1006 && utils.SliceToSet([]int{1, 2, 3})[*event.I1]) || event.Type == 1134 || event.Type == 1135 {
			flag := true
			for _, ev := range events[:idx] {
				if constants.Unblocks[ev.Type] {
					flag = false
					break
				}
			}
			if flag {
				if matchState.Timestamp-bcTimeToTimestamp(event.RegTime) >= 600 {
					alerts = append(alerts, "Источник СТ перенёс/отменил трансляцию")
				}
				return blockAll(settings.TargetEventKind, alerts)
			}
		}
	}

	// Блокировка по constants.Blocks
	for key, value := range constants.Blocks {
		for idx, event := range events {
			if event.Type == key {
				flag := true
				for _, ev := range events[:idx] {
					if value.Cancel[ev.Type] {
						flag = false
						break
					}
				}
				// Разблокировка 1149 и 1176 по опорной точке на матч
				if event.Type == 1149 || event.Type == 1176 {
					for _, bP := range basePoints {
						if bP.EventKind == "1" && bP.TimeStamp > bcTimeToTimestamp(event.RegTime) {
							lenFactors := 0
							for _, factor := range bP.Factors {
								if factor.Value != 0 && !factor.Disabled {
									lenFactors++
								}
							}
							if lenFactors >= 2 {
								flag = false
							}
						}
					}
				}
				if flag {
					for eK := range value.Block {
						result[eK] = true
					}
				}
				break
			}
		}
	}

	// Блокировка {1, 100201} за котировками провайдера
	if settings.FollowProviderBlocks {
		providerBlock := make(map[string]string)
		for eK := range map[string]bool{"1": true, "100201": true} {
			for _, bP := range basePoints {
				if bP.EventKind == eK {
					lenNonZero := 0
					lenDisabled := 0
					for _, factor := range bP.Factors {
						if factor.Value != 0 {
							lenNonZero++
							if factor.Disabled {
								lenDisabled++
							}
						}
					}
					switch {
					case lenNonZero == 0:
						{
							providerBlock[eK] = "empty"
						}
					case lenDisabled == lenNonZero:
						{
							providerBlock[eK] = "blocked"
							result[eK] = true
						}
					default:
						{
							providerBlock[eK] = "unblocked"
						}
					}
				}
			}
		}
		if providerBlock["1"] == "blocked" && providerBlock["100201"] == "empty" {
			result["100201"] = true
		}
	}

	// Блокируем все неактивные eventKind
	for _, eK := range settings.TargetEventKind {
		result[eK] = !isActiveSet[eK] || result[eK]
	}

	return result, alerts

}

func getActiveEventKind(settings structures.SettingsStruct, matchState structures.MatchStateCurrentStruct,
	events []structures.EventStruct, basePoints []structures.BasePointStruct) (
	map[string]structures.EventKindsActive, []string) {

	// Объявляем мапу с указателями на структуры
	result := make(map[string]structures.EventKindsActive)
	isLiveSet := isLiveFun(settings, matchState, events, basePoints)
	isActiveSet := isActiveFun(settings, matchState, isLiveSet)
	isBlockedSet, alerts := isBlockedFun(settings, matchState, events, basePoints, isActiveSet)

	for _, eK := range settings.TargetEventKind {
		tmp := result[eK]
		// tmp.BPID = bPIDFun(eK, basePoints)
		tmp.IsLive = isLiveSet[eK]
		tmp.IsActive = isActiveSet[eK]
		tmp.IsBlocked = isBlockedSet[eK]
		result[eK] = tmp
	}
	return result, alerts
}

func createMatchStateCurrent(events []structures.EventStruct, settings structures.SettingsStruct,
	basePoints []structures.BasePointStruct) structures.MatchStateCurrentStruct {

	var matchState structures.MatchStateCurrentStruct

	matchState.Timestamp = settings.ServerTime
	events = cutEvents(events, matchState.Timestamp)
	matchState.Part, matchState.Timer = partTimer(events, settings.ServerTime, settings)
	matchState.Injury = bcGetInjury(events, settings)
	matchState.EventKinds, matchState.Alerts = getActiveEventKind(settings, matchState, events, basePoints)

	return matchState
}
