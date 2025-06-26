package input

import (
	"fmt"
)

func PrepareInputFun(path string) []Event {

	request := readRequest(path)
	events := parsingEventsFun(request)
	events = bcExclude1020(events)

	fmt.Println(events)

	return events
}
