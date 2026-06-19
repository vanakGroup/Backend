package bootstrap

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type S3 struct {
	Buckets   BucketName
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
}

type BucketName struct {
	LogoPic             string
	ParkingDoc          string
	ParkingPic          string
	RegistrationCardPic string
}

type Env struct {
	PrimaryDB Database
	Logger    Logger
	Redis     Redis
	OTP       OTP
	Jwt       Jwt
	Storage   S3
	Hashed    HashIDConfig
}

type HashIDConfig struct {
	Salt      string
	MinLength int
}

type Jwt struct {
	SecretKey string
}

type OTP struct {
	Length       int
	ExpiryMinute int
	MaxAttempts  int
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

type Redis struct {
	Host            string
	Port            string
	Password        string
	Db              string
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolSize        int
	PoolTimeout     time.Duration
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

func NewEnv() *Env {
	err := godotenv.Load("./.env")
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
			FilePath: os.Getenv("LOGGER_FILE_PATH"),
			Encoding: os.Getenv("LOGGER_ENCODING"),
			Level:    os.Getenv("LOGGER_LEVEL"),
			Logger:   os.Getenv("LOGGER"),
		},
		Redis{
			Host:            os.Getenv("REDIS_HOST_CONFIG"),
			Port:            os.Getenv("REDIS_PORT_CONFIG"),
			Password:        os.Getenv("REDIS_PASSWORD_CONFIG"),
			Db:              os.Getenv("REDIS_DB_CONFIG"),
			DialTimeout:     getEnvDuration("REDIS_DIALTIMEOUT_CONFIG", 5),
			ReadTimeout:     getEnvDuration("REDIS_READTIMEOUT_CONFIG", 5),
			WriteTimeout:    getEnvDuration("REDIS_WRITETIMEOUT_CONFIG", 5),
			PoolTimeout:     getEnvDuration("REDIS_POOLTIMEOUT_CONFIG", 15),
			ConnMaxIdleTime: getEnvDuration("REDIS_CONNMAXIDLETIME_CONFIG", 500),
			ConnMaxLifetime: getEnvDuration("REDIS_CONNMAXLIFETIME_CONFIG", 500),
			PoolSize:        getEnvInt("REDIS_POOLSIZE_CONFIG", 10),
		},
		OTP{
			Length:       getEnvInt("OTP_LENGTH", 6),
			ExpiryMinute: getEnvInt("OTP_EXPIRY_MINUTES", 2),
			MaxAttempts:  getEnvInt("OTP_MAX_ATTEMPTS", 3),
		},
		Jwt{
			SecretKey: os.Getenv("SECRET_KEY"),
		},
		S3{
			Region:    os.Getenv("BUCKET_REGION"),
			AccessKey: os.Getenv("BUCKET_ACCESS_KEY"),
			SecretKey: os.Getenv("BUCKET_SECRET_KEY"),
			Endpoint:  os.Getenv("BUCKET_ENDPOINT"),
			Buckets: BucketName{
				LogoPic:             os.Getenv("LOGO_PIC_BUCKET_NAME"),
				ParkingDoc:          os.Getenv("PARKING_DOC_BUCKET_NAME"),
				ParkingPic:          os.Getenv("PARKING_PIC_BUCKET_NAME"),
				RegistrationCardPic: os.Getenv("REGISTRATION_CARD_PIC_BUCKET_NAME"),
			},
		},
		HashIDConfig{
			Salt:      os.Getenv("HASHID_SALT"),
			MinLength: getEnvInt("HASHID_MIN_LENGTH", 8),
		},
	}
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	i, err := strconv.Atoi(value)

	if err != nil {
		return defaultValue
	}

	return i
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)

	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return time.Duration(i) * time.Second
}
