package input

import (
	"robots-go/utils"
	"time"
)

// Период в матче
type Part struct {
	IsGoing bool // Идет ли игра
	Nmb     int  // Номер периода (Если игра идет: 3 или 4 дополнительное время, 5 серия пенальти;
	// Если игра не идет: 0 матч не начался,  i перерыв после i-ого тайма, -1 матч закончился)
}

func bcTimeToTimestamp(regtime string) int64 {
	t, _ := time.Parse("02.01.2006 15:04:05", regtime)
	return t.Unix()
}

// Вычисление Part и Timer
func partTimer(events []EventStruct, timestamp int64, settings SettingsStruct) (Part, int64) {
	var part Part
	var timer int64
	for i, event := range events {
		switch {
		case utils.IsinSet(event.Type, utils.BcTimer):
			{ // Первое подходящее событие из BcTimer
				part.IsGoing = true
				part.Nmb = *event.I2
				regtime := bcTimeToTimestamp(event.RegTime)
				timer = int64(*event.I3) + timestamp - regtime

				// Ищем следующее подходящее событие из BcTimer для отсечения ошибок
				var timerOld int64 = -1
				for j := i + 1; j < len(events); j++ {
					eventOld := events[j]
					if utils.IsinSet(eventOld.Type, utils.BcTimer) && *eventOld.I2 == part.Nmb {
						regtimeOld := bcTimeToTimestamp(eventOld.RegTime)
						timerOld = int64(*eventOld.I3) + timestamp - regtimeOld
						break
					}
				}
				return part, max(timer, timerOld)
			}
		case event.Type == 1105:
			{ // Первое подходящее событие - начало матча
				part.IsGoing = true
				part.Nmb = 1
				regtime := bcTimeToTimestamp(event.RegTime)
				return part, timestamp - regtime
			}
		case event.Type == 1118:
			{ // Первое подходящее событие - начало тайма
				part.IsGoing = true
				part.Nmb = *event.I1
				regtime := bcTimeToTimestamp(event.RegTime)
				switch part.Nmb {
				case 1:
					return part, timestamp - regtime
				case 2:
					return part, settings.HalfDuration + timestamp - regtime
				case 3:
					return part, settings.MatchDuration + timestamp - regtime
				case 4:
					return part, settings.MatchDuration + 900 + timestamp - regtime
				case 5:
					return part, settings.MatchDuration + 1800
				}
			}
		case event.Type == 1103:
			{ // Первое подходящее событие - окончание матча
				part.IsGoing = false
				if settings.EventGameTypeIdent == "regular" {
					part.Nmb = -1
					return part, settings.MatchDuration
				} else {
					part.Nmb = 2
					return part, settings.MatchDuration
				}
			}
		case event.Type == 1102:
			{ // Первое подходящее событие - конец игрового отрезка
				part.IsGoing = false
				part.Nmb = *event.I1
				switch part.Nmb {
				case 1:
					return part, settings.HalfDuration
				case 2:
					if settings.EventGameTypeIdent == "regular" {
						part.Nmb = -1
						return part, settings.MatchDuration
					} else {
						part.Nmb = 2
						return part, settings.MatchDuration
					}
				case 3:
					return part, settings.MatchDuration + 900
				case 4:
					return part, settings.MatchDuration + 1800
				case 5:
					part.Nmb = -1
					return part, settings.MatchDuration + 1800
				}
			}
		}
	}
	part.IsGoing = false
	part.Nmb = 0
	return part, 0
}
