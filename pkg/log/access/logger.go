package access

import (
	"github.com/shiningrush/grpc-samples/pkg/log/io"

	"gitlab.followme.com/FollowmeGo/utils/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

func init() {
	f := io.NewDailyFileHandler()
	f.SetName("access-%y-%m-%d.log")
	// config encoder config
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeLevel = zapcore.CapitalLevelEncoder
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	// config core
	c := zapcore.AddSync(f)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(ec), c, zap.DebugLevel)
	Logger = zap.New(
		core,
	)

	Logger.With(zap.Int("pid", env.Pid))
}
