package kitchen

type Kitchen struct {
	Conf *Config
}

func NewKitchen(configPath string) *Kitchen {
	return &Kitchen{
		Conf: NewConfig(configPath),
	}
}
