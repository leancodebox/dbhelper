package app

var defaultConfig string

func DefaultConfig(data string) {
	defaultConfig = data
}

func GetDefaultConfig() string {
	return defaultConfig
}
