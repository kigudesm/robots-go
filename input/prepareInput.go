package input

import (
	"fmt"
	"robots-go/utils"
)

func PrepareInputFun(path string) []Event {

	request := utils.ReadRequest(path) // read request
	settings := parsingSettingsFun(request)
	events := parsingEventsFun(request) // parse events
	events = bcExcludeEvents(events)    // exclude 1020 and statistics

	fmt.Println(len(events))
	fmt.Println(settings)

	return events
}
