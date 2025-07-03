package structures

// Начало конец таймов
type PartBeginEnd struct {
	Begin int64
	End   int64
}

// Настройки матча
type SettingsStruct struct {
	MatchType              string               // Тип матча (обычный, товарищеский, 2 по 40 и т.п.)
	EventGameTypeIdent     string               // Тип матча (с добаленным временем, серией пенальти и т.п.)
	HalfDuration           int64                // Длительность тайма
	MatchDuration          int64                // Длительность матча
	PartTimes              map[int]PartBeginEnd // Начало и конец таймов
	SportscastReverseTeams bool                 // реверсивная трансляция
	InjuryDefault          [2]int               // Компенсированное время по умолчанию
	ServerTime             int64                // Время сервера
	TargetEventKind        []string             // Целевые рынки
	Block                  []string             // Заблокированные eventKind
}
