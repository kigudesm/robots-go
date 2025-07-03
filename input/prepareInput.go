package input

import (
	"fmt"
	"robots-go/structures"
)

func PrepareInputFun(path string) []structures.EventStruct {

	request := readRequest(path)            // read request
	Settings := parsingSettingsFun(request) // parse settings
	events, Settings := bcTransformation(request, Settings)

	fmt.Println(partTimer(events, Settings.ServerTime, Settings))
	fmt.Println(Settings)

	return events
}
