package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZapLogger_Basic(t *testing.T) {
	// Simple test to ensure logger can be created
	logger := New(INFO)
	assert.NotNil(t, logger)
	assert.IsType(t, &ZapLogger{}, logger)
}
