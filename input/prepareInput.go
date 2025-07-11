package input

import (
	"fmt"
	"robots-go/structures"
)

func PrepareInputFun(path string) []structures.EventInfo {

	request := readRequest(path)                                                  // read request
	Settings := parsingSettingsFun(request)                                       // settings parsing and transformation
	basePoints := bPTransformation(request, &Settings)                            // basePoints parsing and transformation
	events := bcTransformation(request, &Settings)                                // broadcast parsing and transformation
	MatchStateCurrent := createMatchStateCurrent(&events, &Settings, &basePoints) // creating MatchStateCurrent

	fmt.Println(Settings)
	fmt.Println(MatchStateCurrent)
	// fmt.Println(basePoints)

	return events
}
