package logging

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/logging/logtypes"
	logging_zap "github.com/Vanak/Backend/internal/infrastructure/logger/zap"
	logging_zero "github.com/Vanak/Backend/internal/infrastructure/logger/zero"
)

type Logger interface {
	Init()

	Debug(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{})
	Debugf(template string, args ...interface{})

	Info(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{})
	Infof(template string, args ...interface{})

	Warn(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{})
	Warnf(template string, args ...interface{})

	Error(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{})
	Errorf(template string, args ...interface{})

	Fatal(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{})
	Fatalf(template string, args ...interface{})
}

func NewLogger(cfg *bootstrap.Config) Logger {
	if cfg.Env.Logger.Logger == "zap" {
		return logging_zap.NewZapLogger(cfg)
	} else if cfg.Env.Logger.Logger == "zerolog" {
		return logging_zero.NewZeroLogger(cfg)
	}
	panic("could not create Logger")
}
