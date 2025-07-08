package input

import (
	"robots-go/constants"
	"robots-go/structures"
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
	events []structures.EventStruct) map[string]bool {

	result := make(map[string]bool)
	for _, eK := range settings.TargetEventKind {
		result[eK] = true
	}

	if matchState.Timestamp < settings.StartTime { // Блокировка до начала матча
		return make(map[string]bool)
	}

	if matchState.Part.Nmb == -1 { // Блокировка через 15 сек. после окончания матча
		for _, event := range events {
			if (event.Type == 1103) || (event.Type == 1102) {
				if matchState.Timestamp > bcTimeToTimestamp(event.RegTime)+15 {
					return make(map[string]bool)
				}
			}
		}
	}

	for _, eK := range settings.TargetEventKind { // Блокировка таймов, после окончания первого тайма
		if constants.EventKinds[eK].Name != "match" {
			if !(matchState.Part.Nmb == 0 || (matchState.Part.IsGoing && matchState.Part.Nmb == 1)) {
				result[eK] = false
			}
		}
	}

	if settings.FollowProviderCancels { // Следование за снятиями провайдера
		for _, eK := range settings.TargetEventKind {
			if settings.Providers[eK].MatchClosed {
				result[eK] = false
			}
		}
	}

	return result
}

func bPIDFun(eK string, settings structures.SettingsStruct) string {
	return "0"
}

func getActiveEventKind(settings structures.SettingsStruct, matchState structures.MatchStateCurrentStruct,
	events []structures.EventStruct) map[string]structures.EventKindsActive {

	// Объявляем мапу с указателями на структуры
	result := make(map[string]structures.EventKindsActive)
	isLiveSet := isLiveFun(settings, matchState, events)

	for _, eK := range settings.TargetEventKind {
		tmp := result[eK]
		tmp.BPID = bPIDFun(eK, settings)
		tmp.IsLive = isLiveSet[eK]
		tmp.IsActive = true
		tmp.IsBlocked = true
		result[eK] = tmp
	}
	return result
}

func createMatchStateCurrent(
	events []structures.EventStruct, settings structures.SettingsStruct) structures.MatchStateCurrentStruct {

	var matchState structures.MatchStateCurrentStruct

	matchState.Timestamp = settings.ServerTime
	events = cutEvents(events, matchState.Timestamp)
	matchState.Part, matchState.Timer = partTimer(events, settings.ServerTime, settings)
	matchState.Injury = bcGetInjury(events, settings)
	matchState.EventKinds = getActiveEventKind(settings, matchState, events)

	return matchState
}
