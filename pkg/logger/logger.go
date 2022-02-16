package logger

import (
	"os"

	"github.com/EZ4BRUCE/athena-server/pkg/setting"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(loggerSetting *setting.LOGSettingS) *zap.SugaredLogger {
	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		getLogWriter(loggerSetting),
	)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(loggerSetting *setting.LOGSettingS) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   loggerSetting.LogSavePath + loggerSetting.LogFileName + loggerSetting.LogFileExt,
		MaxSize:    loggerSetting.MaxSize,
		MaxBackups: loggerSetting.MaxBackups,
		MaxAge:     loggerSetting.MaxAge,
		Compress:   loggerSetting.Compress,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
