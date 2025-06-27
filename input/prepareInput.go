package input

import (
	"fmt"
	"robots-go/utils"
)

func PrepareInputFun(path string) []Event {

	request := utils.ReadRequest(path)      // read request
	settings := parsingSettingsFun(request) // parse settings
	events := parsingEventsFun(request)     // parse events
	events = bcExcludeEvents(events)        // exclude 1020 and statistics

	fmt.Println(timerCalc(events, "Half 2", settings.ServerTime))
	fmt.Println(settings)

	return events
}
