package bootstrap

type Config struct {
	Constant *Constants
	Env      *Env
}

func Run() *Config {
	return &Config{
		Constant: NewConstant(),
		Env:      NewEnv(),
	}
}
