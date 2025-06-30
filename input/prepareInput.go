package input

import (
	"fmt"
	"robots-go/utils"
)

func PrepareInputFun(path string) []EventStruct {

	request := utils.ReadRequest(path)      // read request
	Settings := parsingSettingsFun(request) // parse settings
	events := parsingEventsFun(request)     // parse events
	events = bcExcludeEvents(events)        // exclude 1020 and statistics

	fmt.Println(timerPart(events, Settings.ServerTime, Settings))
	fmt.Println(Settings)

	return events
}
