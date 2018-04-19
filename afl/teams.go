package afl

var nicknames map[string]string
var teams []string

func readData() {
	teams = []string{
		"Adelaide Crows",
		"Brisbane Lions",
		"Carlton",
		"Collingwood",
		"Essendon",
		"Fremantle",
		"Geelong Cats",
		"Gold Coast Suns",
		"GWS Giants",
		"Hawthorn",
		"Melbourne",
		"North Melbourne",
		"Port Adelaide Power",
		"Richmond",
		"St Kilda",
		"Sydney Swans",
		"West Coast Eagles",
		"Western Bulldogs",
	}

	nicknames = map[string]string{
		"Crows":         "Adelaide Crows",
		"Adelaide":      "Adelaide Crows",
		"Lions":         "Brisbane Lions",
		"Brisbane":      "Brisbane Lions",
		"Blues":         "Carlton",
		"Carlton Blues": "Carlton",
	}
}
