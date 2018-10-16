package consul

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {

}

func TestRegisterToConsul(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterToConsul("192.168.8.6:8500", "192.168.0.69:5053", "UnitTestService")
	}, "RegisterToConsul should not panic")
}

func TestDeregisterFromConsul(t *testing.T) {
	assert.NotPanics(t, func() {
		DeregisterFromConsul("192.168.8.6:8500", "192.168.0.69:5053", "UnitTestService")
	}, "RegisterToConsul should not panic")
}
