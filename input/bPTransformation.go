package input

import (
	"robots-go/structures"
	"strconv"
)

// func confertFactorsToStruct(ff []any) []structures.Factor {
// 	var factors []structures.Factor
// 	for _, f := range ff {
// 		ft := f.(map[string]any)
// 		var factor structures.Factor
// 		factor.FactorID = ft["f"].(string)
// 		factor.Value = ft["v"]
// 		factors = append(factors, factor)
// 	}
// 	return factors
// }

func filterBPs(bPs []structures.BasePointStruct, settings structures.SettingsStruct) []structures.BasePointStruct {

	var result []structures.BasePointStruct

	// Фильтрация по провайдерам и максимуму timestamp
	for _, eK := range settings.TargetEventKind {
		var timestampMax int64 = 0
		var bPMax structures.BasePointStruct
		for _, bP := range bPs {
			if bP.EventKind == eK && bP.ProviderLayer == settings.Providers[eK].ID && bP.TimeStamp > timestampMax {
				timestampMax = bP.TimeStamp
				bPMax = bP
			}
		}
		result = append(result, bPMax)
	}

	return result
}

func convertBPToStruct(bP any) structures.BasePointStruct {
	baseP := bP.(map[string]any)

	var basePoint structures.BasePointStruct

	// Обязательные поля
	basePoint.ID = baseP["id"].(string)
	basePoint.EventKind = baseP["eventKind"].(string)
	basePoint.ProviderLayer = baseP["providerLayer"].(string)
	num, _ := strconv.ParseInt(baseP["time_abs"].(string), 10, 64)
	basePoint.TimeStamp = num / 1000
	basePoint.ScoreC1 = int(baseP["scoreC1"].(float64))
	basePoint.ScoreC2 = int(baseP["scoreC2"].(float64))
	// basePoint.Factors = confertFactorsToStruct(baseP["factors"].([]any))

	return basePoint
}

func bPTransformation(request map[string]any, settings structures.SettingsStruct) []structures.BasePointStruct {

	bPs := request["basePoints"].([]any)

	var basePoints []structures.BasePointStruct

	// Парсинг
	for _, bP := range bPs {
		basePoints = append(basePoints, convertBPToStruct(bP))
	}

	// Фильтрация
	basePoints = filterBPs(basePoints, settings)

	return basePoints
}
