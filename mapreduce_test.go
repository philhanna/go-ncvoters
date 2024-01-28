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
	ouch := Map(square, inch)

	inputs := []any{1, 2, 3}
	expectedOutputs := []any{1, 4, 9}

	go func() {
		defer close(ouch)
		for _, output := range expectedOutputs {
			ouch <- output
		}
	}()

	result := Map(square, inch)
	
	go func() {
		for _, input := range inputs {
			inch <- input
		}
		close(inch)
	}()

	for expected := range ouch {
		actual, ok := <- result
		assert.True(t, ok)
		assert.Equal(t, expected, actual)
	}

	if _, ok := <-result; ok {
		t.Fatal("Received more outputs than expected")
	}
}
