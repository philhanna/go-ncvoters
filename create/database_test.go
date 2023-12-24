package create

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_estimatedNumberOfVoters(t *testing.T) {
	tests := []struct {
		name string
		size uint64
		want int64
	}{
		{"base", 3913556170, 8470544},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			have := estimatedNumberOfVoters(tt.size)
			diff := math.Abs(float64(want - have))
			delta := diff / float64(want)
			assert.True(t, delta < 3e-4)
		})
	}
}
