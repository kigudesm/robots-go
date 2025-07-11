package input

import (
	"log"
	"robots-go/constants"
	"robots-go/structures"
	"robots-go/utils"
)

func bcGetInjury(eventsPtr *[]structures.EventInfo, settPtr *structures.MatchSettings) [2]int {
	result := settPtr.InjuryDefault
	for i := range 2 {
		for _, event := range *eventsPtr {
			if event.Type == 1104 && *event.I1 != 0 && *event.I2 == i+1 {
				result[i] = *event.I1 * 60
				break
			}
		}
	}
	return result
}

func cutEvents(events []structures.EventInfo, timestamp int64) []structures.EventInfo {
	for i, event := range events {
		if bcTimeToTimestamp(event.RegTime) <= timestamp {
			return events[i:]
		}
	}
	return events[:0]
}

func isLiveFun(settPtr *structures.MatchSettings, matchState structures.MatchStateCurrent,
	eventsPtr *[]structures.EventInfo, basePointsPtr *[]structures.BasePoint) map[string]bool {

	result := make(map[string]bool)
	for _, eK := range settPtr.TargetEventKind {
		result[eK] = true
	}

	// Блокировка до начала матча
	if matchState.Timestamp < settPtr.StartTime {
		return make(map[string]bool)
	}

	// Блокировка через 15 сек. после окончания матча
	if matchState.Half.Nmb == -1 {
		for _, event := range *eventsPtr {
			if (event.Type == 1103) || (event.Type == 1102) {
				if matchState.Timestamp > bcTimeToTimestamp(event.RegTime)+15 {
					return make(map[string]bool)
				}
			}
		}
	}

	// Блокировка таймов, после окончания первого тайма
	if !(matchState.Half.Nmb == 0 || (matchState.Half.IsGoing && matchState.Half.Nmb == 1)) {
		for _, eK := range settPtr.TargetEventKind {
			if constants.EventKinds[eK].Name != "match" {
				result[eK] = false
			}
		}
	}

	// Следование за снятиями провайдера
	if settPtr.FollowProviderCancels {
		for _, eK := range settPtr.TargetEventKind {
			switch {
			case settPtr.Providers[eK].MatchClosed:
				result[eK] = false
			case (eK == "100201" || eK == "100202") && matchState.Timer <= 2400:
				{
					for _, bP := range *basePointsPtr {
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

func isActiveFun(settPtr *structures.MatchSettings, matchState structures.MatchStateCurrent,
	isLiveSet map[string]bool) map[string]bool {

	result := isLiveSet

	// Матч почти завершен
	if matchState.Half.IsGoing && matchState.Half.Nmb == 2 &&
		matchState.Timer >= settPtr.MatchDuration+int64(matchState.InjuryTime[1]) {
		log.Println("Match is almost over")
		return make(map[string]bool)
	}

	// Первый тайм почти завершен
	if matchState.Half.IsGoing && matchState.Half.Nmb == 1 &&
		matchState.Timer >= settPtr.HalfDuration+int64(matchState.InjuryTime[0]) {
		log.Println("First half is almost over")
		for _, eK := range settPtr.TargetEventKind {
			if constants.EventKinds[eK].Name != "match" {
				result[eK] = false
			}
		}
	}
	return result
}

func blockAll(target []string) map[string]bool {
	result := make(map[string]bool)
	for _, eK := range target {
		result[eK] = true
	}
	return result
}

func isBlockedFun(settPtr *structures.MatchSettings, matchState structures.MatchStateCurrent,
	eventsPtr *[]structures.EventInfo, basePointsPtr *[]structures.BasePoint, isActiveSet map[string]bool) map[string]bool {

	result := make(map[string]bool)

	if settPtr.BlockAll { // Если ранее все было заблокировано
		return blockAll(settPtr.TargetEventKind)
	}

	// Блокировка по 1017
	for idx, event := range *eventsPtr {
		if event.Type == 1017 && utils.SliceToSet([]int{1, 4, 5, 6, 8, 10})[*event.I2] {
			flag := true
			setStruct := utils.UnionSets(utils.SliceToSet([]int{1019}), constants.Unblocks)
			for _, ev := range (*eventsPtr)[:idx] {
				if setStruct[ev.Type] {
					flag = false
					break
				}
			}
			if flag {
				return blockAll(settPtr.TargetEventKind)
			}
			break
		}
	}

	// Блокировка по 1006, 1134, 1135
	for idx, event := range *eventsPtr {
		if (event.Type == 1006 && utils.SliceToSet([]int{1, 2, 3})[*event.I1]) || event.Type == 1134 || event.Type == 1135 {
			flag := true
			for _, ev := range (*eventsPtr)[:idx] {
				if constants.Unblocks[ev.Type] {
					flag = false
					break
				}
			}
			if flag {
				if matchState.Timestamp-bcTimeToTimestamp(event.RegTime) >= 600 {
					*matchState.Alerts = append(*matchState.Alerts, "Источник СТ перенёс/отменил трансляцию")
				}
				return blockAll(settPtr.TargetEventKind)
			}
		}
	}

	// Блокировка по constants.Blocks
	for key, value := range constants.Blocks {
		for idx, event := range *eventsPtr {
			if event.Type == key {
				flag := true
				for _, ev := range (*eventsPtr)[:idx] {
					if value.Cancel[ev.Type] {
						flag = false
						break
					}
				}
				// Разблокировка 1149 и 1176 по опорной точке на матч
				if event.Type == 1149 || event.Type == 1176 {
					for _, bP := range *basePointsPtr {
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
	if settPtr.FollowProviderBlocks {
		providerBlock := make(map[string]string)
		for eK := range map[string]bool{"1": true, "100201": true} {
			for _, bP := range *basePointsPtr {
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
	for _, eK := range settPtr.TargetEventKind {
		result[eK] = !isActiveSet[eK] || result[eK]
	}

	return result

}

func getActiveEventKind(settPtr *structures.MatchSettings, matchState structures.MatchStateCurrent,
	eventsPtr *[]structures.EventInfo, basePointsPtr *[]structures.BasePoint) map[string]structures.EventKindsActive {

	// Объявляем мапу с указателями на структуры
	result := make(map[string]structures.EventKindsActive)
	isLiveSet := isLiveFun(settPtr, matchState, eventsPtr, basePointsPtr)
	isActiveSet := isActiveFun(settPtr, matchState, isLiveSet)
	isBlockedSet := isBlockedFun(settPtr, matchState, eventsPtr, basePointsPtr, isActiveSet)

	for _, eK := range settPtr.TargetEventKind {
		tmp := result[eK]
		// tmp.BPID = bPIDFun(eK, basePoints)
		tmp.IsLive = isLiveSet[eK]
		tmp.IsActive = isActiveSet[eK]
		tmp.IsBlocked = isBlockedSet[eK]
		result[eK] = tmp
	}
	return result
}

func suspendedFun(eventsPtr *[]structures.EventInfo) bool {
	for idx, event := range *eventsPtr {
		if event.Type == 1017 && utils.SliceToSet([]int{1, 4})[*event.I2] {
			flag := true
			setStruct := utils.UnionSets(utils.SliceToSet([]int{1019}), constants.Unblocks)
			for _, ev := range (*eventsPtr)[:idx] {
				if setStruct[ev.Type] {
					flag = false
					break
				}
			}
			return flag
		}
	}
	return false
}

func varFun(eventsPtr *[]structures.EventInfo) bool {
	for idx, event := range *eventsPtr {
		if event.Type == 1067 {
			flag := true
			for _, ev := range (*eventsPtr)[:idx] {
				if constants.Blocks[1067].Cancel[ev.Type] {
					flag = false
					break
				}
			}
			return flag
		}
	}
	return false
}

func penaltyFun(eventsPtr *[]structures.EventInfo) int {
	penalty := 0
	for idx, event := range *eventsPtr {
		if event.Type == 1110 {
			penalty = *event.I1
			for _, ev := range (*eventsPtr)[:idx] {
				if constants.PenaltyUnblocks[ev.Type] {
					penalty = 0
					break
				}
			}
			break
		}
	}
	return penalty
}

func removalFun(bcProvider string, eventsPtr *[]structures.EventInfo) [2]int {

	var types map[int]bool

	if bcProvider == "9" {
		types = map[int]bool{1109: true, 1235: true}
	} else {
		types = map[int]bool{1109: true}
	}

	removal := [2]int{0, 0}
	for _, event := range *eventsPtr {
		if types[event.Type] {
			removal[*event.I1-1] += 1
		}
	}
	return removal
}

func createMatchStateCurrent(evPtr *[]structures.EventInfo, settPtr *structures.MatchSettings,
	basePointsPtr *[]structures.BasePoint) structures.MatchStateCurrent {

	var matchState structures.MatchStateCurrent

	matchState.Timestamp = settPtr.ServerTime
	events := cutEvents(*evPtr, matchState.Timestamp)
	matchState.Half, matchState.Timer = partTimer(&events, settPtr.ServerTime, settPtr)
	matchState.InjuryTime = bcGetInjury(&events, settPtr)
	matchState.Alerts = &[]string{}
	matchState.EventKinds = getActiveEventKind(settPtr, matchState, &events, basePointsPtr)
	matchState.Suspended = suspendedFun(&events)
	matchState.Var = varFun(&events)
	matchState.Penalty = penaltyFun(&events)
	matchState.Removal = removalFun(settPtr.SportscastProviderLayerId, &events)
	// matchState.History = historyFun(&events)

	return matchState
}
