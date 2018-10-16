package access

import (
	"github.com/shiningrush/grpc-samples/pkg/log/io"

	"gitlab.followme.com/FollowmeGo/utils/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	l *zap.SugaredLogger
)

func init() {
	f := io.NewDailyFileHandler()
	f.SetName("normal-%y-%m-%d.log")
	// config encoder config
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeLevel = zapcore.CapitalLevelEncoder
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	// config core
	c := zapcore.AddSync(f)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(ec), c, zap.DebugLevel)
	l = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	).Sugar()

	l.With(zap.Int("pid", env.Pid))
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
