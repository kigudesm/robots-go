package input

import (
	"fmt"
	"robots-go/structures"
)

func PrepareInputFun(path string) []structures.EventStruct {

	request := readRequest(path)                                               // read request
	Settings := parsingSettingsFun(request)                                    // settings parsing and transformation
	basePoints := bPTransformation(request, Settings)                          // basePoints parsing and transformation
	events, Settings := bcTransformation(request, Settings)                    // broadcast parsing and transformation
	MatchStateCurrent := createMatchStateCurrent(events, Settings, basePoints) // creating MatchStateCurrent

	fmt.Println(Settings)
	fmt.Println(MatchStateCurrent)
	// fmt.Println(basePoints)

	return events
}
