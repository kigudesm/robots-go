package input

import (
	"robots-go/utils"
	"strconv"
	"time"
)

func bcTimeToTimestamp(regtime string) int64 {
	t, _ := time.Parse("02.01.2006 15:04:05", regtime)
	return t.Unix()
}

// Время таймера, кроме перерыва
func timerCalc(events []Event, half string, timestamp int64) int64 {

	// Преобразуем последний символ half в число (например, "half1" -> 1)
	halfNum, err := strconv.Atoi(string(half[len(half)-1]))
	if err != nil {
		return -1
	}

	for i, event := range events {
		switch {
		case utils.IsinSet(event.Type, utils.BcTimer) && *event.I2 == halfNum:
			{
				regtime := bcTimeToTimestamp(event.RegTime)
				timer := int64(*event.I3) + timestamp - regtime

				// Ищем следующее подходящее событие
				var timerOld int64 = -1
				for j := i + 1; j < len(events); j++ {
					nextEvent := events[j]
					if utils.IsinSet(nextEvent.Type, utils.BcTimer) && *nextEvent.I2 == halfNum {
						regtimeOld := bcTimeToTimestamp(nextEvent.RegTime)
						timerOld = int64(*nextEvent.I3) + timestamp - regtimeOld
						break
					}
				}
				return max(timer, timerOld)
			}
			// case utils.IsinSet(event.Type, utils.BcStart): {

			// }
		}

		// // Проверяем условие для текущего события
		// if utils.IsinSet(event.Type, utils.BcTimer) && *event.I2 == halfNum {
		// 	regtime := bcTimeToTimestamp(event.RegTime)
		// 	timer := int64(*event.I3) + timestamp - regtime

		// 	// Ищем следующее подходящее событие
		// 	var timerOld int64 = -1
		// 	for j := i + 1; j < len(events); j++ {
		// 		nextEvent := events[j]
		// 		if utils.IsinSet(nextEvent.Type, utils.BcTimer) && *nextEvent.I2 == halfNum {
		// 			regtimeOld := bcTimeToTimestamp(nextEvent.RegTime)
		// 			timerOld = int64(*nextEvent.I3) + timestamp - regtimeOld
		// 			break
		// 		}
		// 	}

		// 	return max(timer, timerOld)
		// }
	}
	return 0
}
