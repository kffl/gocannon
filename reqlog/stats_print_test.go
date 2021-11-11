package reqlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	durations := []float64{123456789, 999999, 1, 1100, 0.5, -1}
	expected := []string{"123.4568ms", "999.9990μs", "1.0000ns", "1.1000μs", "0.5000ns", "-"}

	var obtained []string

	for _, duration := range durations {
		obtained = append(obtained, formatDuration(duration))
	}

	assert.ElementsMatch(t, obtained, expected)
}
