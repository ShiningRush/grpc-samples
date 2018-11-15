package log

import (
	"gitlab.followme.com/FollowmeGo/golib/log/io"
	"gitlab.followme.com/FollowmeGo/utils/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(name, dir, level string) *zap.Logger {
	f := io.NewDailyFileHandler()
	if name != "" {
		f.SetName(name)
	}

	if dir != "" {
		f.SetDirectory(dir)
	}

	// config encoder config
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeLevel = zapcore.CapitalLevelEncoder
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	// config core
	ll := zap.DebugLevel
	if level != "" {
		switch level {
		case "info":
			ll = zap.InfoLevel
		case "warn":
			ll = zap.WarnLevel
		case "error":
			ll = zap.ErrorLevel
		}
	}
	c := zapcore.AddSync(f)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(ec), c, ll)
	l := zap.New(
		core,
	)

	l.With(zap.Int("pid", env.Pid))

	return l
}
