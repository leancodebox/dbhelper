package app

var envExample string

func EnvExample(data string) {
	envExample = data
}

func GetEnvExample() string {
	return envExample
}
