package logs

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger 接口
type Logger interface {
	CtxDebug(ctx context.Context, format string, v ...interface{})
	CtxInfo(ctx context.Context, format string, v ...interface{})
	CtxWarn(ctx context.Context, format string, v ...interface{})
	CtxError(ctx context.Context, format string, v ...interface{})
	CtxFatal(ctx context.Context, format string, v ...interface{})
}

// Logger 默认实现依赖于 logrus 实现
type DefaultLogger struct {
	log *logrus.Logger
}

var (
	defaultLogger     *DefaultLogger // 默认的日志实例
	defaultLoggerOnce sync.Once      // 默认的日志实例信号量
)

// Once 信号量实现的单例模式
func DefaultLoggerInstance() *DefaultLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = &DefaultLogger{
			log: logrus.New(),
		}
		//defaultLogger.log.SetReportCaller(true) // 开启日志输出文件名和行号
	})
	return defaultLogger
}

// 所有的实例函数实质为类函数
func (*DefaultLogger) GetLogrus() *logrus.Logger {
	return defaultLogger.log
}

func (*DefaultLogger) CtxFatal(ctx context.Context, format string, v ...interface{}) {
	l := defaultLogger.log.WithField(LEVEL_KEY, FATAL)
	if logId := ctx.Value(LOGID_KEY); logId != nil {
		l = l.WithField(LOGID_KEY, logId)
	}
	l.Fatalf(format, v...)
}

func (*DefaultLogger) CtxError(ctx context.Context, format string, v ...interface{}) {
	l := defaultLogger.log.WithField(LEVEL_KEY, ERROR)
	if logId := ctx.Value(LOGID_KEY); logId != nil {
		l = l.WithField(LOGID_KEY, logId)
	}
	l.Errorf(format, v...)
}

func (*DefaultLogger) CtxWarn(ctx context.Context, format string, v ...interface{}) {
	l := defaultLogger.log.WithField(LEVEL_KEY, WARN)
	if logId := ctx.Value(LOGID_KEY); logId != nil {
		l = l.WithField(LOGID_KEY, logId)
	}
	l.Warnf(format, v...)
}

func (*DefaultLogger) CtxInfo(ctx context.Context, format string, v ...interface{}) {
	l := defaultLogger.log.WithField(LEVEL_KEY, INFO)
	if logId := ctx.Value(LOGID_KEY); logId != nil {
		l = l.WithField(LOGID_KEY, logId)
	}
	//TODO info 日志启用协程异步是否可以加速
	//go l.Infof(format, v...)
	l.Infof(format, v...)
}

func (*DefaultLogger) CtxDebug(ctx context.Context, format string, v ...interface{}) {
	l := defaultLogger.log.WithField(LEVEL_KEY, DEBUG)
	if logId := ctx.Value(LOGID_KEY); logId != nil {
		l = l.WithField(LOGID_KEY, logId)
	}
	//TODO debug 日志启用协程异步是否可以加速
	//go l.Debugf(format, v...)
	l.Debugf(format, v...)
}
