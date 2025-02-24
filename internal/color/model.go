package color

// Scheme
// Represents structure of Color supported by
// Window manager. Knowing, the Color, which can supports
// by the WM not once, in real life will be array segment of Scheme
type Scheme struct {
	Faded  string `json:"faded"`
	Normal string `json:"normal"`
	Vivid  string `json:"vivid"`
}

// SupportedSchemes
// Provides deserialized JSON module, which contains
// loadable in binary colors segments.
var SupportedSchemes = map[string]Scheme{}
