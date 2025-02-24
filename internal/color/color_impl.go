package color

import (
	"encoding/json"
	"syscall/js"
)

// f___ go!

type Provider interface {
	Deserialize(array *map[string]Scheme)
}

// Deserialize
// returns array segment of supported color schemes
func Deserialize(array *map[string]Scheme) {
	raw := js.Global().Get("colors").String()
	if json.Unmarshal([]byte(raw), &array) != nil {
		*array = map[string]Scheme{
			"monochrome": Scheme{
				"#ffffff",
				"#777777",
				"#000000",
			},
		}
	}
}
