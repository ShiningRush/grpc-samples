package access

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	Info("failed to fetch URL", "ahahaha", 321)

	assert.NoError(t, nil, "write log shout not have error")
}
