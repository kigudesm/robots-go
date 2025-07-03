package constants

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
