package logger

import (
	"acat/setting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

var lg *zap.Logger

// InitLogger 初始化Logger
func Init(cfg *setting.LogConf, mode string) (err error) {
	// 获取日志写入器（Writer）
	writeSyncer := getLogWriter(cfg.LogFileName, cfg.LogMaxSize, cfg.LogMaxBackups, cfg.LogMaxAge)
	// 获取编码器（Encoder）
	encoder := getEncoder()

	//解析日志级别
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		return
	}
	var core zapcore.Core
	if mode == "dev" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),                    //输出到文件
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), l), //输出到终端
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg)
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter(...)它的作用是创建一个支持日志轮转（log rotation）的写入器，用于将日志安全地写入文件，并自动管理旧日志文件
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	dir := filepath.Dir(filename)
	_ = os.MkdirAll(dir, 0755)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "logger/" + filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	// 包装成 WriteSyncer
	return zapcore.AddSync(lumberJackLogger)
}
