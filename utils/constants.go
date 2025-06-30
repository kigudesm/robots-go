package utils

var Unblocks = map[int]struct{}{
	1100: {}, 1101: {}, 1102: {}, 1103: {}, 1113: {}, 1115: {}, 1118: {}, 1164: {}, 1168: {}, 1196: {},
	1421: {},
}

var BcTimer = map[int]struct{}{
	1100: {}, 1101: {}, 1106: {}, 1107: {}, 1110: {}, 1113: {}, 1114: {}, 1115: {}, 1116: {}, 1117: {},
	1152: {}, 1163: {}, 1168: {}, 1169: {}, 1188: {}, 1190: {}, 1192: {}, 1194: {}, 1195: {}, 1196: {},
	1197: {}, 1198: {}, 1199: {}, 1228: {}, 1234: {}, 1255: {},
}

var BcStart = map[int]string{
	1105: "match",
	1118: "half",
}

var BcEnd = map[int]string{
	1103: "match",
	1102: "half",
}

var EventKinds = map[string]struct {
	Name string
	Kind string
}{
	"1":        {"match", "goals"},
	"100201":   {"first", "goals"},
	"100202":   {"second", "goals"},
	"400100":   {"match", "corners"},
	"10100201": {"first", "corners"},
	"400200":   {"match", "yellow cards"},
	"10200201": {"first", "yellow cards"},
	"400300":   {"match", "fouls"},
	"10300201": {"first", "fouls"},
	"400400":   {"match", "shots on goal"},
	"10400201": {"first", "shots on goal"},
	"400500":   {"match", "offsides"},
	"10500201": {"first", "offsides"},
	"400700":   {"match", "total shots"},
	"10700201": {"first", "total shots"},
	"400800":   {"match", "hit the woodwork"},
	"10800201": {"first", "hit the woodwork"},
	"401000":   {"match", "throw-ins"},
	"11000201": {"first", "throw-ins"},
	"401100":   {"match", "goal kicks"},
	"11100201": {"first", "goal kicks"},
}
