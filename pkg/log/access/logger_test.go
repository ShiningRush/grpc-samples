package access

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	Logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	assert.NoError(t, nil, "write log shout not have error")
}
