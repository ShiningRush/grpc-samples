package access

import (
	"gitlab.followme.com/FollowmeGo/golib/log"
	"go.uber.org/zap"
)

var (
	l *zap.SugaredLogger
)

func init() {
	l = log.NewLogger("normal-%y-%m-%d.log", "../logs", "info").Sugar()
}

func Debug(v ...interface{}) {
	l.Debug(v)
}

func Debugf(s string, v ...interface{}) {
	l.Debugf(s, v)
}

func Info(v ...interface{}) {
	l.Info(v)
}

func Infof(s string, v ...interface{}) {
	l.Infof(s, v)
}

func Warning(v ...interface{}) {
	l.Warn(v)
}

func Warningf(s string, v ...interface{}) {
	l.Warnf(s, v)
}

func Error(v ...interface{}) {
	l.Error(v)
}

func Errorf(s string, v ...interface{}) {
	l.Errorf(s, v)
}
