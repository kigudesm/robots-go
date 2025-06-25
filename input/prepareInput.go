package input

import (
	"fmt"
)

func PrepareInputFun(path string) []Event {

	request := readRequest(path)
	events := parsingEventsFun(request)

	fmt.Println(events)

	return events
}
