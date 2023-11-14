package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	defaultConfig := defaultConfig()
	assert.IsType(t, &Configuration{}, defaultConfig)
}
