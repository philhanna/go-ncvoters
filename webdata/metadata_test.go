package webdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLayout(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{"bogus", "BOGUS", "", true},
		{"good", URL, "/tmp/voter_layout.txt", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			have, err := GetLayout(tt.url)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, have)
			}
		})
	}
}
