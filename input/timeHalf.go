package input

import (
	"robots-go/constants"
	"robots-go/structures"
	"time"
)

// // Период в матче
// type Part struct {
// 	IsGoing bool // Идет ли игра
// 	Nmb     int  // Номер периода (Если игра идет: 3 или 4 дополнительное время, 5 серия пенальти;
// 	// Если игра не идет: 0 матч не начался,  i перерыв после i-ого тайма, -1 матч закончился)
// }

func bcTimeToTimestamp(regtime string) int64 {
	t, _ := time.Parse("02.01.2006 15:04:05", regtime)
	return t.Unix()
}

// Вычисление Part и Timer
func partTimer(events []structures.EventStruct, timestamp int64, settings structures.SettingsStruct) (structures.Part, int64) {
	var part structures.Part
	var timer int64
	for i, event := range events {
		if _, ok := constants.BcTimer[event.Type]; ok { // Событие из BcTimer
			part.IsGoing = true
			part.Nmb = *event.I2
			regtime := bcTimeToTimestamp(event.RegTime)
			timer = int64(*event.I3) + timestamp - regtime
			// Ищем следующее подходящее событие из BcTimer для отсечения ошибок
			var timerOld int64 = -1
			for j := i + 1; j < len(events); j++ {
				eventOld := events[j]
				if _, ok := constants.BcTimer[eventOld.Type]; ok && *eventOld.I2 == part.Nmb {
					regtimeOld := bcTimeToTimestamp(eventOld.RegTime)
					timerOld = int64(*eventOld.I3) + timestamp - regtimeOld
					break
				}
			}
			return part, max(timer, timerOld)
		}

		switch event.Type {
		case 1105: // Начало матча
			{
				part.IsGoing = true
				part.Nmb = 1
				regtime := bcTimeToTimestamp(event.RegTime)
				return part, timestamp - regtime
			}
		case 1118: // Начало тайма
			{
				part.IsGoing = true
				part.Nmb = *event.I1
				regtime := bcTimeToTimestamp(event.RegTime)
				return part, settings.PartTimes[part.Nmb].Begin + timestamp - regtime
			}
		case 1103: // Окончание матча
			{
				part.IsGoing = false
				if settings.EventGameTypeIdent == "regular" {
					part.Nmb = -1
					return part, settings.MatchDuration
				} else {
					part.Nmb = 2
					return part, settings.MatchDuration
				}
			}
		case 1102: // Конец тайма
			{
				part.IsGoing = false
				part.Nmb = *event.I1
				if (part.Nmb == 2 && settings.EventGameTypeIdent == "regular") || (part.Nmb == 5) {
					part.Nmb = -1
				}
				return part, settings.PartTimes[*event.I1].End
			}
		}
	}
	part.IsGoing = false
	part.Nmb = 0
	return part, 0
}
