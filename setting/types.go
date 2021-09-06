package setting

type Setting struct {
	Hour  string `json:"hour"`
	Image string `json:"image"`
}

type Config struct {
	Setting []Setting `json:"setting"`
}
