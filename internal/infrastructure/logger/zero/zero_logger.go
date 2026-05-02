package logging_zero

import (
	"os"
	"sync"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/logging/logtypes"
	logging_help "github.com/Vanak/Backend/internal/infrastructure/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once
var zeroSinLogger *zerolog.Logger

var logLevelMap = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
}

type zeroLogger struct {
	cfg    *bootstrap.Config
	logger *zerolog.Logger
}

func NewZeroLogger(cfg *bootstrap.Config) *zeroLogger {
	logger := &zeroLogger{cfg: cfg}
	logger.Init()
	return logger
}

func (zl *zeroLogger) getLogLevel() zerolog.Level {
	Level, exits := logLevelMap[zl.cfg.Env.Logger.Level]

	if !exits {
		return zerolog.DebugLevel
	}

	return Level
}

func (zl *zeroLogger) Init() {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.SetGlobalLevel(zl.getLogLevel())

		file, err := os.OpenFile(zl.cfg.Env.Logger.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

		if err != nil {
			panic("could not open the file")
		}

		logger := zerolog.New(file).
			With().
			Timestamp().
			Str("AppName", "MyApp").
			Str("LoggerName", "Zerolog").
			Logger()

		zeroSinLogger = &logger

	})

	zl.logger = zeroSinLogger
}

// func prepareLogKeys(extra map[logtypes.ExtraKey]interface{}, cat logtypes.Category, sub logtypes.SubCategory) []interface{} {
// 	if extra == nil {
// 		extra = make(map[logtypes.ExtraKey]interface{}, 0)
// 	}

// 	extra["Category"] = cat
// 	extra["SubCategory"] = sub

// 	params := logging_help.MapToZeroParams(extra)

//		return params
//	}

func (zl *zeroLogger) Debug(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{}) {
	zl.logger.
		Debug().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logging_help.MapToZeroParams(extra)).
		Msg(msg)
}

func (zl *zeroLogger) Debugf(template string, args ...interface{}) {
	zl.logger.
		Debug().
		Msgf(template, args...)
}

func (zl *zeroLogger) Info(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{}) {
	zl.logger.
		Info().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logging_help.MapToZeroParams(extra)).
		Msg(msg)
}

func (zl *zeroLogger) Infof(template string, args ...interface{}) {
	zl.logger.
		Info().
		Msgf(template, args...)
}

func (zl *zeroLogger) Warn(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{}) {
	zl.logger.
		Warn().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logging_help.MapToZeroParams(extra)).
		Msg(msg)
}

func (zl *zeroLogger) Warnf(template string, args ...interface{}) {
	zl.logger.
		Warn().
		Msgf(template, args...)
}

func (zl *zeroLogger) Error(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{}) {
	zl.logger.
		Error().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logging_help.MapToZeroParams(extra)).
		Msg(msg)
}

func (zl *zeroLogger) Errorf(template string, args ...interface{}) {
	zl.logger.
		Error().
		Msgf(template, args...)
}

func (zl *zeroLogger) Fatal(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{}) {
	zl.logger.
		Error().
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logging_help.MapToZeroParams(extra)).
		Msg(msg)
}

func (zl *zeroLogger) Fatalf(template string, args ...interface{}) {
	zl.logger.
		Error().
		Msgf(template, args...)
}
