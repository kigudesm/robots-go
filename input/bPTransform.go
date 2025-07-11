package input

import (
	"robots-go/structures"
	"strconv"
)

func confertFactorsToStruct(ff []any) []structures.Factor {
	var factors []structures.Factor
	for _, f := range ff {
		ft := f.(map[string]any)

		var factor structures.Factor
		factor.ID = ft["f"].(string)
		num, _ := strconv.ParseFloat(ft["v"].(string), 64)
		factor.Value = num / 1000
		if p, ok := ft["p"].(float64); ok {
			val := p / 1000
			factor.Parameter = &val
		}
		factor.Disabled = ft["disabled"].(bool)
		factors = append(factors, factor)
	}
	return factors
}

func filterBPs(bPsPtr *[]structures.BasePoint, settPtr *structures.MatchSettings) []structures.BasePoint {

	var result []structures.BasePoint

	// Фильтрация по провайдерам и максимуму timestamp
	for _, eK := range settPtr.TargetEventKind {
		var timestampMax int64 = 0
		var bPMax structures.BasePoint
		for _, bP := range *bPsPtr {
			if bP.EventKind == eK && bP.ProviderLayer == settPtr.Providers[eK].ID && bP.TimeStamp > timestampMax {
				timestampMax = bP.TimeStamp
				bPMax = bP
			}
		}
		result = append(result, bPMax)
	}

	return result
}

func convertBPToStruct(bP any) structures.BasePoint {
	baseP := bP.(map[string]any)

	var basePoint structures.BasePoint

	// Обязательные поля
	basePoint.ID = baseP["id"].(string)
	basePoint.EventKind = baseP["eventKind"].(string)
	basePoint.ProviderLayer = baseP["providerLayer"].(string)
	num, _ := strconv.ParseInt(baseP["time_abs"].(string), 10, 64)
	basePoint.TimeStamp = num / 1000
	basePoint.ScoreC1 = int(baseP["scoreC1"].(float64))
	basePoint.ScoreC2 = int(baseP["scoreC2"].(float64))
	basePoint.Factors = confertFactorsToStruct(baseP["factors"].([]any))

	return basePoint
}

func bPTransformation(request map[string]any, settPtr *structures.MatchSettings) []structures.BasePoint {

	bPs := request["basePoints"].([]any)

	var basePoints []structures.BasePoint

	// Парсинг
	for _, bP := range bPs {
		basePoints = append(basePoints, convertBPToStruct(bP))
	}

	// Фильтрация по провайдерам и максимуму timestamp
	basePoints = filterBPs(&basePoints, settPtr)

	return basePoints
}
