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

func onlyMatch(eventKinds []string) []string {
	return eventKinds
}

func getActiveEventKind(settings structures.SettingsStruct, matchState structures.MatchStateCurrentStruct) []string {

	result := settings.TargetEventKind
	switch {
	case matchState.Part.Nmb == -1:
		{
			result = result[:0]
		}
	case matchState.Timer >= settings.MatchDuration+int64(matchState.Injury[1]):
		{
			result = result[:0]
		}
	case matchState.Part.IsGoing && matchState.Part.Nmb == 1 &&
		matchState.Timer >= settings.HalfDuration+int64(matchState.Injury[0]):
		{
			res := make([]string, 0)
			for _, eK := range result {
				name := constants.EventKinds[eK].Name
				if name == "match" {
					res = append(res, eK)
				}
			}
			result = res
		}
	}
	return result
}

func createMatchStateCurrent(
	events []structures.EventStruct, settings structures.SettingsStruct) structures.MatchStateCurrentStruct {

	var matchState structures.MatchStateCurrentStruct

	matchState.Timestamp = settings.ServerTime
	// matchState.Events = cutEvents(events, matchState.Timestamp)
	matchState.Part, matchState.Timer = partTimer(events, settings.ServerTime, settings)
	matchState.Injury = bcGetInjury(events, settings)
	matchState.ActiveEventKinds = getActiveEventKind(settings, matchState)

	return matchState
}
