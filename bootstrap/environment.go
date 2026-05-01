package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PrimaryDB Database
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
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
	}
}
