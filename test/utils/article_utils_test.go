package utils_test

import (
	"testing"

	"github.com/Coteh/gyroid/lib/utils"
	"github.com/stretchr/testify/assert"
)

func TestCalculateExpectedReadTimeCalculatesExpectedReadMinutes(t *testing.T) {
	words := 1481
	expectedWPM := 5
	assert.Equal(t, expectedWPM, utils.CalculateExpectedReadTime(words))
}
