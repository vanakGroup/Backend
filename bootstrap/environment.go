package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PrimaryDB Database
	Logger    Logger
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Logger struct {
	FilePath string
	Encoding string
	Level    string
	Logger   string
}

func NewEnv() *Env {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("fail to load env...")
	}
	return &Env{
		Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Logger{
			FilePath: os.Getenv("Logger_file_path"),
			Encoding: os.Getenv("Logger_encoding"),
			Level:    os.Getenv("Logger_level"),
			Logger:   os.Getenv("Logger"),
		},
	}
}
