package input

import (
	"fmt"
	"robots-go/structures"
)

func PrepareInputFun(path string) []structures.EventStruct {

	request := readRequest(path)                                   // read request
	Settings := parsingSettingsFun(request)                        // parse settings
	events, Settings := bcTransformation(request, Settings)        // broadcast transformation
	MatchStateCurrent := createMatchStateCurrent(events, Settings) // creating MatchStateCurrent

	fmt.Println(Settings)
	fmt.Println(MatchStateCurrent)

	return events
}
