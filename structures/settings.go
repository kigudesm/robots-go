package structures

// Начало конец таймов
type PartBeginEnd struct {
	Begin int64
	End   int64
}

// Настройки матча
type SettingsStruct struct {
	MatchType                     string               // Тип матча (обычный, товарищеский, 2 по 40 и т.п.)
	EventGameTypeIdent            string               // Тип матча (с добаленным временем, серией пенальти и т.п.)
	SwitchToTwoWayBetsProbability float64              // Вероятность для перехода на двухисходку
	HalfDuration                  int64                // Длительность тайма
	MatchDuration                 int64                // Длительность матча
	PartTimes                     map[int]PartBeginEnd // Начало и конец таймов
	SportscastReverseTeams        bool                 // реверсивная трансляция
	InjuryDefault                 [2]int               // Компенсированное время по умолчанию
	ServerTime                    int64                // Время сервера
	StartTime                     int64                // Время начала матча
	TargetEventKind               []string             // Целевые рынки
	Providers                     map[string]string    // Провайдеры
	Block                         []string             // Заблокированные eventKind
}
