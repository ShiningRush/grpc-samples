package io

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {

}

func TestWrite(t *testing.T) {
	f := NewDailyFileHandler()
	_, err := f.Write([]byte("Hello"))
	assert.NoError(t, err, "write file shout not have error")
}

func TestDailyName(t *testing.T) {
	f := NewDailyFileHandler().dailyName()
	assert.Regexp(t, "logs-.*.log", f, "Default daily name should start with logs")
}
