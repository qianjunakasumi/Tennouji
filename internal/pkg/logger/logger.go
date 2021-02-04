package logger

import (
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

// Init 初始化 logger
func Init() {

	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.RFC3339TimeEncoder
	conf.EncodeCaller = zapcore.FullCallerEncoder

	lumber := &lumberjack.Logger{
		Filename: time.Now().Format(".tennouji/logs/2006-01-02.log"), MaxAge: 7, LocalTime: true, Compress: true,
	}

	logger = zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(conf), zapcore.AddSync(lumber), zapcore.DebugLevel),
		zap.AddCaller(),
	)
	job(lumber)

}

// job 日志轮转
func job(l *lumberjack.Logger) {

	c := cron.New()
	_, _ = c.AddFunc("0 0 * * *", func() {
		if err := l.Rotate(); err != nil {
			logger.Named("日志").Warn("无法轮转", zap.Error(err))
			return
		}
	})
	c.Start()

}
