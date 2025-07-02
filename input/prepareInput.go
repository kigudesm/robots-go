package input

import (
	"fmt"
)

func PrepareInputFun(path string) []EventStruct {

	request := readRequest(path)            // read request
	Settings := parsingSettingsFun(request) // parse settings
	events, Settings := bcTransformation(request, Settings)

	fmt.Println(partTimer(events, Settings.ServerTime, Settings))
	fmt.Println(Settings)

	return events
}
