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

	tests := []struct {
		name     string
		bufsize  int
		limit    int
		expected []any
	}{
		{
			name:     "Unbuffered",
			bufsize:  0,
			limit:    3,
			expected: []any{1, 4, 9},
		},
		{
			name:     "Buffered",
			bufsize:  10,
			limit:    3,
			expected: []any{1, 4, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inch := make(chan any)
			var ouch chan any
			switch tt.bufsize {
			case 0:
				ouch = Map(square, inch)
			default:
				ouch = Map(square, inch, tt.bufsize)
			}

			go func() {
				for i := 1; i <= tt.limit; i++ {
					inch <- i
				}
				close(inch)
			}()

			for _, expected := range tt.expected {
				actual, ok := <-ouch
				assert.Truef(t, ok, "Not enough values found in channel")
				assert.Equal(t, expected, actual)
			}

			_, ok := <-ouch
			assert.Falsef(t, ok, "Too many values found in channel")

		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name   string
		inputs []any
		want   int
	}{
		{
			name:   "Happy path",
			inputs: []any{1, 2, 3},
			want:   6,
		},
		{
			name:   "Just one value",
			inputs: []any{1},
			want:   1,
		},
		{
			name:   "Empty",
			inputs: []any{},
			want:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan any)
			go func() {
				defer close(ch)
				for _, input := range tt.inputs {
					ch <- input
				}
			}()
			add := func(x, y any) any {
				return x.(int) + y.(int)
			}
			have := Reduce(add, ch, 0)
			assert.Equal(t, tt.want, have)
		})
	}
}
