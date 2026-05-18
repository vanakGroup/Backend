package logging_zap

import (
	"sync"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/logging/logtypes"
	logging_help "github.com/Vanak/Backend/internal/infrastructure/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once
var zapSinLogger *zap.SugaredLogger

var logLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

type zapLogger struct {
	cfg    *bootstrap.Config
	logger *zap.SugaredLogger
}

func NewZapLogger(cfg *bootstrap.Config) *zapLogger {
	logger := &zapLogger{cfg: cfg}
	logger.Init()
	return logger
}

func (zl *zapLogger) getLogLevel() zapcore.Level {
	Level, exits := logLevelMap[zl.cfg.Env.Logger.Level]

	if !exits {
		return zapcore.DebugLevel
	}

	return Level
}

func (zl *zapLogger) Init() {
	once.Do(func() {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   zl.cfg.Env.Logger.FilePath,
			MaxSize:    1,
			MaxAge:     5,
			LocalTime:  true,
			MaxBackups: 10,
			Compress:   true,
		})

		config := zap.NewProductionEncoderConfig()

		config.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			w,
			zl.getLogLevel(),
		)

		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()

		logger.With("AppName", "MyApp", "LoggerName", "Zerolog")
		zapSinLogger = logger
	})

	zl.logger = zapSinLogger
}

func prepareLogKeys(extra map[logtypes.ExtraKey]interface{}, cat logtypes.Category, sub logtypes.SubCategory) []interface{} {
	if extra == nil {
		extra = make(map[logtypes.ExtraKey]interface{}, 0)
	}

	extra["Category"] = cat
	extra["SubCategory"] = sub

	params := logging_help.MapToZapParams(extra)

	return params
}
func (zl *zapLogger) Debug(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{}) {
	params := prepareLogKeys(extra, cat, sub)
	zl.logger.Debugw(msg, params...)
}

func (zl *zapLogger) Debugf(template string, args ...interface{}) {
	zl.logger.Debugf(template, args...)
}
func (zl *zapLogger) Info(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{}) {
	params := prepareLogKeys(extra, cat, sub)
	zl.logger.Infow(msg, params...)
}

func (zl *zapLogger) Infof(template string, args ...interface{}) {
	zl.logger.Infof(template, args...)
}

func (zl *zapLogger) Warn(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{}) {
	params := prepareLogKeys(extra, cat, sub)
	zl.logger.Warnw(msg, params...)
}

func (zl *zapLogger) Warnf(template string, args ...interface{}) {
	zl.logger.Warnf(template, args...)
}
func (zl *zapLogger) Error(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{}) {
	params := prepareLogKeys(extra, cat, sub)
	zl.logger.Errorw(msg, params...)
}

func (zl *zapLogger) Errorf(template string, args ...interface{}) {
	zl.logger.Errorf(template, args...)
}
func (zl *zapLogger) Fatal(cat logtypes.Category, sub logtypes.SubCategory, msg string, extra map[logtypes.ExtraKey]interface{}) {
	params := prepareLogKeys(extra, cat, sub)
	zl.logger.Fatalw(msg, params...)
}

func (zl *zapLogger) Fatalf(template string, args ...interface{}) {
	zl.logger.Fatalf(template, args...)
}
