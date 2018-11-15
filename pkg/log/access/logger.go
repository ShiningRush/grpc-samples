package access

import (
	"gitlab.followme.com/FollowmeGo/golib/log"
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

func init() {
	Logger = log.NewLogger("access-%y-%m-%d.log", "../logs", "info")
}
