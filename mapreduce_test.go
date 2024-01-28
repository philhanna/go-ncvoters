package goncvoters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	square := func(input any) any {
		n := input.(int)
		return n * n
	}

	inch := make(chan any)
	inputs := []any{1, 2, 3, 4, 5}
	ouch := Map(square, inch)

	go func() {
		for _, input := range inputs {
			inch <- input
		}
		close(inch)
	}()

	for _, expected := range []any{1, 4, 9, 16, 25} {
		actual, ok := <-ouch
		assert.True(t, ok)
		assert.Equal(t, expected, actual)
	}

	if _, ok := <-ouch; ok {
		t.Fatal("Received more outputs than expected")
	}
}
