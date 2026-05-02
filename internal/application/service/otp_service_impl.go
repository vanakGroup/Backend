package service

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/redis/go-redis/v9"
)

type OtpService struct {
	logger logging.Logger
	cfg    *bootstrap.Config
	redis  *redis.Client
}
